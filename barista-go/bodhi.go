package barista

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/log"
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

func init() {
	commandlib.RegisterTag(commandlib.Tag{
		Name:     I18n("Bodhi Updates"),
		Usage:    I18n("Link to Bodhi updates"),
		Examples: `FEDORA-2020-81f9f75f04, FEDORA-2020-217b6928cc, FEDORA-2020-06730065a6, FEDORA-EPEL-2020-a729ac8728`,
		Samples: []commandlib.TagSample{
			{
				Tag:  "FEDORA-*",
				Desc: I18n("Fedora Bodhi Updates"),
			},
		},
		ID:     "bodhi-updates",
		Match:  regexp.MustCompile("FEDORA-"),
		Action: Bodhi,
	})
}

func Bodhi(c commandlib.Context) {
	var embeds []commandlib.Embed
	for _, word := range c.Args() {
		resp, err := http.Get(fmt.Sprintf("https://bodhi.fedoraproject.org/updates/%s", word))
		if err != nil {
			log.Error("%+v", err)
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error("%+v", err)
			continue
		}

		var updateWrapper UpdateWrapper
		err = json.Unmarshal(body, &updateWrapper)
		if err != nil {
			log.Error("%+v", err)
			continue
		}

		update := updateWrapper.Update

		var builds []string

		for _, build := range update.Builds {
			builds = append(
				builds,
				c.GenerateLink(build.NVR, fmt.Sprintf("https://koji.fedoraproject.org/koji/search?terms=%s&type=build&match=glob", build.NVR)),
			)
		}

		var karmamsg string
		if update.Karma == 0 {
			karmamsg = fmt.Sprintf("%d", update.Karma)
		} else if update.Karma > 0 {
			karmamsg = fmt.Sprintf("%d ✅", update.Karma)
		} else if update.Karma < 0 {
			karmamsg = fmt.Sprintf("%d ❌", update.Karma)
		}

		embeds = append(embeds, commandlib.Embed{
			Colour: 0x3C6EB4,
			Title: commandlib.EmbedHeader{
				Text: update.Name,
				URL:  update.URL,
			},
			Body: update.Description,
			Header: commandlib.EmbedHeader{
				Text: update.Author.Name,
				Icon: update.Author.Avatar,
			},
			Footer: commandlib.EmbedHeader{
				Text: c.I18n("Fedora Updates System"),
				Icon: "https://bodhi.fedoraproject.org/static/v5.1.0/ico/favicon.ico",
			},
			Fields: []commandlib.EmbedField{
				{
					Title:  c.I18n("Distro"),
					Body:   update.Release.Display,
					Inline: true,
				},
				{
					Title:  c.I18n("Karma"),
					Body:   karmamsg,
					Inline: true,
				},
				{
					Title:  c.I18n("Builds"),
					Body:   strings.Join(builds, ", "),
					Inline: true,
				},
				{
					Title: c.I18n("Status"),
					Body:  strings.Title(update.Status),
				},
			},
		})
	}
	if len(embeds) > 0 {
		c.SendTags("primary", embeds)
	}
}
