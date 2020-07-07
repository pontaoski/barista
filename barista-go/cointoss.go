package barista

import (
	"math/rand"
	"time"

	"github.com/appadeia/barista/barista-go/commandlib"
)

func init() {
	commandlib.RegisterCommand(commandlib.Command{
		Name:     I18n("Coin Toss"),
		Usage:    I18n("Flip a coin."),
		Matches:  []string{"sudo cointoss", "o cointoss"},
		Examples: `sudo cointoss`,
		ID:       "coinflip",
		Action:   CoinToss,
	})
}

func CoinToss(c commandlib.Context) {
	rand.Seed(time.Now().UnixNano())
	side := ""
	if rand.Intn(2) == 0 {
		side = c.I18n("You got some head!")
	} else {
		side = c.I18n("You got Tails The Fox!")
	}
	c.SendMessage("primary", side)
}
