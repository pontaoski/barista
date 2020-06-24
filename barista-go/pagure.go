package barista

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/log"
)

type PagureUser struct {
	FullName  string `json:"fullname"`
	UserName  string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type PagureUserData struct {
	User PagureUser `json:"user"`
}

func (self PagureUser) GetAvatarURL(instance PagureInstance) string {
	user, err := http.Get(instance.Path("user/" + self.UserName))
	if err != nil {
		return ""
	}

	defer user.Body.Close()
	body, err := ioutil.ReadAll(user.Body)
	if err != nil {
		return ""
	}

	var userData PagureUserData
	err = json.Unmarshal(body, &userData)
	if err != nil {
		return ""
	}

	return userData.User.AvatarURL
}

type PagureIssue struct {
	Title       string     `json:"title"`
	Description string     `json:"content"`
	Status      string     `json:"status"`
	Tags        []string   `json:"tags"`
	User        PagureUser `json:"user"`
}

type PagureRepo struct {
	Fullname string `json:"fullname"`
}

type PagurePullRequest struct {
	Title          string     `json:"title"`
	InitialComment string     `json:"initial_comment"`
	Status         string     `json:"status"`
	User           PagureUser `json:"user"`
	FromRepo       PagureRepo `json:"repo_from"`
	FromBranch     string     `json:"branch_from"`
	ToBranch       string     `json:"branch"`
}

type PagureInstance struct {
	Name    string
	Matches []string
	Icon    string
	URL     string
	UserURL string
	Colour  int
}

func (self PagureInstance) UserPath(path string) string {
	return self.UserURL + path
}

func (self PagureInstance) Path(path string) string {
	return self.URL + path
}

var PagureInstances []PagureInstance = []PagureInstance{
	{
		Name:    "Fedora Package Sources",
		Matches: []string{"srcfpo"},
		Icon:    "https://fedoraproject.org/w/uploads/archive/e/e5/20110717032101%21Fedora_infinity.png",
		Colour:  0x0b57a4,
		URL:     "https://src.fedoraproject.org/api/0/",
		UserURL: "https://src.fedoraproject.org/",
	},
	{
		Name:    "pagure.io",
		Matches: []string{"pagureio"},
		Icon:    "https://avatars0.githubusercontent.com/u/44003730?s=400&v=4",
		Colour:  0x0b57a4,
		URL:     "https://pagure.io/api/0/",
		UserURL: "https://pagure.io/",
	},
	{
		Name:    "CentOS Git Server",
		Matches: []string{"gitco"},
		Icon:    "https://upload.wikimedia.org/wikipedia/commons/thumb/6/63/CentOS_color_logo.svg/1024px-CentOS_color_logo.svg.png",
		Colour:  0x951C7A,
		URL:     "https://git.centos.org/api/0/",
		UserURL: "https://git.centos.org/",
	},
}

func init() {
	var matches []string
	for _, instance := range PagureInstances {
		matches = append(matches, instance.Matches...)
	}
	commandlib.RegisterTag(commandlib.Tag{
		Name:  I18n("Pagure"),
		Usage: I18n("Link to Pagure issues or pull requests"),
		Examples: `srcfpo#rpms/python-flask-babel#PR5
pagureio/fedora-comps#492`,
		Samples: func() []commandlib.TagSample {
			var ret []commandlib.TagSample
			for _, pagure := range PagureInstances {
				for _, match := range pagure.Matches {
					ret = append(ret, commandlib.TagSample{Tag: fmt.Sprintf("%s#repo#1234", match), Desc: pagure.Name + " Pull Requests"})
					ret = append(ret, commandlib.TagSample{Tag: fmt.Sprintf("%s#repo#PR1234", match), Desc: pagure.Name + " Issues"})
				}
			}
			return ret
		}(),
		ID:     "pagure",
		Match:  regexp.MustCompile(fmt.Sprintf("(%s)", strings.Join(matches, "|"))),
		Action: Pagure,
	})
}

func Pagure(c commandlib.Context) {
	var embeds []commandlib.Embed
	for _, word := range c.Args() {
		for _, pagure := range PagureInstances {
			for _, match := range pagure.Matches {
				if strings.HasPrefix(word, match+"#") {
					items := strings.Split(word, "#")
					if len(items) < 3 { // instance#project#number
						continue
					}
					if strings.HasPrefix(strings.ToLower(items[2]), "pr") {
						id := strings.TrimPrefix(strings.ToLower(items[2]), "pr")
						pr, err := http.Get(pagure.Path(items[1] + "/pull-request/" + id))
						if err != nil {
							log.Error("%+v", err)
							continue
						}
						defer pr.Body.Close()
						body, err := ioutil.ReadAll(pr.Body)
						if err != nil {
							log.Error("%+v", err)
							continue
						}

						var pullRequest PagurePullRequest
						err = json.Unmarshal(body, &pullRequest)
						if err != nil {
							log.Error("%+v", err)
							continue
						}

						embeds = append(embeds, commandlib.Embed{
							Title: commandlib.EmbedHeader{
								Text: fmt.Sprintf("PR#%s: %s", id, pullRequest.Title),
								URL:  pagure.UserPath(items[1] + "/pull-request/" + id),
							},
							Body:   pullRequest.InitialComment,
							Colour: pagure.Colour,
							Fields: []commandlib.EmbedField{
								{
									Title:  c.I18n("Status"),
									Body:   pullRequest.Status,
									Inline: true,
								},
								{
									Title: c.I18n("Pull Request"),
									Body: fmt.Sprintf(
										c.I18n("Merge [%s/%s](%s) into [%s](%s)"),
										pullRequest.FromRepo.Fullname,
										pullRequest.FromBranch,
										strings.ReplaceAll(pagure.UserPath(pullRequest.FromRepo.Fullname+"/tree/"+pullRequest.FromBranch), "forks", "fork"),
										pullRequest.ToBranch,
										pagure.UserPath(items[1]+"/tree/"+pullRequest.ToBranch),
									),
									Inline: true,
								},
							},
							Header: commandlib.EmbedHeader{
								Text: pullRequest.User.FullName,
								Icon: pullRequest.User.GetAvatarURL(pagure),
							},
							Footer: commandlib.EmbedHeader{
								Text: pagure.Name,
								Icon: pagure.Icon,
							},
						})
					} else {
						issue, err := http.Get(pagure.Path(items[1] + "/issue/" + items[2]))
						if err != nil {
							log.Error("%+v", err)
							continue
						}
						defer issue.Body.Close()
						body, err := ioutil.ReadAll(issue.Body)
						if err != nil {
							log.Error("%+v", err)
							continue
						}

						var pIssue PagureIssue
						err = json.Unmarshal(body, &pIssue)
						if err != nil {
							log.Error("%+v", err)
							continue
						}

						embeds = append(embeds, commandlib.Embed{
							Title: commandlib.EmbedHeader{
								Text: fmt.Sprintf(c.I18n("Issue #%s: %s - %s"), items[2], pIssue.Title, items[1]),
								URL:  pagure.UserPath(items[1] + "/issue/" + items[2]),
							},
							Body:   pIssue.Description,
							Colour: pagure.Colour,
							Fields: []commandlib.EmbedField{
								{
									Title:  c.I18n("Status"),
									Body:   pIssue.Status,
									Inline: true,
								},
							},
							Header: commandlib.EmbedHeader{
								Text: pIssue.User.FullName,
								Icon: pIssue.User.GetAvatarURL(pagure),
							},
						})
					}
				}
			}
		}
	}
	if len(embeds) > 0 {
		c.SendTags("primary", embeds)
	}
}
