package barista

import (
	"fmt"

	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/util"
	"github.com/godbus/dbus"
)

func init() {
	commandlib.RegisterCommand(commandlib.Command{
		Name:  "DNF Repository Query",
		Usage: "Query repositories with DNF",
		ID:    "dnf-repoquery",
		Examples: `dnf repoquery -d=opensuse --provides dnf
dnf repoquery -d=fedora --whatprovides cmake(KF5Kirigami2)
dnf repoquery -d=mageia -l chromium`,
		Match: [][]string{
			{"dnf", "rq"},
			{"dnf", "repoquery"},
		},
		Flags: commandlib.FlagList{
			commandlib.StringFlag{
				LongFlag:  "distro",
				ShortFlag: "d",
				FlagUsage: "Which distro you want to query repos for",
				Value:     "",
			},
			// What does X flags
			commandlib.StringFlag{LongFlag: "whatconflicts"},
			commandlib.StringFlag{LongFlag: "whatobsoletes"},
			commandlib.StringFlag{LongFlag: "whatprovides"},
			commandlib.StringFlag{LongFlag: "whatrecommends"},
			commandlib.StringFlag{LongFlag: "whatenhances"},
			commandlib.StringFlag{LongFlag: "whatsupplements"},
			commandlib.StringFlag{LongFlag: "whatsuggests"},
			commandlib.StringFlag{LongFlag: "whatrequires"},
			commandlib.StringFlag{LongFlag: "file", ShortFlag: "f"},
			// Relationship flags
			commandlib.BoolFlag{LongFlag: "provides"},
			commandlib.BoolFlag{LongFlag: "requires"},
			commandlib.BoolFlag{LongFlag: "recommends"},
			commandlib.BoolFlag{LongFlag: "suggests"},
			commandlib.BoolFlag{LongFlag: "supplements"},
			commandlib.BoolFlag{LongFlag: "enhances"},
			commandlib.BoolFlag{LongFlag: "conflicts"},
			commandlib.BoolFlag{LongFlag: "obsoletes"},
			// Flag to list files of a package
			commandlib.BoolFlag{LongFlag: "list", ShortFlag: "l"},
		},
		Action: DnfRepoquery,
	})
}

var whatFlags = []string{"file", "whatconflicts", "whatobsoletes", "whatprovides", "whatrecommends", "whatenhances", "whatsupplements", "whatsuggests", "whatrequires"}
var relFlags = []string{"provides", "requires", "recommends", "suggests", "supplements", "enhances", "conflicts", "obsoletes"}

func DnfRepoquery(c commandlib.Context) {
	def := schemas["default-distro"].ReadValue(c)

	if def == "" {
		def = c.FlagValue("distro")
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
	obj := conn.Object("com.github.Appadeia.QueryKit", "/com/github/Appadeia/QueryKit")

	if c.AnySet(relFlags...) || c.IsFlagSet("list") {
		if c.Content() == "" {
			c.SendMessage("primary", commandlib.ErrorEmbed("Please provide a query with your flag."))
			return
		}
	}
	if c.AnySet(whatFlags...) {
		var pkgs [][]interface{}
		queries := make(map[string]string)
		for _, flag := range whatFlags {
			if c.IsFlagSet(flag) {
				queries[flag] = c.FlagValue(flag)
			}
		}
		err = obj.Call("com.github.Appadeia.QueryKit.QueryRepo", 0, queries, distro.queryKitName).Store(&pkgs)
		if err != nil {
			util.OutputError(err)
			c.SendMessage("primary", commandlib.ErrorEmbed("There was an issue searching for packages: "+err.Error()))
			return
		}
		if len(pkgs) == 0 {
			c.SendMessage("primary", commandlib.ErrorEmbed("No packages were found."))
			return
		}
		c.SendMessage("primary", pkgListToUnionEmbed(toPackageList(pkgs), distro, c))
	} else if c.AnySet(relFlags...) {
		qkAction := c.ChoiceFlags(relFlags...)
		var reldeps []string

		err = obj.Call("com.github.Appadeia.QueryKit.QueryRepoPackage", 0, c.Content(), qkAction, distro.queryKitName).Store(&reldeps)
		if err != nil {
			util.OutputError(err)
			c.SendMessage("primary", commandlib.ErrorEmbed("There was an error querying packages: "+err.Error()))
			return
		}

		titleText := fmt.Sprintf("Query %s for %s", qkAction, c.Content())
		header := commandlib.EmbedHeader{
			Text: fmt.Sprintf("%s Repoquery", distro.displayName),
			Icon: distro.iconURL,
		}

		embed, list := commandlib.PaginateList(c, reldeps)
		if embed != nil {
			embed.Title.Text = titleText
			embed.Colour = distro.colour
			embed.Header = header
			c.SendMessage("primary", *embed)
		} else {
			var embedList commandlib.EmbedList
			embedList.ItemTypeName = "Page"
			for _, embed := range *list {
				embedList.Embeds = append(embedList.Embeds, commandlib.Embed{
					Title: commandlib.EmbedHeader{
						Text: titleText,
					},
					Colour: distro.colour,
					Body:   embed.Body,
					Header: header,
				})
			}
			c.SendMessage("primary", embedList)
		}
	} else if c.IsFlagSet("list") {
		var files []string
		err = obj.Call("com.github.Appadeia.QueryKit.ListFiles", 0, c.Content(), distro.queryKitName).Store(&files)
		if err != nil {
			util.OutputError(err)
			c.SendMessage("primary", commandlib.ErrorEmbed("There was an issue listing files: "+err.Error()))
			return
		}
		titleText := fmt.Sprintf("Filelist for %s", c.Content())
		header := commandlib.EmbedHeader{
			Text: fmt.Sprintf("%s Repoquery", distro.displayName),
			Icon: distro.iconURL,
		}

		embed, list := commandlib.PaginateList(c, files)
		if embed != nil {
			embed.Title.Text = titleText
			embed.Colour = distro.colour
			embed.Header = header
			c.SendMessage("primary", *embed)
		} else {
			var embedList commandlib.EmbedList
			embedList.ItemTypeName = "Page"
			for _, embed := range *list {
				embedList.Embeds = append(embedList.Embeds, commandlib.Embed{
					Title: commandlib.EmbedHeader{
						Text: titleText,
					},
					Colour: distro.colour,
					Body:   embed.Body,
					Header: header,
				})
			}
			c.SendMessage("primary", embedList)
		}
	} else {
		c.SendMessage("primary", commandlib.ErrorEmbed("Please provide a query."))
		return
	}
}
