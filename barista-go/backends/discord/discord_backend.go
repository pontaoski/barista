package discord

import (
	"context"
	"fmt"

	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/config"
	"github.com/appadeia/barista/barista-go/log"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

// The DiscordBackend handles Discord connections
type DiscordBackend struct {
	token string
	name  string
	s     *state.State
	me    *discord.User
}

var backends = map[discord.UserID]*DiscordBackend{}

func init() {
	for _, token := range config.BotConfig.Services.Discord {
		backend := new(DiscordBackend)
		backend.token = token.Token
		backend.name = token.Name
		commandlib.RegisterBackend(backend)
	}
}

func (d *DiscordBackend) Stats() (r *commandlib.BackendStats) {
	r = &commandlib.BackendStats{}
	guilds, err := d.s.Client.Guilds(0)
	if err != nil {
		return r
	}
	r.Communities = uint64(len(guilds))
	var users uint64
	var all map[discord.UserID]struct{}
	for _, guild := range guilds {
		guildUsers, _ := d.s.Members(guild.ID)
		for _, user := range guildUsers {
			all[user.User.ID] = struct{}{}
		}
	}
	users = uint64(len(all))
	r.Users = users
	return
}

func (d *DiscordBackend) CanGiveStats() bool {
	return true
}

// Name is the name of the Discord backend
func (d *DiscordBackend) Name() string {
	return fmt.Sprintf("Discord (%s)", d.name)
}

func (d *DiscordBackend) ID() string {
	return "discord"
}

func (d *DiscordBackend) IsBotOwner(c commandlib.Context) bool {
	var ctx interface{} = c
	casted := ctx.(*DiscordContext)
	return casted.tm.Author.ID.String() == config.BotConfig.Owner.Discord
}

// Start starts the Discord backend
func (d *DiscordBackend) Start(cancel chan struct{}) error {
	discord := state.New("Bot " + d.token)
	defer discord.Close()

	discord.AddIntents(gateway.IntentGuilds)
	discord.AddIntents(gateway.IntentGuildMembers)
	discord.AddIntents(gateway.IntentGuildMessages)
	discord.AddIntents(gateway.IntentMessageContent)

	discord.AddHandler(d.discordMessageCreate)
	discord.AddHandler(d.discordMessageEdit)
	discord.AddHandler(d.discordMessageDelete)

	d.s = discord
	err := discord.Connect(context.Background())
	if err != nil {
		return err
	}

	d.token = ""
	me, err := d.s.Me()
	if err != nil {
		return err
	}
	backends[me.ID] = d
	d.me = me

	log.Info("Discord (%s) session started", d.name)

	<-cancel
	return nil
}

func (d *DiscordBackend) discordMessageCreate(m *gateway.MessageCreateEvent) {
	if m.Author.ID == d.me.ID {
		return
	}
	DiscordMessage(d, &m.Message)
}

func (d *DiscordBackend) discordMessageEdit(m *gateway.MessageUpdateEvent) {
	if m.Author.ID == d.me.ID {
		return
	}
	DiscordMessage(d, &m.Message)
}

func (d *DiscordBackend) discordMessageDelete(m *gateway.MessageDeleteEvent) {
	DeleteDiscordMessage(d, m)
}

func (d *DiscordBackend) paginator(m *gateway.InteractionCreateEvent) {
	switch e := m.InteractionEvent.Data.(type) {
	case *discord.ButtonInteraction:
		if val, ok := paginatorCache.Get(m.Message.ID); ok {
			if e.CustomID == "previous" {
				val.(*paginator).Prev()
			} else {
				val.(*paginator).Next()
			}
		}
	}
}
