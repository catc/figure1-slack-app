package main

import (
	"fmt"
	"log"
	"net/http"

	// ...
	_ "encoding/json"
)

const TOKEN = "q51eCdHKw2uYBxVVbKrXMyJs"
const PORT = "3200"
const ADDRESS = "localhost:" + PORT

type query struct {
	token string
	text  string
}

func main() {
	fmt.Println("hello world", PORT)

	http.HandleFunc("/sla", caseHandler)

	if err := http.ListenAndServe(ADDRESS, nil); err != nil {
		log.Fatal(err)
	}
}

func caseHandler(res http.ResponseWriter, req *http.Request) {
	q := &query{}

	if err := req.ParseForm(); err != nil {
		log.Printf("Error parsing form: %s", err)
		return
	}

	q.token = req.Form.Get("token")
	q.text = req.Form.Get("text")

	if q.token != TOKEN {
		fmt.Println("Token does not match")
		return
	}

	fmt.Println("all is well", q.token)
	/*
		TODO
		- validate request
			- token
			- text
	*/
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
*/
