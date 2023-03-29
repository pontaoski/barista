package telegram

import (
	"fmt"
	"strconv"
	"time"

	"github.com/appadeia/barista/barista-go/commandlib"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type MessageTelegramContext struct {
	BaseTelegramContext

	tm *tgbotapi.Message
}

func (t MessageTelegramContext) AuthorName() string {
	return t.tm.From.String()
}

func (t MessageTelegramContext) AuthorIdentifier() string {
	return strconv.FormatInt(int64(t.tm.From.ID), 16)
}

func (t MessageTelegramContext) RoomIdentifier() string {
	return strconv.FormatInt(t.tm.Chat.ID, 16)
}

func (t MessageTelegramContext) CommunityIdentifier() string {
	return strconv.FormatInt(t.tm.Chat.ID, 16)
}

func (t MessageTelegramContext) I18n(message string) string {
	return t.I18nInternal(commandlib.GetLanguage(&t), message)
}

func (t MessageTelegramContext) I18nc(context, message string) string {
	return t.I18n(message)
}

func (t *MessageTelegramContext) keyboard() tgbotapi.InlineKeyboardMarkup {
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
func (t MessageTelegramContext) SendTags(_ string, tags []commandlib.Embed) {
	for _, tag := range tags {
		msg := telegramEmbed(tag)
		msg.ChatID = t.tm.Chat.ID
		t.bot.Send(msg)
	}
}

func (t *MessageTelegramContext) SendMessage(_ string, content interface{}) {
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

func (t *MessageTelegramContext) NextResponse() (out chan string) {
	out = make(chan string)
	go func() {
		for usermsg := range waitForTelegramMessage() {
			if usermsg.Chat.ID == t.tm.Chat.ID && usermsg.From.ID == t.tm.From.ID {
				out <- usermsg.Text
				return
			}
		}
	}()
	return out
}

func (t *MessageTelegramContext) AwaitResponse(tenpo time.Duration) (response string, ok bool) {
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
