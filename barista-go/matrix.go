package barista

import (
	"fmt"
	"time"

	"github.com/appadeia/barista/barista-go/commandlib/matrix"
	"github.com/matrix-org/gomatrix"
)

func MatrixMain() {
	fmt.Println("Initializing Barista Matrix...")

	sec := Cfg.Section("Matrix")

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
		matrix.MatrixMessage(client, ev)
	})
	syncer.OnEventType("m.room.member", func(ev *gomatrix.Event) {
		if val, ok := ev.Content["membership"]; ok {
			if val.(string) == "invite" {
				client.JoinRoom(ev.RoomID, "", nil)
			}
		}
	})
	go func() {
		println("Barista Matrix is now running.")
		for {
			client.Sync()
			time.Sleep(time.Millisecond * 500)
		}
	}()
}
