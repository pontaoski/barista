package telegram

import (
	"strconv"

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

func Stalk(tg *TelegramContext) {
	commandlib.StoreData(tg, "__telegram_id"+strconv.FormatInt(int64(tg.tm.From.ID), 10), tg.tm.From.String(), commandlib.Global)
	if tg.tm.From.UserName != "" {
		commandlib.StoreData(tg, "__telegram_user"+tg.tm.From.UserName, strconv.FormatInt(int64(tg.tm.From.ID), 10), commandlib.Global)
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
			cmd.Action(&tg)
		})
	} else {
		for _, tc := range commandlib.LexTags(m.Text) {
			tg := TelegramContext{}
			tg.ContextMixin = tc.Context
			tg.bot = b
			tg.tm = m
			go log.CanPanic(func() {
				tc.Tag.Action(&tg)
			})
		}
	}
}
