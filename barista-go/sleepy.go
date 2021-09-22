package barista

import (
	"time"

	"github.com/appadeia/barista/barista-go/commandlib"
)

func init() {
	commandlib.RegisterCommand(commandlib.Command{
		Name:  I18n("Sleepy"),
		Usage: I18n("Make Barista sleep for a bit"),
		ID:    "sleepy",
		Matches: []string{
			"ilo, sleep for",
		},
		Action: func(c commandlib.Context) {
			dur, err := time.ParseDuration(c.Content())
			if err != nil {
				c.SendMessage("primary", commandlib.ErrorEmbed(err.Error()))
			}
			time.Sleep(dur)
			c.SendMessage("primary", "mu")
		},
	})
}
