package barista

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/appadeia/barista/barista-go/commandlib"
)

var languages = map[string]int{
	"c#":          1,
	"vb.net":      2,
	"f#":          3,
	"java":        4,
	"python2":     5,
	"c-gcc":       6,
	"c++-gcc":     7,
	"php":         8,
	"pascal":      9,
	"objective-C": 10,
	"haskell":     11,
	"ruby":        12,
	"perl":        13,
	"lua":         14,
	"nasm":        15,
	"javascript":  17,
	"lisp":        18,
	"prolog":      19,
	"go":          20,
	"scala":       21,
	"scheme":      22,
	"node.js":     23,
	"python3":     24,
	"octave":      25,
	"c-clang":     26,
	"c++-clang":   27,
	"c++-vc++":    28,
	"c-vc":        29,
	"d":           30,
	"r":           31,
	"tcl":         32,
	"oracle":      35,
	"swift":       37,
	"bash":        38,
	"ada":         39,
	"erlang":      40,
	"elixir":      41,
	"ocaml":       42,
	"kotlin":      43,
	"brainfuck":   44,
	"fortran":     45,
	"rust":        46,
	"clojure":     47,
}

func init() {
	commandlib.RegisterCommand(commandlib.Command{
		Name:    I18n("Programming"),
		Usage:   I18n("Program in a language."),
		ID:      "programming",
		Matches: []string{"go build"},
		Action:  Programming,
	})
}

func Programming(c commandlib.Context) {
	if c.Arg(0) == "list" || c.NArgs() == 0 {
		var langs []string
		for key := range languages {
			langs = append(langs, key)
		}
		c.SendMessage("primary", c.WrapCodeBlock(strings.Join(langs, "\n")))
		return
	}
	if val, ok := languages[c.Arg(0)]; ok {
		resp, err := http.Get(fmt.Sprintf("https://rextester.com/rundotnet/api?LanguageChoice=%d&Program=%s", val, url.QueryEscape(strings.TrimPrefix(c.RawContent(), c.Arg(0)))))
		if err != nil {
			c.SendMessage("primary", commandlib.ErrorEmbed(fmt.Sprintf(c.I18n("There was an error accessing the API: %s"), err.Error())))
			return
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.SendMessage("primary", commandlib.ErrorEmbed(fmt.Sprintf(c.I18n("There was an error accessing data from the API: %s"), err.Error())))
			return
		}
		var parsedData struct {
			Result  string
			Warning string
			Errors  string
			Stats   string
		}
		json.Unmarshal(data, &parsedData)
		embed := commandlib.Embed{
			Title: commandlib.EmbedHeader{
				Text: c.I18n("Program Results"),
			},
		}
		appendEmbed := func(embed *commandlib.Embed, title, data string) {
			if data != "" {
				embed.Fields = append(embed.Fields, commandlib.EmbedField{
					Title: title,
					Body:  c.WrapCodeBlock(data),
				})
			}
		}
		appendEmbed(&embed, c.I18n("Output"), parsedData.Result)
		appendEmbed(&embed, c.I18n("Warnings"), parsedData.Warning)
		appendEmbed(&embed, c.I18n("Errors"), parsedData.Errors)
		embed.Footer.Text = parsedData.Stats
		if len(embed.Fields) == 0 {
			embed.Body = c.I18n("Your program produced no output.")
		}
		c.SendMessage("primary", embed)
	} else {
		c.SendMessage("primary", commandlib.ErrorEmbed(fmt.Sprintf(c.I18n("%s is not a valid programming language. Use go build list to see a list of programming languages."), c.Arg(0))))
	}
}
