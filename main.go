package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const PORT = "3200"
const ADDRESS = ":" + PORT

type Oembed struct {
	Email       string
	Password    string
	Token       string
	BearerToken string
}

func main() {
	oembed := initOembed()
	_ = oembed.getBearerToken()

	http.HandleFunc("/case", oembed.caseHandler)

	fmt.Println("Figure 1 case oembed listening on " + ADDRESS)
	if err := http.ListenAndServe(ADDRESS, nil); err != nil {
		log.Fatal(err)
	}
}

func initOembed() Oembed {
	file, err := os.Open("conf.json")
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	o := Oembed{}
	if err = decoder.Decode(&o); err != nil {
		log.Fatal("error loading config.json", err)
	}
	return o
}

func (o *Oembed) caseHandler(res http.ResponseWriter, req *http.Request) {
	type query struct {
		token string
		text  string
	}
	q := &query{}

	if err := req.ParseForm(); err != nil {
		log.Printf("Error parsing form: %s", err)
		http.Error(res, "Internal error", http.StatusInternalServerError)
		return
	}
	q.token = req.Form.Get("token")
	q.text = req.Form.Get("text")

	if q.token != o.Token {
		http.Error(res, "Tokens did not match", http.StatusInternalServerError)
		return
	}

	// get case id
	var id string
	if id = getCaseId(q.text); id == "" {
		fmt.Println("Could not parse case id: ", q.text)
		http.Error(res, "Could not find the case id, try again", 400)
		return
	}

	// get case
	data, err := o.getCase(id)
	if err != nil {
		http.Error(res, err.Error(), 400)
		return
	}

	slackResponse(res, &data)
}

func getCaseId(text string) (id string) {
	if len(text) == 24 {
		id = text
	} else {
		u, err := url.Parse(text)

		if err != nil || u.Host == "" {
			fmt.Println("Failed to parse url: ", q.text)
			return ""
		}

		// parse url, 2 types of links (page url and share link)
		query := u.Query()
		imageid := query.Get("imageid")
		image := query.Get("image")

		switch {
		case len(imageid) == 24:
			// is modal link
			id = imageid
		case len(image) == 24:
			// is web app url
			id = image
		default:
			// is share link
			path := strings.Split(u.EscapedPath(), "/")
			idAttempt := path[len(path)-1]
			if len(idAttempt) == 24 {
				id = idAttempt
			}
		}
	}

	return id
}
