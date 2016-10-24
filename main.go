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
	PORT    = "3200"
	ADDRESS = ":" + PORT
)

type Oembed struct {
	Email       string
	Password    string
	BearerToken string

	// integration tokens
	CaseIntegrationToken string `json:"case_integration_token"`
	UserIntegrationToken string `json:"user_integration_token"`
}

func main() {
	oembed := initOembed()
	if err := oembed.getBearerToken(); err != nil {
		log.Fatal("Failed to get bearer token, credentials are probably incorrect")
	}

	mux := http.NewServeMux()

	// add any routes
	mux.HandleFunc("/case", oembed.caseHandler)
	mux.HandleFunc("/user", oembed.userHandler)

	server := &http.Server{
		Addr:           ADDRESS,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("Figure 1 case oembed listening on " + ADDRESS)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func initOembed() Oembed {
	file, err := os.Open("conf.json")
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	o := Oembed{}
	if err = decoder.Decode(&o); err != nil {
		log.Fatal("error loading config.json", err)
	}
	return o
}
