package barista

import (
	"strings"

	"github.com/appadeia/barista/barista-go/commandlib"
)

func init() {
	commandlib.RegisterCommand(commandlib.Command{
		Name:  I18n("Notes"),
		Usage: I18n("Store and recall notes."),
		ID:    "notes",
		Matches: []string{
			"sudo notes",
		},
		Examples: `sudo notes store id this is a very cool note`,
		Action:   Notes,
	})
	commandlib.RegisterCommand(commandlib.Command{
		Hidden:  true,
		Action:  SilentRecall,
		Matches: []string{".."},
	})
	commandlib.RegisterCommand(commandlib.Command{
		Hidden:  true,
		Action:  SilentStore,
		Matches: []string{"!!"},
	})
}

func SilentRecall(c commandlib.Context) {
	if c.Arg(0) == "" {
		return
	}
	data, ok := commandlib.GetNote(c, c.Arg(0), commandlib.Community)
	if ok {
		c.SendMessage("primary", data)
	}
}

func SilentStore(c commandlib.Context) {
	if c.Arg(0) == "" {
		return
	}
	trimmed := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(c.RawContent()), c.Arg(0)))
	commandlib.StoreNote(c, c.Arg(0), trimmed, commandlib.Community)
	c.SendMessage("primary", "Note saved!")
}

func Notes(c commandlib.Context) {
	switch c.Arg(0) {
	case "store":
		if c.Arg(1) == "" {
			c.SendMessage("primary", CommandHelp(c, c.Command()))
			return
		}
		trimmed := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(c.RawContent()), "store")), c.Arg(1)))
		commandlib.StoreNote(c, c.Arg(1), trimmed, commandlib.Community)
		c.SendMessage("primary", "Note saved!")
	case "read":
		if c.Arg(1) == "" {
			c.SendMessage("primary", CommandHelp(c, c.Command()))
			return
		}
		data, ok := commandlib.GetNote(c, c.Arg(1), commandlib.Community)
		if !ok {
			c.SendMessage("primary", commandlib.ErrorEmbed(c.I18n("There was an error recalling the note. Does it exist?")))
		}
		c.SendMessage("primary", data)
	default:
		c.SendMessage("primary", CommandHelp(c, c.Command()))
		return
	}
}
