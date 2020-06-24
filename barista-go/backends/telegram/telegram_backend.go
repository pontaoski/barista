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

func (t TelegramBackend) IsBotOwner(c commandlib.Context) bool {
	var ctx interface{} = c
	casted := ctx.(*TelegramContext)
	return casted.tm.From.ID == config.BotConfig.Owner.Telegram
}

func (t TelegramBackend) Name() string {
	return "Telegram"
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
			if update.Message != nil {
				TelegramMessage(bot, update.Message)
			}
		case <-cancel:
			break forLoop
		}
	}
	return nil
}
