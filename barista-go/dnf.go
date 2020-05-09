package barista

import (
	"fmt"
	"strings"

	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/dustin/go-humanize"
)

type Package struct {
	name        string
	desc        string
	vers        string
	downsize    string
	installsize string
	url         string
}

type Distro struct {
	matches      []string
	displayName  string
	queryKitName string
	iconURL      string
	colour       int
}

var Distros []Distro = []Distro{
	{
		displayName:  "openSUSE Tumbleweed",
		queryKitName: "tumbleweed",
		matches:      []string{"tumbleweed", "tw", "opensuse", "os"},
		iconURL:      "https://en.opensuse.org/images/c/cd/Button-colour.png",
		colour:       0x73ba25,
	},
	{
		displayName:  "openSUSE Leap",
		queryKitName: "leap",
		matches:      []string{"leap", "opensuse-leap", "os-leap"},
		iconURL:      "https://en.opensuse.org/images/c/cd/Button-colour.png",
		colour:       0x73ba25,
	},
	{
		displayName:  "Fedora",
		queryKitName: "fedora",
		matches:      []string{"fedora"},
		iconURL:      "https://fedoraproject.org/w/uploads/archive/e/e5/20110717032101%21Fedora_infinity.png",
		colour:       0x0b57a4,
	},
	{
		displayName:  "Mageia",
		queryKitName: "mageia",
		matches:      []string{"mageia"},
		iconURL:      "https://pbs.twimg.com/profile_images/553311070215892992/lf8QV6oJ_400x400.png",
		colour:       0x2397d4,
	},
	{
		displayName:  "OpenMandriva",
		queryKitName: "openmandriva",
		matches:      []string{"openmandriva"},
		iconURL:      "https://pbs.twimg.com/profile_images/1140547712208822272/dG9610ZK_400x400.jpg",
		colour:       0x40a5da,
	},
	{
		displayName:  "CentOS",
		queryKitName: "centos",
		matches:      []string{"centos"},
		iconURL:      "https://upload.wikimedia.org/wikipedia/commons/thumb/6/63/CentOS_color_logo.svg/1024px-CentOS_color_logo.svg.png",
		colour:       0x951C7A,
	},
	{
		displayName:  "Packman for openSUSE Leap",
		queryKitName: "packman-leap",
		matches:      []string{"packman-leap", "pm-l", "pm-leap"},
		colour:       0x2C74CC,
	},
	{
		displayName:  "Packman for openSUSE Tumbleweed",
		queryKitName: "packman-tumbleweed",
		matches:      []string{"packman", "pm-tw", "pm-tumbleweed"},
		colour:       0x2C74CC,
	},
	{
		displayName:  "RPMFusion",
		queryKitName: "rpmfusion",
		iconURL:      "https://rpmfusion.org/moin_static1910/rpmfusion/logo.png",
		matches:      []string{"rpmfusion", "rf", "rpmf"},
		colour:       0x0855A7,
	},
}

func distroList() string {
	var distros []string
	for _, dist := range Distros {
		distros = append(distros, dist.matches[0])
	}
	return strings.Join(distros, ", ")
}

func toPackageList(pkgs [][]interface{}) []Package {
	var ret []Package
	for _, pkg := range pkgs {
		ret = append(ret, Package{
			name:        pkg[0].(string),
			desc:        pkg[1].(string),
			vers:        pkg[2].(string),
			downsize:    humanize.Bytes(uint64(pkg[3].(int32))),
			installsize: humanize.Bytes(uint64(pkg[4].(int32))),
			url:         pkg[5].(string),
		})
	}
	return ret
}

func pkgListToUnionEmbed(pkgs []Package, distro Distro, c commandlib.Context) commandlib.UnionEmbed {
	var embeds []commandlib.Embed
	var tableData [][]string
	for _, pkg := range pkgs {
		embeds = append(embeds, commandlib.Embed{
			Colour: distro.colour,
			Title: commandlib.EmbedHeader{
				Text: pkg.name,
				URL:  pkg.url,
			},
			Header: commandlib.EmbedHeader{
				Text: fmt.Sprintf(l10n(c, "%s Package Search"), distro.displayName),
				Icon: distro.iconURL,
			},
			Body: pkg.desc,
			Fields: []commandlib.EmbedField{
				{
					Title:  l10n(c, "Version"),
					Body:   pkg.vers,
					Inline: true,
				},
				{
					Title:  l10n(c, "Download Size"),
					Body:   pkg.downsize,
					Inline: true,
				},
				{
					Title:  l10n(c, "Install Size"),
					Body:   pkg.downsize,
					Inline: true,
				},
			},
		})
		tableData = append(tableData, []string{pkg.name, pkg.vers, pkg.desc, pkg.downsize, pkg.installsize})
	}
	return commandlib.UnionEmbed{
		EmbedList: commandlib.EmbedList{
			ItemTypeName: l10n(c, "Package"),
			Embeds:       embeds,
		},
		EmbedTable: commandlib.EmbedTable{
			Heading:  fmt.Sprintf(l10n(c, "Search results for %s in %s"), c.Content(), distro.displayName),
			Subtitle: fmt.Sprintf(l10n(c, "%d packages found"), len(tableData)),
			Headers:  []string{l10n(c, "Name"), l10n(c, "Version"), l10n(c, "Description"), l10n(c, "Download Size"), l10n(c, "Install Size")},
			Data:     tableData,
		},
	}
}

func resolveDistro(name string) (Distro, bool) {
	var distro Distro
	set := false

	for _, dist := range Distros {
		for _, match := range dist.matches {
			if strings.ToLower(name) == strings.ToLower(match) {
				distro = dist
				set = true
			}
		}
	}

	return distro, set
}
