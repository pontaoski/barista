package barista

import (
	"strings"

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
