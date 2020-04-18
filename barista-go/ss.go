package barista

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type ScreenshotUpvotes struct {
	UpvotedUsers   []string
	DownvotedUsers []string
}

func (self *ScreenshotUpvotes) Upvote(id string) bool {
	self.RemoveDownvote(id)
	for _, user := range self.UpvotedUsers {
		if user == id {
			return false
		}
	}
	self.UpvotedUsers = append(self.UpvotedUsers, id)
	return true
}

func (self *ScreenshotUpvotes) RemoveUpvote(id string) bool {
	idx := -1
	for index, user := range self.UpvotedUsers {
		if user == id {
			idx = index
		}
	}
	if idx != -1 {
		self.UpvotedUsers[idx] = self.UpvotedUsers[len(self.UpvotedUsers)-1]
		self.UpvotedUsers[len(self.UpvotedUsers)-1] = ""
		self.UpvotedUsers = self.UpvotedUsers[:len(self.UpvotedUsers)-1]
		return true
	}
	return false
}

func (self *ScreenshotUpvotes) RemoveDownvote(id string) bool {
	idx := -1
	for index, user := range self.DownvotedUsers {
		if user == id {
			idx = index
		}
	}
	if idx != -1 {
		self.DownvotedUsers[idx] = self.DownvotedUsers[len(self.DownvotedUsers)-1]
		self.DownvotedUsers[len(self.DownvotedUsers)-1] = ""
		self.DownvotedUsers = self.DownvotedUsers[:len(self.DownvotedUsers)-1]
		return true
	}
	return false
}

func (self *ScreenshotUpvotes) Downvote(id string) bool {
	self.RemoveUpvote(id)
	for _, user := range self.DownvotedUsers {
		if user == id {
			return false
		}
	}
	self.DownvotedUsers = append(self.DownvotedUsers, id)
	return true
}

func (self *ScreenshotUpvotes) Serialize() string {
	data, err := json.Marshal(self)
	if err != nil {
		return ""
	}
	return string(data)
}

func DeserializeScreenshotUpvotes(data string) ScreenshotUpvotes {
	var votes ScreenshotUpvotes
	json.Unmarshal([]byte(data), &votes)
	return votes
}

func Screenshot(s *discordgo.Session, cmd *LexedCommand) {
	if !commandEnabled(cmd, "profl") {
		return
	}
	helpmsg := "```dsconfig\n" + sshelp + "```"

	used := false
	updated := false
	var user string
	if val := cmd.GetFlagPair("-d", "--set-description"); val != "" {
		SetGlobalKey(cmd.MemKey("ss-description"), val)
		used, updated = true, true
	}
	if val := cmd.GetFlagPair("-s", "--set-screenshot"); val != "" {
		if strings.HasPrefix(val, "http") && strings.HasSuffix(val, "png") {
			SetGlobalKey(cmd.MemKey("screenshot"), val)
			used, updated = true, true
		} else {
			embed := NewEmbed().
				SetColor(0xff0000).
				SetTitle("Please specify a URL to a .png file.")
			msgSend := discordgo.MessageSend{
				Embed: embed.MessageEmbed,
			}
			cmd.SendMessage(&msgSend)
			return
		}
	}
	if val := cmd.GetFlagPair("-u", "--user"); val != "" {
		used = true
		user = val
	}
	if cmd.GetFlagPair("", "--upvote") == "nil" || cmd.GetFlagPair("", "--downvote") == "nil" {
		if val := cmd.GetFlagPair("-u", "--user"); val == "" {
			cmd.SendErrorEmbed("Please specify a user in your command.", "")
			return
		}
		used = true
	}
	if !used {
		cmd.SendErrorEmbed("Please specify some arguments", helpmsg)
		return
	}
	if updated == true {
		embed := NewEmbed().
			SetColor(cmd.Author.Colour).
			SetTitle("Screenshot updated!").
			SetDescription(GetGlobalKey(cmd.MemKey("ss-description"))).
			SetImage(Default(GetGlobalKey(cmd.MemKey("screenshot")), ""))
		SetGlobalKey(cmd.MemKey("ss-votes"), "")
		msgSend := discordgo.MessageSend{
			Embed: embed.MessageEmbed,
		}
		cmd.SendMessage(&msgSend)
		return
	}
	if val := cmd.MatchUser(user); val != "" {
		mem, _ := cmd.Session.GuildMember(cmd.CommandMessage.GuildID, val)
		colour := cmd.Session.State.UserColor(val, cmd.CommandMessage.ChannelID)
		var disp string
		if mem.Nick == "" {
			disp = mem.User.Username
		} else {
			disp = mem.Nick
		}

		votes := DeserializeScreenshotUpvotes(GetGlobalKey(val + "ss-votes"))

		if cmd.GetFlagPair("", "--upvote") == "nil" {
			votes.Upvote(cmd.CommandMessage.Author.ID)
			SetGlobalKey(val+"ss-votes", votes.Serialize())
		}
		if cmd.GetFlagPair("", "--downvote") == "nil" {
			votes.Downvote(cmd.CommandMessage.Author.ID)
			SetGlobalKey(val+"ss-votes", votes.Serialize())
		}

		if cmd.GetFlagPair("", "--upvote") == "nil" || cmd.GetFlagPair("", "--downvote") == "nil" {
			embed := NewEmbed().
				SetColor(colour).
				SetTitle("Voted on " + disp + "'s Screenshot").
				SetDescription(GetGlobalKey(val + "ss-description")).
				SetImage(Default(GetGlobalKey(val+"screenshot"), "")).
				SetFooter(fmt.Sprintf("Score: %d", len(votes.UpvotedUsers)-len(votes.DownvotedUsers)))

			msgSend := discordgo.MessageSend{
				Embed: embed.MessageEmbed,
			}
			cmd.SendMessage(&msgSend)
			return
		}

		embed := NewEmbed().
			SetColor(colour).
			SetTitle(disp + "'s Screenshot").
			SetThumbnail(mem.User.AvatarURL("")).
			SetDescription(GetGlobalKey(val + "ss-description")).
			SetImage(Default(GetGlobalKey(val+"screenshot"), "")).
			SetFooter(fmt.Sprintf("Score: %d", len(votes.UpvotedUsers)-len(votes.DownvotedUsers)))

		msgSend := discordgo.MessageSend{
			Embed: embed.MessageEmbed,
		}
		cmd.SendMessage(&msgSend)
		return
	} else {
		embed := NewEmbed().
			SetColor(0xff0000).
			SetDescription(helpmsg).
			SetTitle("No user specified.")
		msgSend := discordgo.MessageSend{
			Embed: embed.MessageEmbed,
		}
		cmd.SendMessage(&msgSend)
		return
	}
}
