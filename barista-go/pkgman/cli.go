package pkgman

import "github.com/urfave/cli/v2"

func Cli(c *cli.Context) {
	data := RPMBackend{
		rDistros: map[string]RPMDistro{
			"fedora": {
				Name: "fedora",
				Repos: []RPMRepo{
					{
						Name:    "core",
						BaseURL: "http://download.fedoraproject.org/pub/fedora/linux/releases/32/Everything/x86_64/os",
					},
				},
			},
		},
	}
	err := data.Refresh("fedora")
	if err != nil {
		panic(err)
	}
}
