package barista

import (
	"encoding/json"

	"github.com/bwmarrin/discordgo"
)

func EmbedCmd(s *discordgo.Session, cmd *LexedCommand) {
	var embed discordgo.MessageEmbed
	json.Unmarshal([]byte(cmd.Query.Content), &embed)

	if !cmd.Author.IsAdmin {
		return
	}

	if val := cmd.GetFlagPair("-m", "--message"); val != "" {
		s.ChannelMessageEditEmbed(
			cmd.CommandMessage.ChannelID,
			val,
			&embed,
		)
	} else {
		s.ChannelMessageSendEmbed(
			cmd.CommandMessage.ChannelID,
			&embed,
		)
	}
}
