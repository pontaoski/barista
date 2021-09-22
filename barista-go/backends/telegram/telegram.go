package telegram

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/appadeia/barista/barista-go/commandlib"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramContext struct {
	commandlib.ContextMixin

	bot *tgbotapi.BotAPI
	tm  *tgbotapi.Message
}

type telegramPaginator struct {
	Pages []tgbotapi.MessageConfig
	Index int

	context *TelegramContext
	message *tgbotapi.Message
	bot     *tgbotapi.BotAPI
}

func newPaginator(bot *tgbotapi.BotAPI, context *TelegramContext) telegramPaginator {
	return telegramPaginator{
		bot:     bot,
		context: context,
	}
}

func (p *telegramPaginator) AddPage(msg tgbotapi.MessageConfig) {
	p.Pages = append(p.Pages, msg)
}

func (t TelegramContext) AuthorName() string {
	return t.tm.From.String()
}

func (t TelegramContext) AuthorIdentifier() string {
	return strconv.FormatInt(int64(t.tm.From.ID), 16)
}

func (t TelegramContext) RoomIdentifier() string {
	return strconv.FormatInt(t.tm.Chat.ID, 16)
}

func (t TelegramContext) CommunityIdentifier() string {
	return strconv.FormatInt(t.tm.Chat.ID, 16)
}

func (t TelegramContext) I18n(message string) string {
	return t.I18nInternal(commandlib.GetLanguage(&t), message)
}

func (t TelegramContext) I18nc(context, message string) string {
	return t.I18n(message)
}

func (t *TelegramContext) keyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				t.I18n("Previous"),
				"previous",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				t.I18n("Next"),
				"next",
			),
		),
	)
}

func (t TelegramContext) Backend() commandlib.Backend {
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
	var fields []string
	for _, field := range d.Fields {
		fields = append(fields, fmt.Sprintf("%s: %s", field.Title, field.Body))
	}
	msg := tgbotapi.NewMessage(0, "")
	var sb strings.Builder
	embedTemplate.Execute(&sb, d)
	msg.Text = sb.String()
	msg.ParseMode = tgbotapi.ModeHTML
	return msg
}

func (t TelegramContext) SendTags(_ string, tags []commandlib.Embed) {
	for _, tag := range tags {
		msg := telegramEmbed(tag)
		msg.ChatID = t.tm.Chat.ID
		t.bot.Send(msg)
	}
}

func (t TelegramContext) WrapCodeBlock(code string) string {
	return "<pre>\n" + code + "\n</pre>"
}

func (t TelegramContext) GenerateLink(text, URL string) string {
	return fmt.Sprintf(`<a href="%s">%s</a>`, URL, text)
}

func (t *TelegramContext) SendMessage(_ string, content interface{}) {
	switch a := content.(type) {
	case string:
		msg := tgbotapi.NewMessage(t.tm.Chat.ID, content.(string))
		t.bot.Send(msg)
	case commandlib.Embed:
		msg := telegramEmbed(content.(commandlib.Embed))
		msg.ChatID = t.tm.Chat.ID
		t.bot.Send(msg)
	case commandlib.EmbedList:
		paginator := newPaginator(t.bot, t)
		title := "Item"
		if content.(commandlib.EmbedList).ItemTypeName != "" {
			title = content.(commandlib.EmbedList).ItemTypeName
		}
		for idx, page := range content.(commandlib.EmbedList).Embeds {
			page.Footer.Text = fmt.Sprintf(t.I18n("%s %d out of %d"), title, idx+1, len(content.(commandlib.EmbedList).Embeds))
			msg := telegramEmbed(page)
			msg.ChatID = t.tm.Chat.ID
			paginator.AddPage(msg)
		}
		paginator.Send()
	case commandlib.UnionEmbed:
		t.SendMessage("", content.(commandlib.UnionEmbed).EmbedList)
		return
	case commandlib.File:
		t.bot.Send(tgbotapi.NewDocumentUpload(t.tm.Chat.ID, tgbotapi.FileReader{
			Name:   a.Name,
			Reader: a.Reader,
			Size:   -1,
		}))
	}
}

func (t *TelegramContext) NextResponse() (out chan string) {
	out = make(chan string)
	go func() {
		for {
			select {
			case usermsg := <-waitForTelegramMessage():
				if usermsg.Chat.ID == t.tm.Chat.ID && usermsg.From.ID == t.tm.From.ID {
					out <- usermsg.Text
					return
				}
			}
		}
	}()
	return out
}

func (t *TelegramContext) AwaitResponse(tenpo time.Duration) (response string, ok bool) {
	timeoutChan := make(chan struct{})
	go func() {
		time.Sleep(tenpo)
		timeoutChan <- struct{}{}
	}()
	for {
		select {
		case msg := <-t.NextResponse():
			return msg, true
		case <-timeoutChan:
			return "", false
		}
	}
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
