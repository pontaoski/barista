package barista

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/godbus/dbus"
)

func TelegramMain() {
	bot, err := tgbotapi.NewBotAPI(Cfg.Section("Bot").Key("telegramtoken").String())

	if err != nil {
		return
	}

	log.Printf("Authorised on account %s", bot.Self.UserName)
	log.Printf("Barista Telegram is now running")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() { // ignore any non-Message and non-Command Updates
			continue
		}
		switch update.Message.CommandWithAt() {
		case "dnfsearch", "dnfsearch@cafeterabot":
			args := strings.Split(update.Message.CommandArguments(), " ")
			if len(args) < 2 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please specify arguments like so: `dnfsearch distro query`")
				bot.Send(msg)
				continue
			}
			dist, set := resolveDistro(args[0])
			if !set {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "That is not a distro supported by Barista.")
				bot.Send(msg)
				continue
			}
			pkgs, err := TelegramDnfSearch(args[1], dist.queryKitName)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "There was an error searching for packages.")
				bot.Send(msg)
				continue
			}
			if len(pkgs) == 0 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "No matches found.")
				bot.Send(msg)
				continue
			}
			slice := 10
			if len(pkgs) < 10 {
				slice = len(pkgs)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join(pkgs[:slice], "\n"))
			msg.ParseMode = "HTML"
			bot.Send(msg)
			continue
		}
	}
}

func TelegramDnfSearch(search string, distro string) ([]string, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return []string{}, err
	}
	var pkgs [][]interface{}
	obj := conn.Object("com.github.Appadeia.QueryKit", "/com/github/Appadeia/QueryKit")
	err = obj.Call("com.github.Appadeia.QueryKit.SearchPackages", 0, search, distro).Store(&pkgs)
	if err != nil {
		return []string{}, err
	}
	pkgStrings := []string{}
	for _, pkg := range pkgs {
		pkgStrings = append(pkgStrings, fmt.Sprintf("<b>%s</b> — %s", pkg[0].(string), pkg[1].(string)))
	}
	return pkgStrings, nil
}
