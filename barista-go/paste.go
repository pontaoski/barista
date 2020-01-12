package barista

import "github.com/bwmarrin/discordgo"

func Paste(s *discordgo.Session, cmd *LexedCommand) {
	if cmd.Author.IsOwner {
		embed := NewEmbed().
			SetColor(cmd.Author.Colour).
			SetTitle("Here's a link to pastebin.com").
			SetURL("https://pastebin.com/")
		msgSend := discordgo.MessageSend{
			Embed: embed.MessageEmbed,
		}
		cmd.SendMessage(&msgSend)
		return
	} else {
		embed := NewEmbed().
			SetColor(0xff0000).
			SetTitle("tf you doing bruh, go tell someone to use a pastebin yourself instead of being a lazy ass and having a robot do it for you")
		msgSend := discordgo.MessageSend{
			Embed: embed.MessageEmbed,
		}
		cmd.SendMessage(&msgSend)
		return
	}
}
