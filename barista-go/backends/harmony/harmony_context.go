package harmony

import (
	"fmt"
	"strconv"
	"time"

	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/log"
	"github.com/harmony-development/shibshib"
	chatv1 "github.com/harmony-development/shibshib/gen/chat/v1"
)

// Context is a context holding information about a Harmony command
type Context struct {
	commandlib.ContextMixin
	b  *Backend
	c  *shibshib.Client
	tm *shibshib.LocatedMessage
}

type mgid struct {
	chatv1.MessageWithId

	guildID   uint64
	channelID uint64
}

func buildContext(ctx commandlib.ContextMixin, b *Backend, c *shibshib.Client, m *shibshib.LocatedMessage) Context {
	dc := Context{
		ContextMixin: ctx,
	}
	dc.b = b
	dc.c = c
	dc.tm = m
	return dc
}

func (c *Context) AuthorIdentifier() string {
	return strconv.FormatUint(c.tm.Message.AuthorId, 10)
}

func (c *Context) AuthorName() string {
	return c.c.UsernameFor(c.tm.Message)
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
	return strconv.FormatUint(c.tm.GuildID, 10)
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

func waitForMessage(c *shibshib.Client) chan *shibshib.LocatedMessage {
	channel := make(chan *shibshib.LocatedMessage)
	var f func()
	f = func() {
		c.HandleOnce(func(ev *shibshib.LocatedMessage) {
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
				if usermsg.Message.AuthorId == c.tm.Message.AuthorId && usermsg.ChannelID == c.tm.ChannelID {
					out <- usermsg.GetMessage().GetContent().GetTextMessage().String()
					return
				}
			}
		}
	}()
	return out
}

func (c *Context) RoomIdentifier() string {
	return strconv.FormatUint(c.tm.ChannelID, 10)
}

func ftext(s string) *chatv1.FormattedText {
	return &chatv1.FormattedText{Text: s}
}

func convert(embed commandlib.Embed) *chatv1.Embed {
	c := int32(embed.Colour)
	return &chatv1.Embed{
		Body:  ftext(embed.Body),
		Color: &c,
		Title: embed.Title.Text,
		Header: &chatv1.Embed_EmbedHeading{
			Text: embed.Header.Text,
			Url:  &embed.Header.URL,
			Icon: &embed.Header.Icon,
		},
		Footer: &chatv1.Embed_EmbedHeading{
			Text: embed.Footer.Text,
			Url:  &embed.Footer.URL,
			Icon: &embed.Footer.Icon,
		},
		Fields: func() (f []*chatv1.Embed_EmbedField) {
			for _, field := range embed.Fields {
				f = append(f, &chatv1.Embed_EmbedField{
					Title:        field.Title,
					Body:         ftext(field.Body),
					Presentation: chatv1.Embed_EmbedField_PRESENTATION_DATA_UNSPECIFIED,
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
			GuildId:   c.tm.GuildID,
			ChannelId: c.tm.ChannelID,
			Content:   &chatv1.Content{Content: &chatv1.Content_TextMessage{TextMessage: &chatv1.Content_TextContent{Content: ftext(content)}}},
		})
		if err != nil {
			log.Error("%+v", err)
		}
	case commandlib.Embed:
		_, err := c.c.ChatKit.SendMessage(&chatv1.SendMessageRequest{
			GuildId:   c.tm.GuildID,
			ChannelId: c.tm.ChannelID,
			Content: &chatv1.Content{
				Content: &chatv1.Content_EmbedMessage{
					EmbedMessage: &chatv1.Content_EmbedContent{
						Embed: convert(content),
					},
				},
			},
		})
		if err != nil {
			log.Error("%+v", err)
		}
	case commandlib.EmbedList:
		log.Error("unhandled embed list")
		return
		_, err := c.c.ChatKit.SendMessage(&chatv1.SendMessageRequest{
			GuildId:   c.tm.GuildID,
			ChannelId: c.tm.ChannelID,
			// Embeds: func() (r []*types.Embed) {
			// 	for _, embed := range content.Embeds {
			// 		r = append(r, convert(embed))
			// 	}
			// 	return
			// }(),
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
	log.Error("unhandled embed list")
	return
	c.c.ChatKit.SendMessage(&chatv1.SendMessageRequest{
		GuildId:   c.tm.GuildID,
		ChannelId: c.tm.ChannelID,
		// Embeds: func() (r []*types.Embed) {
		// 	for _, embed := range tags {
		// 		r = append(r, convert(embed))
		// 	}
		// 	return
		// }(),
	})
}

func (c *Context) WrapCodeBlock(code string) string {
	return fmt.Sprintf("```%s```", code)
}
