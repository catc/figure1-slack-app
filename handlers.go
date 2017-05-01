package main

import (
	"fmt"
	"net/http"
)

type slashCommandRequestBody struct {
	Token     string
	ChannelID string
	Username  string
	Text      string
}

func (app *SlackApp) slashCommandHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var handler func(http.ResponseWriter, *slashCommandRequestBody)

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

	// parse form
	if err := req.ParseForm(); err != nil {
		msg := "Failed to parse body"
		(&slackError{msg, msg, err}).handleError(res)
		return
	}

	// body
	body := slashCommandRequestBody{
		Token:     req.FormValue("token"),
		ChannelID: req.FormValue("channel_id"),
		Username:  req.FormValue("user_name"),
		Text:      req.FormValue("text"),
	}

	// check request token is valid
	if body.Token != app.VerificationToken {
		msg := fmt.Sprintf("Token provided did not match (token: %v)", body.Token)
		(&slackError{"Token provided did not match", msg, nil}).handleError(res)
		return
	}

	// more basic body validation
	if body.ChannelID == "" || body.Username == "" || body.Text == "" {
		msg := fmt.Sprintf("Invalid request body (channel: %v, username: %v, text: %v)", body.Token, body.Username, body.Text)
		(&slackError{"Invalid body", msg, nil}).handleError(res)
		return
	}

	// handle request
	handler(res, &body)
}

func (app *SlackApp) handleCase(res http.ResponseWriter, body *slashCommandRequestBody) {
	// validate case id
	var id string
	if id = getCaseID(body.Text); id == "" {
		msg := fmt.Sprintf("Failed to parse url/id (text: %v)", body.Text)
		(&slackError{"Invalid case id/url, please try again", msg, nil}).handleError(res)
		return
	}

	// get case
	f1Case, err := app.getCase(id)
	if err != nil {
		msg := fmt.Sprintf("Failed retrieve case (id: %v)", id)
		(&slackError{"Failed to retrieve case", msg, err}).handleError(res)
		return
	}

	// send data to slack service to format and post message
	generateCaseContent(res, &f1Case, body.ChannelID, body.Username, app.OAuthAccessToken)

	// send 204
	res.WriteHeader(http.StatusNoContent)
}

func (app *SlackApp) handleUser(res http.ResponseWriter, body *slashCommandRequestBody) {
	// get username
	var username string
	if username = getUsername(body.Text); username == "" {
		msg := fmt.Sprintf("Failed to parse username (text: %v)", body.Text)
		(&slackError{"Invalid user id/url, please try again", msg, nil}).handleError(res)
		return
	}

	// get user data
	f1User, err := app.getUser(username)
	if err != nil {
		msg := fmt.Sprintf("Failed retrieve user data (username: %v)", username)
		(&slackError{"Failed to retrieve user", msg, err}).handleError(res)
		return
	}

	// send data to slack service to format and post message
	generateUserContent(res, &f1User, body.ChannelID, body.Username, app.OAuthAccessToken)

	// send 204
	res.WriteHeader(http.StatusNoContent)
}

func (app *SlackApp) handleCollection(res http.ResponseWriter, body *slashCommandRequestBody) {
	/*
		TODO
		- validate case id/url
		- get case
		- send 204
		- send case to slack service
			- send channel_id
			- send user_name
			- send case data
			- send token
		- slack service:
			- formats content
			- posts it
	*/
}
