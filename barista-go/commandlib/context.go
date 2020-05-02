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
	ChoiceFlags(flags ...string) string
	AnySet(flags ...string) bool
	IsFlagSet(name string) bool
	NArgs() int
	// Flags needed by implementations
	SendMessage(id string, content interface{})
	SendTags(id string, tags []Embed)
	WrapCodeBlock(code string) string
}

type contextImpl struct {
	flagSet  flag.FlagSet
	lastUsed time.Time
	isTag    bool
	words    []string
}

func (c contextImpl) FlagValue(name string) string {
	if c.isTag {
		return ""
	}
	return c.flagSet.Lookup(name).Value.String()
}

func (c contextImpl) Arg(i int) string {
	if c.isTag {
		if len(c.words) >= i {
			return ""
		}
		return c.words[i]
	}
	return c.flagSet.Arg(i)
}

func (c contextImpl) Args() []string {
	if c.isTag {
		return c.words
	}
	return c.flagSet.Args()
}

func (c contextImpl) IsFlagSet(name string) bool {
	if c.isTag {
		return false
	}
	val := false
	c.flagSet.Visit(func(f *flag.Flag) {
		if f.Name == name {
			val = true
		}
	})
	return val
}

func (c contextImpl) NArgs() int {
	if c.isTag {
		return len(c.words)
	}
	return c.flagSet.NArg()
}

func (c contextImpl) Content() string {
	if c.isTag {
		return strings.Join(c.words, " ")
	}
	return strings.Join(c.flagSet.Args(), " ")
}

func (c contextImpl) ChoiceFlags(flags ...string) string {
	for _, flag := range flags {
		if c.IsFlagSet(flag) {
			return flag
		}
	}
	return ""
}

func (c contextImpl) AnySet(flags ...string) bool {
	for _, flag := range flags {
		if c.IsFlagSet(flag) {
			return true
		}
	}
	return false
}
