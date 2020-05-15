package commandlib

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramContext struct {
	contextImpl

	bot *tgbotapi.BotAPI
	tm  *tgbotapi.Message
}

type telegramPaginator struct {
	Pages []tgbotapi.MessageConfig
	Index int

	context  *TelegramContext
	message  *tgbotapi.Message
	bot      *tgbotapi.BotAPI
	lastused time.Time
}

func newPaginator(bot *tgbotapi.BotAPI, context *TelegramContext) telegramPaginator {
	return telegramPaginator{
		bot:     bot,
		context: context,
	}
}

func TelegramPaginatorHandler(messageID int, direction string) {
	if val, ok := telegramPaginators[messageID]; ok {
		if direction == "previous" {
			val.Prev()
		} else {
			val.Next()
		}
	}
}

func init() {
	go telegramCleaner()
}

func telegramCleaner() {
	for {
		time.Sleep(5 * time.Minute)
		var rmkeys []int
		for key, cmd := range telegramPaginators {
			if time.Now().Sub(cmd.lastused) >= 10*time.Minute {
				rmkeys = append(rmkeys, key)
			}
		}
		for _, key := range rmkeys {
			telegramMutex.Lock()
			delete(telegramPaginators, key)
			telegramMutex.Unlock()
		}
	}
}

func (p *telegramPaginator) AddPage(msg tgbotapi.MessageConfig) {
	p.Pages = append(p.Pages, msg)
}

func (t TelegramContext) RoomIdentifier() string {
	return strconv.FormatInt(t.tm.Chat.ID, 16)
}

var i18nschema = Schema{
	Name:           "Preferred Locale",
	Description:    "The preferred language of this channel.",
	ID:             "locale",
	DefaultValue:   "en",
	PossibleValues: []string{"en", "de", "es", "fr", "it", "nl", "pl", "tpo"},
}

func (t TelegramContext) I18n(message string) string {
	return t.I18nInternal(i18nschema.ReadValue(&t), message)
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

func (p *telegramPaginator) Send() {
	p.Index = 0
	send := p.Pages[p.Index]
	send.ReplyMarkup = p.context.keyboard()
	msg, err := p.bot.Send(send)
	if err == nil {
		p.message = &msg
		p.lastused = time.Now()
		telegramPaginators[msg.MessageID] = p
	}
}

func (p *telegramPaginator) Prev() {
	p.Index--
	if p.Index < 0 {
		p.Index = len(p.Pages) - 1
	}
	send := p.Pages[p.Index]
	send.ReplyMarkup = p.context.keyboard()
	edit := tgbotapi.NewEditMessageText(p.message.Chat.ID, p.message.MessageID, "")
	edit.Text = p.Pages[p.Index].Text
	edit.ParseMode = p.Pages[p.Index].ParseMode
	kb := p.context.keyboard()
	edit.ReplyMarkup = &kb
	msg, err := p.bot.Send(edit)
	if err != nil {
		p.message = &msg
		p.lastused = time.Now()
	}
}

func (p *telegramPaginator) Next() {
	p.Index++
	if p.Index+1 > len(p.Pages) {
		p.Index = 0
	}
	send := p.Pages[p.Index]
	send.ReplyMarkup = p.context.keyboard()
	edit := tgbotapi.NewEditMessageText(p.message.Chat.ID, p.message.MessageID, "")
	edit.Text = p.Pages[p.Index].Text
	edit.ParseMode = p.Pages[p.Index].ParseMode
	kb := p.context.keyboard()
	edit.ReplyMarkup = &kb
	msg, err := p.bot.Send(edit)
	if err != nil {
		p.message = &msg
		p.lastused = time.Now()
	}
}

var telegramMutex = &sync.Mutex{}
var telegramPaginators map[int]*telegramPaginator = make(map[int]*telegramPaginator)

func telegramEmbed(d Embed) tgbotapi.MessageConfig {
	d.Truncate()
	var fields []string
	for _, field := range d.Fields {
		fields = append(fields, fmt.Sprintf("%s: %s", field.Title, field.Body))
	}
	msg := tgbotapi.NewMessage(0, "")
	msg.Text = fmt.Sprintf(`<i>%s</i>

<b>%s</b>
<i>%s</i>
%s

%s`, d.Header.Text, d.Title.Text, d.Body, strings.Join(fields, "\n"), d.Footer.Text)
	msg.ParseMode = tgbotapi.ModeHTML
	return msg
}

func (t TelegramContext) SendTags(_ string, tags []Embed) {
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
	switch content.(type) {
	case string:
		msg := tgbotapi.NewMessage(t.tm.Chat.ID, content.(string))
		t.bot.Send(msg)
	case Embed:
		msg := telegramEmbed(content.(Embed))
		msg.ChatID = t.tm.Chat.ID
		t.bot.Send(msg)
	case EmbedList:
		telegramMutex.Lock()
		defer telegramMutex.Unlock()
		paginator := newPaginator(t.bot, t)
		title := "Item"
		if content.(EmbedList).ItemTypeName != "" {
			title = content.(EmbedList).ItemTypeName
		}
		for idx, page := range content.(EmbedList).Embeds {
			page.Footer.Text = fmt.Sprintf(t.I18n("%s %d out of %d"), title, idx+1, len(content.(EmbedList).Embeds))
			msg := telegramEmbed(page)
			msg.ChatID = t.tm.Chat.ID
			paginator.AddPage(msg)
		}
		paginator.Send()
	case UnionEmbed:
		t.SendMessage("", content.(UnionEmbed).EmbedList)
		return
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

func TelegramMessage(b *tgbotapi.BotAPI, m *tgbotapi.Message) {
	for _, handler := range tgHandlers {
		handler.handler(m)
		removeTelegramHandler(handler)
	}
	if cmd, contextImpl, ok := lexContent(m.Text); ok {
		tg := TelegramContext{}
		tg.contextImpl = contextImpl
		tg.contextImpl.contextType = CreateCommand
		tg.bot = b
		tg.tm = m
		go cmd.Action(&tg)
	} else {
		for _, tc := range lexTags(m.Text) {
			tg := TelegramContext{}
			tg.contextImpl = tc.Context
			tg.bot = b
			tg.tm = m
			go tc.Tag.Action(&tg)
		}
	}
}
