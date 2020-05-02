package barista

import (
	"fmt"

	"github.com/appadeia/barista/barista-go/commandlib"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func TelegramMain() {
	fmt.Println("Initializing Barista Telegram...")

	bot, err := tgbotapi.NewBotAPI(Cfg.Section("Bot").Key("telegramtoken").String())
	if err != nil {
		fmt.Println("Error creating Telegram session: ", err.Error())
		return
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	fmt.Println("Barista Telegram is now running.")
	for update := range updates {
		if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
			commandlib.TelegramPaginatorHandler(update.CallbackQuery.Message.MessageID, update.CallbackQuery.Data)
		}
		if update.Message != nil {
			commandlib.TelegramMessage(bot, update.Message)
		}
	}
}
