package barista

import (
	"fmt"
	"strings"

	"github.com/appadeia/barista/barista-go/commandlib"
)

func init() {
	commandlib.RegisterCommand(commandlib.Command{
		Name:  "Help",
		Usage: "See help for commands.",
		ID:    "help",
		Matches: []string{
			"sudo help",
			"o help",
		},
		Action: Help,
	})
}

func CommandHelp(c commandlib.Context, command commandlib.Command) commandlib.Embed {
	return commandlib.Embed{
		Title: commandlib.EmbedHeader{
			Text: fmt.Sprintf("%s (%s)", c.I18n(command.Name), command.ID),
		},
		Body: c.I18n(command.Usage),
		Fields: func() []commandlib.EmbedField {
			ret := []commandlib.EmbedField{
				{
					Title: c.I18n("Aliases"),
					Body:  strings.Join(command.Matches, ", "),
				},
			}
			exmps := command.Examples
			if exmps != "" {
				ret = append(ret, commandlib.EmbedField{
					Title: c.I18n("Examples"),
					Body:  c.WrapCodeBlock(exmps),
				})
			}
			flags := command.Flags.GetFlagSet().FlagUsages()
			if flags != "" {
				ret = append(ret, commandlib.EmbedField{
					Title: c.I18n("Flags"),
					Body:  c.WrapCodeBlock(flags),
				})
			}
			return ret
		}(),
	}
}

func Help(c commandlib.Context) {
	var commandEmbeds []commandlib.Embed
	var tagEmbeds []commandlib.Embed
	for _, command := range commandlib.Commands() {
		if command.Hidden {
			continue
		}
		commandEmbeds = append(commandEmbeds, CommandHelp(c, command))
	}
	for _, tag := range commandlib.Tags() {
		tagEmbeds = append(tagEmbeds, commandlib.Embed{
			Title: commandlib.EmbedHeader{
				Text: fmt.Sprintf("%s (%s)", tag.Name, tag.ID),
			},
			Body: tag.Usage,
			Fields: []commandlib.EmbedField{
				{
					Title: c.I18n("Tags"),
					Body: c.WrapCodeBlock(func() string {
						var ret []string
						for _, match := range tag.Samples {
							ret = append(ret, fmt.Sprintf("%s\t%s", match.Tag, match.Desc))
						}
						return strings.Join(ret, "\n")
					}()),
				},
				{
					Title: c.I18n("Examples"),
					Body:  c.WrapCodeBlock(tag.Examples),
				},
			},
		})
	}
	c.SendMessage("cmds", commandlib.EmbedList{
		ItemTypeName: c.I18n("Command"),
		Embeds:       commandEmbeds,
	})
	c.SendMessage("tags", commandlib.EmbedList{
		ItemTypeName: c.I18n("Tag"),
		Embeds:       tagEmbeds,
	})
}
