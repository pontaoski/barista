package telegram

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/ini.v1"
)

func Main() {
	fmt.Println("Initializing Barista Telegram...")

	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Failed to load config.ini")
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Section("Bot").Key("telegramtoken").String())

	log.Printf("Authorised on account %s", bot.Self.UserName)
	log.Printf("Barista Telegram is now running")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	go cleaner()

	for update := range updates {
		if update.CallbackQuery != nil {
			if update.CallbackQuery.Message != nil {
				paginatorHandler(update.CallbackQuery.Message.MessageID, update.CallbackQuery.Data)
			}
		}
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}
		switch update.Message.CommandWithAt() {
		case "dnfsearch":
			dnf(update.Message, bot)
		}
	}
}
