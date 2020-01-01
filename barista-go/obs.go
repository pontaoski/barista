package barista

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/bwmarrin/discordgo"
)

func isInt(s string) bool {
	_, err := strconv.ParseInt(s, 0, 0)
	if err == nil {
		return true
	} else {
		return false
	}
}

func Obs(s *discordgo.Session, cmd *LexedCommand) {
	ctnt := strings.ToLower(cmd.CommandMessage.Content)
	words := strings.Split(ctnt, " ")
	var ids []int64
	var urlIDS []int64
	if strings.Contains(ctnt, "sr") {
		for index, word := range words {
			if strings.Contains(word, "sr") {
				if len(word) != 2 {
					if isInt(word[2:]) {
						// SR1234
						inty, _ := strconv.ParseInt(word[2:], 0, 0)
						ids = append(ids, inty)
					} else if isInt(word[3:]) {
						// SR#1234
						inty, _ := strconv.ParseInt(word[3:], 0, 0)
						ids = append(ids, inty)
					}
				} else {
					if isInt(words[index+1]) {
						inty, _ := strconv.ParseInt(words[index+1], 0, 0)
						ids = append(ids, inty)
					}
				}
			}
		}
	}
	r, _ := regexp.Compile(`http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\(\), ]|(?:%[0-9a-fA-F][0-9a-fA-F]))+`)

	for _, url := range r.FindAllString(ctnt, -1) {
		components := strings.Split(url, "/")
		if strings.Contains(url, "build.opensuse.org") && strings.Contains(url, "request/show") {
			str := components[len(components)-1]
			if isInt(str) {
				inty, _ := strconv.ParseInt(str, 0, 0)
				ids = append(ids, inty)
				urlIDS = append(urlIDS, inty)
			}
		}
	}

	client := &http.Client{}
	var embeds []*Embed

	for _, id := range ids {
		s.ChannelTyping(cmd.CommandMessage.ChannelID)
		request, _ := http.NewRequest("GET", fmt.Sprintf("https://api.opensuse.org/request/%d", id), nil)
		request.SetBasicAuth("zyp_user", "zyp_pw_1")
		resp, err := client.Do(request)
		if err != nil {
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		s := string(body)

		doc, err := xmlquery.Parse(strings.NewReader(s))
		if err != nil {
			continue
		}

		list, err := xmlquery.QueryAll(doc, "/request/description")
		if err != nil {
			continue
		}
		desc := list[0].InnerText()

		list, err = xmlquery.QueryAll(doc, "/request/state/@name")
		if err != nil {
			continue
		}
		state := list[0].InnerText()

		list, err = xmlquery.QueryAll(doc, "/request/action/@type")
		if err != nil {
			continue
		}
		action := list[0].InnerText()

		if action == "maintenance_incident" {
			action = "Submit update from"
		}

		list, err = xmlquery.QueryAll(doc, "/request/action/source/@project")
		if err != nil {
			continue
		}
		sourceProj := list[0].InnerText()

		list, err = xmlquery.QueryAll(doc, "/request/action/source/@package")
		if err != nil {
			continue
		}
		sourcePkg := list[0].InnerText()

		list, err = xmlquery.QueryAll(doc, "/request/action/target/@project")
		if err != nil {
			continue
		}
		targetProj := list[0].InnerText()

		list, err = xmlquery.QueryAll(doc, "/request/action/target/@package")
		if err != nil {
			continue
		}
		var targetPkg string
		if len(list) > 0 {
			targetPkg = list[0].InnerText()
		}

		embed := NewEmbed().
			SetColor(0x73ba25).
			SetTitle(fmt.Sprintf("SR#%d — %s — **%s**", id, desc, strings.Title(state))).
			SetAuthor("openSUSE Build Service", "https://en.opensuse.org/images/c/cd/Button-colour.png", fmt.Sprintf("https://build.opensuse.org/request/show/%d", id))

		if targetPkg == "" {
			embed.SetDescription(fmt.Sprintf("%s **%s**:**%s** → **%s**", strings.Title(action), sourceProj, sourcePkg, targetProj))
		} else {
			embed.SetDescription(fmt.Sprintf("%s **%s**:**%s** → **%s**:**%s**", strings.Title(action), sourceProj, sourcePkg, targetProj, targetPkg))
		}

		embeds = append(embeds, embed)
	}

	cmd.SendTags(embeds)
}
