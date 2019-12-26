package barista

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Setting struct {
	name         string
	parent       string
	defaultValue string
	summary      string
}

type Schema struct {
	name        string
	description string
	settings    []Setting
}

var Schemas []Schema = []Schema{
	{
		name:        "dnf",
		description: "The package management command of choice.",
		settings: []Setting{
			{
				parent:       "dnf",
				name:         "defaultDistro",
				defaultValue: "fedora",
				summary:      "The default distro of this guild",
			},
		},
	},
}

func (set *Setting) getValue(cmd *LexedCommand) string {
	val := cmd.GetGuildKey(fmt.Sprintf("gsettings_%s_%s", set.parent, set.name))
	if val == "" {
		return set.defaultValue
	}
	return val
}

func (set *Setting) setValue(cmd *LexedCommand, val string) {
	cmd.SetGuildKey(fmt.Sprintf("gsettings_%s_%s", set.parent, set.name), val)
}

func schemaExists(name string) bool {
	for _, schema := range Schemas {
		if schema.name == name {
			return true
		}
	}
	return false
}

func getSchema(name string) Schema {
	var schema Schema
	for _, schem := range Schemas {
		if schem.name == name {
			return schem
		}
	}
	return schema
}

func getSetting(schemaName string, settingName string) Setting {
	var set Setting
	schema := getSchema(schemaName)
	for _, setting := range schema.settings {
		if setting.name == settingName {
			return setting
		}
	}
	return set
}

func settingExists(schemaName string, settingName string) bool {
	if !schemaExists(schemaName) {
		return false
	}
	for _, schema := range Schemas {
		if schema.name == schemaName {
			for _, settings := range schema.settings {
				if settings.name == settingName {
					return true
				}
			}
		}
	}
	return false
}

func Gsettings(s *discordgo.Session, cmd *LexedCommand) {
	helpmsg := "```dsconfig\n" + gsettingshelp + "\n```"

	if !cmd.Author.IsAdmin {
		embed := NewEmbed().
			SetColor(0xff0000).
			SetTitle("You need admin permissions to use this command.")
		msgSend := discordgo.MessageSend{
			Embed: embed.MessageEmbed,
		}
		cmd.SendMessage(&msgSend)
		return
	}
	if cmd.GetFlagPair("-g", "--get") == "" &&
		cmd.GetFlagPair("-s", "--set") == "" &&
		cmd.GetFlagPair("-l", "--list-schemas") != "nil" &&
		cmd.GetFlagPair("-k", "--list-settings") != "nil" {
		embed := NewEmbed().
			SetColor(0xff0000).
			SetDescription(helpmsg).
			SetTitle("Please specify an action to take in your command.")
		msgSend := discordgo.MessageSend{
			Embed: embed.MessageEmbed,
		}
		cmd.SendMessage(&msgSend)
		return
	}
	if cmd.GetFlagPair("-g", "--get") != "" ||
		cmd.GetFlagPair("-s", "--set") != "" {
		if cmd.GetFlagPair("", "--schema") == "" ||
			cmd.GetFlagPair("", "--setting") == "" {
			embed := NewEmbed().
				SetColor(0xff0000).
				SetDescription(helpmsg).
				SetTitle("Please specify a schema and setting.")
			msgSend := discordgo.MessageSend{
				Embed: embed.MessageEmbed,
			}
			cmd.SendMessage(&msgSend)
			return
		}
		if !settingExists(cmd.GetFlagPair("", "--schema"), cmd.GetFlagPair("", "--setting")) {
			embed := NewEmbed().
				SetColor(0xff0000).
				SetDescription(helpmsg).
				SetTitle("Please list a valid schema and setting.")
			msgSend := discordgo.MessageSend{
				Embed: embed.MessageEmbed,
			}
			cmd.SendMessage(&msgSend)
			return
		}
	}
	if cmd.GetFlagPair("-k", "--list-settings") != "" {
		if cmd.GetFlagPair("", "--schema") == "" {
			embed := NewEmbed().
				SetColor(0xff0000).
				SetDescription(helpmsg).
				SetTitle("Please specify a schema to list settings for.")
			msgSend := discordgo.MessageSend{
				Embed: embed.MessageEmbed,
			}
			cmd.SendMessage(&msgSend)
			return
		}
		if !schemaExists(cmd.GetFlagPair("", "--schema")) {
			embed := NewEmbed().
				SetColor(0xff0000).
				SetDescription(helpmsg).
				SetTitle("Please specify a valid schema.")
			msgSend := discordgo.MessageSend{
				Embed: embed.MessageEmbed,
			}
			cmd.SendMessage(&msgSend)
			return
		}
	}
	if cmd.GetFlagPair("-l", "--list-schemas") != "" {
		embedmsg := []string{}
		for _, schema := range Schemas {
			embedmsg = append(embedmsg, fmt.Sprintf("%s\n\t%s", schema.name, schema.description))
		}
		embed := NewEmbed().
			SetTitle("List of schemas").
			SetDescription("```\n" + strings.Join(embedmsg, "\n") + "\n```").
			SetColor(0xc12fb7)
		msgSend := discordgo.MessageSend{
			Embed: embed.MessageEmbed,
		}
		cmd.SendMessage(&msgSend)
		return
	}
	if cmd.GetFlagPair("-k", "--list-settings") != "" {
		schema := getSchema(cmd.GetFlagPair("", "--schema"))
		embedmsg := []string{}
		for _, set := range schema.settings {
			embedmsg = append(embedmsg, fmt.Sprintf("%s\n\t%s", set.name, set.summary))
		}
		embed := NewEmbed().
			SetTitle("List of settings in schema " + schema.name).
			SetDescription("```\n" + strings.Join(embedmsg, "\n") + "\n```").
			SetColor(0xc12fb7)
		msgSend := discordgo.MessageSend{
			Embed: embed.MessageEmbed,
		}
		cmd.SendMessage(&msgSend)
		return
	}
	if cmd.GetFlagPair("-g", "--get") != "" {
		setting := getSetting(cmd.GetFlagPair("", "--schema"), cmd.GetFlagPair("", "--setting"))
		embed := NewEmbed().
			SetTitle("Value of " + setting.name + ":").
			SetDescription(setting.getValue(cmd)).
			SetColor(0xc12fb7)
		msgSend := discordgo.MessageSend{
			Embed: embed.MessageEmbed,
		}
		cmd.SendMessage(&msgSend)
		return
	}
	if cmd.GetFlagPair("-s", "--set") != "" {
		if cmd.GetFlagPair("", "--value") == "" || cmd.GetFlagPair("", "--value") == "nil" {
			embed := NewEmbed().
				SetColor(0xff0000).
				SetDescription(helpmsg).
				SetTitle("Please specify a value.")
			msgSend := discordgo.MessageSend{
				Embed: embed.MessageEmbed,
			}
			cmd.SendMessage(&msgSend)
			return
		}
		setting := getSetting(cmd.GetFlagPair("", "--schema"), cmd.GetFlagPair("", "--setting"))
		setting.setValue(cmd, cmd.GetFlagPair("", "--value"))
		embed := NewEmbed().
			SetTitle("Value of " + setting.name + " set:").
			SetDescription(setting.getValue(cmd)).
			SetColor(0xc12fb7)
		msgSend := discordgo.MessageSend{
			Embed: embed.MessageEmbed,
		}
		cmd.SendMessage(&msgSend)
		return
	}
}
