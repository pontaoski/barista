package barista

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/appadeia/barista/barista-go/commandlib"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func init() {
	commandlib.RegisterCommand(commandlib.Command{
		Name:     I18n("sitelen pona"),
		Usage:    I18n("Write with sitelen pona."),
		Matches:  []string{"ilo o sitelen pona", ",sp"},
		Examples: `ilo o sitelen pona mu`,
		ID:       "sitelenpona",
		Action: func(c commandlib.Context) {
			filename := "/tmp/" + randSeq(10) + ".png"
			cmd := exec.Command("pango-view", "--no-display", "-t", c.Content(), "--font", "linja sike 50", "-o", filename, "--align=center", "--hinting=full", "--margin=10px")

			var b bytes.Buffer
			cmd.Stdout = &b
			cmd.Stderr = &b

			if err := cmd.Run(); err != nil {
				c.SendMessage("main", commandlib.ErrorEmbed("ilo li pakala a! \n"+b.String()))
				return
			}

			mu, err := os.Open(filename)
			if err != nil {
				c.SendMessage("main", commandlib.ErrorEmbed("ilo li pakala a! \n"+err.Error()))
				return
			}

			url, err := uploadFile(mu)
			if err != nil {
				c.SendMessage("main", commandlib.ErrorEmbed("ilo li pakala a! \n"+err.Error()))
				return
			}

			c.SendMessage("main", url)
		},
	})
}

func uploadFile(file *os.File) (string, error) {
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
