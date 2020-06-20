package commandlib

import (
	"regexp"
	"strings"
)

var tags []Tag

func RegisterTag(tag Tag) {
	tags = append(tags, tag)
}

func Tags() []Tag {
	return tags
}

type TagSample struct {
	Tag  string
	Desc string
}

type Tag struct {
	Name  string
	Usage string

	Examples string
	Samples  []TagSample

	ID string

	Match  *regexp.Regexp
	Action Action
}

type TagContext struct {
	Tag     Tag
	Context ContextMixin
}

func LexTags(content string) []TagContext {
	var ret []TagContext
	for _, tag := range tags {
		if tag.Match.Match([]byte(content)) {
			ret = append(ret, TagContext{
				Tag: tag,
				Context: ContextMixin{
					Words: strings.Fields(content),
					IsTag: true,
				},
			})
		}
	}
	return ret
}
