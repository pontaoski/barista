package web

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/appadeia/barista/barista-go"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}
func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "https://discordapp.com/api/oauth2/authorize?client_id=594613582696677416&redirect_uri=http%%3A%%2F%%2Flinuxcafe.ddns.net%%3A9000%%2Foauth&response_type=code&scope=identify%%20guilds")
}

func oauthHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "barista_session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vals := r.URL.Query()
	code := vals["code"][0]

	reqURL := fmt.Sprintf("https://discordapp.com/api/oauth2/token")
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", barista.Cfg.Section("App").Key("redir").String())

	client := &http.Client{}
	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer([]byte(form.Encode())))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.SetBasicAuth(barista.Cfg.Section("App").Key("id").String(), barista.Cfg.Section("App").Key("secret").String())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	text, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bod, _ := req.GetBody()
	reqBody, err := ioutil.ReadAll(ioutil.NopCloser(bod))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	println(formatRequest(req))
	println("Request body:")
	println(string(reqBody))
	fmt.Fprintf(w, "<br>===================================<br>")
	fmt.Fprintf(w, form.Encode())
	fmt.Fprintf(w, "<br>===================================<br>")
	fmt.Fprintf(w, string(text))

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Main : The main function of the web server.
func Main() {
	if barista.GetGlobalKey("cookie_key") == "" {
		rand.Seed(time.Now().UnixNano())
		barista.SetGlobalKey("cookie_key", randomString(32))
		store = sessions.NewCookieStore([]byte(barista.GetGlobalKey("cookie_key")))
	} else {
		store = sessions.NewCookieStore([]byte(barista.GetGlobalKey("cookie_key")))
	}
	http.HandleFunc("/", handler)
	http.HandleFunc("/oauth", oauthHandler)
	go http.ListenAndServe(":9000", nil)
	println("HTTP server now running.")
}
