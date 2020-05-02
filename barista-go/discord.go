package barista

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/bwmarrin/discordgo"
)

func discordMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Message.WebhookID != "" {
		return
	}
	commandlib.DiscordMessage(s, m.Message)
}

func discordMessageEdit(s *discordgo.Session, m *discordgo.MessageUpdate) {
	msg, err := s.ChannelMessage(m.ChannelID, m.ID)
	if err != nil {
		println(err.Error())
		return
	}
	if msg.Author.ID == s.State.User.ID {
		return
	}
	commandlib.DiscordMessage(s, m.Message)
}

func DiscordMain() {
	fmt.Println("Initializing Barista Discord...")

	discord, err := discordgo.New("Bot " + Cfg.Section("Bot").Key("token").String())
	if err != nil {
		fmt.Println("Error creating Discord session: ", err.Error())
		return
	}

	discord.AddHandler(discordMessageCreate)
	discord.AddHandler(discordMessageEdit)

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}

	fmt.Println("Barista Discord is now running.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}
