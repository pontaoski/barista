package harmony

import (
	"fmt"

	"github.com/appadeia/barista/barista-go/backends/harmony/client"
	corev1 "github.com/appadeia/barista/barista-go/backends/harmony/gen/core"
	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/config"
	"github.com/appadeia/barista/barista-go/log"
)

// Backend is an instance of a connection to a Harmony homeserver
type Backend struct {
	email      string
	password   string
	homeserver string
}

func init() {
	backend := new(Backend)
	backend.email = config.BotConfig.Services.Harmony.Email
	backend.password = config.BotConfig.Services.Harmony.Password
	backend.homeserver = config.BotConfig.Services.Harmony.Homeserver

	commandlib.RegisterBackend(backend)
}

// Stats gives the stats of the Harmony backend
func (b *Backend) Stats() (r *commandlib.BackendStats) {
	return
}

// CanGiveStats indicates whether this backend can give stats
func (b *Backend) CanGiveStats() bool {
	return false
}

// Name is the name of the Discord backend
func (b *Backend) Name() string {
	return fmt.Sprintf("Harmony (%s on %s)", b.email, b.homeserver)
}

// IsBotOwner checks whether the context is of the bot owner
func (b *Backend) IsBotOwner(c commandlib.Context) bool {
	panic("Implement me")
}

// Start starts the Harmony backend
func (b *Backend) Start(cancel chan struct{}) error {
	client, err := client.NewClient(b.homeserver, b.email, b.password)
	if err != nil {
		return err
	}

	channel, err := client.Start()
	if err != nil {
		return err
	}

	log.Info("%s session started", b.Name())

	for {
		select {
		case ev := <-channel:
			switch a := ev.Event.Event.(type) {
			case *corev1.Event_SentMessage:
				b.Message(&client.Client, a.SentMessage.Message)
			}
		case <-cancel:
			return nil
		}
	}
}
