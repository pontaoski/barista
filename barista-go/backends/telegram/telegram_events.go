package telegram

import (
	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func TelegramPaginatorHandler(messageID int, direction string) {
	if val, ok := paginatorCache.Get(messageID); ok {
		if direction == "previous" {
			val.(*telegramPaginator).Prev()
		} else {
			val.(*telegramPaginator).Next()
		}
	}
}

func TelegramMessage(b *tgbotapi.BotAPI, m *tgbotapi.Message) {
	for _, handler := range tgHandlers {
		handler.handler(m)
		removeTelegramHandler(handler)
	}
	if cmd, ContextMixin, ok := commandlib.LexCommand(m.Text); ok {
		tg := TelegramContext{}
		tg.ContextMixin = ContextMixin
		tg.ContextMixin.ContextType = commandlib.CreateCommand
		tg.bot = b
		tg.tm = m
		go log.CanPanic(func() {
			b.Send(tgbotapi.NewChatAction(
				m.Chat.ID,
				"typing",
			))
			cmd.Action(&tg)
		})
	} else {
		for _, tc := range commandlib.LexTags(m.Text) {
			tg := TelegramContext{}
			tg.ContextMixin = tc.Context
			tg.bot = b
			tg.tm = m
			go log.CanPanic(func() {
				b.Send(tgbotapi.NewChatAction(
					m.Chat.ID,
					"typing",
				))
				tc.Tag.Action(&tg)
			})
		}
	}
}
