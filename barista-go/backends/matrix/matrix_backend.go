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

	client.SetCredentials(config.BotConfig.Services.Matrix.Username, config.BotConfig.Services.Matrix.Token)
	client.UserID = config.BotConfig.Services.Matrix.Username

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
			err := client.Sync()
			if err != nil {
				log.Error("matrix syncer failed: %v", err)
				return
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()

	log.Info("Matrix session started")
	<-cancel
	return nil
}
