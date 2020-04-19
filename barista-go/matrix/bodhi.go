package matrix

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/matrix-org/gomatrix"
)

type Release struct {
	Display string `json:"long_name"`
}

type Build struct {
	NVR       string `json:"nvr"`
	ReleaseId int    `json:"release_id"`
	Signed    bool   `json:"signed"`
	Type      string `json:"rpm"`
	Epoch     int    `json:"epoch"`
}

type User struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type Update struct {
	Karma       int     `json:"karma"`
	Name        string  `json:"title"`
	Description string  `json:"notes"`
	Release     Release `json:"release"`
	Builds      []Build `json:"builds"`
	Author      User    `json:"user"`
	URL         string  `json:"url"`
	Status      string  `json:"status"`
}

type UpdateWrapper struct {
	Update  Update `json:"update"`
	CanEdit bool   `json:"can_edit"`
}

func Bodhi(client *gomatrix.Client, ev *gomatrix.Event, body []string) {
	words := body

	if strings.Contains(strings.Join(body, " "), "FEDORA-") {
		for _, word := range words {
			if strings.HasPrefix(word, "FEDORA-") {
				resp, err := http.Get(fmt.Sprintf("https://bodhi.fedoraproject.org/updates/%s", word))
				if err != nil {
					continue
				}
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					continue
				}

				var updateWrapper UpdateWrapper
				err = json.Unmarshal(body, &updateWrapper)
				if err != nil {
					continue
				}

				update := updateWrapper.Update

				builds := fmt.Sprintf(`<a href="https://koji.fedoraproject.org/koji/search?terms=%s&type=build&match=glob">%s</a>`, update.Builds[0].NVR, update.Builds[0].NVR)

				for _, build := range update.Builds[1:] {
					builds = fmt.Sprintf(`%s, <a href="https://koji.fedoraproject.org/koji/search?terms=%s&type=build&match=glob">%s</a>`, builds, build.NVR, build.NVR)
				}
				var karmamsg string
				if update.Karma == 0 {
					karmamsg = fmt.Sprintf("%d", update.Karma)
				} else if update.Karma > 0 {
					karmamsg = fmt.Sprintf("%d ✅", update.Karma)
				} else if update.Karma < 0 {
					karmamsg = fmt.Sprintf("%d ❌", update.Karma)
				}
				output := fmt.Sprintf(`
<table>
	<tr>
		<td>
			<a href="%s">%s</a> by %s <br>
			%s
			<table>
				<tr>
					<th>Distro</th>
					<th>Karma</th>
					<th>Builds</th>
					<th>Status</th>
				</tr>
				<tr>
					<td>%s</th>
					<td>%s</th>
					<td>%s</th>
					<td>%s</th>
				</tr>
			</table>
		</td>
	</tr>
</table>

				`, update.URL, update.Name, update.Author.Name, update.Description, update.Release.Display, karmamsg, builds, strings.Title(update.Status))

				go SendHTMLMessage(client, ev.RoomID, output, "")
			}
		}
	}
}
