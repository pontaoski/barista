package barista

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/alecthomas/repr"
	"github.com/appadeia/barista/barista-go/commandlib"
)

const (
	SefariaSearchAPI = "https://www.sefaria.org/api/search-wrapper"
	SefariaTextAPI   = "https://www.sefaria.org/api/texts"
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
		"<i>", "",
		"</i>", "",
	)
)

func searchSefaria(query string) (res []sefariaResult, err error) {
	resp, err := http.Post(SefariaSearchAPI, "application/json", toJSONReader(sefariaQuery{
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

type verseKind interface {
	String() string
	isVerse()
}

type singleVerse int

func (s singleVerse) isVerse() {}
func (s singleVerse) String() string {
	return strconv.Itoa(int(s))
}

type singleVerseWithComment struct {
	verse   int
	comment int
}

func (s singleVerseWithComment) isVerse() {}
func (s singleVerseWithComment) String() string {
	return strconv.Itoa(int(s.verse)) + ":" + strconv.Itoa(int(s.comment))
}

type verseRange struct {
	from int
	to   int
}

func (s verseRange) isVerse() {}
func (s verseRange) String() string {
	return strconv.Itoa(int(s.from)) + "-" + strconv.Itoa(int(s.to))
}

type SefText string

func (p *SefText) UnmarshalJSON(data []byte) error {
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	switch a := v.(type) {
	case []interface{}:
		var b []string
		for _, item := range a {
			b = append(b, item.(string))
		}
		*p = SefText(strings.Join(b, " "))
	case string:
		*p = SefText(a)
	default:
		return errors.New("Unrecognised value: " + repr.String(v))
	}
	return nil
}

func locateText(book, chapter string, verse verseKind) (string, error) {
	book = strings.ReplaceAll(book, " ", "_")

	var url string

	switch a := verse.(type) {
	case singleVerse:
		url = fmt.Sprintf("%s/%s.%s.%d", SefariaTextAPI, book, chapter, a)
	case singleVerseWithComment:
		url = fmt.Sprintf("%s/%s.%s.%d.%d", SefariaTextAPI, book, chapter, a.verse, a.comment)
	case verseRange:
		url = fmt.Sprintf("%s/%s.%s.%d-%d", SefariaTextAPI, book, chapter, a.from, a.to)
	default:
		panic("unimplemented")
	}

	url += "?context=0"

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Text SefText `json:"text"`
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return "", err
	}

	return replacer.Replace(string(result.Text)), nil
}

func init() {
	commandlib.RegisterCommand(commandlib.Command{
		Name:  I18n("Jewish Text Lookup"),
		Usage: I18n("Look up a specific Jewish text using Sefaria"),
		ID:    "sefaria",
		Matches: []string{
			"o pana e lipu Ju ni tawa mi:",
		},
		Action: func(c commandlib.Context) {
			book := c.Arg(0)
			chapterVerse := strings.SplitN(c.Arg(1), ":", 2)
			chapter, rawVerse := chapterVerse[0], chapterVerse[1]
			var verse verseKind

			badInt := commandlib.ErrorEmbed("Your number doesn't appear to be a number")

			switch {
			case strings.ContainsRune(rawVerse, ':'):
				data := strings.Split(rawVerse, ":")
				nanpa, comment := data[0], data[1]
				nanpaInt, err := strconv.ParseInt(nanpa, 10, 64)
				if err != nil {
					c.SendMessage("primary", badInt)
					return
				}
				commentInt, err := strconv.ParseInt(comment, 10, 64)
				if err != nil {
					c.SendMessage("primary", badInt)
					return
				}
				verse = singleVerseWithComment{
					verse:   int(nanpaInt),
					comment: int(commentInt),
				}
			case strings.ContainsRune(rawVerse, '-'):
				data := strings.Split(rawVerse, "-")
				nanpa, comment := data[0], data[1]
				nanpaInt, err := strconv.ParseInt(nanpa, 10, 64)
				if err != nil {
					c.SendMessage("primary", badInt)
					return
				}
				commentInt, err := strconv.ParseInt(comment, 10, 64)
				if err != nil {
					c.SendMessage("primary", badInt)
					return
				}
				verse = verseRange{
					from: int(nanpaInt),
					to:   int(commentInt),
				}
			default:
				nanpa, err := strconv.ParseInt(rawVerse, 10, 64)
				if err != nil {
					c.SendMessage("primary", badInt)
					return
				}
				verse = singleVerse(nanpa)
			}

			data, err := locateText(book, chapter, verse)
			if err != nil {
				c.SendMessage("primary", commandlib.ErrorEmbed("There was an error looking up your desired text: "+err.Error()))
				return
			}

			c.SendMessage("primary", commandlib.Embed{
				Title: commandlib.EmbedHeader{
					Text: fmt.Sprintf("%s %s:%s", book, chapter, verse.String()),
				},
				Body: data,
			})
		},
	})
}
