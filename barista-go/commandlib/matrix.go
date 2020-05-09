package commandlib

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/matrix-org/gomatrix"
)

type MatrixContext struct {
	contextImpl

	client       *gomatrix.Client
	triggerEvent *gomatrix.Event
}

type matrixMessage struct {
	Format string `json:"format"`
	HTML   string `json:"formatted_body"`
	Body   string `json:"body"`
	Type   string `json:"msgtype"`
}

func sendHTMLMessage(client *gomatrix.Client, roomID, html, text string) {
	message := matrixMessage{
		Format: "org.matrix.custom.html",
		Body:   text,
		HTML:   html,
		Type:   "m.text",
	}
	if message.Body == "" {
		message.Body = "This message can only be viewed in Matrix clients that support HTML."
	}
	_, err := client.SendMessageEvent(roomID, "m.room.message", message)
	if err != nil {
		println(err.Error())
	}
}

func sendMessage(client *gomatrix.Client, roomID, text string) {
	message := matrixMessage{
		Body: text,
		Type: "m.text",
	}
	_, err := client.SendMessageEvent(roomID, "m.room.message", message)
	if err != nil {
		println(err.Error())
	}
}

func wrapCode(code string) string {
	return "<pre><code>" + code + "</code></pre>"
}

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("barista-go/commandlib/template/*.html"))
}

func (m MatrixContext) SendTags(_ string, tags []Embed) {
	var sb strings.Builder
	for _, tag := range tags {
		tmpl.ExecuteTemplate(&sb, "single", tag)
		sb.WriteString("<br>\n")
	}
	sendHTMLMessage(m.client, m.triggerEvent.RoomID, sb.String(), "")
}

func (m MatrixContext) RoomIdentifier() string {
	return m.triggerEvent.RoomID
}

func (m MatrixContext) SendMessage(_ string, content interface{}) {
	switch content.(type) {
	case string:
		sendMessage(m.client, m.triggerEvent.RoomID, content.(string))
	case Embed:
		var sb strings.Builder
		tmpl.ExecuteTemplate(&sb, "single", content.(Embed))
		sendHTMLMessage(m.client, m.triggerEvent.RoomID, sb.String(), "")
	case EmbedList:
		var sb strings.Builder
		tmpl.ExecuteTemplate(&sb, "multiple", content.(EmbedList))
		sendHTMLMessage(m.client, m.triggerEvent.RoomID, sb.String(), "")
	case EmbedTable:
		var sb strings.Builder
		tmpl.ExecuteTemplate(&sb, "table", content.(EmbedTable))
		sendHTMLMessage(m.client, m.triggerEvent.RoomID, sb.String(), "")
	case UnionEmbed:
		m.SendMessage("", content.(UnionEmbed).EmbedTable)
		return
	}
}

func (m MatrixContext) WrapCodeBlock(code string) string {
	return fmt.Sprintf(`<pre><code>
%s
</code></pre>`, code)
}

func (m MatrixContext) GenerateLink(text, URL string) string {
	return fmt.Sprintf(`<a href="%s">%s</a>`, URL, text)
}

func MatrixMessage(client *gomatrix.Client, ev *gomatrix.Event) {
	if val, ok := ev.Content["body"]; ok {
		if cmd, contextImpl, ok := lexContent(val.(string)); ok {
			mc := MatrixContext{}
			mc.contextImpl = contextImpl
			mc.client = client
			mc.triggerEvent = ev
			go cmd.Action(&mc)
		} else {
			for _, tc := range lexTags(val.(string)) {
				mc := MatrixContext{}
				mc.contextImpl = tc.Context
				mc.client = client
				mc.triggerEvent = ev
				go tc.Tag.Action(&mc)
			}
		}
	}
}
