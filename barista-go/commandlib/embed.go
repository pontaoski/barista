package commandlib

import (
	"io"
	"strings"
)

type EmbedHeader struct {
	Icon string
	Text string
	URL  string
}

type EmbedField struct {
	Title  string
	Body   string
	Inline bool
}

//Embed ...
type Embed struct {
	Header EmbedHeader
	Title  EmbedHeader
	Footer EmbedHeader

	Fields []EmbedField
	Body   string
	Colour int
}

type File struct {
	Mimetype string
	Name     string
	Reader   io.ReadCloser
}

type EmbedList struct {
	ItemTypeName string
	Embeds       []Embed
}

type EmbedTable struct {
	Heading  string
	Subtitle string
	Headers  []string
	Data     [][]string
}

type UnionEmbed struct {
	EmbedList  EmbedList
	EmbedTable EmbedTable
}

// Constants for message embed character limits
const (
	EmbedLimitTitle       = 256
	EmbedLimitDescription = 2048
	EmbedLimitFieldValue  = 1024
	EmbedLimitFieldName   = 256
	EmbedLimitField       = 25
	EmbedLimitFooter      = 2048
	EmbedLimit            = 4000
)

func PaginateList(c Context, l []string) (*Embed, *[]Embed) {
	totalLen := 0
	for _, str := range l {
		totalLen += len(str)
		totalLen += 1
	}
	if totalLen < 1800 {
		return &Embed{Body: c.WrapCodeBlock(strings.Join(l, "\n"))}, nil
	} else {
		var strs [][]string
		var currentStr []string
		currentChunkSize := 0
		for _, str := range l {
			currentStr = append(currentStr, str)
			currentChunkSize += len(str)
			currentChunkSize += 1
			if currentChunkSize > 1800 {
				strs = append(strs, currentStr)
				currentStr = make([]string, 0)
				currentChunkSize = 0
			}
		}
		var embeds []Embed
		for _, strarr := range strs {
			embeds = append(embeds, Embed{
				Body: c.WrapCodeBlock(strings.Join(strarr, "\n")),
			})
		}
		return nil, &embeds
	}
}

// Truncate truncates any embed value over the character limit.
func (e *Embed) Truncate() {
	if len(e.Body) > EmbedLimitDescription {
		e.Body = e.Body[:EmbedLimitDescription]
	}
	for _, v := range e.Fields {
		if len(v.Title) > EmbedLimitFieldName {
			v.Title = v.Title[:EmbedLimitFieldName]
		}

		if len(v.Body) > EmbedLimitFieldValue {
			v.Body = v.Body[:EmbedLimitFieldValue]
		}
	}
	if len(e.Title.Text) > EmbedLimitTitle {
		e.Title.Text = e.Title.Text[:EmbedLimitTitle]
	}
	if len(e.Footer.Text) > EmbedLimitFooter {
		e.Footer.Text = e.Footer.Text[:EmbedLimitFooter]
	}
}
