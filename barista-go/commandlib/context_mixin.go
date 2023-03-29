package commandlib

import (
	"fmt"
	"strings"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/leonelquinteros/gotext"
)

type ContextMixin struct {
	FlagSet      flag.FlagSet
	LastUsed     time.Time
	IsTag        bool
	Words        []string
	RawData      string
	ContextType  ContextType
	Action       Command
	CacheHintTTL int
	CacheHintSet bool
	Data         map[string]interface{}
}

var pofiles map[string]*gotext.Po = make(map[string]*gotext.Po)

func (c *ContextMixin) CacheHint(ttl int) {
	c.CacheHintTTL = ttl
	c.CacheHintSet = true
}

func (c *ContextMixin) GetTTL() int {
	if c.CacheHintSet {
		return c.CacheHintTTL
	}
	return 300
}

func (c *ContextMixin) Command() Command {
	return c.Action
}

func (c *ContextMixin) SetData(key string, v interface{}) {
	c.Data[key] = v
}

func (c *ContextMixin) RecallData(key string) (val interface{}, ok bool) {
	val, ok = c.Data[key]
	return
}

func (c ContextMixin) Type() ContextType {
	return c.ContextType
}

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

func (c ContextMixin) RawContent() string {
	return c.RawData
}

func (c ContextMixin) I18nInternal(locale, message string) string {
	return grabPo(locale).Get(message)
}

func (c ContextMixin) Usage() string {
	return c.FlagSet.FlagUsages()
}

func (c ContextMixin) FlagValue(name string) string {
	if c.IsTag {
		return ""
	}
	return c.FlagSet.Lookup(name).Value.String()
}

func (c ContextMixin) Arg(i int) string {
	if c.IsTag {
		if len(c.Words) >= i {
			return ""
		}
		return c.Words[i]
	}
	return c.FlagSet.Arg(i)
}

func (c ContextMixin) Args() []string {
	if c.IsTag {
		return c.Words
	}
	return c.FlagSet.Args()
}

func (c ContextMixin) IsFlagSet(name string) bool {
	if c.IsTag {
		return false
	}
	val := false
	c.FlagSet.Visit(func(f *flag.Flag) {
		if f.Name == name {
			val = true
		}
	})
	return val
}

func (c ContextMixin) NArgs() int {
	if c.IsTag {
		return len(c.Words)
	}
	return c.FlagSet.NArg()
}

func (c ContextMixin) Content() string {
	if c.IsTag {
		return strings.Join(c.Words, " ")
	}
	return strings.Join(c.FlagSet.Args(), " ")
}

func (c ContextMixin) ChoiceFlags(flags ...string) string {
	for _, flag := range flags {
		if c.IsFlagSet(flag) {
			return flag
		}
	}
	return ""
}

func (c ContextMixin) AnySet(flags ...string) bool {
	for _, flag := range flags {
		if c.IsFlagSet(flag) {
			return true
		}
	}
	return false
}
