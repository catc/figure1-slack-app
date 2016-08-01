package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var bearer string = "Bearer 52f85c59d2a5aa387eb5f1b83f6eae941c0c60709d4b57363d9aa966052e3570"

type f1Case struct {
	Caption string
}

func getCase(id string) {
	fmt.Printf("getting case with id %v", id)

	url := "https://app.figure1.com/s/case/" + id
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", bearer) // TODO - replace

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
}

func getBearerToken(id string) {

}
