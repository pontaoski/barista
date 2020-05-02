package commandlib

import (
	"strings"
	"time"

	flag "github.com/spf13/pflag"
)

type Context interface {
	// Flags handled by contextImpl
	FlagValue(name string) string
	Arg(i int) string
	Args() []string
	Content() string
	IsFlagSet(name string) bool
	NArgs() int
	// Flags needed by implementations
	SendMessage(id string, content interface{})
}

type contextImpl struct {
	flagSet  flag.FlagSet
	lastUsed time.Time
}

func (c contextImpl) FlagValue(name string) string {
	return c.flagSet.Lookup(name).Value.String()
}

func (c contextImpl) Arg(i int) string {
	return c.flagSet.Arg(i)
}

func (c contextImpl) Args() []string {
	return c.flagSet.Args()
}

func (c contextImpl) IsFlagSet(name string) bool {
	val := false
	c.flagSet.Visit(func(f *flag.Flag) {
		if f.Name == name {
			val = true
		}
	})
	return val
}

func (c contextImpl) NArgs() int {
	return c.flagSet.NArg()
}

func (c contextImpl) Content() string {
	return strings.Join(c.flagSet.Args(), " ")
}
