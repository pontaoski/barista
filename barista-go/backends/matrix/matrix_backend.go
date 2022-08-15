package matrix

import (
	"time"

	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/config"
	"github.com/appadeia/barista/barista-go/log"
	"github.com/matrix-org/gomatrix"
)

var backend = MatrixBackend{}

func init() {
	commandlib.RegisterBackend(&backend)
}

func (m MatrixBackend) Stats() (r *commandlib.BackendStats) {
	return
}

func (m MatrixBackend) CanGiveStats() bool {
	return false
}

type MatrixBackend struct{}

func (m MatrixBackend) Name() string {
	return "Matrix"
}

func (m MatrixBackend) ID() string {
	return "matrix"
}

func (m MatrixBackend) IsBotOwner(c commandlib.Context) bool {
	var ctx interface{} = c
	casted := ctx.(*MatrixContext)
	return casted.triggerEvent.Sender == config.BotConfig.Owner.Matrix
}

func (m MatrixBackend) Start(cancel chan struct{}) error {
	client, err := gomatrix.NewClient(config.BotConfig.Services.Matrix.Homeserver, "", "")
	if err != nil {
		return err
	}

	resp, err := client.Login(&gomatrix.ReqLogin{
		Type:     "m.login.password",
		User:     config.BotConfig.Services.Matrix.Username,
		Password: config.BotConfig.Services.Matrix.Password,
	})
	if err != nil {
		return err
	}

	client.SetCredentials(resp.UserID, resp.AccessToken)
	client.UserID = resp.UserID

	syncer := client.Syncer.(*gomatrix.DefaultSyncer)
	syncer.OnEventType("m.room.message", func(ev *gomatrix.Event) {
		MatrixMessage(client, ev)
	})
	syncer.OnEventType("m.room.member", func(ev *gomatrix.Event) {
		if val, ok := ev.Content["membership"]; ok {
			if val.(string) == "invite" {
				client.JoinRoom(ev.RoomID, "", nil)
			}
		}
	})
	go func() {
		for {
			client.Sync()
			time.Sleep(time.Millisecond * 500)
		}
	}()

	log.Info("Matrix session started")
	<-cancel
	return nil
}
