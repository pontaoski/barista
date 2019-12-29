package barista

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Release struct {
	Display string `json:"long_name"`
}

type Build struct {
	NVR       string `json:"nvr"`
	ReleaseId int    `json:"release_id"`
	Signed    bool   `json:"signed"`
	Type      string `json:"rpm"`
	Epoch     int    `json:"epoch"`
}

type User struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type Update struct {
	Karma       int     `json:"karma"`
	Name        string  `json:"title"`
	Description string  `json:"notes"`
	Release     Release `json:"release"`
	Builds      []Build `json:"builds"`
	Author      User    `json:"user"`
	URL         string  `json:"url"`
	Status      string  `json:"status"`
}

type UpdateWrapper struct {
	Update  Update `json:"update"`
	CanEdit bool   `json:"can_edit"`
}

func Bodhi(s *discordgo.Session, cmd *LexedCommand) {
	ctnt := cmd.CommandMessage.Content
	words := strings.Split(ctnt, " ")
	var embeds []*Embed

	if strings.Contains(ctnt, "FEDORA-") {
		for _, word := range words {
			if strings.HasPrefix(word, "FEDORA-") {
				s.ChannelTyping(cmd.CommandMessage.ChannelID)
				resp, err := http.Get(fmt.Sprintf("https://bodhi.fedoraproject.org/updates/%s", word))
				if err != nil {
					continue
				}
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					continue
				}

				var updateWrapper UpdateWrapper
				err = json.Unmarshal(body, &updateWrapper)
				if err != nil {
					continue
				}

				update := updateWrapper.Update

				builds := fmt.Sprintf("[%s](https://koji.fedoraproject.org/koji/search?terms=%s&type=build&match=glob)", update.Builds[0].NVR, update.Builds[0].NVR)

				for _, build := range update.Builds[1:] {
					builds = fmt.Sprintf("%s, [%s](https://koji.fedoraproject.org/koji/search?terms=%s&type=build&match=glob)", builds, build.NVR, build.NVR)
				}
				var karmamsg string
				if update.Karma == 0 {
					karmamsg = fmt.Sprintf("%d", update.Karma)
				} else if update.Karma > 0 {
					karmamsg = fmt.Sprintf("%d ✅", update.Karma)
				} else if update.Karma < 0 {
					karmamsg = fmt.Sprintf("%d ❌", update.Karma)
				}
				embed := NewEmbed().
					SetTitle(update.Name).
					SetURL(update.URL).
					SetDescription(update.Description).
					AddField("Distro", update.Release.Display, true).
					AddField("Karma", karmamsg, true).
					AddField("Builds", builds, true).
					AddField("Status", strings.Title(update.Status), true).
					SetAuthor(update.Author.Name, update.Author.Avatar).
					SetColor(0x3C6EB4).
					SetFooter("Fedora Updates System", "https://bodhi.fedoraproject.org/static/v5.1.0/ico/favicon.ico")

				embeds = append(embeds, embed)
			}
		}
		if len(embeds) > 0 {
			cmd.SendTags(embeds)
		}
	}
}
