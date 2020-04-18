package barista

import (
	"github.com/ajanicij/goduckgo/goduckgo"
	"github.com/bwmarrin/discordgo"
)

func Ddg(s *discordgo.Session, cmd *LexedCommand) {
	if !commandEnabled(cmd, "ddgse") {
		return
	}
	if cmd.Query.Content == "" {
		embed := NewEmbed().
			SetColor(0xff0000).
			SetTitle("Please specify a search term.")
		msgSend := discordgo.MessageSend{
			Embed: embed.MessageEmbed,
		}
		cmd.SendMessage(&msgSend)
		return
	}
	gdgQuery, err := goduckgo.Query(cmd.Query.Content)
	if err != nil {
		embed := NewEmbed().
			SetColor(0xff0000).
			SetTitle("There was an error getting search results.")
		msgSend := discordgo.MessageSend{
			Embed: embed.MessageEmbed,
		}
		cmd.SendMessage(&msgSend)
		return
	}
	embed := NewEmbed().
		SetColor(0xd34b2b).
		SetTitle(gdgQuery.Heading).
		SetDescription(gdgQuery.Abstract).
		SetThumbnail(gdgQuery.Image).
		SetURL(gdgQuery.AbstractURL).
		SetAuthor("DuckDuckGo", "https://duckduckgo.com/favicon.png", "https://duckduckgo.com")
	if gdgQuery.Answer != "" {
		embed = embed.SetDescription(gdgQuery.Answer)
	}
	msgSend := discordgo.MessageSend{
		Embed: embed.MessageEmbed,
	}
	cmd.SendMessage(&msgSend)
}
