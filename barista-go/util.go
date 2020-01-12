package barista

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sahilm/fuzzy"
)

// Default : Default string handling.
func Default(in string, def string) string {
	if in != "" {
		return in
	} else {
		return def
	}
}

// GetGuildMembers : Get all members in a guild
func (cmd *LexedCommand) GetGuildMembers(id string) []*discordgo.Member {
	i := 1000
	var prev string
	var members []*discordgo.Member
	for i == 1000 {
		mems, _ := cmd.Session.GuildMembers(id, prev, 1000)
		if len(mems) == 0 {
			continue
		}
		members = append(members, mems...)
		prev = mems[len(mems)-1].User.ID
		i = len(mems)
	}
	return members
}

// GetChannelMembers : Get all members that can read messages
func (cmd *LexedCommand) GetChannelMembers(chanID, guildID string) []*discordgo.Member {
	members := cmd.GetGuildMembers(guildID)
	var canSee []*discordgo.Member
	for _, member := range members {
		perms, _ := cmd.Session.State.UserChannelPermissions(member.User.ID, chanID)
		if perms&discordgo.PermissionReadMessages == discordgo.PermissionReadMessages {
			canSee = append(canSee, member)
		}
	}
	return canSee
}

// MatchUser : Get a user matching a string and returns an ID.
func (cmd *LexedCommand) MatchUser(user string) string {
	mems := cmd.GetChannelMembers(cmd.CommandMessage.ChannelID, cmd.CommandMessage.GuildID)
	var memNames []string
	for _, mem := range mems {
		if mem.User.Bot {
			continue
		}
		memNames = append(memNames, mem.User.Username+"#"+mem.User.Discriminator)
	}
	matches := fuzzy.Find(user, memNames)
	matchesToUserID := make(map[string]string)
	for _, val := range matches {
		for _, mem := range mems {
			if mem.User.Username+"#"+mem.User.Discriminator == val.Str {
				matchesToUserID[val.Str] = mem.User.ID
			}
		}
	}
	matchesSlice := matches[:]
	if len(matches) == 1 {
		return matchesToUserID[matches[0].Str]
	} else if len(matches) == 0 {
		return ""
	}
	sort.Slice(matchesSlice, func(i, j int) bool {
		return matches[i].Score < matches[j].Score
	})
	embed := NewEmbed().
		SetTitle("Select which user you want.")
	if len(matchesSlice) > 5 {
		matchesSlice = matchesSlice[:5]
	}
	for index, val := range matchesSlice {
		embed = embed.AddField(strconv.Itoa(index+1), val.Str, false)
		for _, mem := range mems {
			if mem.User.Username+"#"+mem.User.Discriminator == val.Str {
				matchesToUserID[val.Str] = mem.User.ID
			}
		}
	}
	timeoutChan := make(chan int)
	go func() {
		time.Sleep(10 * time.Second)
		timeoutChan <- 0
	}()
	msg, _ := cmd.Session.ChannelMessageSendEmbed(cmd.CommandMessage.ChannelID, embed.MessageEmbed)
	for {
		select {
		case usermsg := <-waitForMessage(cmd.Session):
			if usermsg.Author.ID != cmd.CommandMessage.Author.ID {
				continue
			}
			defer cmd.Session.ChannelMessageDelete(cmd.CommandMessage.ChannelID, msg.ID)
			defer cmd.Session.ChannelMessageDelete(cmd.CommandMessage.ChannelID, usermsg.ID)
			if strings.Contains("1", usermsg.Content) {
				return matchesToUserID[matchesSlice[0].Str]
			} else if strings.Contains("2", usermsg.Content) {
				return matchesToUserID[matchesSlice[1].Str]
			} else if strings.Contains("3", usermsg.Content) {
				return matchesToUserID[matchesSlice[2].Str]
			} else if strings.Contains("4", usermsg.Content) {
				return matchesToUserID[matchesSlice[3].Str]
			} else if strings.Contains("5", usermsg.Content) {
				return matchesToUserID[matchesSlice[4].Str]
			}
			break
		case <-timeoutChan:
			cmd.Session.ChannelMessageDelete(cmd.CommandMessage.ChannelID, msg.ID)
			embed := NewEmbed().
				SetColor(0xff0000).
				SetTitle("User selection timed out.")
			msgSend := discordgo.MessageSend{
				Embed: embed.MessageEmbed,
			}
			cmd.SendMessage(&msgSend)
			return ""
		}
	}
}

func waitForMessage(s *discordgo.Session) chan *discordgo.MessageCreate {
	channel := make(chan *discordgo.MessageCreate)
	s.AddHandlerOnce(func(_ *discordgo.Session, e *discordgo.MessageCreate) {
		channel <- e
	})
	return channel
}
