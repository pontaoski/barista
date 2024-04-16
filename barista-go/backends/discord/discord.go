package discord

import (
	"fmt"
	"time"

	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	stripmd "github.com/writeas/go-strip-markdown"
)

// DiscordContext is a Discord context
type DiscordContext struct {
	commandlib.ContextMixin
	author     discord.User
	pm         map[string]*discord.Message
	paginators map[string]*paginator
	tags       map[string][]*discord.Message
	s          *DiscordBackend
	me         *discord.User
	tm         *discord.Message
}

func (d *DiscordContext) Backend() commandlib.Backend {
	return backends[d.me.ID]
}

func (d *DiscordContext) cleanID(id string) {
	if val, ok := d.paginators[id]; ok {
		val.Inactive()
		delete(d.paginators, id)
	}
	if val, ok := d.tags[id]; ok {
		for _, msg := range val {
			d.s.s.Client.DeleteMessage(msg.ChannelID, msg.ID, api.AuditLogReason(""))
		}
	}
}

func (t *DiscordContext) keyboard() *discord.ActionRowComponent {
	return &discord.ActionRowComponent{
		&discord.ButtonComponent{
			Label:    t.I18n("Previous"),
			CustomID: "previous",
		},
		&discord.ButtonComponent{
			Label:    t.I18n("Next"),
			CustomID: "next",
		},
	}
}

func (d *DiscordContext) SendTags(id string, tags []commandlib.Embed) {
	d.cleanID(id)
	for _, tag := range tags {
		msg, _ := d.s.s.SendEmbeds(d.tm.ChannelID, discordEmbed(tag))
		if msg != nil {
			d.tags[id] = append(d.tags[id], msg)
		}
	}
}

func (d *DiscordContext) GenerateLink(text, URL string) string {
	return fmt.Sprintf("[%s](%s)", text, URL)
}

func (d *DiscordContext) WrapCodeBlock(code string) string {
	return "```\n" + code + "\n```"
}

func (d *DiscordContext) SendMessage(id string, content interface{}) {
	if val, ok := d.pm[id]; ok {
		switch a := content.(type) {
		case string:
			d.pm[id], _ = d.s.s.EditMessage(val.ChannelID, val.ID, content.(string))
			goto clean
		case commandlib.Embed:
			d.pm[id], _ = d.s.s.EditEmbeds(val.ChannelID, val.ID, discordEmbed(content.(commandlib.Embed)))
			goto clean
		case commandlib.EmbedList:
			goto paginator
		case commandlib.UnionEmbed:
			d.SendMessage(id, content.(commandlib.UnionEmbed).EmbedList)
			return
		case commandlib.File:
			d.s.s.DeleteMessage(val.ChannelID, d.pm[id].ID, api.AuditLogReason(""))
			d.pm[id], _ = d.s.s.SendMessageComplex(val.ChannelID, api.SendMessageData{
				Files: []sendpart.File{
					{
						Name:   a.Name,
						Reader: a.Reader,
					},
				},
			})
			a.Reader.Close()
		}
	} else {
		switch a := content.(type) {
		case string:
			d.pm[id], _ = d.s.s.SendMessage(d.tm.ChannelID, content.(string))
		case commandlib.Embed:
			d.pm[id], _ = d.s.s.SendEmbeds(d.tm.ChannelID, discordEmbed(content.(commandlib.Embed)))
		case commandlib.EmbedList:
			goto paginator
		case commandlib.UnionEmbed:
			d.SendMessage(id, content.(commandlib.UnionEmbed).EmbedList)
			return
		case commandlib.File:
			d.pm[id], _ = d.s.s.SendMessageComplex(d.tm.ChannelID, api.SendMessageData{
				Files: []sendpart.File{
					{
						Name:   a.Name,
						Reader: a.Reader,
					},
				},
			})
			a.Reader.Close()
		}
	}
	return
clean:
	d.cleanID(id)
	return
paginator:
	embedList := content.(commandlib.EmbedList)
	if val, ok := d.pm[id]; ok {
		d.s.s.DeleteMessage(val.ChannelID, val.ID, api.AuditLogReason(""))
		delete(d.pm, id)
	}
	if val, ok := d.paginators[id]; ok {
		val.Inactive()
	}
	paginator := newPaginator(d)
	d.paginators[id] = paginator
	title := "Item"
	if embedList.ItemTypeName != "" {
		title = embedList.ItemTypeName
	}
	for index, page := range embedList.Embeds {
		page.Footer.Text = fmt.Sprintf("%s %d out of %d", title, index+1, len(embedList.Embeds))
		paginator.AddPage(discordEmbed(page))
	}
}

func (d DiscordContext) AuthorName() string {
	return d.tm.Author.Username
}

func (d DiscordContext) AuthorIdentifier() string {
	return d.tm.Author.ID.String()
}

func (d DiscordContext) RoomIdentifier() string {
	return d.tm.ChannelID.String()
}

func (d DiscordContext) CommunityIdentifier() string {
	return d.tm.GuildID.String()
}

func (d DiscordContext) I18n(message string) string {
	return d.I18nInternal(commandlib.GetLanguage(&d), message)
}

func (d DiscordContext) I18nc(context, message string) string {
	return d.I18n(message)
}

func waitForMessage(s *state.State) chan *gateway.MessageCreateEvent {
	channel := make(chan *gateway.MessageCreateEvent)
	var rm func()
	rm = s.PreHandler.AddHandler(func(c *gateway.MessageCreateEvent) {
		channel <- c
		rm()
	})
	return channel
}

func (d *DiscordContext) NextResponse() (out chan string) {
	out = make(chan string)
	go func() {
		for {
			select {
			case usermsg := <-waitForMessage(d.s.s):
				if usermsg.Author.ID == d.tm.Author.ID && usermsg.ChannelID == d.tm.ChannelID {
					out <- stripmd.Strip(usermsg.Content)
					return
				}
			}
		}
	}()
	return out
}

func (d *DiscordContext) AwaitResponse(tenpo time.Duration) (response string, ok bool) {
	timeoutChan := make(chan struct{})
	go func() {
		time.Sleep(tenpo)
		timeoutChan <- struct{}{}
	}()
	for {
		select {
		case msg := <-d.NextResponse():
			return msg, true
		case <-timeoutChan:
			return "", false
		}
	}
}

func discordEmbed(d commandlib.Embed) discord.Embed {
	d.Truncate()
	var fields []discord.EmbedField
	for _, field := range d.Fields {
		fields = append(fields, discord.EmbedField{
			Name:   field.Title,
			Value:  field.Body,
			Inline: field.Inline,
		})
	}
	return discord.Embed{
		Title:       d.Title.Text,
		URL:         d.Title.URL,
		Description: d.Body,
		Author: &discord.EmbedAuthor{
			Name: d.Header.Text,
			URL:  d.Header.URL,
			Icon: d.Header.Icon,
		},
		Footer: &discord.EmbedFooter{
			Text: d.Footer.Text,
			Icon: d.Footer.URL,
		},
		Fields: fields,
		Color:  discord.Color(d.Colour),
	}
}

func buildContext(c commandlib.ContextMixin, s *DiscordBackend, m *discord.Message) DiscordContext {
	dc := DiscordContext{
		ContextMixin: c,
	}
	dc.author = m.Author
	dc.s = s
	dc.tm = m
	dc.pm = make(map[string]*discord.Message)
	dc.paginators = make(map[string]*paginator)
	dc.tags = make(map[string][]*discord.Message)
	return dc
}
