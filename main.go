package main

import (
	"time"

	"github.com/appadeia/barista/barista-go"
	"github.com/appadeia/barista/barista-go/web"
)

func main() {
	go func() {
		time.Sleep(5 * time.Second)
		web.Main()
	}()
	barista.Main()
}
