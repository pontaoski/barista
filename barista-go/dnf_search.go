package barista

import (
	"fmt"
	"sort"

	"github.com/agext/levenshtein"
	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/log"
	"github.com/godbus/dbus"
)

func init() {
	commandlib.RegisterCommand(commandlib.Command{
		Name:  I18n("DNF Package Search"),
		Usage: I18n("Search packages with DNF"),
		ID:    "dnf-search",
		Matches: []string{
			"ilo, dnf search",
			"ilo, dnf se",
			"ilo, zypper search",
			"ilo, zypper se",
		},
		Examples: `ilo, dnf se -d=fedora chromium
ilo, dnf se -d=rpmfusion nvidia`,
		Flags: commandlib.FlagList{
			commandlib.StringFlag{
				LongFlag:  "distro",
				ShortFlag: "d",
				FlagUsage: I18n("Which distro you want to use."),
				Value:     "",
			},
		},
		Action: DnfSearch,
	})
}

func DnfSearch(c commandlib.Context) {
	def := schemas["default-distro"].ReadValue(c)
	if !c.IsFlagSet("distro") && def == "" {
		c.SendMessage("primary", commandlib.ErrorEmbed(c.I18n("Please provide a distro with the flag `--distro`.")))
		return
	}
	if def == "" {
		def = c.FlagValue("distro")
	}
	if c.NArgs() < 1 {
		c.SendMessage("primary", commandlib.ErrorEmbed(c.I18n("Please provide a search term.")))
		return
	}
	var distro Distro
	var ok bool
	if distro, ok = resolveDistro(def); !ok {
		c.SendMessage(
			"primary",
			commandlib.ErrorEmbed(fmt.Sprintf(c.I18n("Please provide a distro from the following list: %s"), distroList())),
		)
		return
	}
	conn, err := dbus.SystemBus()
	if err != nil {
		log.Error("%+v", err)
		c.SendMessage("primary", commandlib.ErrorEmbed(c.I18n("There was an issue connecting to QueryKit, the package search service.")))
		return
	}
	var pkgs [][]interface{}
	obj := conn.Object("com.github.Appadeia.QueryKit", "/com/github/Appadeia/QueryKit")
	err = obj.Call("com.github.Appadeia.QueryKit.SearchPackages", 0, c.Content(), distro.queryKitName).Store(&pkgs)
	if err != nil {
		log.Error("%+v", err)
		c.SendMessage("primary", commandlib.ErrorEmbed(c.I18n("There was an issue searching for packages: ")+err.Error()))
		return
	}
	if len(pkgs) == 0 {
		c.SendMessage("primary", commandlib.ErrorEmbed(c.I18n("No packages were found.")))
		return
	}
	packs := toPackageList(pkgs)
	sort.Slice(packs, func(i, j int) bool {
		return levenshtein.Distance(c.Content(), packs[i].name, nil) < levenshtein.Distance(c.Content(), packs[j].name, nil)
	})
	c.SendMessage("primary", pkgListToUnionEmbed(packs, distro, c))
}
