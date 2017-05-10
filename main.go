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
	port    = "3400"
	address = ":" + port
)

// SlackApp contains figure1 and slack tokens/secrets
type SlackApp struct {
	Email       string
	Password    string
	BearerToken string

	// slack tokens/secrets
	OAuthAccessToken  string `json:"oauth_access_token"`
	VerificationToken string `json:"verification_token"`
}

func main() {
	slackApp := newSlackApp()
	if err := slackApp.getBearerToken(); err != nil {
		log.Fatal("Failed to get bearer token, credentials are probably incorrect")
	}

	mux := http.NewServeMux()

	// add routes
	mux.HandleFunc("/case", slackApp.slashCommandHandler)
	mux.HandleFunc("/user", slackApp.slashCommandHandler)
	mux.HandleFunc("/collection", slackApp.slashCommandHandler)

	server := &http.Server{
		Addr:           address,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("Figure 1 slack app listening on " + address)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func newSlackApp() SlackApp {
	file, err := os.Open("conf.json")
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	app := SlackApp{}
	if err = decoder.Decode(&app); err != nil {
		log.Fatal("error loading config.json", err)
	}
	return app
}
