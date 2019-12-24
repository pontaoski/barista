package barista

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Profile(s *discordgo.Session, cmd *LexedCommand) {
	helpmsg := "```" + `dsconfig
# Syntax: sudo profile --flag value
( --user | -u )
	Get the user specified.
( --set-desktop-environment | -w )
	Set your desktop environment or window manager.
( --set-distro | -d )
	Set your distro.
( --set-shell | -s )
	Set your command line shell.
( --set-editor | -e )
	Set your editor.
( --set-languages | -p )
	Set your programming languages.
( --set-blurb | -b )
	Set your profile blurb.
( --set-screenshot | -i )
	Set your screenshot to a PNG.` + "```"
	used := false
	updated := false
	var user string
	if val := cmd.GetFlagPair("-w", "--set-desktop-environment"); val != "" {
		SetGlobalKey(cmd.MemKey("de"), val)
		used, updated = true, true
	}
	if val := cmd.GetFlagPair("-s", "--set-shell"); val != "" {
		SetGlobalKey(cmd.MemKey("shell"), val)
		used, updated = true, true
	}
	if val := cmd.GetFlagPair("-d", "--set-distro"); val != "" {
		SetGlobalKey(cmd.MemKey("distro"), val)
		used, updated = true, true
	}
	if val := cmd.GetFlagPair("-e", "--set-editor"); val != "" {
		SetGlobalKey(cmd.MemKey("editor"), val)
		used, updated = true, true
	}
	if val := cmd.GetFlagPair("-p", "--set-languages"); val != "" {
		SetGlobalKey(cmd.MemKey("langs"), val)
		used, updated = true, true
	}
	if val := cmd.GetFlagPair("-b", "--set-blurb"); val != "" {
		SetGlobalKey(cmd.MemKey("blurb"), val)
		used, updated = true, true
	}
	if val := cmd.GetFlagPair("-i", "--set-screenshot"); val != "" {
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
	if !used {
		embed := NewEmbed().
			SetColor(0xff0000).
			SetDescription(helpmsg).
			SetTitle("Please specify some arguments.")
		msgSend := discordgo.MessageSend{
			Embed: embed.MessageEmbed,
		}
		cmd.SendMessage(&msgSend)
		return
	}
	if updated {
		embed := NewEmbed().
			SetColor(cmd.Author.Colour).
			SetTitle("Profile updated!").
			SetDescription(GetGlobalKey(cmd.MemKey("blurb"))).
			AddField("Distro", Default(GetGlobalKey(cmd.MemKey("distro")), "No distro set."), true).
			AddField("Shell", Default(GetGlobalKey(cmd.MemKey("shell")), "No shell set."), true).
			AddField("Editor", Default(GetGlobalKey(cmd.MemKey("editor")), "No editor set."), true).
			AddField("Programming Languages", Default(GetGlobalKey(cmd.MemKey("langs")), "No programming languages set."), true).
			AddField("DE/WM", Default(GetGlobalKey(cmd.MemKey("de")), "No DE/WM set."), true).
			SetImage(Default(GetGlobalKey(cmd.MemKey("screenshot")), ""))
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
		embed := NewEmbed().
			SetColor(colour).
			SetTitle(disp+"'s Profile").
			SetDescription(GetGlobalKey(val+"blurb")).
			AddField("Distro", Default(GetGlobalKey(val+"distro"), "No distro set."), true).
			AddField("Shell", Default(GetGlobalKey(val+"shell"), "No shell set."), true).
			AddField("Editor", Default(GetGlobalKey(val+"editor"), "No editor set."), true).
			AddField("Programming Languages", Default(GetGlobalKey(val+"langs"), "No programming languages set."), true).
			AddField("DE/WM", Default(GetGlobalKey(val+"de"), "No DE/WM set."), true)

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
