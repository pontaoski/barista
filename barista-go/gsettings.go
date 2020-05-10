package barista

import (
	"fmt"
	"strings"

	"github.com/appadeia/barista/barista-go/commandlib"
)

func init() {
	commandlib.RegisterCommand(commandlib.Command{
		Name:  I18n("GSettings"),
		Usage: I18n("Configure the bot"),
		ID:    "gsettings",
		Match: [][]string{
			{"sudo", "gsettings"},
			{"o", "gsettings"},
		},
		Action: GSettings,
	})
}

var schemas = map[string]commandlib.Schema{
	"locale": {
		Name:           I18n("Preferred Locale"),
		Description:    I18n("The preferred language of this channel."),
		ID:             "locale",
		DefaultValue:   "en",
		PossibleValues: []string{"en", "de", "es", "fr", "it", "nl", "pl", "tpo"},
	},
	"default-distro": {
		Name:           I18n("Default Distro"),
		Description:    I18n("The default distro to use for package searches."),
		ID:             "default-distro",
		DefaultValue:   "",
		PossibleValues: []string{"tumbleweed", "leap", "fedora", "mageia", "openmandriva", "centos", "packman-leap", "packman-tumbleweed", "rpmfusion"},
	},
}

func GSettings(c commandlib.Context) {
	if key := c.Arg(0); key != "" {
		if value := c.Arg(1); value != "" {
			if schema, ok := schemas[key]; ok {
				if value == "reset" {
					schema.WriteValue(c, schema.DefaultValue)
					c.SendMessage(
						"primary",
						commandlib.Embed{
							Title: commandlib.EmbedHeader{
								Text: c.I18n("Setting reset!"),
							},
						},
					)
					return
				}
				if schema.ValueValid(value) {
					schema.WriteValue(c, value)
					c.SendMessage(
						"primary",
						commandlib.Embed{
							Title: commandlib.EmbedHeader{
								Text: c.I18n("Setting updated!"),
							},
						},
					)
				} else {
					c.SendMessage(
						"primary",
						commandlib.ErrorEmbed(
							fmt.Sprintf(
								c.I18n("%s is not an accepted value. Please use a value from the following list:\n%s"),
								value,
								c.WrapCodeBlock(strings.Join(schema.PossibleValues, ", ")),
							),
						),
					)
				}
			} else {
				c.SendMessage("primary", commandlib.ErrorEmbed(fmt.Sprintf(c.I18n("%s is not a valid setting key"), key)))
			}
		} else {
			if schema, ok := schemas[key]; ok {
				c.SendMessage("primary", schema.ToEmbed(c))
			} else {
				c.SendMessage("primary", commandlib.ErrorEmbed(fmt.Sprintf(c.I18n("%s is not a valid setting key"), key)))
			}
		}
	} else {
		var schemaEmbeds []commandlib.Embed
		for _, schema := range schemas {
			schemaEmbeds = append(schemaEmbeds, schema.ToEmbed(c))
		}
		c.SendMessage("primary", commandlib.EmbedList{
			ItemTypeName: c.I18n("Command"),
			Embeds:       schemaEmbeds,
		})
	}
}
