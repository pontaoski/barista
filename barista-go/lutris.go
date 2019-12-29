package barista

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Necroforger/dgwidgets"
	"github.com/bwmarrin/discordgo"
)

type LutrisPlatform struct {
	Name string `json:"name"`
}

type LutrisGenre struct {
	Name string `json:"name"`
}

type LutrisGame struct {
	Name          string           `json:"name"`
	ID            string           `json:"slug"`
	Year          int              `json:"year"`
	Platforms     []LutrisPlatform `json:"platforms"`
	Genres        []LutrisGenre    `json:"genres"`
	BannerURL     string           `json:"banner_url"`
	IconURL       string           `json:"icon_url"`
	SteamID       int              `json:"steamid"`
	GOGSlug       string           `json:"gogslug"`
	HumbleStoreID string           `json:"humblestoreid"`
}

type LutrisGames struct {
	Games []LutrisGame
	Count int
}

type LutrisRequest struct {
	Count    int          `json:"count"`
	Next     string       `json:"next"`
	Previous string       `json:"previous"`
	Results  []LutrisGame `json:"results"`
}

var games LutrisGames = LutrisGames{}
var lutrisRunning bool = false

func updateLutrisGames() {
	for {
		wg.Add(1)
		var tmpGames LutrisGames
		next := "https://lutris.net/api/games"
		for {
			if next == "" {
				break
			}
			resp, err := http.Get(next)
			if err != nil {
				continue
			}

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				continue
			}

			var lr LutrisRequest
			err = json.Unmarshal(body, &lr)
			if err != nil {
				continue
			}

			tmpGames.Games = append(tmpGames.Games, lr.Results...)
			if lr.Count != 0 {
				tmpGames.Count = lr.Count
			}

			next = lr.Next
		}
		games = tmpGames
		wg.Done()
		time.Sleep(5 * time.Minute)
	}
}

func startLutrisUpdates() {
	if !lutrisRunning {
		lutrisRunning = true
		go updateLutrisGames()
	}
}

var wg sync.WaitGroup

func Lutris(s *discordgo.Session, cmd *LexedCommand) {
	startLutrisUpdates()
	wg.Wait()
	cmd.PaginatorPageName = "Game"
	var filteredGames []LutrisGame

	for _, game := range games.Games {
		if cmd.Query.Content == "*" {
			filteredGames = games.Games
			break
		}
		if strings.Contains(strings.ToLower(game.Name), strings.ToLower(cmd.Query.Content)) || strings.Contains(strings.ToLower(game.ID), strings.ToLower(cmd.Query.Content)) {
			filteredGames = append(filteredGames, game)
		}
	}

	if cmd.Query.Content == "" {
		embed := NewEmbed().
			SetColor(0xff0000).
			SetTitle("Please provide a search term.")
		msgSend := discordgo.MessageSend{
			Embed: embed.MessageEmbed,
		}
		cmd.SendMessage(&msgSend)
		return
	}

	if len(filteredGames) == 0 {
		embed := NewEmbed().
			SetColor(0xff0000).
			SetTitle(fmt.Sprintf("No results were found for %s", cmd.Query.Content))
		msgSend := discordgo.MessageSend{
			Embed: embed.MessageEmbed,
		}
		cmd.SendMessage(&msgSend)
		return
	}

	paginator := dgwidgets.NewPaginator(cmd.Session, cmd.CommandMessage.ChannelID)

	for _, game := range filteredGames {
		embed := NewEmbed().
			SetTitle(game.Name).
			SetURL("https://lutris.net/games/" + game.ID + "/").
			SetColor(0xff9900).
			SetThumbnail("https://lutris.net" + game.IconURL).
			SetImage("https://lutris.net" + game.BannerURL)

		if len(game.Platforms) > 0 {
			var tmp []string
			for _, platform := range game.Platforms {
				tmp = append(tmp, platform.Name)
			}
			embed.AddField("Platforms", strings.Join(tmp, ", "), true)
		}
		if len(game.Genres) > 0 {
			var tmp []string
			for _, genre := range game.Genres {
				tmp = append(tmp, genre.Name)
			}
			embed.AddField("Genres", strings.Join(tmp, ", "), true)
		}
		if game.Year != 0 {
			embed.AddField("Year", fmt.Sprintf("%d", game.Year), true)
		}
		if game.SteamID != 0 {
			embed.AddField("Steam ID", fmt.Sprintf("[%d](https://store.steampowered.com/app/%d/)", game.SteamID, game.SteamID), true)
		}
		if game.HumbleStoreID != "" {
			embed.AddField("Humble Store ID", fmt.Sprintf("[%s](https://www.humblebundle.com/store/%s)", game.HumbleStoreID, game.HumbleStoreID), true)
		}

		paginator.Add(embed.MessageEmbed)
	}

	cmd.SendPaginator(paginator)
}
