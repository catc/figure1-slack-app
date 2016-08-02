package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type SlackResponse struct {
	ResponseType string       `json:"response_type"`
	Text         string       `json:"text,omitempty"`
	Attachments  []Attachment `json:"attachments"`
}

type Attachment struct {
	Title      string   `json:"title,omitempty"`
	TitleLink  string   `json:"title_link,omitempty"`
	Text       string   `json:"text,omitempty"`
	PreText    string   `json:"pretext,omitempty"`
	ThumbUrl   string   `json:"thumb_url,omitempty"`
	Footer     string   `json:"footer,omitempty"`
	FooterIcon string   `json:"footer_icon,omitempty"`
	Color      string   `json:"color,omitempty"`
	Markdown   []string `json:"mrkdwn_in,omitempty"`
}

func slackResponse(res http.ResponseWriter, data *f1Case) {
	resp := &SlackResponse{
		ResponseType: "in_channel",
	}

	// author
	authorSection := Attachment{
		Title:     data.Author.Username,
		TitleLink: linkgen("user", data.Author.Username),
		PreText:   "*Figure 1 Case*",
		Markdown:  []string{"pretext"},
	}
	if data.Author.TopContributor {
		authorSection.Footer = "Top Contributor"
		authorSection.FooterIcon = "http://i.imgur.com/oYpmgwF.jpg"
	} else if data.Author.Verified {
		authorSection.Footer = "Verified"
		authorSection.FooterIcon = "http://i.imgur.com/9eyI61P.jpg"
	}
	resp.Attachments = append(resp.Attachments, authorSection)

	// case info
	caseInfoSection := Attachment{
		ThumbUrl: linkgen("image", data.Id),
	}
	split := strings.Split(data.Caption, " ")
	const limit int = 36
	if len(split) > limit {
		caseInfoSection.Text = strings.Join(split[0:limit], " ") + "..."
	} else {
		caseInfoSection.Text = data.Caption
	}
	caseInfoSection.Footer = strings.Join([]string{
		data.ImageViews,
		strconv.Itoa(data.VoteCount) + " stars",
		strconv.Itoa(data.CommentCount) + " comments",
		strconv.Itoa(data.Followers) + " followers",
	}, ", ")
	resp.Attachments = append(resp.Attachments, caseInfoSection)

	// share links
	shareSection := Attachment{
		Title: "Share case link",
		Text:  linkgen("case", data.Id),
		Color: "#8bcaf1",
	}
	resp.Attachments = append(resp.Attachments, shareSection)

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(resp)
}

func linkgen(linkType, val string) string {
	switch linkType {
	case "user":
		return "https://app.figure1.com/rd/publicprofile?username=" + val
	case "case":
		return "https://app.figure1.com/rd/image?imageid=" + val
	case "image":
		return "https://s3.amazonaws.com/static.figure1.com/img/share/" + val
	}
	return ""
}
