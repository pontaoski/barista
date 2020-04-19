package matrix

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/matrix-org/gomatrix"
)

type BugzillaInstance struct {
	Name    string
	Matches []string
	Icon    string
	URL     string
	Colour  int
}

var BugzillaInstances []BugzillaInstance = []BugzillaInstance{
	{
		Name:    "Red Hat Bugzilla",
		Matches: []string{"rh", "rhbz"},
		Icon:    "https://lcp.world/images/trackers/RedHat.png",
		URL:     "https://bugzilla.redhat.com",
		Colour:  0xEE0000,
	},
	{
		Name:    "SUSE Bugzilla",
		Matches: []string{"bsc", "susebz"},
		Icon:    "https://lcp.world/images/trackers/SUSE.png",
		URL:     "https://bugzilla.suse.com",
		Colour:  0x2D35F,
	},
	{
		Name:    "openSUSE Bugzilla",
		Matches: []string{"boo"},
		Icon:    "https://lcp.world/images/trackers/openSUSE.png",
		URL:     "https://bugzilla.opensuse.org",
		Colour:  0x73BA25,
	},
	{
		Name:    "Novell Bugzilla",
		Matches: []string{"bnc"},
		Icon:    "https://lcp.world/images/trackers/Novell.png",
		URL:     "https://bugzilla.novell.com",
		Colour:  0xFF1E1E,
	},
	{
		Name:    "GNOME Bugzilla",
		Matches: []string{"bgo"},
		Icon:    "https://lcp.world/images/trackers/GNOME.png",
		URL:     "https://bugzilla.gnome.org",
		Colour:  0x4A86CF,
	},
	{
		Name:    "Kernel Bugzilla",
		Matches: []string{"bko"},
		Icon:    "https://lcp.world/images/trackers/Linux.png",
		URL:     "https://bugzilla.kernel.org",
		Colour:  0xFFD133,
	},
	{
		Name:    "Mozilla Bugzilla",
		Matches: []string{"bmo"},
		Icon:    "https://lcp.world/images/trackers/Mozilla.png",
		URL:     "https://bugzilla.mozilla.org",
		Colour:  0xFFFFFF,
	},
	{
		Name:    "Samba Bugzilla",
		Matches: []string{"bso"},
		Icon:    "https://lcp.world/images/trackers/Samba.png",
		URL:     "13172736",
		Colour:  0xC90000,
	},
	{
		Name:    "Xfce Bugzilla",
		Matches: []string{"bxo"},
		Icon:    "https://lcp.world/images/trackers/Xfce.png",
		URL:     "https://bugzilla.xfce.org",
		Colour:  0x2284F2,
	},
	{
		Name:    "KDE Bugzilla",
		Matches: []string{"kde"},
		Icon:    "https://lcp.world/images/trackers/KDE.png",
		URL:     "https://bugs.kde.org",
		Colour:  1939955,
	},
	{
		Name:    "Freedesktop Bugzilla",
		Matches: []string{"fdo"},
		Icon:    "https://lcp.world/images/trackers/Freedesktop.png",
		URL:     "https://bugs.freedesktop.org",
		Colour:  3899566,
	},
	{
		Name:    "GCC Bugzilla",
		Matches: []string{"gcc"},
		Icon:    "https://lcp.world/images/trackers/GCC.png",
		URL:     "https://gcc.gnu.org/bugzilla/",
		Colour:  16764843,
	},
	{
		Name:    "Mageia Bugzilla",
		Matches: []string{"mga", "mgabz"},
		Icon:    "https://lcp.world/images/trackers/Mageia.png",
		URL:     "https://bugs.mageia.org/",
		Colour:  2332628,
	},
}

func Bugzilla(client *gomatrix.Client, ev *gomatrix.Event, body []string) {
	words := body

	for _, word := range words {
	BugzillaLoop:
		for _, bugzilla := range BugzillaInstances {
			for _, match := range bugzilla.Matches {
				if strings.HasPrefix(word, match+"#") {
					tag := strings.TrimPrefix(word, match+"#")
					bug, err := http.Get(fmt.Sprintf("%s/show_bug.cgi?id=%s&ctype=xml", bugzilla.URL, tag))
					if err != nil {
						continue BugzillaLoop
					}
					body, err := ioutil.ReadAll(bug.Body)
					if err != nil {
						continue BugzillaLoop
					}

					doc, err := xmlquery.Parse(strings.NewReader(string(body)))
					if err != nil {
						continue BugzillaLoop
					}

					bugs, err := xmlquery.QueryAll(doc, "/bugzilla/bug")
					if err != nil {
						continue BugzillaLoop
					}

					if len(bugs) == 0 {
						continue BugzillaLoop
					}

					if bugs[0].SelectElement("short_desc") == nil {
						continue BugzillaLoop
					}

					table := fmt.Sprintf(`
<table>
	<tr>
		<th>Bug #%s at %s</th>
		<th>
			<a href="%s">%s</a>
		</th>
		<th> </th>
	</tr>
	<tr>
		<td>
			<b>Product</b><br/>
			%s
		</td>
		<td>
			<b>Version</b><br/>
			%s
		</td>
		<td>
			<b>Component</b><br/>
			%s
		</td>
	</tr>
	<tr>
		<td>
			<b>Priority</b><br/>
			%s
		</td>
		<td>
			<b>Severity</b><br/>
			%s
		</td>
		<td>
			<b>Status</b><br/>
			%s
		</td>
	</tr>
</table>
					`, tag,
						bugzilla.Name,
						fmt.Sprintf(
							"%s/show_bug.cgi?id=%s",
							bugzilla.URL,
							tag,
						),
						bugs[0].SelectElement("short_desc").InnerText(),
						bugs[0].SelectElement("product").InnerText(),
						bugs[0].SelectElement("version").InnerText(),
						bugs[0].SelectElement("component").InnerText(),
						bugs[0].SelectElement("priority").InnerText(),
						bugs[0].SelectElement("bug_severity").InnerText(),
						bugs[0].SelectElement("bug_status").InnerText(),
					)

					go SendHTMLMessage(client, ev.RoomID, table, "")
				}
			}
		}
	}
}
