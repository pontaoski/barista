package barista

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/appadeia/barista/barista-go/commandlib"
)

const (
	SefariaAPI = "https://www.sefaria.org/api/search-wrapper"
)

type sefariaQuery struct {
	Query string `json:"query"`
	Type  string `json:"type"`
}

type sefariaQueryResponse struct {
	Hits struct {
		Hits []struct {
			ID        string `json:"_id"`
			Highlight struct {
				Exact []string `json:"exact"`
			} `json:"highlight"`
		} `json:"hits"`
	} `json:"hits"`
}

type sefariaResult struct {
	Location string
	Matches  []string
}

func toJSONReader(v interface{}) io.Reader {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(data)
}

var (
	replacer = strings.NewReplacer(
		"<b>", "",
		"</b>", "",
	)
)

func searchSefaria(query string) (res []sefariaResult, err error) {
	resp, err := http.Post(SefariaAPI, "application/json", toJSONReader(sefariaQuery{
		Query: query,
		Type:  "text",
	}))
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var q sefariaQueryResponse
	err = json.Unmarshal(data, &q)
	if err != nil {
		return
	}
	for _, data := range q.Hits.Hits {
		res = append(res, sefariaResult{
			Location: data.ID,
			Matches: func() (ret []string) {
				for _, item := range data.Highlight.Exact {
					ret = append(ret, replacer.Replace(item))
				}
				return
			}(),
		})
	}
	return
}

func init() {
	commandlib.RegisterCommand(commandlib.Command{
		Name:  I18n("Jewish Text Search"),
		Usage: I18n("Search a heck tonne of Jewish texts using Sefaria"),
		ID:    "sefaria",
		Matches: []string{
			"o alasa lon lipu suli Ju e ni:",
		},
		Action: func(c commandlib.Context) {
			data, err := searchSefaria(c.Content())
			if err != nil {
				c.SendMessage("primary", commandlib.ErrorEmbed("There was an error searching texts"))
			}
			c.SendMessage("primary", commandlib.EmbedList{
				ItemTypeName: "Result",
				Embeds: func() (ret []commandlib.Embed) {
					for _, item := range data {
						ret = append(ret, commandlib.Embed{
							Title: commandlib.EmbedHeader{
								Text: item.Location,
							},
							Body: strings.Join(item.Matches, "\n\n"),
						})
					}
					return
				}(),
			})
		},
	})
}
