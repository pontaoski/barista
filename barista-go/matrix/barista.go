package matrix

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/matrix-org/gomatrix"
	"gopkg.in/ini.v1"
)

func prettyPrint(v interface{}) {
	out, _ := json.MarshalIndent(v, "", "    ")
	println(string(out))
}

func Main() {
	fmt.Println("Initializing Barista Matrix...")

	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Failed to load config.ini, stopping Barista Matrix...")
		return
	}
	sec := cfg.Section("Matrix")

	client, err := gomatrix.NewClient(sec.Key("Homeserver").String(), "", "")
	if err != nil {
		fmt.Printf("Failed to create connection to Matrix, stopping Barista Matrix...")
		return
	}

	resp, err := client.Login(&gomatrix.ReqLogin{
		Type:     "m.login.password",
		User:     sec.Key("User").String(),
		Password: sec.Key("Password").String(),
	})
	if err != nil {
		fmt.Printf("Failed to authenticate with Matrix, stoppig Barista Matrix...")
		return
	}

	client.SetCredentials(resp.UserID, resp.AccessToken)
	client.UserID = resp.UserID

	syncer := client.Syncer.(*gomatrix.DefaultSyncer)
	syncer.OnEventType("m.room.message", func(ev *gomatrix.Event) {
		if val, ok := ev.Content["body"]; ok {
			words := strings.Fields(val.(string))
			if len(words) > 2 {
				switch words[0] + " " + words[1] {
				case "dnf search", "dnf se":
					go DnfSearch(client, ev, words[2:])
				case "dnf repoquery", "dnf rq":
					go DnfRepoquery(client, ev, words[2:])
				}
			}
			go Bugzilla(client, ev, words)
			go Bodhi(client, ev, words)
		}
	})
	syncer.OnEventType("m.room.member", func(ev *gomatrix.Event) {
		if val, ok := ev.Content["membership"]; ok {
			if val.(string) == "invite" {
				_, err := client.JoinRoom(ev.RoomID, "", nil)
				if err != nil {
					println("Failed to join", ev.RoomID)
				}
			}
		}
		prettyPrint(ev)
	})

	go func() {
		println("Barista Matrix is now running.")
		for {
			client.Sync()
			time.Sleep(time.Millisecond * 500)
		}
	}()
}

type Msg struct {
	Format string `json:"format"`
	HTML   string `json:"formatted_body"`
	Body   string `json:"body"`
	Type   string `json:"msgtype"`
}

func SendHTMLMessage(client *gomatrix.Client, roomID, html, text string) {
	message := Msg{
		Format: "org.matrix.custom.html",
		Body:   text,
		HTML:   html,
		Type:   "m.text",
	}
	if message.Body == "" {
		message.Body = "This message can only be viewed in Matrix clients that support HTML."
	}
	client.SendMessageEvent(roomID, "m.room.message", message)
}

func SendMessage(client *gomatrix.Client, roomID, text string) {
	message := Msg{
		Body: text,
		Type: "m.text",
	}
	client.SendMessageEvent(roomID, "m.room.message", message)
}

func WrapCode(code string) string {
	return "<pre><code>" + code + "</code></pre>"
}
