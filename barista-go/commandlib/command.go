package commandlib

import (
	"regexp"
	"strings"
)

var commands []Command
var tags []Tag

func RegisterCommand(command Command) {
	commands = append(commands, command)
}

func Commands() []Command {
	return commands
}

func RegisterTag(tag Tag) {
	tags = append(tags, tag)
}

func Tags() []Tag {
	return tags
}

type Action func(c Context)

type Command struct {
	Name  string
	Usage string

	Examples string

	ID    string
	Match [][]string

	Flags  FlagList
	Action Action
	Hidden bool
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

type tagContext struct {
	Tag     Tag
	Context contextImpl
}

func lexTags(content string) []tagContext {
	var ret []tagContext
	for _, tag := range tags {
		if tag.Match.Match([]byte(content)) {
			ret = append(ret, tagContext{
				Tag: tag,
				Context: contextImpl{
					words: strings.Fields(content),
					isTag: true,
				},
			})
		}
	}
	return ret
}

func lexContent(content string) (Command, contextImpl, bool) {
	if content == "" {
		return Command{}, contextImpl{}, false
	}
	for _, command := range commands {
		var sliced []string
		fields := strings.Fields(content)
		fieldsLen := len(fields)
	outerLoop:
		for _, match := range command.Match {
			matchLen := len(match)
			if fieldsLen < matchLen {
				continue
			}
			for i := 0; i < matchLen; i++ {
				if fields[i] != match[i] {
					continue outerLoop
				}
			}
			sliced = fields[matchLen:]
			goto matched
		}
		continue
	matched:
		ctxt := contextImpl{}
		ctxt.flagSet = *command.Flags.GetFlagSet()
		err := ctxt.flagSet.Parse(sliced)
		if err != nil {
			println(err.Error())
		}
		return command, ctxt, true
	}
	return Command{}, contextImpl{}, false
}
