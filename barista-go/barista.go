package barista

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	badger "github.com/dgraph-io/badger"
	"gopkg.in/ini.v1"
)

// Barista's config
var Cfg *ini.File

var cmds map[string]CommandFunc = map[string]CommandFunc{
	"sudo echo":       Echo,
	"dnf search":      Dnf,
	"zypper search":   Dnf,
	"dnf se":          Dnf,
	"zypper se":       Dnf,
	"pacman -Ss":      Dnf,
	"apt search":      Dnf,
	"apt se":          Dnf,
	"apt-get search":  Dnf,
	"apt-get se":      Dnf,
	"sudo profile":    Profile,
	"sudo help":       Help,
	"sudo ddg":        Ddg,
	"dnf repoquery":   DnfRepoQuery,
	"dnf rq":          DnfRepoQuery,
	"sudo about":      About,
	"sudo gsettings":  Gsettings,
	"lutris search":   Lutris,
	"sudo ss":         Screenshot,
	"sudo screenshot": Screenshot,
	"sudo paste":      Paste,
	"sudo embed":      EmbedCmd,
}

var handlers []CommandFunc = []CommandFunc{
	Obs,
	Bodhi,
	Bugzilla,
	Pagure,
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Message.WebhookID != "" {
		return
	}
	for key, val := range cmds {
		if strings.HasPrefix(m.Content, key) {
			go LexedCommandFunction(s, m.Message, val)
		}
	}
	for _, val := range handlers {
		go LexedCommandFunction(s, m.Message, val)
	}
}
func messageEdit(s *discordgo.Session, m *discordgo.MessageUpdate) {
	msg, err := s.ChannelMessage(m.ChannelID, m.ID)
	if err != nil {
		println(err.Error())
		return
	}
	if msg.Author.ID == s.State.User.ID {
		return
	}
	for key, val := range cmds {
		if strings.HasPrefix(m.Content, key) {
			go LexedCommandFunction(s, msg, val)
		}
	}
}

var db *badger.DB

// Main : Call this function to start the bot's main loop.
func Main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("There was a panic!\n\t%s", err)
		}
	}()

	var err error
	Cfg, err = ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Failed to load config.ini")
		os.Exit(1)
	}

	db, err = badger.Open(badger.DefaultOptions("./storage/db"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	discord, err := discordgo.New("Bot " + Cfg.Section("Bot").Key("token").String())
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
	}

	discord.AddHandler(messageCreate)
	discord.AddHandler(messageEdit)

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}

	go Cleaner()

	fmt.Println("Barista Discord is now running.")
	go TelegramMain()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}
