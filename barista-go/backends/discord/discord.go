package discord

import (
	"fmt"
	"time"

	"github.com/Necroforger/dgwidgets"
	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/bwmarrin/discordgo"
	stripmd "github.com/writeas/go-strip-markdown"
)

// DiscordContext is a Discord context
type DiscordContext struct {
	commandlib.ContextMixin
	author     *discordgo.User
	pm         map[string]*discordgo.Message
	paginators map[string]*dgwidgets.Paginator
	tags       map[string][]*discordgo.Message
	s          *discordgo.Session
	tm         *discordgo.Message
}

func (d *DiscordContext) cleanID(id string) {
	if val, ok := d.paginators[id]; ok {
		val.Widget.Close <- true
		delete(d.paginators, id)
	}
	if val, ok := d.tags[id]; ok {
		for _, msg := range val {
			d.s.ChannelMessageDelete(msg.ChannelID, msg.ID)
		}
	}
}

func (d *DiscordContext) SendTags(id string, tags []commandlib.Embed) {
	d.cleanID(id)
	for _, tag := range tags {
		msg, _ := d.s.ChannelMessageSendEmbed(d.tm.ChannelID, discordEmbed(tag))
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
		switch content.(type) {
		case string:
			d.pm[id], _ = d.s.ChannelMessageEdit(val.ChannelID, val.ID, content.(string))
			goto clean
		case commandlib.Embed:
			d.pm[id], _ = d.s.ChannelMessageEditEmbed(val.ChannelID, val.ID, discordEmbed(content.(commandlib.Embed)))
			goto clean
		case commandlib.EmbedList:
			goto paginator
		case commandlib.UnionEmbed:
			d.SendMessage(id, content.(commandlib.UnionEmbed).EmbedList)
			return
		}
	} else {
		switch content.(type) {
		case string:
			d.pm[id], _ = d.s.ChannelMessageSend(d.tm.ChannelID, content.(string))
		case commandlib.Embed:
			d.pm[id], _ = d.s.ChannelMessageSendEmbed(d.tm.ChannelID, discordEmbed(content.(commandlib.Embed)))
		case commandlib.EmbedList:
			goto paginator
		case commandlib.UnionEmbed:
			d.SendMessage(id, content.(commandlib.UnionEmbed).EmbedList)
			return
		}
	}
	return
clean:
	d.cleanID(id)
	return
paginator:
	embedList := content.(commandlib.EmbedList)
	if val, ok := d.pm[id]; ok {
		d.s.ChannelMessageDelete(val.ChannelID, val.ID)
		delete(d.pm, id)
	}
	if val, ok := d.paginators[id]; ok {
		val.Widget.Close <- true
	}
	paginator := dgwidgets.NewPaginator(d.s, d.tm.ChannelID)
	d.paginators[id] = paginator
	title := "Item"
	if embedList.ItemTypeName != "" {
		title = embedList.ItemTypeName
	}
	for index, page := range embedList.Embeds {
		page.Footer.Text = fmt.Sprintf("%s %d out of %d", title, index+1, len(embedList.Embeds))
		paginator.Add(discordEmbed(page))
	}
	paginator.DeleteMessageWhenDone = true
	go paginator.Spawn()
}

func (d DiscordContext) AuthorName() string {
	return d.tm.Author.Username
}

func (d DiscordContext) AuthorIdentifier() string {
	return d.tm.Author.ID
}

func (d DiscordContext) RoomIdentifier() string {
	return d.tm.ChannelID
}

func (d DiscordContext) CommunityIdentifier() string {
	return d.tm.GuildID
}

func (d DiscordContext) I18n(message string) string {
	return d.I18nInternal(commandlib.GetLanguage(&d), message)
}

func (d DiscordContext) I18nc(context, message string) string {
	return d.I18n(message)
}

func waitForMessage(s *discordgo.Session) chan *discordgo.MessageCreate {
	channel := make(chan *discordgo.MessageCreate)
	s.AddHandlerOnce(func(_ *discordgo.Session, e *discordgo.MessageCreate) {
		channel <- e
	})
	return channel
}

func (d *DiscordContext) NextResponse() (out chan string) {
	out = make(chan string)
	go func() {
		for {
			select {
			case usermsg := <-waitForMessage(d.s):
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

func discordEmbed(d commandlib.Embed) *discordgo.MessageEmbed {
	d.Truncate()
	var fields []*discordgo.MessageEmbedField
	for _, field := range d.Fields {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   field.Title,
			Value:  field.Body,
			Inline: field.Inline,
		})
	}
	return &discordgo.MessageEmbed{
		Title:       d.Title.Text,
		URL:         d.Title.URL,
		Description: d.Body,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    d.Header.Text,
			URL:     d.Header.URL,
			IconURL: d.Header.Icon,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    d.Footer.Text,
			IconURL: d.Footer.URL,
		},
		Fields: fields,
		Color:  d.Colour,
	}
}

func buildContext(c commandlib.ContextMixin, s *discordgo.Session, m *discordgo.Message) DiscordContext {
	dc := DiscordContext{
		ContextMixin: c,
	}
	dc.author = m.Author
	dc.s = s
	dc.tm = m
	dc.pm = make(map[string]*discordgo.Message)
	dc.paginators = make(map[string]*dgwidgets.Paginator)
	dc.tags = make(map[string][]*discordgo.Message)
	return dc
}
