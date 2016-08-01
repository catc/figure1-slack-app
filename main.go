package main

import (
	"fmt"
	"log"
	"net/http"

	// ...
	"encoding/json"
	"net/url"
	"os"
	"strings"
)

const PORT = "3200"
const ADDRESS = "localhost:" + PORT

var config configuration

type configuration struct {
	Username string
	Password string
	Token    string
}

func main() {
	setupConf()

	http.HandleFunc("/case", caseHandler)

	fmt.Println("Figure 1 case oembed listening on " + ADDRESS)
	if err := http.ListenAndServe(ADDRESS, nil); err != nil {
		log.Fatal(err)
	}
}

func setupConf() {
	file, err := os.Open("conf.json")
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	c := configuration{}
	if err = decoder.Decode(&c); err != nil {
		log.Fatal("error loading config.json", err)
	}
	config = c
}

func caseHandler(res http.ResponseWriter, req *http.Request) {
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

	if q.token != config.Token {
		fmt.Println("Token does not match")
		return
	}

	// get case id
	var id string
	if id = getCaseId(q.text); id == "" {
		fmt.Println("Could not find the case id, try again")
	} else {
		fmt.Println("id is ", id)
	}

	// get case
	getCase(id)
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
