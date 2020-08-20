package irc

import (
	"strings"
	"time"

	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/lrstanley/girc"
)

type IRCContext struct {
	commandlib.ContextMixin
	tm string
	cl *girc.Client
	ev girc.Event
}

func (i IRCContext) AuthorIdentifier() string {
	return i.ev.Source.String()
}

func (i IRCContext) AuthorName() string {
	return i.ev.Source.String()
}

func (i IRCContext) AwaitResponse(tenpo time.Duration) (content string, ok bool) {
	timeoutChan := make(chan struct{})
	go func() {
		time.Sleep(tenpo)
		timeoutChan <- struct{}{}
	}()
	for {
		select {
		case msg := <-NextMessage():
			if msg.Source.Equals(i.ev.Source) {
				return msg.Last(), true
			}
		case <-timeoutChan:
			return "", false
		}
	}
}

func (i IRCContext) Backend() commandlib.Backend {
	return backend
}

func (i IRCContext) CommunityIdentifier() string {
	return i.ev.Params[0]
}

func (i IRCContext) GenerateLink(text string, URL string) string {
	return URL
}

func (i IRCContext) I18n(message string) string {
	return message
}

func (i IRCContext) I18nc(context, message string) string {
	return message
}

func (i IRCContext) NextResponse() chan string {
	retChan := make(chan string)
	go func() {
		for {
			select {
			case msg := <-NextMessage():
				if msg.Source.Equals(i.ev.Source) {
					retChan <- msg.Last()
					return
				}
			}
		}
	}()
	return retChan
}

func (i IRCContext) RoomIdentifier() string {
	return i.ev.Params[0]
}

func (i IRCContext) SendMessage(_ string, content interface{}) {
	switch content.(type) {
	case string:
		i.cl.Cmd.ReplyTo(i.ev, content.(string))
	case commandlib.Embed:
		msg := ircEmbed(content.(commandlib.Embed))
		for _, str := range msg {
			i.cl.Cmd.ReplyTo(i.ev, str)
		}
	case commandlib.EmbedList:
		for _, page := range content.(commandlib.EmbedList).Embeds {
			msg := ircEmbed(page)
			for _, str := range msg {
				i.cl.Cmd.ReplyTo(i.ev, str)
			}
		}
	case commandlib.UnionEmbed:
		i.SendMessage("", content.(commandlib.UnionEmbed).EmbedList)
	}
}

func (i IRCContext) SendTags(_ string, tags []commandlib.Embed) {
	for _, page := range tags {
		msg := ircEmbed(page)
		for _, str := range msg {
			i.cl.Cmd.ReplyTo(i.ev, str)
		}
	}
}

func (i IRCContext) Mentions() (ret []string) {
	channel := i.cl.LookupChannel(i.ev.Params[len(i.ev.Params)-2])
	for _, mem := range channel.UserList {
		for _, word := range strings.Fields(i.RawData) {
			if mem == word {
				ret = append(ret, mem)
			}
		}
	}
	return
}

func (i IRCContext) DisplayNameForID(id string) string {
	return id
}

func (i IRCContext) WrapCodeBlock(code string) string {
	return "[Barista IRC does not support codeblocks. Please use Barista from another service to view codeblocks.]"
}
