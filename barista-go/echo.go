package barista

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Echo : The echo command.
func Echo(s *discordgo.Session, cmd *LexedCommand) {
	embed := NewEmbed().
		SetColor(cmd.Author.Colour).
		SetDescription(cmd.Query.Content)
	embed.MessageEmbed.Author = &discordgo.MessageEmbedAuthor{
		Name:    fmt.Sprintf("%s said...", cmd.Author.DisplayName),
		IconURL: cmd.Author.Avatar,
	}
	msgSend := discordgo.MessageSend{
		Embed: embed.MessageEmbed,
	}
	cmd.SendMessage(&msgSend)
}
