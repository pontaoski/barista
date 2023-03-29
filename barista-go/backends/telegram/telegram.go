package telegram

import (
	"fmt"
	"strings"
	"sync"
	"text/template"

	"github.com/appadeia/barista/barista-go/commandlib"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BaseTelegramContext struct {
	commandlib.ContextMixin

	bot *tgbotapi.BotAPI
}

type telegramPaginator struct {
	Pages []tgbotapi.MessageConfig
	Index int

	context *MessageTelegramContext
	message *tgbotapi.Message
	bot     *tgbotapi.BotAPI
}

func newPaginator(bot *tgbotapi.BotAPI, context *MessageTelegramContext) telegramPaginator {
	return telegramPaginator{
		bot:     bot,
		context: context,
	}
}

func (p *telegramPaginator) AddPage(msg tgbotapi.MessageConfig) {
	p.Pages = append(p.Pages, msg)
}

func (t BaseTelegramContext) Backend() commandlib.Backend {
	return backend
}

type templater map[string]string

func (t templater) eval() *template.Template {
	tmpl := template.Must(template.New("").Parse(t[""]))
	delete(t, "")
	for key, val := range t {
		template.Must(tmpl.New(key).Parse(val))
	}
	return tmpl
}

var (
	embedTemplate = templater{
		"header-strong": `{{ if .URL }}<a href="{{ .URL }}"><b><u>{{ .Text }}</u></b></a>{{ else }}<b>{{ .Text }}</b>{{ end }}`,
		//
		//
		//
		"header": `{{ if .URL }}<a href="{{ .URL }}"><i><u>{{ .Text }}</u></i></a>{{ else }}<i>{{ .Text }}</i>{{ end }}`,
		//
		//
		//
		"field": `<u>{{ .Title }}</u>
{{ .Body }}`,
		//
		//
		//
		"": `{{ if .Header }}{{ template "header" .Header }}
{{ end }}
{{ template "header-strong" .Title }}

{{ .Body }}
{{ if .Fields }}
{{ range .Fields }}{{ template "field" . }}

{{ end }}{{ end }}{{ if .Footer }}{{ template "header" .Footer }}{{ end }}`,
	}.eval()
)

func telegramEmbed(d commandlib.Embed) tgbotapi.MessageConfig {
	d.Truncate()
	msg := tgbotapi.NewMessage(0, "")
	var sb strings.Builder
	embedTemplate.Execute(&sb, d)
	msg.Text = sb.String()
	msg.ParseMode = tgbotapi.ModeHTML
	return msg
}

func (t BaseTelegramContext) WrapCodeBlock(code string) string {
	return "<pre>\n" + code + "\n</pre>"
}

func (t BaseTelegramContext) GenerateLink(text, URL string) string {
	return fmt.Sprintf(`<a href="%s">%s</a>`, URL, text)
}

var tgHandlers []*tgEventHandlerInstance
var tgHandlerMutex = sync.Mutex{}

type tgEventHandlerInstance struct {
	handler func(m *tgbotapi.Message)
}

func removeTelegramHandler(ehi *tgEventHandlerInstance) {
	tgHandlerMutex.Lock()
	defer tgHandlerMutex.Unlock()
	for idx, handler := range tgHandlers {
		if handler == ehi {
			tgHandlers = append(tgHandlers[:idx], tgHandlers[idx+1:]...)
		}
	}
}

func addTelegramHandlerOnce(input func(m *tgbotapi.Message)) {
	tgHandlerMutex.Lock()
	defer tgHandlerMutex.Unlock()
	ehi := tgEventHandlerInstance{input}
	tgHandlers = append(tgHandlers, &ehi)
}

func waitForTelegramMessage() chan *tgbotapi.Message {
	channel := make(chan *tgbotapi.Message)
	addTelegramHandlerOnce(func(m *tgbotapi.Message) {
		channel <- m
	})
	return channel
}
