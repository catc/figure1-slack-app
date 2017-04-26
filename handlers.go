package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type slashCommandRequestBody struct {
	Token     string `json:"token"`
	ChannelID string `json:"channel_id"`
	Username  string `json:"user_name"`
	Text      string `json:"text"`
}

func (sa *SlackApp) slashCommandHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var handler func(http.ResponseWriter, *slashCommandRequestBody)

	// check if handler exists for path
	switch req.URL.Path {
	case "/case":
		handler = handleCase
	case "/user":
		handler = handleUser
	case "/collection":
		handler = handleCollection
	default:
		http.Error(res, "Not found", http.StatusNotFound)
		return
	}

	// decode body
	var body slashCommandRequestBody
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		fmt.Println("Failed to decode slash command json, ", err)
		http.Error(res, "Failed to parse body", http.StatusBadRequest)
		return
	}

	// check request token is valid
	if body.Token != sa.VerificationToken {
		http.Error(res, "Tokens did not match", http.StatusInternalServerError)
		return
	}

	// more basic body validation
	if body.ChannelID == "" || body.Username == "" || body.Text == "" {
		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	// handle request
	handler(res, &body)
}

func handleCase(res http.ResponseWriter, body *slashCommandRequestBody) {
	fmt.Println("handling case")
}

func handleUser(res http.ResponseWriter, body *slashCommandRequestBody) {

}

func handleCollection(res http.ResponseWriter, body *slashCommandRequestBody) {

}
