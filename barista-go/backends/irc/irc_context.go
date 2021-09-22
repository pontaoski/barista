package irc

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/lrstanley/girc"
)

type IRCContext struct {
	commandlib.ContextMixin
	tm string
	cl *girc.Client
	ev girc.Event
}

func (i IRCContext) AuthorIdentifier() string {
	return i.ev.Source.String()
}

func (i IRCContext) AuthorName() string {
	return i.ev.Source.String()
}

func (i IRCContext) AwaitResponse(tenpo time.Duration) (content string, ok bool) {
	timeoutChan := make(chan struct{})
	go func() {
		time.Sleep(tenpo)
		timeoutChan <- struct{}{}
	}()
	for {
		select {
		case msg := <-NextMessage():
			if msg.Source.Equals(i.ev.Source) {
				return msg.Last(), true
			}
		case <-timeoutChan:
			return "", false
		}
	}
}

func (i IRCContext) Backend() commandlib.Backend {
	return backend
}

func (i IRCContext) CommunityIdentifier() string {
	return i.ev.Params[0]
}

func (i IRCContext) GenerateLink(text string, URL string) string {
	return URL
}

func (i IRCContext) I18n(message string) string {
	return message
}

func (i IRCContext) I18nc(context, message string) string {
	return message
}

func (i IRCContext) NextResponse() chan string {
	retChan := make(chan string)
	go func() {
		for {
			select {
			case msg := <-NextMessage():
				if msg.Source.Equals(i.ev.Source) {
					retChan <- msg.Last()
					return
				}
			}
		}
	}()
	return retChan
}

func (i IRCContext) RoomIdentifier() string {
	return i.ev.Params[0]
}

func uploadFile(file io.Reader) (string, error) {
	const (
		url = "https://0x0.st"
	)
	var (
		err    error
		client http.Client
		b      bytes.Buffer
	)

	values := map[string]io.Reader{
		"file": file,
	}

	writer := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = writer.CreateFormFile(key, x.Name()); err != nil {
				return "", err
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return "", err
		}

	}
	writer.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return strings.Replace(bodyString, "\n", "", -1), nil
	}
	return "", fmt.Errorf("bad status: %s", resp.Status)
}

func (i IRCContext) SendMessage(_ string, content interface{}) {
	switch a := content.(type) {
	case string:
		i.cl.Cmd.ReplyTo(i.ev, content.(string))
	case commandlib.Embed:
		msg := ircEmbed(content.(commandlib.Embed))
		for _, str := range msg {
			i.cl.Cmd.ReplyTo(i.ev, str)
		}
	case commandlib.EmbedList:
		for _, page := range content.(commandlib.EmbedList).Embeds {
			msg := ircEmbed(page)
			for _, str := range msg {
				i.cl.Cmd.ReplyTo(i.ev, str)
			}
		}
	case commandlib.UnionEmbed:
		i.SendMessage("", content.(commandlib.UnionEmbed).EmbedList)
	case commandlib.File:
		fi, err := uploadFile(a.Reader)
		if err != nil {
			i.cl.Cmd.ReplyTo(i.ev, "failed to upload a file")
		}
		i.cl.Cmd.ReplyTo(i.ev, fi)
	}
}

func (i IRCContext) SendTags(_ string, tags []commandlib.Embed) {
	for _, page := range tags {
		msg := ircEmbed(page)
		for _, str := range msg {
			i.cl.Cmd.ReplyTo(i.ev, str)
		}
	}
}

func (i IRCContext) WrapCodeBlock(code string) string {
	return "[Barista IRC does not support codeblocks. Please use Barista from another service to view codeblocks.]"
}
