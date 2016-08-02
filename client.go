package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type f1Case struct {
	Caption string
}

func getCase(id string) f1Case {
	fmt.Printf("getting case with id %v", id)

	url := "https://app.figure1.com/s/case/" + id
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	// req.Header.Add("Authorization", bearer) // TODO - replace
	req.Header.Add("Authorization", config.BearerToken) // TODO - replace

	// make the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		// TODO - deal with error
		fmt.Println(err)
	}
	defer res.Body.Close()

	var body f1Case
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		fmt.Println(err)
	}

	fmt.Println("GOT THE CASE!", body.Caption)
	/*
		TODO - handle errors
	*/
	return body
}

func getBearerToken() {
	reqBody := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		config.Email,
		config.Password,
	}
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println("Failed to marshal json, ", err)
	}

	url := "https://app.figure1.com/s/auth/login"
	req, err := http.NewRequest("POST", url, bytes.NewReader(reqJSON))
	req.Header.Add("Content-Type", "application/json")

	// make the request
	client := &http.Client{}
	res, err := client.Do(req)
	defer res.Body.Close()

	if err != nil {
		fmt.Println("Failed to connect to Figure 1 API, ", err)
	}

	// handle response
	type resp struct {
		Token string
	}
	var resBody resp
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		fmt.Println("Failed to decode response body, ", err)
	}

	config.BearerToken = resBody.Token
	fmt.Println(config.BearerToken)

	/*
		NEED:
		- incorrect credentials handling
	*/
}
