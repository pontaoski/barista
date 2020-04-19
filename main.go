package main

import (
	"github.com/appadeia/barista/barista-go"
	"github.com/appadeia/barista/barista-go/matrix"
	"github.com/appadeia/barista/barista-go/telegram"
)

func main() {
	go barista.Main()
	go matrix.Main()
	telegram.Main()
}
