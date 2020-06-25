package irc

import (
	"sync"

	"github.com/lrstanley/girc"
)

var IRCMutex sync.Mutex
var IRCChans []chan girc.Event

func NextMessage() chan girc.Event {
	IRCMutex.Lock()
	defer IRCMutex.Unlock()

	channel := make(chan girc.Event)
	IRCChans = append(IRCChans, channel)
	return channel
}

func IRCOnetimeHandler(c *girc.Client, e girc.Event) {
	IRCMutex.Lock()
	defer IRCMutex.Unlock()

	for len(IRCChans) > 0 {
		var channel chan girc.Event
		channel, IRCChans = IRCChans[len(IRCChans)-1], IRCChans[:len(IRCChans)-1]
		channel <- e
	}
}
