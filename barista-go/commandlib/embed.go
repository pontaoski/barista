package commandlib

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

// Truncate truncates any embed value over the character limit.
func (e *Embed) Truncate() *Embed {
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
	return e
}
