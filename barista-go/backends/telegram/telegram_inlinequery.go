package telegram

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/appadeia/barista/barista-go/commandlib"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type InlineQueryContext struct {
	BaseTelegramContext

	iq *tgbotapi.InlineQuery
}

var _ commandlib.Context = &InlineQueryContext{}

func (t *InlineQueryContext) AuthorIdentifier() string {
	return strconv.FormatInt(int64(t.iq.From.ID), 16)
}

func (t *InlineQueryContext) AuthorName() string {
	return t.iq.From.String()
}

func (t *InlineQueryContext) AwaitResponse(time time.Duration) (content string, ok bool) {
	return "", false
}

func (t *InlineQueryContext) CommunityIdentifier() string {
	return t.AuthorIdentifier()
}

func (t *InlineQueryContext) I18n(message string) string {
	return t.I18nInternal(commandlib.GetLanguage(t), message)
}

func (t *InlineQueryContext) I18nc(context string, message string) string {
	return t.I18n(message)
}

func (t *InlineQueryContext) NextResponse() chan string {
	return nil
}

func (t *InlineQueryContext) RoomIdentifier() string {
	return t.AuthorIdentifier()
}

func (t *InlineQueryContext) SendMessage(id string, content interface{}) {
	switch content := content.(type) {
	case string:
		t.bot.AnswerInlineQuery(tgbotapi.InlineConfig{
			InlineQueryID: t.iq.ID,
			Results: []interface{}{
				tgbotapi.NewInlineQueryResultArticle(
					randSeq(16),
					t.Command().Name,
					content,
				),
			},
		})
	case commandlib.Embed:
		msg := telegramEmbed(content)
		t.bot.AnswerInlineQuery(tgbotapi.InlineConfig{
			InlineQueryID: t.iq.ID,
			Results: []interface{}{
				tgbotapi.NewInlineQueryResultArticleHTML(
					randSeq(16),
					t.Command().Name,
					msg.Text,
				),
			},
		})
	case commandlib.EmbedList:
		var results []interface{}
		for _, page := range content.Embeds {
			msg := telegramEmbed(page)
			results = append(results, tgbotapi.NewInlineQueryResultArticleHTML(
				randSeq(16),
				content.ItemTypeName+" - "+page.Title.Text,
				msg.Text,
			))
		}
		t.bot.AnswerInlineQuery(tgbotapi.InlineConfig{
			InlineQueryID: t.iq.ID,
			Results:       results,
		})
	case commandlib.UnionEmbed:
		t.SendMessage(id, content.EmbedList)
		return
	case commandlib.File:
		t.SendMessage(id, "Sorry, files aren't supported with Barista inline.")
	}
}

func (t *InlineQueryContext) SendTags(id string, tags []commandlib.Embed) {
	panic("unimplemented")
}
