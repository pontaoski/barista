package pkgman

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"runtime"
	"strings"

	"github.com/adrg/xdg"
	"github.com/appadeia/barista/barista-go/log"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/ulikunitz/xz"
)

var (
	chttp = httpcache.NewTransport(diskcache.New(xdg.CacheHome)).Client()
)

type RPMRepo struct {
	Name    string
	BaseURL string
}

type RPMDistro struct {
	Name  string
	Repos []RPMRepo
}

type RepoMD struct {
	Datas []struct {
		Type     string `xml:"type,attr"`
		Location struct {
			URL string `xml:"href,attr"`
		} `xml:"location"`
	} `xml:"data"`
}

type RPMEntry struct {
	Name    string `xml:"name,attr"`
	Kind    string `xml:"kind,attr"`
	Flags   string `xml:"flags,attr"`
	Epoch   string `xml:"epoch,attr"`
	Version string `xml:"ver,attr"`
	Release string `xml:"rel,attr"`
}

type Version struct {
	Epoch   string `xml:"epoch,attr"`
	Version string `xml:"ver,attr"`
	Release string `xml:"rel,attr"`
}

func (v Version) String() string {
	return fmt.Sprintf("%s:%s-%s", v.Epoch, v.Version, v.Release)
}

type RPMSize struct {
	Package   uint64 `xml:"package,attr"`
	Installed uint64 `xml:"installed,attr"`
	Archive   uint64 `xml:"archive,attr"`
}

type RPMLocation struct {
	URL string `xml:"href,attr"`
}

type RPMEntryList struct {
	Entries []RPMEntry `xml:"rpm__entry"`
}

func (r RPMEntryList) ToNames() (ret []string) {
	for _, ent := range r.Entries {
		ret = append(ret, ent.Name)
	}
	return
}

type RPMFormat struct {
	License     string       `xml:"rpm__license"`
	Requires    RPMEntryList `xml:"rpm__requires"`
	Obsoletes   RPMEntryList `xml:"rpm__obsoletes"`
	Provides    RPMEntryList `xml:"rpm__provides"`
	Recommends  RPMEntryList `xml:"rpm__recommends"`
	Supplements RPMEntryList `xml:"rpm__supplements"`
	Conflicts   RPMEntryList `xml:"rpm__conflicts"`
	Enhances    RPMEntryList `xml:"rpm__enhances"`
	Suggests    RPMEntryList `xml:"rpm__suggests"`
}

type RPMPackage struct {
	Name string `xml:"name"`
	Arch string `xml:"arch"`

	Version Version `xml:"version"`

	Summary     string `xml:"summary"`
	Description string `xml:"description"`
	URL         string `xml:"url"`

	Size RPMSize `xml:"size"`

	Location RPMLocation `xml:"location"`

	Data RPMFormat `xml:"format"`
}

type Primary struct {
	Package []RPMPackage `xml:"package"`
}

type PackageFileList struct {
	Name string   `xml:"name,attr"`
	File []string `xml:"file"`
}

type FileList struct {
	Packages []PackageFileList `xml:"package"`
}

func Get(url string) (data []byte, err error) {
	resp, err := chttp.Get(url)
	if err != nil {
		return
	}
	data, err = ioutil.ReadAll(resp.Body)

	log.Info("Fetching %s...", url)

	var reader io.Reader
	switch {
	case strings.HasSuffix(url, "gz"):
		reader, err = gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return
		}
	case strings.HasSuffix(url, "xz"):
		reader, err = xz.NewReader(bytes.NewReader((data)))
		if err != nil {
			return
		}
	}

	if reader != nil {
		log.Info("Decompressing %s...", url)
		data, err = ioutil.ReadAll(reader)
	}

	data = []byte(strings.ReplaceAll(string(data), "<rpm:", "<rpm__"))
	runtime.GC()
	data = []byte(strings.ReplaceAll(string(data), "</rpm:", "</rpm__"))
	runtime.GC()

	return
}

type PrimaryList struct {
}

func (r RPMDistro) isDistro() {}

type RPMBackend struct {
	rDistros map[string]RPMDistro
}

func (r *RPMBackend) Distros() (ret []string) {
	for _, distro := range r.rDistros {
		ret = append(ret, distro.Name)
	}
	return
}

func (r *RPMBackend) Refresh(target string) error {
	var distro RPMDistro
	var ok bool
	if distro, ok = r.rDistros[target]; !ok {
		return ErrBadDistro
	}
	for _, repo := range distro.Repos {
		data, err := Get(fmt.Sprintf("%s/repodata/repomd.xml", repo.BaseURL))
		if err != nil {
			return err
		}

		var md RepoMD
		err = xml.Unmarshal(data, &md)
		if err != nil {
			return err
		}

		var primaryURL, filelistURL string
		for _, item := range md.Datas {
			switch item.Type {
			case "primary":
				primaryURL = item.Location.URL
			case "filelists":
				filelistURL = item.Location.URL
			}
		}

		primaryData, err := Get(fmt.Sprintf("%s/%s", repo.BaseURL, primaryURL))
		if err != nil {
			return err
		}

		var pd Primary
		log.Info("Unmarshalling primary data...")
		err = xml.Unmarshal(primaryData, &pd)
		if err != nil {
			return err
		}

		filelistData, err := Get(fmt.Sprintf("%s/%s", repo.BaseURL, filelistURL))
		if err != nil {
			return err
		}
		var fl FileList

		log.Info("Unmarshalling file lists...")
		err = xml.Unmarshal(filelistData, &fl)
		if err != nil {
			return err
		}

		log.Info("Clearing distro %s...", target)
		err = ClearDistro(target)
		if err != nil {
			return err
		}
		log.Info("Inserting packages...")
		for _, pkg := range pd.Package {
			err := InsertPackage(DBPackage{
				Distro:       target,
				Name:         pkg.Name,
				Description:  pkg.Description,
				Version:      pkg.Version.String(),
				DownloadSize: pkg.Size.Archive,
				InstallSize:  pkg.Size.Installed,
				Relations: map[string][]string{
					"Conflicts":   pkg.Data.Conflicts.ToNames(),
					"Provides":    pkg.Data.Provides.ToNames(),
					"Requires":    pkg.Data.Requires.ToNames(),
					"Recommends":  pkg.Data.Recommends.ToNames(),
					"Suggests":    pkg.Data.Suggests.ToNames(),
					"Supplements": pkg.Data.Supplements.ToNames(),
					"Enhances":    pkg.Data.Enhances.ToNames(),
					"Obsoletes":   pkg.Data.Obsoletes.ToNames(),
				},
				Files: func() (ret []string) {
					for _, file := range fl.Packages {
						if file.Name == pkg.Name {
							return file.File
						}
					}
					return
				}(),
			})
			if err != nil {
				return err
			}
		}
	}
	log.Info("Refreshed distro %s", target)
	return nil
}
