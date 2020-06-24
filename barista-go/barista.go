package barista

import (
	_ "github.com/appadeia/barista/barista-go/backends/discord"
	_ "github.com/appadeia/barista/barista-go/backends/matrix"
	_ "github.com/appadeia/barista/barista-go/backends/telegram"
	"github.com/appadeia/barista/barista-go/commandlib"
)

// Main : Call this function to start the bot's main loop.
func Main() {
	commandlib.StartBackends()
}
