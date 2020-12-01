package barista

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/config"
	"github.com/appadeia/barista/barista-go/log"
	"github.com/gomarkdown/markdown"
	"github.com/xanzy/go-gitlab"
)

const validSlug = `([a-zA-Z_\-\.]*)`
const slash = `\/`
const excl = `\!`
const number = `([0-9]+)`

var inventRegex = regexp.MustCompile(validSlug + slash + validSlug + excl + number)

var glClient *gitlab.Client

var htmlReplacer = strings.NewReplacer(
	`<img src="`, `<img src`,
)

func init() {
	var err error
	glClient, err = gitlab.NewClient(config.BotConfig.Tokens.InventKDEOrg, gitlab.WithBaseURL("https://invent.kde.org/api/v4"))
	if err != nil {
		log.Fatal(log.UnknownReason, "%+v", err)
	}
	if !inventRegex.MatchString("plasma/breeze!26") {
		panic("does not match")
	}
	commandlib.RegisterTag(commandlib.Tag{
		Name:     I18n("invent.kde.org"),
		Usage:    I18n("Tag merge requests on invent.kde.org"),
		Examples: `plasma/breeze!26`,
		ID:       "invent",
		Match:    inventRegex,
		Action:   invent,
	})
}

func invent(c commandlib.Context) {
	var embeds []commandlib.Embed

	for _, word := range c.Args() {
		matches := inventRegex.FindAllStringSubmatch(word, -1)
		if len(matches) == 1 && len(matches[0]) == 4 {
			grp := matches[0][1:]

			group := grp[0]
			repo := grp[1]
			number, err := strconv.ParseInt(grp[2], 10, 64)
			if err != nil {
				panic(err)
			}

			mr, _, err := glClient.MergeRequests.GetMergeRequest(group+"/"+repo, int(number), &gitlab.GetMergeRequestsOptions{})
			if err != nil {
				continue
			}

			if strings.Contains(strings.ToLower(c.Backend().Name()), "matrix") {
				mr.Description = strings.ReplaceAll(
					string(markdown.ToHTML([]byte(mr.Description), nil, nil)),
					`<img src="`,
					`<img src="`+fmt.Sprintf("https://invent.kde.org/%s/%s", group, repo),
				)
			}

			embeds = append(embeds, commandlib.Embed{
				Title: commandlib.EmbedHeader{
					Text: fmt.Sprintf("%s", mr.Title),
				},
				Body: mr.Description,
				Fields: []commandlib.EmbedField{
					{
						Title:  c.I18n("Upvotes"),
						Body:   fmt.Sprintf("%d", mr.Upvotes),
						Inline: true,
					},
					{
						Title:  c.I18n("Downvotes"),
						Body:   fmt.Sprintf("%d", mr.Downvotes),
						Inline: true,
					},
					{
						Title:  c.I18n("Status"),
						Body:   mr.MergeStatus,
						Inline: true,
					},
					{
						Title:  c.I18n("Author"),
						Body:   mr.Author.Name,
						Inline: true,
					},
				},
			})
			slice := &embeds[len(embeds)-1].Fields
			if mr.Assignee != nil {
				*slice = append(*slice, commandlib.EmbedField{
					Title:  c.I18n("Assignee"),
					Body:   mr.Assignee.Name,
					Inline: true,
				})
			}
		}
	}

	if len(embeds) > 0 {
		c.SendTags("tag-primary", embeds)
	}
}
