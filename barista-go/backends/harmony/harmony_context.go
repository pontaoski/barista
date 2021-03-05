package harmony

import (
	"fmt"
	"strconv"
	"time"

	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/log"
	"github.com/harmony-development/shibshib"
	chatv1 "github.com/harmony-development/shibshib/gen/chat/v1"
	types "github.com/harmony-development/shibshib/gen/harmonytypes/v1"
)

// Context is a context holding information about a Harmony command
type Context struct {
	commandlib.ContextMixin
	b  *Backend
	c  *shibshib.Client
	tm *types.Message
}

func buildContext(ctx commandlib.ContextMixin, b *Backend, c *shibshib.Client, m *types.Message) Context {
	dc := Context{
		ContextMixin: ctx,
	}
	dc.b = b
	dc.c = c
	dc.tm = m
	return dc
}

func (c *Context) AuthorIdentifier() string {
	return strconv.FormatUint(c.tm.AuthorId, 10)
}

func (c *Context) AuthorName() string {
	return c.c.UsernameFor(c.tm)
}

func (c *Context) AwaitResponse(tenpo time.Duration) (content string, ok bool) {
	timeoutChan := make(chan struct{})
	go func() {
		time.Sleep(tenpo)
		timeoutChan <- struct{}{}
	}()
	for {
		select {
		case msg := <-c.NextResponse():
			return msg, true
		case <-timeoutChan:
			return "", false
		}
	}
}

func (c *Context) Backend() commandlib.Backend {
	return c.b
}

func (c *Context) CommunityIdentifier() string {
	return strconv.FormatUint(c.tm.GuildId, 10)
}

func (c *Context) GenerateLink(text string, URL string) string {
	return fmt.Sprintf("(%s)[%s]", text, URL)
}

func (c *Context) I18n(message string) string {
	return c.I18nInternal(commandlib.GetLanguage(c), message)
}

func (c *Context) I18nc(context, message string) string {
	return c.I18n(message)
}

func waitForMessage(c *shibshib.Client) chan *types.Message {
	channel := make(chan *types.Message)
	var f func()
	f = func() {
		c.HandleOnce(func(ev *types.Message) {
			channel <- ev
		})
	}
	f()
	return channel
}

func (c *Context) NextResponse() chan string {
	out := make(chan string)
	go func() {
		for {
			select {
			case usermsg := <-waitForMessage(c.c):
				if usermsg.AuthorId == c.tm.AuthorId && usermsg.ChannelId == usermsg.ChannelId {
					out <- usermsg.Content
					return
				}
			}
		}
	}()
	return out
}

func (c *Context) RoomIdentifier() string {
	return strconv.FormatUint(c.tm.ChannelId, 10)
}

func convert(embed commandlib.Embed) *types.Embed {
	return &types.Embed{
		Body:  embed.Body,
		Color: int64(embed.Colour),
		Title: embed.Title.Text,
		Header: &types.EmbedHeading{
			Text: embed.Header.Text,
			Url:  embed.Header.URL,
			Icon: embed.Header.Icon,
		},
		Footer: &types.EmbedHeading{
			Text: embed.Footer.Text,
			Url:  embed.Footer.URL,
			Icon: embed.Footer.Icon,
		},
		Fields: func() (f []*types.EmbedField) {
			for _, field := range embed.Fields {
				f = append(f, &types.EmbedField{
					Title:        field.Title,
					Body:         field.Body,
					Presentation: types.FieldPresentation_Data,
				})
			}
			return
		}(),
	}
}

func (c *Context) SendMessage(id string, content interface{}) {
	switch content := content.(type) {
	case string:
		_, err := c.c.ChatKit.SendMessage(&chatv1.SendMessageRequest{
			GuildId:   c.tm.GuildId,
			ChannelId: c.tm.ChannelId,
			Content:   content,
		})
		if err != nil {
			log.Error("%+v", err)
		}
	case commandlib.Embed:
		_, err := c.c.ChatKit.SendMessage(&chatv1.SendMessageRequest{
			GuildId:   c.tm.GuildId,
			ChannelId: c.tm.ChannelId,
			Embeds: []*types.Embed{
				convert(content),
			},
		})
		if err != nil {
			log.Error("%+v", err)
		}
	case commandlib.EmbedList:
		_, err := c.c.ChatKit.SendMessage(&chatv1.SendMessageRequest{
			GuildId:   c.tm.GuildId,
			ChannelId: c.tm.ChannelId,
			Embeds: func() (r []*types.Embed) {
				for _, embed := range content.Embeds {
					r = append(r, convert(embed))
				}
				return
			}(),
		})
		if err != nil {
			log.Error("%+v", err)
		}
	case commandlib.UnionEmbed:
		c.SendMessage(id, content.EmbedList)
		return
	}
}

func (c *Context) SendTags(id string, tags []commandlib.Embed) {
	c.c.ChatKit.SendMessage(&chatv1.SendMessageRequest{
		GuildId:   c.tm.GuildId,
		ChannelId: c.tm.ChannelId,
		Embeds: func() (r []*types.Embed) {
			for _, embed := range tags {
				r = append(r, convert(embed))
			}
			return
		}(),
	})
}

func (c *Context) WrapCodeBlock(code string) string {
	return fmt.Sprintf("```%s```", code)
}
