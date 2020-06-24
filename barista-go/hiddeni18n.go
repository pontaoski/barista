package barista

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/log"
)

func init() {
	commandlib.RegisterCommand(commandlib.Command{
		Matches: []string{"o kama e sitelen toki ni:"},
		Hidden:  true,
		Action:  HiddenI18n,
	})
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func HiddenI18n(c commandlib.Context) {
	if !c.Backend().IsBotOwner(c) {
		c.SendMessage("primary", commandlib.ErrorEmbed("You are not the bot owner."))
	}
	if c.Arg(0) == "" {
		c.SendMessage("primary", commandlib.ErrorEmbed("Please provide a language to download."))
	}
	in := fmt.Sprintf("https://raw.githubusercontent.com/pontaoski/barista/master/messages/barista_%s.po", c.Arg(0))
	log.Info(in)
	err := DownloadFile(fmt.Sprintf("messages/barista_%s.po", c.Arg(0)), in)
	if err != nil {
		c.SendMessage("primary", commandlib.ErrorEmbed("Failed to download: "+err.Error()))
	} else {
		commandlib.DropPo(c.Arg(0))
		c.SendMessage("primary", commandlib.ErrorEmbed("Updated languages."))
	}
}
