package barista

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/appadeia/barista/barista-go/commandlib"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

var sheet *sheets.Service

func init() {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	sheet, err = sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
}

func init() {
	commandlib.RegisterCommand(commandlib.Command{
		Matches: []string{"o lanpan e nimi ale pona nanpa tu"},
		Hidden:  true,
		Action:  HiddenNAP,
	})
}

const sheetID = "1t-pjAgZDyKPXcCRnEdATFQOxGbQFMjZm-8EvXiQd2Po"

var words []Word

type Word struct {
	Names          []string
	Category       string
	Definition     string
	SourceLanguage string
	Etymology      string
	Tag            string
}

func init() {
	err := lanpan()
	if err != nil {
		panic(err)
	}
}

func lanpan() error {
	resp, err := sheet.Spreadsheets.Values.Get(sheetID, "nimi (inli)!A1:F1000").Do()
	if err != nil {
		return err
	}

	words = []Word{
		{
			Names:          []string{"akesun"},
			Category:       "post-pu",
			Definition:     "The toki pona alphabet",
			SourceLanguage: "toki pona",
			Etymology:      "akesi + esun",
		},
		{
			Names:          []string{"kese"},
			Category:       "post-pu",
			Definition:     "queer; LGBT+",
			SourceLanguage: "Hebrew",
			Etymology:      "késhet 'rainbow' (also name of Jewish LGBT org)",
		},
		{
			Names:          []string{"tajan"},
			Category:       "post-pu",
			Definition:     "triangular, spikey, triplet",
			SourceLanguage: "English",
			Etymology:      "kinda sounds like triangle",
		},
		{
			Names:          []string{"lipamanka"},
			Category:       "post-pu",
			Definition:     "jewish, jan olin pi jan opasan",
			SourceLanguage: "jew",
			Etymology:      "lipamanka is jew",
		},
		{
			Names:          []string{"opasan"},
			Category:       "post-pu",
			Definition:     "jan olin pi jan lipamanka",
			SourceLanguage: "pona",
			Etymology:      "opasan is pona",
		},
	}

	for _, row := range resp.Values {
		word := Word{
			Names:          strings.Split(row[0].(string), ","),
			Category:       row[1].(string),
			Definition:     row[2].(string),
			SourceLanguage: row[3].(string),
		}
		switch {
		case len(row) >= 6:
			word.Etymology = row[5].(string)
		case len(row) >= 5:
			word.Etymology = row[4].(string)
		}
		words = append(words, word)
	}

	return nil
}

func HiddenNAP(c commandlib.Context) {
	if !c.Backend().IsBotOwner(c) {
		c.SendMessage("primary", commandlib.ErrorEmbed("You are not the bot owner."))
	}

	err := lanpan()
	if err != nil {
		c.SendMessage("primary", commandlib.ErrorEmbed("There was an error: "+err.Error()))
	}

	c.SendMessage("primary", "mi lanpan e nimi ale pona nanpa tu!")
}
