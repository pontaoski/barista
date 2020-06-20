package commandlib

import (
	"time"
)

type Context interface {
	// Flags handled by ContextMixin
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
	Type() ContextType
	RecallData(key string) (interface{}, bool)
	SetData(key string, v interface{})
	Command() Command
	// Flags needed by implementations
	I18n(message string) string
	I18nc(context, message string) string
	AuthorName() string
	AuthorIdentifier() string
	RoomIdentifier() string
	CommunityIdentifier() string
	SendMessage(id string, content interface{})
	SendTags(id string, tags []Embed)
	WrapCodeBlock(code string) string
	GenerateLink(text string, URL string) string
	NextResponse() chan string
	AwaitResponse(time time.Duration) (content string, ok bool)
}

type ContextType int

const (
	InvalidContextType ContextType = iota
	CreateCommand
	EditCommand
	DeleteCommand
)
