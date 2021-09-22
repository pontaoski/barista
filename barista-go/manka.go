package barista

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Tnze/go-mc/bot"
	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/google/uuid"
	"google.golang.org/api/chat/v1"
)

type status struct {
	Description chat.Message
	Players     struct {
		Max    int
		Online int
		Sample []struct {
			ID   uuid.UUID
			Name string
		}
	}
	Version struct {
		Name     string
		Protocol int
	}
	Favicon string
}

func init() {
	commandlib.RegisterCommand(commandlib.Command{
		Name:  "musi manka",
		Usage: "o lukin e manka",
		ID:    "manka",
		Matches: []string{
			"ilo o manka",
		},
		Action: func(c commandlib.Context) {
			ip, port := "musi.musimanka.com", 25565
			if strings.Contains(c.Content(), "ante") {
				ip, port = "5.9.30.86", 25565
			}
			c.SendMessage("primary", commandlib.Embed{
				Title: commandlib.EmbedHeader{
					Text: "mi pali...",
				},
			})
			resp, _, err := bot.PingAndList(ip, port)
			if err != nil {
				c.SendMessage("primary", commandlib.ErrorEmbed("There was an error getting server info: "+err.Error()))
				return
			}

			var s status
			err = json.Unmarshal(resp, &s)
			if err != nil {
				c.SendMessage("primary", commandlib.ErrorEmbed("There was an error getting server info: "+err.Error()))
				return
			}

			c.SendMessage("primary", commandlib.Embed{
				Title: commandlib.EmbedHeader{
					Text: fmt.Sprintf("jan %d li musi lon %s:%d", s.Players.Online, ip, port),
				},
				Fields: []commandlib.EmbedField{
					{
						Title: "jan musi",
						Body: func() string {
							data := strings.Join(func() (ret []string) {
								for _, item := range s.Players.Sample {
									ret = append(ret, "- "+item.Name)
								}
								return
							}(), "\n")
							if data == "" {
								return "jan ala li musi. o musi!"
							}
							return data
						}(),
					},
					{
						Title: "nanpa musi",
						Body:  s.Version.Name,
					},
				},
			})
		},
	})
}
