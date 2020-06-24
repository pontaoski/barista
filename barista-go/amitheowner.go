package barista

import "github.com/appadeia/barista/barista-go/commandlib"

func init() {
	commandlib.RegisterCommand(commandlib.Command{
		Matches: []string{"am i the owner"},
		Hidden:  true,
		Action: func(c commandlib.Context) {
			if c.Backend().IsBotOwner(c) {
				c.SendMessage("primary", "yes")
			} else {
				c.SendMessage("primary", "no")
			}
		},
	})
}
