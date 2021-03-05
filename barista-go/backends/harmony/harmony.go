package harmony

import (
	"fmt"

	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/config"
	"github.com/appadeia/barista/barista-go/log"
	"github.com/harmony-development/shibshib"
)

// Backend is an instance of a connection to a Harmony homeserver
type Backend struct {
	homeserver string
	userID     uint64
	token      string
}

func init() {
	backend := new(Backend)
	backend.homeserver = config.BotConfig.Services.Harmony.Homeserver
	backend.userID = config.BotConfig.Services.Harmony.UserID
	backend.token = config.BotConfig.Services.Harmony.Token

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
	return fmt.Sprintf("Harmony (%d on %s)", b.userID, b.homeserver)
}

// IsBotOwner checks whether the context is of the bot owner
func (b *Backend) IsBotOwner(c commandlib.Context) bool {
	panic("Implement me")
}

// Start starts the Harmony backend
func (b *Backend) Start(cancel chan struct{}) error {
	client, err := shibshib.NewFederatedClient(b.homeserver, b.token, b.userID)
	if err != nil {
		return err
	}

	log.Info("%s session started", b.Name())

	evs, err := client.Start()
	if err != nil {
		return err
	}

	for {
		select {
		case ev := <-evs:
			b.Message(ev.Client, ev.Event)
		case <-cancel:
			return nil
		}
	}
}
