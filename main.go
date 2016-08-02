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
const ADDRESS = "localhost:" + PORT

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
		return
	}
	q.token = req.Form.Get("token")
	q.text = req.Form.Get("text")

	if q.token != o.Token {
		fmt.Println("Token does not match")
		return
	}

	// get case id
	var id string
	if id = getCaseId(q.text); id == "" {
		fmt.Println("Could not find the case id, try again")
		// TODO - send back invalid response
	}

	// get case
	data, err := o.getCase(id)
	if err != nil {
		// TODO - send back invalid response
		fmt.Println("ERROR TRYING TO GET CASE", err)
		return
	}

	formatSlackResponse(&data)

	// TODO - send back valid response
}

func getCaseId(text string) (id string) {
	if len(text) == 24 {
		id = text
	} else {
		u, err := url.Parse(text)

		if err != nil || u.Host == "" {
			fmt.Println("Could not parse url")
		}

		// parse url, 2 types of links (page url and share link)
		query := u.Query()
		idAttempt := query.Get("imageid")
		if idAttempt != "" && len(idAttempt) == 24 {
			// is share link
			id = idAttempt
		} else {
			// is web app url
			path := strings.Split(u.EscapedPath(), "/")
			idAttempt = path[len(path)-1]
			if len(idAttempt) == 24 {
				id = idAttempt
			}
		}
	}
	return id
}

/*
	TODO
	- need error handler
	- wildcard route
	- account for > 3 sec delay (slack timeout?)
	- error handler for response
		- if:
			- could not find case
			- imageid is incorrect
			- cant parse link
			- etc
	- add credentials/config file
		- to store passwords for fig1
		- to store slack token
*/
