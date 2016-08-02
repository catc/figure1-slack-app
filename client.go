package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type f1Case struct {
	Id           string `json:"_id"`
	Caption      string `json:"caption"`
	IsPagingCase bool   `json:"isPagingCase"`

	// stats
	ImageViews   string `json:"imageViews"`
	Followers    int    `json:"followers"`
	CommentCount int    `json:"CommentCount"`
	VoteCount    int    `json:"voteCount"`

	// author
	Author struct {
		Username       string `json:"username"`
		TopContributor bool   `json:"topContributor"`
		Verified       bool   `json:"verified"`
	}
}

func (o *Oembed) getCase(id string) (f1Case, error) {
	var body f1Case

	url := "https://app.figure1.com/s/case/" + id
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", o.BearerToken)

	// make the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return body, errors.New("Failed to create case http request")
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		fmt.Println("Need to relog")
		if err = o.getBearerToken(); err == nil {
			return o.getCase(id)
		}

		return body, errors.New("Failed to refresh bearer token")
	}

	if res.StatusCode != http.StatusOK {
		fmt.Println("Failed to retrieve case")
		return body, errors.New("Failed to retrieve case, please try again later")
	}

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		fmt.Println("Failed to decode json, ", err)
		return body, err
	}

	return body, nil
}

func (o *Oembed) getBearerToken() error {
	reqBody := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		o.Email,
		o.Password,
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
		return errors.New("Failed to connect to Figure 1 API. Try again later.")
	}

	if res.StatusCode != http.StatusOK {
		log.Fatal("Failed to retrieve bearer token: Incorrect credentials")
	}

	// handle response
	type resp struct {
		Token string
	}
	var resBody resp
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		fmt.Println("Failed to decode response body, ", err)
		return errors.New("Failed to decode bearer response body")
	}

	o.BearerToken = resBody.Token
	return nil
}
