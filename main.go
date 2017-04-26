package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	PORT    = "3400"
	ADDRESS = ":" + PORT
)

type SlackApp struct {
	Email       string
	Password    string
	BearerToken string

	// integration tokens
	OAuthAccessToken  string `json:"oauth_access_token"`
	VerificationToken string `json:"verification_token"`
}

func main() {
	slackApp := NewSlackApp()
	if err := slackApp.getBearerToken(); err != nil {
		log.Fatal("Failed to get bearer token, credentials are probably incorrect")
	}

	mux := http.NewServeMux()

	// add any routes
	mux.HandleFunc("/case", slackApp.slashCommandHandler)
	mux.HandleFunc("/user", slackApp.slashCommandHandler)
	mux.HandleFunc("/collection", slackApp.slashCommandHandler)

	server := &http.Server{
		Addr:           ADDRESS,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("Figure 1 slack app listening on " + ADDRESS)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func NewSlackApp() SlackApp {
	file, err := os.Open("conf.json")
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	sa := SlackApp{}
	if err = decoder.Decode(&sa); err != nil {
		log.Fatal("error loading config.json", err)
	}
	return sa
}
