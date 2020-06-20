package barista

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/appadeia/barista/barista-go/commandlib/discord"
	"github.com/bwmarrin/discordgo"
)

func discordMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Message.WebhookID != "" {
		return
	}
	discord.DiscordMessage(s, m.Message, m)
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
	discord.DiscordMessage(s, m.Message, m)
}

func discordMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	discord.DeleteDiscordMessage(s, m)
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
	discord.AddHandler(discordMessageDelete)

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
