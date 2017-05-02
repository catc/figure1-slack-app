package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type f1Case struct {
	ID           string `json:"_id"`
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

type f1User struct {
	ID             string `json:"_id"`
	Username       string `json:"username"`
	Verified       bool   `json:"verified"`
	TopContributor bool   `json:"topContributor"`
	Category       string
	Specialty      string

	// extra info
	Country        string `json:"country"`
	ProfileCountry string `json:"profileCountry"`
	FullName       string `json:"fullName"`
	Institution    string `json:"institution"`
	Bio            string `json:"bio"`
	Link           string `json:"link"`

	// specialty object
	SpecialtyObject struct {
		Category struct {
			Strings struct {
				Label string `json:"label"`
			} `json:"strings"`
		} `json:"category"`
		Strings struct {
			Label string `json:"label"`
		} `json:"strings"`
	} `json:"specialtyObject"`

	// stats
	CommentsCount  int `json:"profileCommentsCount"`
	FavoritesCount int `json:"profileFavoritesCount"`
	FollowersCount int `json:"profileFollowersCount"`
	FollowingCount int `json:"profileFollowingCount"`
	UploadsCount   int `json:"profileUploadsCount"`
}

type f1Collection struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ID          string `json:"id"`
	Size        int    `json:"size"`
	Embedded    struct {
		Items []struct {
			ID           string `json:"_id"`
			Caption      string `json:"caption"`
			Title        string `json:"title"`
			ContentType  int    `json:"contentType"`
			CommentCount int    `json:"commentCount"`
			Followers    int    `json:"followers"`
			VoteCount    int    `json:"voteCount"`
			Links        struct {
				Image struct {
					Href string `json:"href"`
				} `json:"image"`
			} `json:"_links"`
		} `json:"items"`

		Authors []struct {
			Username          string `json:"username"`
			ID                string `json:"_id"`
			Verified          bool   `json:"verified"`
			SpecialtyName     string `json:"specialtyName"`
			SpecialtyCategory string `json:"specialtyCategory"`
			TopContributor    bool   `json:"topContributor"`
		} `json:"authors"`
	} `json:"_embedded"`
}

type f1Response interface {
	decode(io.Reader) error
}

func (f *f1Case) decode(body io.Reader) error {
	if err := json.NewDecoder(body).Decode(&f); err != nil {
		return err
	}
	return nil
}
func (f *f1User) decode(body io.Reader) error {
	if err := json.NewDecoder(body).Decode(&f); err != nil {
		return err
	}
	return nil
}
func (f *f1Collection) decode(body io.Reader) error {
	if err := json.NewDecoder(body).Decode(&f); err != nil {
		return err
	}
	return nil
}

func (app *SlackApp) fig1Request(url string, marsh f1Response) error {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", app.BearerToken)

	// make the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return errors.New("Failed to create http request")
	}
	defer res.Body.Close()

	// check if request is authorized
	if res.StatusCode == http.StatusUnauthorized {
		fmt.Println("Need to relog")
		if err = app.getBearerToken(); err != nil {
			return errors.New("Failed to refresh auth token, please try again.")
		}
		return app.fig1Request(url, marsh)
	}

	if res.StatusCode != http.StatusOK {
		fmt.Println("Failed to retrieve case", res.Status)
		return errors.New("Failed to retrieve case, please try again later")
	}

	if err := marsh.decode(res.Body); err != nil {
		return err
	}

	return nil
}

func (app *SlackApp) getCase(id string) (f1Case, error) {
	var body f1Case
	url := "https://app.figure1.com/s/case/" + id

	err := app.fig1Request(url, &body)
	if err != nil {
		return body, err
	}

	return body, nil
}

func (app *SlackApp) getUser(username string) (f1User, error) {
	var body f1User
	url := "https://app.figure1.com/s/profile/public/" + username

	err := app.fig1Request(url, &body)
	if err != nil {
		return body, err
	}

	body.Category = body.SpecialtyObject.Category.Strings.Label
	body.Specialty = body.SpecialtyObject.Strings.Label

	return body, nil
}

func (app *SlackApp) getCollection(id string) (f1Collection, error) {
	var body f1Collection
	url := "https://api.figure1.com/collections/" + id

	err := app.fig1Request(url, &body)
	if err != nil {
		return body, err
	}

	return body, nil
}

func (app *SlackApp) getBearerToken() error {
	reqBody := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		app.Email,
		app.Password,
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
	if err != nil {
		fmt.Println("Failed to connect to Figure 1 API, ", err)
		return errors.New("Failed to connect to Figure 1 API, try again later")
	}
	defer res.Body.Close()

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

	app.BearerToken = resBody.Token
	return nil
}
