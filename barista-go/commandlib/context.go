package commandlib

import (
	"fmt"
	"strings"
	"time"

	"github.com/leonelquinteros/gotext"
	flag "github.com/spf13/pflag"
)

type Context interface {
	// Flags handled by contextImpl
	FlagValue(name string) string
	Arg(i int) string
	Args() []string
	Content() string
	RawContent() string
	ChoiceFlags(flags ...string) string
	AnySet(flags ...string) bool
	IsFlagSet(name string) bool
	NArgs() int
	Usage() string
	I18nInternal(locale, message string) string
	// Flags needed by implementations
	I18n(message string) string
	I18nc(context, message string) string
	RoomIdentifier() string
	SendMessage(id string, content interface{})
	SendTags(id string, tags []Embed)
	WrapCodeBlock(code string) string
	GenerateLink(text string, URL string) string
}

type contextImpl struct {
	flagSet    flag.FlagSet
	lastUsed   time.Time
	isTag      bool
	words      []string
	rawContent string
}

var pofiles map[string]*gotext.Po = make(map[string]*gotext.Po)

func grabPo(locale string) *gotext.Po {
	if val, ok := pofiles[locale]; ok {
		return val
	}
	po := new(gotext.Po)
	po.ParseFile(fmt.Sprintf("messages/barista_%s.po", locale))
	pofiles[locale] = po
	return po
}

func DropPo(locale string) {
	delete(pofiles, locale)
}

func (c contextImpl) RawContent() string {
	return c.rawContent
}

func (c contextImpl) I18nInternal(locale, message string) string {
	return grabPo(locale).Get(message)
}

func (c contextImpl) Usage() string {
	return c.flagSet.FlagUsages()
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
