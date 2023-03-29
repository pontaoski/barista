package barista

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"

	"github.com/appadeia/barista/barista-go/commandlib"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func init() {
	commandlib.RegisterCommand(commandlib.Command{
		Name:     I18n("sitelen pona"),
		Usage:    I18n("Write with sitelen pona."),
		Matches:  []string{"ilo o sitelen pona", ",os"},
		Examples: `ilo o sitelen pona mu`,
		ID:       "sitelenpona",
		Action: func(c commandlib.Context) {
			filename := "/tmp/" + randSeq(10) + ".png"
			cmd := exec.Command("pango-view", "--no-display", "-t", c.RawContent(), "--font", "linja lipamanka 50", "-o", filename, "--align=center", "--hinting=full", "--margin=10px")
			fmt.Printf("%+v", cmd.Args)

			var b bytes.Buffer
			cmd.Stdout = &b
			cmd.Stderr = &b

			if err := cmd.Run(); err != nil {
				c.SendMessage("main", commandlib.ErrorEmbed("ilo li pakala a! \n"+b.String()))
				return
			}

			mu, err := os.Open(filename)
			if err != nil {
				c.SendMessage("main", commandlib.ErrorEmbed("ilo li pakala a! \n"+err.Error()))
				return
			}

			c.SendMessage("main", commandlib.File{
				Name:     "image.png",
				Mimetype: "image/png",
				Reader:   mu,
			})
		},
	})
}
