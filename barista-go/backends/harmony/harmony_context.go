package harmony

import (
	"fmt"
	"strconv"
	"time"

	"github.com/appadeia/barista/barista-go/backends/harmony/client"
	corev1 "github.com/appadeia/barista/barista-go/backends/harmony/gen/core"
	profilev1 "github.com/appadeia/barista/barista-go/backends/harmony/gen/profile"
	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/log"
)

// Context is a context holding information about a Harmony command
type Context struct {
	commandlib.ContextMixin
	b  *Backend
	c  *client.Client
	tm *corev1.Message
}

func buildContext(ctx commandlib.ContextMixin, b *Backend, c *client.Client, m *corev1.Message) Context {
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
	resp, err := c.c.Profilekit.GetUser(c.c.Context(), &profilev1.GetUserRequest{
		UserId: c.tm.AuthorId,
	})
	if err != nil {
		return c.AuthorIdentifier()
	}

	return resp.UserName
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
	return strconv.FormatUint(c.tm.Location.GuildId, 10)
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

func waitForMessage(c *client.Client) chan *corev1.Message {
	channel := make(chan *corev1.Message)
	var f func()
	f = func() {
		c.HandleOnce(func(ev *corev1.GuildEvent) {
			if val, ok := ev.Event.(*corev1.GuildEvent_SentMessage); ok {
				channel <- val.SentMessage.Message
			} else {
				f()
			}
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
				if usermsg.AuthorId == c.tm.AuthorId && usermsg.Location.ChannelId == usermsg.Location.ChannelId {
					out <- usermsg.Content
					return
				}
			}
		}
	}()
	return out
}

func (c *Context) RoomIdentifier() string {
	return strconv.FormatUint(c.tm.Location.ChannelId, 10)
}

func convert(embed commandlib.Embed) *corev1.Embed {
	return &corev1.Embed{
		Body:  embed.Body,
		Color: int64(embed.Colour),
		Title: embed.Title.Text,
		Header: &corev1.EmbedHeading{
			Text: embed.Header.Text,
			Url:  embed.Header.URL,
			Icon: embed.Header.Icon,
		},
		Footer: &corev1.EmbedHeading{
			Text: embed.Footer.Text,
			Url:  embed.Footer.URL,
			Icon: embed.Footer.Icon,
		},
		Fields: func() (f []*corev1.EmbedField) {
			for _, field := range embed.Fields {
				f = append(f, &corev1.EmbedField{
					Title:        field.Title,
					Body:         field.Body,
					Presentation: corev1.FieldPresentation_Data,
				})
			}
			return
		}(),
	}
}

func (c *Context) SendMessage(id string, content interface{}) {
	switch content := content.(type) {
	case string:
		_, err := c.c.CoreKit.SendMessage(c.c.Context(), &corev1.SendMessageRequest{
			Location: c.tm.Location,
			Content:  content,
		})
		if err != nil {
			log.Error("%+v", err)
		}
	case commandlib.Embed:
		_, err := c.c.CoreKit.SendMessage(c.c.Context(), &corev1.SendMessageRequest{
			Location: c.tm.Location,
			Embeds: []*corev1.Embed{
				convert(content),
			},
		})
		if err != nil {
			log.Error("%+v", err)
		}
	case commandlib.EmbedList:
		_, err := c.c.CoreKit.SendMessage(c.c.Context(), &corev1.SendMessageRequest{
			Location: c.tm.Location,
			Embeds: func() (r []*corev1.Embed) {
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
	c.c.CoreKit.SendMessage(c.c.Context(), &corev1.SendMessageRequest{
		Location: c.tm.Location,
		Embeds: func() (r []*corev1.Embed) {
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
