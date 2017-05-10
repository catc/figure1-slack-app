package main

import (
	"fmt"
	"net/http"
)

type slashCommandRequestBody struct {
	Token       string
	ChannelID   string
	Username    string
	Text        string
	ResponseURL string
}

func (app *SlackApp) slashCommandHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var handler func(*slashCommandRequestBody)

	// check if handler exists for path
	switch req.URL.Path {
	case "/case":
		handler = app.handleCase
	case "/user":
		handler = app.handleUser
	case "/collection":
		handler = app.handleCollection
	default:
		http.Error(res, "Not found", http.StatusNotFound)
		return
	}

	logErr("Request made to '%v'", req.URL.Path)

	// parse form
	if err := req.ParseForm(); err != nil {
		msg := "Failed to parse body"
		(&requestError{msg, msg, err}).handleError(res)
		return
	}

	// body
	body := slashCommandRequestBody{
		Token:       req.FormValue("token"),
		ChannelID:   req.FormValue("channel_id"),
		Username:    req.FormValue("user_name"),
		Text:        req.FormValue("text"),
		ResponseURL: req.FormValue("response_url"),
	}

	// check request token is valid
	if body.Token != app.VerificationToken {
		msg := fmt.Sprintf("Token provided did not match (token: %v)", body.Token)
		(&requestError{"Token provided did not match", msg, nil}).handleError(res)
		return
	}

	// more basic body validation
	if body.ChannelID == "" || body.Username == "" || body.Text == "" {
		msg := fmt.Sprintf("Invalid request body (channel: %v, username: %v, text: %v)", body.Token, body.Username, body.Text)
		(&requestError{"Invalid body", msg, nil}).handleError(res)
		return
	}

	// assume everything is fine, any further errors will be sent via the `response_url`
	res.Write([]byte("Fetching content..."))

	// spawn slack response
	go func() {
		handler(&body)
	}()
}

func (app *SlackApp) handleCase(body *slashCommandRequestBody) {
	// validate case id
	var id string
	if id = getCaseID(body.Text); id == "" {
		msg := fmt.Sprintf("Failed to parse case url/id (text: %v)", body.Text)
		(&slackError{"Invalid case id/url, please try again", msg, nil}).handleError(body.ResponseURL)
		return
	}

	// get case
	f1Case, err := app.getCase(id)
	if err != nil {
		msg := fmt.Sprintf("Failed retrieve case (id: %v)", id)
		(&slackError{"Failed to retrieve case", msg, err}).handleError(body.ResponseURL)
		return
	}

	// generate content
	attachments := generateCaseContent(&f1Case)

	// respond
	respondToSlashCommand(body.ResponseURL, attachments)
}

func (app *SlackApp) handleUser(body *slashCommandRequestBody) {
	// get username
	var username string
	if username = getUsername(body.Text); username == "" {
		msg := fmt.Sprintf("Failed to parse username (text: %v)", body.Text)
		(&slackError{"Invalid user id/url, please try again", msg, nil}).handleError(body.ResponseURL)
		return
	}

	// get user data
	f1User, err := app.getUser(username)
	if err != nil {
		msg := fmt.Sprintf("Failed retrieve user data (username: %v)", username)
		(&slackError{"Failed to retrieve user", msg, err}).handleError(body.ResponseURL)
		return
	}

	// generate content
	attachments := generateUserContent(&f1User)

	// respond
	respondToSlashCommand(body.ResponseURL, attachments)
}

func (app *SlackApp) handleCollection(body *slashCommandRequestBody) {
	// get id
	var id string
	if id = getCollectionID(body.Text); id == "" {
		msg := fmt.Sprintf("Failed to parse collection url/id (text: %v)", body.Text)
		(&slackError{"Invalid collection id/url, please try again", msg, nil}).handleError(body.ResponseURL)
		return
	}

	// get user data
	f1Collection, err := app.getCollection(id)
	if err != nil {
		msg := fmt.Sprintf("Failed retrieve collection (id: %v)", id)
		(&slackError{"Failed to retrieve collection", msg, err}).handleError(body.ResponseURL)
		return
	}

	// generate content
	attachments := generateCollectionContent(&f1Collection)

	// respond
	respondToSlashCommand(body.ResponseURL, attachments)
}
