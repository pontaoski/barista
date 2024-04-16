package discord

import (
	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/log"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

// DeleteDiscordMessage handles a deleted Discord message
func DeleteDiscordMessage(s *DiscordBackend, m *gateway.MessageDeleteEvent) {
	if val, ok := commandCache.Get(m.ID); ok {
		tmp := val.(*DiscordContext)
		tmp.ContextMixin.ContextType = commandlib.DeleteCommand
		if tmp.Action.DeleteAction != nil {
			log.CanPanic(func() {
				tmp.Action.DeleteAction(tmp)
			})
		}
		for _, paginator := range tmp.paginators {
			paginator.Inactive()
			s.s.Client.DeleteMessage(paginator.message.ChannelID, paginator.message.ID, api.AuditLogReason(""))
		}
		for _, message := range tmp.pm {
			if message != nil {
				s.s.Client.DeleteMessage(message.ChannelID, message.ID, api.AuditLogReason(""))
			}
		}
	}
}

// DiscordMessage handles a created or edited Discord message
func DiscordMessage(d *DiscordBackend, m *discord.Message) {
	strip := m.Content
	if val, ok := commandCache.Get(m.ID); ok {
		if cmd, contextMixin, ok := commandlib.LexCommand(strip); ok {
			contextMixin.ContextType = commandlib.EditCommand
			tmp := val.(*DiscordContext)
			tmp.ContextMixin = contextMixin
			go log.CanPanic(func() {
				cmd.Action(tmp)
			})
		}
	} else {
		if cmd, contextMixin, ok := commandlib.LexCommand(strip); ok {
			dc := buildContext(contextMixin, d, m)
			contextMixin.ContextType = commandlib.CreateCommand
			commandCache.Add(m.ID, &dc)
			go log.CanPanic(func() {
				cmd.Action(&dc)
			})
		}
	}

	if val, ok := tagCache.Get(m.ID); ok {
		for _, tag := range commandlib.LexTags(strip) {
			tmp := val.(*DiscordContext)
			tmp.ContextMixin = tag.Context
			go log.CanPanic(func() {
				tag.Tag.Action(tmp)
			})
		}
	} else {
		for _, tag := range commandlib.LexTags(strip) {
			dc := buildContext(tag.Context, d, m)
			tagCache.Add(m.ID, &dc)
			go log.CanPanic(func() {
				tag.Tag.Action(&dc)
			})
		}
	}
}
