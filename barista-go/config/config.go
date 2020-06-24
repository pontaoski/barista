package config

import (
	"github.com/BurntSushi/toml"
	"github.com/appadeia/barista/barista-go/log"
)

// Config defines what the bot configuration looks like
type Config struct {
	Services struct {
		Discord struct {
			Token string `toml:"token"`
		} `toml:"discord"`
		Telegram struct {
			Token string `toml:"token"`
		} `toml:"telegram"`
		Matrix struct {
			Homeserver string `toml:"homeserver"`
			Username   string `toml:"username"`
			Password   string `toml:"password"`
		} `toml:"matrix"`
	} `toml:"services"`
}

// BotConfig holds an instance of Config with values loaded
var BotConfig Config

func init() {
	log.Info("Reading config...")
	_, err := toml.DecodeFile("config.toml", &BotConfig)
	if err != nil {
		log.Fatal(log.ConfigFailure, "Failed to read config: %+v", err)
	}
}
