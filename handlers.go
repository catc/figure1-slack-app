package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

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

	if q.token != o.VerificationToken {
		fmt.Println(q.token, o.VerificationToken)
		http.Error(res, "Tokens did not match (case integration token)", http.StatusInternalServerError)
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

	slackCaseResponse(res, &data)
}

func getCaseId(text string) (id string) {
	if len(text) == 24 {
		id = text
	} else {
		u, err := url.Parse(text)

		if err != nil || u.Host == "" {
			fmt.Println("Failed to parse url: ", text)
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

func (o *Oembed) userHandler(res http.ResponseWriter, req *http.Request) {
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

	if q.token != o.VerificationToken {
		http.Error(res, "Tokens did not match (user integration token)", http.StatusInternalServerError)
		return
	}

	// get username
	var username string
	if username = getUsername(q.text); username == "" {
		fmt.Println("Could not parse username: ", q.text)
		http.Error(res, "Could not parse username, try again", 400)
		return
	}

	// get user data
	data, err := o.getUser(username)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to get user data for user %v (parsed from %v)", username, q.text))
		http.Error(res, err.Error(), 400)
		return
	}

	slackUserResponse(res, &data)
}

func getUsername(text string) string {
	// try parsing as url
	u, err := url.Parse(text)
	if err != nil {
		return ""
	}

	if u.Host != "" {
		query := u.Query()

		uq := query.Get("username")
		if uq != "" {
			// eg: 	https://app.figure1.com/rd/publicprofile?username=penguinophile
			return uq
		} else {
			// eg: 	https://app.figure1.com/user/penguinophile
			path := strings.Split(u.EscapedPath(), "/")
			return path[len(path)-1]
		}
	} else {
		// is regular user
		return text
	}
}
