package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func getCaseID(text string) (id string) {
	if len(text) == 24 {
		id = text
	} else {
		u, err := url.Parse(text)

		if err != nil || u.Host == "" {
			fmt.Println("Failed to parse case url: ", text)
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
	if len(text) == 24 {
		id = text
	} else {
		u, err := url.Parse(text)

		if err != nil || u.Host == "" {
			fmt.Println("Failed to parse collection url: ", text)
			return ""
		}

		// parse url, 2 types of links (page url and share link)
		query := u.Query()
		collectionid := query.Get("id")

		if len(collectionid) == 24 {
			// is share link
			id = collectionid
		} else {
			// is web app route
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

/*
	error handling
*/

type requestError struct {
	ClientResp string
	Msg        string
	Err        error
}

func (err *requestError) handleError(res http.ResponseWriter) {
	// log
	logErr("%v %v", err.Msg, err.Err)

	// send response to slack
	if err.ClientResp == "" {
		err.ClientResp = "Internal error"
	}
	res.Write([]byte(err.ClientResp))
}

type slackError struct {
	ClientResp string
	Msg        string
	Err        error
}

type slackErrorRequestBody struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text,omitempty"`
}

func (se *slackError) handleError(link string) {
	// log
	logErr("%v %v", se.Msg, se.Err)

	// post error to `response_url`
	if se.ClientResp == "" {
		se.ClientResp = "Error fetching content"
	}

	body := &slackErrorRequestBody{
		Text:         se.ClientResp,
		ResponseType: "ephemeral",
	}
	reqBody := new(bytes.Buffer)
	if err := json.NewEncoder(reqBody).Encode(body); err != nil {
		logErr("Error marshaling slack error body: %v", err)
		return
	}

	req, err := http.NewRequest("POST", link, reqBody)
	if err != nil {
		logErr("Error creating slack error request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logErr("Error making slack error request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(resp.Body)
		logErr("Slack error request was not OK (status %v): %v", resp.StatusCode, string(data))
		return
	}
}

func logErr(format string, a ...interface{}) {
	now := time.Now()
	timestamp := now.Format(time.RFC822)

	vals := append([]interface{}{timestamp}, a...)
	str := "%v: " + format
	fmt.Println(fmt.Sprintf(str, vals...))
}
