package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func getCaseID(text string) (id string) {
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
		}
		// eg: 	https://app.figure1.com/user/penguinophile
		path := strings.Split(u.EscapedPath(), "/")
		return path[len(path)-1]
	}
	// is regular user
	return text
}

func getCollectionID(text string) (id string) {
	// TODO - finish this + add tests
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

func truncateString(text string) string {
	split := strings.Split(text, " ")
	limit := 36
	if len(split) > limit {
		return strings.Join(split[0:limit], " ") + "..."
	}
	return text
}

type slackError struct {
	ClientResp string
	Msg        string
	Err        error
}

func (err *slackError) handleError(res http.ResponseWriter) {
	res.Write([]byte(err.ClientResp))
	fmt.Println(err.Msg, err.Err)
}