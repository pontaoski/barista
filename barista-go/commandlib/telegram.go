package commandlib

import (
	"fmt"
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

	message  *tgbotapi.Message
	bot      *tgbotapi.BotAPI
	lastused time.Time
}

func newPaginator(bot *tgbotapi.BotAPI) telegramPaginator {
	return telegramPaginator{
		bot: bot,
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

var keyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Previous", "previous"),
		tgbotapi.NewInlineKeyboardButtonData("Next", "next"),
	),
)

func (p *telegramPaginator) Send() {
	p.Index = 0
	send := p.Pages[p.Index]
	send.ReplyMarkup = keyboard
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
	send.ReplyMarkup = keyboard
	edit := tgbotapi.NewEditMessageText(p.message.Chat.ID, p.message.MessageID, "")
	edit.Text = p.Pages[p.Index].Text
	edit.ParseMode = p.Pages[p.Index].ParseMode
	edit.ReplyMarkup = &keyboard
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
	send.ReplyMarkup = keyboard
	edit := tgbotapi.NewEditMessageText(p.message.Chat.ID, p.message.MessageID, "")
	edit.Text = p.Pages[p.Index].Text
	edit.ParseMode = p.Pages[p.Index].ParseMode
	edit.ReplyMarkup = &keyboard
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
		paginator := newPaginator(t.bot)
		title := "Item"
		if content.(EmbedList).ItemTypeName != "" {
			title = content.(EmbedList).ItemTypeName
		}
		for idx, page := range content.(EmbedList).Embeds {
			page.Footer.Text = fmt.Sprintf("%s %d out of %d", title, idx+1, len(content.(EmbedList).Embeds))
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

func TelegramMessage(b *tgbotapi.BotAPI, m *tgbotapi.Message) {
	if cmd, contextImpl, ok := lexContent(m.Text); ok {
		tg := TelegramContext{}
		tg.contextImpl = contextImpl
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
