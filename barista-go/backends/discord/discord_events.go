package discord

import (
	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/log"
	"github.com/bwmarrin/discordgo"
)

// DeleteDiscordMessage handles a deleted Discord message
func DeleteDiscordMessage(s *discordgo.Session, m *discordgo.MessageDelete) {
	if val, ok := commandCache.Get(m.ID); ok {
		tmp := val.(*DiscordContext)
		tmp.ContextMixin.ContextType = commandlib.DeleteCommand
		if tmp.Action.DeleteAction != nil {
			log.CanPanic(func() {
				tmp.Action.DeleteAction(tmp)
			})
		}
		for _, paginator := range tmp.paginators {
			paginator.DeleteMessageWhenDone = true
			paginator.Widget.Close <- true
		}
		for _, message := range tmp.pm {
			if message != nil {
				s.ChannelMessageDelete(message.ChannelID, message.ID)
			}
		}
	}
}

// DiscordMessage handles a created or edited Discord message
func DiscordMessage(s *discordgo.Session, m *discordgo.Message, ev interface{}) {
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
			dc := buildContext(contextMixin, s, m)
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
			dc := buildContext(tag.Context, s, m)
			tagCache.Add(m.ID, &dc)
			go log.CanPanic(func() {
				tag.Tag.Action(&dc)
			})
		}
	}
}
