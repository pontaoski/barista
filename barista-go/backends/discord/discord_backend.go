package discord

import (
	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/config"
	"github.com/appadeia/barista/barista-go/log"
	"github.com/bwmarrin/discordgo"
)

// The DiscordBackend handles Discord connections
type DiscordBackend struct{}

func init() {
	commandlib.RegisterBackend(DiscordBackend{})
}

// Name is the name of the Discord backend
func (d DiscordBackend) Name() string {
	return "Discord"
}

// Start starts the Discord backend
func (d DiscordBackend) Start(cancel chan struct{}) error {
	discord, err := discordgo.New("Bot " + config.BotConfig.Services.Discord.Token)
	defer discord.Close()
	if err != nil {
		return err
	}

	err = discord.Open()
	if err != nil {
		return err
	}

	log.Info("Discord session started")
	discord.AddHandler(discordMessageCreate)
	discord.AddHandler(discordMessageEdit)
	discord.AddHandler(discordMessageDelete)

	<-cancel
	return nil
}

func discordMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Message.WebhookID != "" {
		return
	}
	DiscordMessage(s, m.Message, m)
}

func discordMessageEdit(s *discordgo.Session, m *discordgo.MessageUpdate) {
	msg, err := s.ChannelMessage(m.ChannelID, m.ID)
	if err != nil {
		return
	}
	if msg.Author.ID == s.State.User.ID {
		return
	}
	DiscordMessage(s, m.Message, m)
}

func discordMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	DeleteDiscordMessage(s, m)
}
