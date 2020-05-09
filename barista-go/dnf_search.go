package barista

import (
	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/util"
	"github.com/godbus/dbus"
)

func init() {
	commandlib.RegisterCommand(commandlib.Command{
		Name:  "DNF Package Search",
		Usage: "Search packages with DNF",
		ID:    "dnf-search",
		Match: [][]string{
			{"dnf", "search"},
			{"dnf", "se"},
			{"zypper", "search"},
			{"zypper", "se"},
		},
		Examples: `dnf se -d=fedora chromium
dnf se -d=rpmfusion nvidia`,
		Flags: commandlib.FlagList{
			commandlib.StringFlag{
				LongFlag:  "distro",
				ShortFlag: "d",
				FlagUsage: "Which distro you want to use.",
				Value:     "",
			},
		},
		Action: DnfSearch,
	})
}

func DnfSearch(c commandlib.Context) {
	def := schemas["default-distro"].ReadValue(c)
	if !c.IsFlagSet("distro") && def == "" {
		c.SendMessage("primary", commandlib.ErrorEmbed("Please provide a distro with the flag `--distro`."))
		return
	}
	if def == "" {
		def = c.FlagValue("distro")
	}
	if c.NArgs() < 1 {
		c.SendMessage("primary", commandlib.ErrorEmbed("Please provide a search term."))
		return
	}
	var distro Distro
	var ok bool
	if distro, ok = resolveDistro(def); !ok {
		c.SendMessage(
			"primary",
			commandlib.ErrorEmbed("Please provide a distro from the following list: `"+distroList()+"`"),
		)
		return
	}
	conn, err := dbus.SessionBus()
	if err != nil {
		util.OutputError(err)
		c.SendMessage("primary", commandlib.ErrorEmbed("There was an issue connecting to QueryKit, the package search service."))
		return
	}
	var pkgs [][]interface{}
	obj := conn.Object("com.github.Appadeia.QueryKit", "/com/github/Appadeia/QueryKit")
	err = obj.Call("com.github.Appadeia.QueryKit.SearchPackages", 0, c.Content(), distro.queryKitName).Store(&pkgs)
	if err != nil {
		util.OutputError(err)
		c.SendMessage("primary", commandlib.ErrorEmbed(l10n(c, "There was an issue searching for packages: ")+err.Error()))
		return
	}
	if len(pkgs) == 0 {
		c.SendMessage("primary", commandlib.ErrorEmbed(l10n(c, "No packages were found.")))
		return
	}
	c.SendMessage("primary", pkgListToUnionEmbed(toPackageList(pkgs), distro, c))
}
