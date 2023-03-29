package config

import (
	"github.com/BurntSushi/toml"
	"github.com/appadeia/barista/barista-go/log"
)

// Config defines what the bot configuration looks like
type Config struct {
	Services struct {
		Backends []string `toml:"active"`

		Discord []struct {
			Token string `toml:"token"`
			Name  string `toml:"name"`
		} `toml:"discord"`
		Telegram struct {
			Token string `toml:"token"`
		} `toml:"telegram"`
		Matrix struct {
			Homeserver string `toml:"homeserver"`
			Username   string `toml:"username"`
			Password   string `toml:"password"`
		} `toml:"matrix"`
		IRC struct {
			Server   string   `toml:"server"`
			Nickname string   `toml:"nickname"`
			Username string   `toml:"username"`
			Channels []string `toml:"channels"`
		} `toml:"irc"`
		Harmony struct {
			Homeserver string `toml:"homeserver"`
			UserID     uint64 `toml:"userID"`
			Token      string `toml:"token"`
		} `toml:"harmony"`
	} `toml:"services"`
	Owner struct {
		Discord  string `toml:"discord"`
		Matrix   string `toml:"matrix"`
		Telegram int    `toml:"telegram"`
	} `toml:"owner"`
	Tokens struct {
		InventKDEOrg string `toml:"invent.kde.org"`
		OpenAI       string `toml:"openai"`
	} `toml:"tokens"`
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
