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
		Match: [][]string{
			{"sudo", "help"},
			{"o", "help"},
		},
		Action: Help,
	})
}

func join(strs [][]string) string {
	var ret []string
	for _, str := range strs {
		ret = append(ret, strings.Join(str, " "))
	}
	return strings.Join(ret, ", ")
}

func Help(c commandlib.Context) {
	var commandEmbeds []commandlib.Embed
	var tagEmbeds []commandlib.Embed
	for _, command := range commandlib.Commands() {
		commandEmbeds = append(commandEmbeds, commandlib.Embed{
			Title: commandlib.EmbedHeader{
				Text: fmt.Sprintf("%s (%s)", command.Name, command.ID),
			},
			Body: command.Usage,
			Fields: []commandlib.EmbedField{
				{
					Title: "Aliases",
					Body:  join(command.Match),
				},
				{
					Title: "Examples",
					Body:  c.WrapCodeBlock(command.Examples),
				},
			},
		})
	}
	for _, tag := range commandlib.Tags() {
		tagEmbeds = append(tagEmbeds, commandlib.Embed{
			Title: commandlib.EmbedHeader{
				Text: fmt.Sprintf("%s (%s)", tag.Name, tag.ID),
			},
			Body: tag.Usage,
			Fields: []commandlib.EmbedField{
				{
					Title: "Tags",
					Body: c.WrapCodeBlock(func() string {
						var ret []string
						for _, match := range tag.Samples {
							ret = append(ret, fmt.Sprintf("%s\t%s", match.Tag, match.Desc))
						}
						return strings.Join(ret, "\n")
					}()),
				},
				{
					Title: "Examples",
					Body:  c.WrapCodeBlock(tag.Examples),
				},
			},
		})
	}
	c.SendMessage("cmds", commandlib.EmbedList{
		ItemTypeName: "Command",
		Embeds:       commandEmbeds,
	})
	c.SendMessage("tags", commandlib.EmbedList{
		ItemTypeName: "Tag",
		Embeds:       tagEmbeds,
	})
}
