package telegram

import (
	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/config"
	"github.com/appadeia/barista/barista-go/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var backend = TelegramBackend{}

func init() {
	commandlib.RegisterBackend(&backend)
}

type TelegramBackend struct{}

func (t TelegramBackend) Stats() (r *commandlib.BackendStats) {
	return
}

func (t TelegramBackend) CanGiveStats() bool {
	return false
}

func (t TelegramBackend) IsBotOwner(c commandlib.Context) bool {
	var ctx interface{} = c
	switch casted := ctx.(type) {
	case *MessageTelegramContext:
		return casted.tm.From.ID == config.BotConfig.Owner.Telegram
	case *InlineQueryContext:
		return casted.iq.From.ID == config.BotConfig.Owner.Telegram
	default:
		return false
	}
}

func (t TelegramBackend) Name() string {
	return "Telegram"
}

func (t TelegramBackend) ID() string {
	return "telegram"
}

func (t TelegramBackend) Start(cancel chan struct{}) error {
	bot, err := tgbotapi.NewBotAPI(config.BotConfig.Services.Telegram.Token)
	if err != nil {
		return err
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	log.Info("Telegram session started")
forLoop:
	for {
		select {
		case update := <-updates:
			if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
				TelegramPaginatorHandler(update.CallbackQuery.Message.MessageID, update.CallbackQuery.Data)
			}
			if update.InlineQuery != nil {
				TelegramInlineQuery(bot, update.InlineQuery)
			}
			if update.Message != nil {
				TelegramMessage(bot, update.Message)
			}
		case <-cancel:
			break forLoop
		}
	}
	return nil
}
