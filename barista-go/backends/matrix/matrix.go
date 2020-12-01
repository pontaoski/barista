package matrix

import (
	"fmt"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/log"
	"github.com/matrix-org/gomatrix"
)

type MatrixContext struct {
	commandlib.ContextMixin

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
		log.Error(err.Error())
	}
}

func sendMessage(client *gomatrix.Client, roomID, text string) {
	message := matrixMessage{
		Body: text,
		Type: "m.text",
	}
	_, err := client.SendMessageEvent(roomID, "m.room.message", message)
	if err != nil {
		log.Error(err.Error())
	}
}

func wrapCode(code string) string {
	return "<pre><code>" + code + "</code></pre>"
}

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("barista-go/commandlib/template/*.html"))
}

func (m MatrixContext) Backend() commandlib.Backend {
	return backend
}

func (m MatrixContext) SendTags(_ string, tags []commandlib.Embed) {
	var sb strings.Builder
	for _, tag := range tags {
		tmpl.ExecuteTemplate(&sb, "single", tag)
		sb.WriteString("<br>\n")
	}
	sendHTMLMessage(m.client, m.triggerEvent.RoomID, sb.String(), "")
}

func (m MatrixContext) AuthorName() string {
	resp, _ := m.client.GetDisplayName(m.triggerEvent.Sender)
	if resp == nil {
		return m.triggerEvent.Sender
	}
	return resp.DisplayName
}

func (m MatrixContext) AuthorIdentifier() string {
	return m.triggerEvent.Sender
}

func (m MatrixContext) RoomIdentifier() string {
	return m.triggerEvent.RoomID
}

func (m MatrixContext) CommunityIdentifier() string {
	return m.triggerEvent.RoomID
}

func (m MatrixContext) I18n(message string) string {
	return m.I18nInternal(commandlib.GetLanguage(&m), message)
}

func (m MatrixContext) NextResponse() (out chan string) {
	out = make(chan string)
	go func() {
		for {
			select {
			case usermsg := <-waitForMatrixMessage():
				if val, ok := usermsg.Content["body"]; ok {
					if usermsg.Sender == m.triggerEvent.Sender && usermsg.RoomID == m.triggerEvent.RoomID {
						out <- val.(string)
						return
					}
				}
			}
		}
	}()
	return out
}

func (m *MatrixContext) AwaitResponse(tenpo time.Duration) (response string, ok bool) {
	timeoutChan := make(chan struct{})
	go func() {
		time.Sleep(tenpo * 3)
		timeoutChan <- struct{}{}
	}()
	for {
		select {
		case msg := <-m.NextResponse():
			return msg, true
		case <-timeoutChan:
			return "", false
		}
	}
}

func (m MatrixContext) I18nc(context, message string) string {
	return m.I18n(message)
}

func (m MatrixContext) SendMessage(_ string, content interface{}) {
	switch content.(type) {
	case string:
		sendMessage(m.client, m.triggerEvent.RoomID, content.(string))
	case commandlib.Embed:
		var sb strings.Builder
		tmpl.ExecuteTemplate(&sb, "single", content.(commandlib.Embed))
		sendHTMLMessage(m.client, m.triggerEvent.RoomID, sb.String(), "")
	case commandlib.EmbedList:
		var sb strings.Builder
		tmpl.ExecuteTemplate(&sb, "multiple", content.(commandlib.EmbedList))
		sendHTMLMessage(m.client, m.triggerEvent.RoomID, sb.String(), "")
	case commandlib.EmbedTable:
		var sb strings.Builder
		tmpl.ExecuteTemplate(&sb, "table", content.(commandlib.EmbedTable))
		sendHTMLMessage(m.client, m.triggerEvent.RoomID, sb.String(), "")
	case commandlib.UnionEmbed:
		m.SendMessage("", content.(commandlib.UnionEmbed).EmbedTable)
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

var handlers []*eventHandlerInstance
var handlerMutex = sync.Mutex{}

type eventHandlerInstance struct {
	handler func(ev *gomatrix.Event)
}

func removeMatrixHandler(ehi *eventHandlerInstance) {
	handlerMutex.Lock()
	defer handlerMutex.Unlock()
	for idx, handler := range handlers {
		if handler == ehi {
			handlers = append(handlers[:idx], handlers[idx+1:]...)
		}
	}
}

func addMatrixHandlerOnce(input func(ev *gomatrix.Event)) {
	handlerMutex.Lock()
	defer handlerMutex.Unlock()
	ehi := eventHandlerInstance{input}
	handlers = append(handlers, &ehi)
}

func waitForMatrixMessage() chan *gomatrix.Event {
	channel := make(chan *gomatrix.Event)
	addMatrixHandlerOnce(func(ev *gomatrix.Event) {
		channel <- ev
	})
	return channel
}

func MatrixMessage(client *gomatrix.Client, ev *gomatrix.Event) {
	for _, handler := range handlers {
		handler.handler(ev)
		removeMatrixHandler(handler)
	}
	if val, ok := ev.Content["body"]; ok {
		if cmd, ContextMixin, ok := commandlib.LexCommand(val.(string)); ok {
			mc := MatrixContext{}
			mc.ContextMixin = ContextMixin
			mc.ContextMixin.ContextType = commandlib.CreateCommand
			mc.client = client
			mc.triggerEvent = ev
			go log.CanPanic(func() {
				cmd.Action(&mc)
			})
		} else {
			for _, tc := range commandlib.LexTags(val.(string)) {
				mc := MatrixContext{}
				mc.ContextMixin = tc.Context
				mc.client = client
				mc.triggerEvent = ev
				go log.CanPanic(func() {
					tc.Tag.Action(&mc)
				})
			}
		}
	}
}
