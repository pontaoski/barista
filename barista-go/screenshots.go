package barista

import (
	"strings"

	"github.com/appadeia/barista/barista-go/commandlib"
)

func init() {
	commandlib.RegisterCommand(commandlib.Command{
		Name:    I18n("Screenshots"),
		Usage:   I18n("User screenshots."),
		ID:      "screenshots",
		Matches: []string{"sudo ss", "sudo screenshot"},
		Action:  Screenshots,
	})
}

func Screenshots(c commandlib.Context) {
	switch c.Arg(0) {
	case "set":
		switch c.Arg(1) {
		case "ss", "screenshot":
			if !strings.HasPrefix(c.Arg(2), ".png") {
				c.SendMessage("primary", commandlib.ErrorEmbed(c.I18n("The given screenshot must be a PNG file.")))
				return
			}
			commandlib.StoreData(c, c.AuthorIdentifier()+"user-screenshot", c.Arg(2), commandlib.Global)
			c.SendMessage("primary", "Updated screenshot!")
		case "desc", "description":
			strings.Join(c.Args()[3:], " ")
			commandlib.StoreData(c, c.AuthorIdentifier()+"user-screenshot", c.Arg(2), commandlib.Global)
			c.SendMessage("primary", "Updated screenshot!")
		}
	case "get":
	}
}
