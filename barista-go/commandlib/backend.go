package commandlib

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/appadeia/barista/barista-go/log"
)

// A Backend represents a service that Barista can chat on
type Backend interface {
	Name() string
	Start(chan struct{}) error
}

var backends []Backend

// RegisterBackend registers a backend for Barista
func RegisterBackend(b Backend) {
	backends = append(backends, b)
}

// StartBackends starts all backends and waits on them to exit
func StartBackends() {
	wg := sync.WaitGroup{}
	for _, backend := range backends {
		wg.Add(1)
		log.Info("Starting backend %s", backend.Name())

		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		ic := make(chan struct{}, 1)
		go func() {
			<-sc
			ic <- struct{}{}
		}()

		go func(b Backend) {
			defer wg.Done()
			err := b.Start(ic)
			if err != nil {
				log.Fatal(log.BackendFailure, "Error starting backend %s: %s", b.Name(), err.Error())
			} else {
				log.Info("Backend %s exited", b.Name())
			}
		}(backend)
	}
	wg.Wait()
}
