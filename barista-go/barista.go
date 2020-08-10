package barista

import (
	"os"

	_ "github.com/appadeia/barista/barista-go/backends/discord"
	_ "github.com/appadeia/barista/barista-go/backends/irc"
	_ "github.com/appadeia/barista/barista-go/backends/matrix"
	_ "github.com/appadeia/barista/barista-go/backends/telegram"
	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/pkgman"
	"github.com/urfave/cli/v2"
)

// Main : Call this function to start the bot's main loop.
func Main() {
	app := cli.App{
		Commands: []*cli.Command{
			{
				Name: "pkgman",
				Action: func(c *cli.Context) error {
					pkgman.Cli(c)
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			commandlib.StartBackends()
			return nil
		},
	}
	app.Run(os.Args)
}
