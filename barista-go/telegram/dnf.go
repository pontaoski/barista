package telegram

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/godbus/dbus"
)

type Package struct {
	name        string
	desc        string
	vers        string
	downsize    string
	installsize string
}

type Distro struct {
	matches      []string
	displayName  string
	queryKitName string
}

var Distros []Distro = []Distro{
	{
		displayName:  "openSUSE Tumbleweed",
		queryKitName: "tumbleweed",
		matches:      []string{"opensuse", "os", "tumbleweed", "tw"},
	},
	{
		displayName:  "openSUSE Leap",
		queryKitName: "leap",
		matches:      []string{"leap", "opensuse-leap", "os-leap"},
	},
	{
		displayName:  "Fedora",
		queryKitName: "fedora",
		matches:      []string{"fedora"},
	},
	{
		displayName:  "Mageia",
		queryKitName: "mageia",
		matches:      []string{"mageia"},
	},
	{
		displayName:  "OpenMandriva",
		queryKitName: "openmandriva",
		matches:      []string{"openmandriva"},
	},
	{
		displayName:  "CentOS",
		queryKitName: "centos",
		matches:      []string{"centos"},
	},
}

func resolveDistro(name string) (Distro, bool) {
	var distro Distro
	set := false

	for _, dist := range Distros {
		for _, match := range dist.matches {
			if strings.ToLower(name) == strings.ToLower(match) {
				distro = dist
				set = true
			}
		}
	}

	return distro, set
}

func searchDnf(search, distro string) ([]Package, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return []Package{}, err
	}
	obj := conn.Object("com.github.Appadeia.QueryKit", "/com/github/Appadeia/QueryKit")
	var pkgs [][]interface{}
	err = obj.Call("com.github.Appadeia.QueryKit.SearchPackages", 0, search, distro).Store(&pkgs)
	if err != nil {
		return []Package{}, err
	}
	var packages []Package
	for _, pkg := range pkgs {
		packages = append(packages, Package{
			name:        pkg[0].(string),
			desc:        pkg[1].(string),
			vers:        pkg[2].(string),
			downsize:    humanize.Bytes(uint64(pkg[3].(int32))),
			installsize: humanize.Bytes(uint64(pkg[4].(int32))),
		})
	}
	return packages, nil
}

func dnf(msg *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	args := strings.Fields(msg.CommandArguments())
	if len(args) < 2 {
		msg := tgbotapi.NewMessage(msg.Chat.ID, "Please specify arguments like so: `dnfsearch distro query`")
		bot.Send(msg)
		return
	}
	dist, set := resolveDistro(args[0])
	if !set {
		msg := tgbotapi.NewMessage(msg.Chat.ID, "That is not a distro supported by Barista.")
		bot.Send(msg)
		return
	}
	pkgs, err := searchDnf(strings.Join(args[1:], " "), dist.queryKitName)
	if err != nil {
		msg := tgbotapi.NewMessage(msg.Chat.ID, "That is not a distro supported by Barista.")
		bot.Send(msg)
		return
	}
	paginator := NewPaginator(bot)
	for idx, pkg := range pkgs {
		msg := tgbotapi.NewMessage(msg.Chat.ID, "")
		msg.Text = fmt.Sprintf(`<b>%s</b> — %s
<i>%s</i>
Download Size: %s
Install Size: %s

Package %d of %d`, pkg.name, pkg.vers, pkg.desc, pkg.downsize, pkg.installsize, idx+1, len(pkgs))
		msg.ParseMode = tgbotapi.ModeHTML
		paginator.AddPage(msg)
	}
	paginator.Send()
}
