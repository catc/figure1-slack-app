package main

import (
	"fmt"
	"testing"
)

func TestCaseIDParse(t *testing.T) {
	valid := []string{
		"59076d6324d11b594b2dff1d",                                                // id
		"https://app.figure1.com/?image=59076d6324d11b594b2dff1d&t=0",             // web app modal
		"https://app.figure1.com/rd/image?imageid=59076d6324d11b594b2dff1d",       // rd share link
		"https://app.figure1.com/images/59076d6324d11b594b2dff1d?imageType=0&t=0", // web app standalone/url
	}
	for _, str := range valid {
		id := getCaseID(str)
		if id == "" {
			t.Error(fmt.Sprintf("Expected case \"%v\" to be valid", str))
		}
	}

	invalid := []string{
		"abcd",
	}
	for _, str := range invalid {
		id := getCaseID(str)
		if id != "" {
			t.Error(fmt.Sprintf("Expected case \"%v\" to be invalid", str))
		}
	}
}

func TestUsernameParse(t *testing.T) {
	valid := []string{
		"59076d6324d11b594b2dff1d",                                 // username
		"https://app.figure1.com/rd/publicprofile?username=ccovic", // rd share link
		"https://app.figure1.com/user/richardpenner",               // web app route
	}
	for _, str := range valid {
		id := getUsername(str)
		if id == "" {
			t.Error(fmt.Sprintf("Expected user \"%v\" to be valid", str))
		}
	}
}

func TestCollectionIDParse(t *testing.T) {
	valid := []string{
		"5907549889c89eef5b1b3511",                                           // username
		"https://app.figure1.com/rd/collections?id=5907549889c89eef5b1b3511", // rd share link
		"https://app.figure1.com/collections/5907549889c89eef5b1b3511",       // web app route
	}
	for _, str := range valid {
		id := getCollectionID(str)
		if id == "" {
			t.Error(fmt.Sprintf("Expected collection \"%v\" to be valid", str))
		}
	}
}
