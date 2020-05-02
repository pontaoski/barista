package commandlib

import (
	"strings"
)

var commands []Command

func RegisterCommand(command Command) {
	commands = append(commands, command)
}

type Action func(c Context)

type Command struct {
	Name  string
	Usage string

	ID    string
	Match [][]string

	Flags  FlagList
	Action Action
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
		ctxt.flagSet = command.Flags.GetFlagSet()
		err := ctxt.flagSet.Parse(sliced)
		if err != nil {
			println(err.Error())
		}
		return command, ctxt, true
	}
	return Command{}, contextImpl{}, false
}
