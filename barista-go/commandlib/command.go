package commandlib

import (
	"strings"

	iradix "github.com/hashicorp/go-immutable-radix"
)

var commandRadix = iradix.New()
var commandList []Command

func RegisterCommand(command Command) {
	for _, match := range command.Matches {
		commandRadix, _, _ = commandRadix.Insert([]byte(match), command)
	}
	commandList = append(commandList, command)
}

func Commands() []Command {
	return commandList
}

type Action func(c Context)

type Command struct {
	Name  string
	Usage string

	Examples string

	ID      string
	Matches []string

	Flags        FlagList
	Action       Action
	DeleteAction Action
	Hidden       bool
}

func LexCommand(content string) (Command, ContextMixin, bool) {
	if content == "" {
		return Command{}, ContextMixin{}, false
	}
	prefix, value, ok := commandRadix.Root().LongestPrefix([]byte(content))
	if !ok {
		return Command{}, ContextMixin{}, false
	}
	content = strings.TrimSpace(strings.TrimPrefix(content, string(prefix)))
	cmd := value.(Command)
	ctx := ContextMixin{}
	ctx.Action = cmd
	ctx.Data = make(map[string]interface{})
	ctx.FlagSet = *cmd.Flags.GetFlagSet()
	ctx.RawData = content
	ctx.FlagSet.Parse(strings.Fields(content))
	return cmd, ctx, true
}
