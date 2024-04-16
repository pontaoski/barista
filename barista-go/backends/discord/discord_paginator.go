package discord

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	lru "github.com/hashicorp/golang-lru"
)

var paginatorCache *lru.ARCCache

func init() {
	var err error
	paginatorCache, err = lru.NewARC(1 << 16)
	if err != nil {
		panic(err)
	}
}

type paginator struct {
	Pages []discord.Embed
	Index int

	context *DiscordContext
	message *discord.Message
}

func newPaginator(context *DiscordContext) *paginator {
	return &paginator{
		context: context,
	}
}

func (p *paginator) AddPage(msg discord.Embed) {
	p.Pages = append(p.Pages, msg)
}

func (p *paginator) Send() {
	p.Index = 0
	send := p.Pages[p.Index]
	msg, err := p.context.s.s.Client.SendMessageComplex(p.context.tm.ChannelID, api.SendMessageData{
		Embeds: []discord.Embed{send},
		Components: discord.ContainerComponents{
			p.context.keyboard(),
		},
	})
	if err == nil {
		p.message = msg
		paginatorCache.Add(msg.ID, p)
	}
}

func (p *paginator) Inactive() {
	send := p.Pages[p.Index]
	msg, err := p.context.s.s.Client.EditMessageComplex(p.message.ChannelID, p.message.ID, api.EditMessageData{
		Embeds:     &[]discord.Embed{send},
		Components: &discord.ContainerComponents{},
	})
	if err != nil {
		p.message = msg
	}
}

func (p *paginator) Prev(e *discord.InteractionEvent) {
	p.Index--
	if p.Index < 0 {
		p.Index = len(p.Pages) - 1
	}
	send := p.Pages[p.Index]
	p.context.s.s.Client.RespondInteraction(e.ID, e.Token, api.InteractionResponse{
		Type: api.UpdateMessage,
		Data: &api.InteractionResponseData{
			Embeds: &[]discord.Embed{send},
		},
	})
}

func (p *paginator) Next(e *discord.InteractionEvent) {
	p.Index++
	if p.Index+1 > len(p.Pages) {
		p.Index = 0
	}
	send := p.Pages[p.Index]
	p.context.s.s.Client.RespondInteraction(e.ID, e.Token, api.InteractionResponse{
		Type: api.UpdateMessage,
		Data: &api.InteractionResponseData{
			Embeds: &[]discord.Embed{send},
		},
	})
}
