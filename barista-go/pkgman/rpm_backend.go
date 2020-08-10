package pkgman

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alecthomas/repr"
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
		repomd, err := http.Get(fmt.Sprintf("%s/repodata/repomd.xml", repo.BaseURL))
		if err != nil {
			return err
		}
		data, err := ioutil.ReadAll(repomd.Body)
		if err != nil {
			return err
		}
		var md RepoMD
		println(string(data))
		err = xml.Unmarshal(data, &md)
		if err != nil {
			return err
		}
		repr.Println(md)
	}
	return nil
}
