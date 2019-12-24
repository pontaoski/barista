package barista

import (
	"fmt"
	"strings"
	"time"

	"github.com/Necroforger/dgwidgets"
	"github.com/bwmarrin/discordgo"
)

// CommandFunc : The type definition for a command function.
type CommandFunc func(*discordgo.Session, *LexedCommand)

// LexedAuthor : information about a command's author.
type LexedAuthor struct {
	DisplayName string
	Colour      int
	Avatar      string
}

// LexedQuery : information about what a command requests
type LexedQuery struct {
	Array      []string
	TextLength []string
	Content    string
}

// LexedCommand : A lexed command with editable stuff.
type LexedCommand struct {
	Author LexedAuthor
	Query  LexedQuery
	Flags  map[string]string

	PaginatorPageName string

	CommandMessage  *discordgo.Message
	SentMessage     *discordgo.Message
	SentTagMessages []*discordgo.Message
	SentPaginator   *dgwidgets.Paginator

	Session *discordgo.Session

	LastUsed time.Time
}

// SendMessage : Send a message.
func (cmd *LexedCommand) SendMessage(s *discordgo.MessageSend) {
	var err error
	if cmd.SentPaginator != nil {
		cmd.SentPaginator.Widget.Close <- true
		cmd.SentPaginator = nil
	}
	if cmd.SentMessage == nil {
		cmd.SentMessage, err = cmd.Session.ChannelMessageSendComplex(cmd.CommandMessage.ChannelID, s)
	} else {
		msgedit := discordgo.MessageEdit{
			Content: &s.Content,
			Embed:   s.Embed,
			ID:      cmd.SentMessage.ID,
			Channel: cmd.CommandMessage.ChannelID,
		}
		cmd.SentMessage, err = cmd.Session.ChannelMessageEditComplex(&msgedit)
	}
	if err != nil {
		println(err.Error())
	}
}

// SendTags : Send a bunch of embeds
func (cmd *LexedCommand) SendTags(embeds []*Embed) {
	for _, old := range cmd.SentTagMessages {
		cmd.Session.ChannelMessageDelete(old.ChannelID, old.ID)
	}
	cmd.SentTagMessages = []*discordgo.Message{}
	for _, new := range embeds {
		msgSend := discordgo.MessageSend{
			Embed: new.MessageEmbed,
		}
		msg, err := cmd.Session.ChannelMessageSendComplex(cmd.CommandMessage.ChannelID, &msgSend)
		if err == nil {
			cmd.SentTagMessages = append(cmd.SentTagMessages, msg)
		}
	}
}

// SendPaginator : Send paginator
func (cmd *LexedCommand) SendPaginator(paginator *dgwidgets.Paginator) {
	if cmd.SentPaginator != nil {
		cmd.SentPaginator.Widget.Close <- true
		cmd.SentPaginator = nil
	}
	if cmd.SentMessage != nil {
		cmd.Session.ChannelMessageDelete(cmd.SentMessage.ChannelID, cmd.SentMessage.ID)
		cmd.SentMessage = nil
	}
	for index, page := range paginator.Pages {
		page.Footer = &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%s %d out of %d", cmd.PaginatorPageName, index+1, len(paginator.Pages)),
		}
	}
	paginator.DeleteMessageWhenDone = true
	cmd.SentPaginator = paginator
	cmd.SentPaginator.Spawn()
}

// GetFlagPair : Get a flag pair
func (cmd *LexedCommand) GetFlagPair(short string, long string) string {
	if val, ok := cmd.Flags[short]; ok {
		return val
	} else if val, ok := cmd.Flags[long]; ok {
		return val
	} else {
		return ""
	}
}

func (cmd *LexedCommand) lex() {
	if len(strings.Split(cmd.CommandMessage.Content, " ")) < 2 {
		return
	}
	arr := strings.Split(cmd.CommandMessage.Content, " ")[2:]
	var toremove []int
	hasFlag := false
	for index, i := range arr {
		if i == "--" {
			cmd.Query.Content = strings.Join(arr[index+1:], " ")
			break
		}
		if strings.HasPrefix(i, "-") {
			if strings.Contains(i, "=") {
				flag := strings.Split(i, "=")
				cmd.Flags[flag[0]] = flag[1]
				hasFlag = true
			} else {
				val := ""
				for indx, ii := range arr[index+1:] {
					toremove = append(toremove, index+1+indx)
					if !strings.HasPrefix(ii, "-") {
						val = val + " " + ii
					} else {
						break
					}
				}
				if strings.TrimSpace(val) != "" {
					cmd.Flags[i] = strings.TrimSpace(val)
				} else {
					cmd.Flags[i] = "nil"
				}
				toremove = append(toremove, index)
				hasFlag = true
			}
		}
	}
	if !hasFlag {
		cmd.Query.Content = strings.Join(arr, " ")
	} else {
		var queryarr []string
		for index, word := range arr {
			remove := false
			for _, toremoveindex := range toremove {
				if index == toremoveindex {
					remove = true
				}
			}
			if !remove {
				queryarr = append(queryarr, word)
			}
		}
		cmd.Query.Content = strings.Join(queryarr, " ")
	}
	blank := LexedAuthor{}
	if cmd.Author == blank {
		if cmd.CommandMessage.GuildID == "" {
			cmd.Author.DisplayName = cmd.CommandMessage.Author.Username
			cmd.Author.Colour = 0xc12fb7
			cmd.Author.Avatar = cmd.CommandMessage.Author.AvatarURL("")
		} else {
			member, _ := cmd.Session.GuildMember(cmd.CommandMessage.GuildID, cmd.CommandMessage.Author.ID)
			if member.Nick != "" {
				cmd.Author.DisplayName = member.Nick
			} else {
				cmd.Author.DisplayName = cmd.CommandMessage.Author.Username
			}
			cmd.Author.Avatar = cmd.CommandMessage.Author.AvatarURL("")
			cmd.Author.Colour = cmd.Session.State.UserColor(cmd.CommandMessage.Author.ID, cmd.CommandMessage.ChannelID)
		}
	}
}

var lexedcommands map[string]*LexedCommand = map[string]*LexedCommand{}

// Cleaner : Cleans old lexed commands.
func Cleaner() {
	for {
		time.Sleep(5 * time.Minute)
		var rmkeys []string
		for key, cmd := range lexedcommands {
			if time.Now().Sub(cmd.LastUsed) >= 10*time.Minute {
				rmkeys = append(rmkeys, key)
			}
		}
		for _, key := range rmkeys {
			delete(lexedcommands, key)
		}
	}
}

// NewLexedCommandForMessageAndSession : A function that creates a new lexed command for a message.
func NewLexedCommandForMessageAndSession(m *discordgo.Message, s *discordgo.Session) *LexedCommand {
	if val, ok := lexedcommands[m.ID]; ok {
		cmd := val
		cmd.CommandMessage = m
		cmd.LastUsed = time.Now()
		cmd.lex()
		return cmd
	} else {
		cmd := LexedCommand{
			CommandMessage:    m,
			Session:           s,
			Flags:             make(map[string]string),
			PaginatorPageName: "Page",
			LastUsed:          time.Now(),
		}
		cmd.lex()
		lexedcommands[m.ID] = &cmd
		return &cmd
	}
}

// LexedCommandFunction : Calls a command function with arguments lexed from discord events
func LexedCommandFunction(s *discordgo.Session, m *discordgo.Message, f CommandFunc) {
	cmd := NewLexedCommandForMessageAndSession(m, s)
	f(s, cmd)
}
