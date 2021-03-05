package harmony

import (
	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/log"
	"github.com/harmony-development/shibshib"
	types "github.com/harmony-development/shibshib/gen/harmonytypes/v1"
)

// DeleteMessage handles a deleted message
func (b *Backend) DeleteMessage(c *shibshib.Client, m *types.Message) {
	if val, ok := commandCache.Get(m.MessageId); ok {
		tmp := val.(*Context)
		tmp.ContextMixin.ContextType = commandlib.DeleteCommand
		if tmp.Action.DeleteAction != nil {
			log.CanPanic(func() {
				tmp.Action.DeleteAction(tmp)
			})
		}

	}
}

// Message handles a created or edited message
func (b *Backend) Message(c *shibshib.Client, m *types.Message) {
	strip := m.Content
	if val, ok := commandCache.Get(m.MessageId); ok {
		if cmd, contextMixin, ok := commandlib.LexCommand(strip); ok {
			contextMixin.ContextType = commandlib.EditCommand
			tmp := val.(*Context)
			tmp.ContextMixin = contextMixin
			go log.CanPanic(func() {
				cmd.Action(tmp)
			})
		}
	} else {
		if cmd, contextMixin, ok := commandlib.LexCommand(strip); ok {
			dc := buildContext(contextMixin, b, c, m)
			contextMixin.ContextType = commandlib.CreateCommand
			commandCache.Add(m.MessageId, &dc)
			go log.CanPanic(func() {
				cmd.Action(&dc)
			})
		}
	}

	if val, ok := tagCache.Get(m.MessageId); ok {
		for _, tag := range commandlib.LexTags(strip) {
			tmp := val.(*Context)
			tmp.ContextMixin = tag.Context
			go log.CanPanic(func() {
				tag.Tag.Action(tmp)
			})
		}
	} else {
		for _, tag := range commandlib.LexTags(strip) {
			dc := buildContext(tag.Context, b, c, m)
			tagCache.Add(m.MessageId, &dc)
			go log.CanPanic(func() {
				tag.Tag.Action(&dc)
			})
		}
	}
}
