package irc

import (
	"strings"

	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/config"
	"github.com/appadeia/barista/barista-go/log"
	"github.com/lrstanley/girc"
)

// The IRCBackend handles Discord connections
type IRCBackend struct {
}

var backend = IRCBackend{}

func init() {
	commandlib.RegisterBackend(&backend)
}

// Name is the name of the Discord backend
func (i IRCBackend) Name() string {
	return "IRC"
}

// IsBotOwner always returns false because IRC is insecure
func (i IRCBackend) IsBotOwner(c commandlib.Context) bool {
	return false
}

// Start starts the Discord backend
func (i IRCBackend) Start(cancel chan struct{}) error {
	client := girc.New(girc.Config{
		Server:     config.BotConfig.Services.IRC.Server,
		Port:       6667,
		Nick:       config.BotConfig.Services.IRC.Nickname,
		User:       config.BotConfig.Services.IRC.Username,
		AllowFlood: false,
	})
	client.Handlers.Add(girc.CONNECTED, func(c *girc.Client, e girc.Event) {
		c.Cmd.Join(config.BotConfig.Services.IRC.Channels...)
	})
	client.Handlers.Add(girc.PRIVMSG, IRCOnetimeHandler)
	client.Handlers.Add(girc.PRIVMSG, func(c *girc.Client, e girc.Event) {
		msg := e.Last()
		if strings.HasPrefix(msg, "!join") {
			channel := strings.TrimSpace(strings.TrimPrefix(msg, "!join"))
			c.Cmd.Join(channel)
		}
		if cmd, ContextMixin, ok := commandlib.LexCommand(msg); ok {
			irc := IRCContext{}
			irc.ContextMixin = ContextMixin
			irc.ContextMixin.ContextType = commandlib.CreateCommand
			irc.tm = msg
			irc.ev = e
			irc.cl = client
			go log.CanPanic(func() {
				cmd.Action(&irc)
			})
		} else {
			for _, tc := range commandlib.LexTags(msg) {
				irc := IRCContext{}
				irc.ContextMixin = tc.Context
				irc.tm = msg
				irc.ev = e
				irc.cl = client
				go log.CanPanic(func() {
					tc.Tag.Action(&irc)
				})
			}
		}
	})
	canal := make(chan error, 1)
	go func() {
		err := client.Connect()
		canal <- err
	}()
	log.Info("IRC session started")
	for {
		select {
		case err := <-canal:
			return err
		case <-cancel:
			return nil
		}
	}
}
