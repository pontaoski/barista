package barista

import (
	"fmt"

	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/godbus/dbus"
)

func init() {
	commandlib.RegisterCommand(commandlib.Command{
		Name:  "DNF Package Search",
		Usage: "Search packages with DNF",
		Match: [][]string{
			{"dnf", "search"},
			{"dnf", "se"},
		},
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
	if !c.IsFlagSet("distro") {
		c.SendMessage("primary", commandlib.ErrorEmbed("Please provide a distro with the flag `--distro`."))
		return
	}
	if c.NArgs() < 1 {
		c.SendMessage("primary", commandlib.ErrorEmbed("Please provide a search term."))
		return
	}
	var distro Distro
	var ok bool
	if distro, ok = resolveDistro(c.FlagValue("distro")); !ok {
		c.SendMessage(
			"primary",
			commandlib.ErrorEmbed("Please provide a distro from the following list: `"+distroList()+"`"),
		)
		return
	}
	conn, err := dbus.SessionBus()
	if err != nil {
		c.SendMessage("primary", commandlib.ErrorEmbed("There was an issue connecting to QueryKit, the package search service."))
		return
	}
	var pkgs [][]interface{}
	obj := conn.Object("com.github.Appadeia.QueryKit", "/com/github/Appadeia/QueryKit")
	err = obj.Call("com.github.Appadeia.QueryKit.SearchPackages", 0, c.Content(), distro.queryKitName).Store(&pkgs)
	if err != nil {
		c.SendMessage("primary", commandlib.ErrorEmbed("There was an issue searching for packages: "+err.Error()))
	}
	packages := toPackageList(pkgs)
	var embeds []commandlib.Embed
	var tableData [][]string
	for _, pkg := range packages {
		embeds = append(embeds, commandlib.Embed{
			Colour: distro.colour,
			Title: commandlib.EmbedHeader{
				Text: pkg.name,
				URL:  pkg.url,
			},
			Header: commandlib.EmbedHeader{
				Text: fmt.Sprintf("%s Package Search", distro.displayName),
				Icon: distro.iconURL,
			},
			Body: pkg.desc,
			Fields: []commandlib.EmbedField{
				{
					Title:  "Version",
					Body:   pkg.vers,
					Inline: true,
				},
				{
					Title:  "Download Size",
					Body:   pkg.downsize,
					Inline: true,
				},
				{
					Title:  "Install Size",
					Body:   pkg.downsize,
					Inline: true,
				},
			},
		})
		tableData = append(tableData, []string{pkg.name, pkg.vers, pkg.desc, pkg.downsize, pkg.installsize})
	}
	c.SendMessage("primary", commandlib.UnionEmbed{
		EmbedList: commandlib.EmbedList{
			ItemTypeName: "Package",
			Embeds:       embeds,
		},
		EmbedTable: commandlib.EmbedTable{
			Heading:  fmt.Sprintf("Search results for %s in %s", c.Content(), distro.displayName),
			Subtitle: fmt.Sprintf("%d packages found", len(tableData)),
			Headers:  []string{"Name", "Version", "Description", "Download Size", "Install Size"},
			Data:     tableData,
		},
	})
}
