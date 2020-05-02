package barista

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/ini.v1"
)

// Barista's config
var Cfg *ini.File

// Main : Call this function to start the bot's main loop.
func Main() {
	var err error
	Cfg, err = ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Failed to load config.ini")
		os.Exit(1)
	}
	go DiscordMain()
	go TelegramMain()
	go MatrixMain()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	println("Closing connections...")
	time.Sleep(3 * time.Second)
}
