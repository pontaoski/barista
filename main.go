package main

import (
	"github.com/appadeia/barista/barista-go/telegram"
)

func main() {
	go barista.Main()
	telegram.Main()
}
