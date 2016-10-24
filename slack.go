package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type SlackResponse struct {
	ResponseType string        `json:"response_type"`
	Text         string        `json:"text,omitempty"`
	Attachments  []*Attachment `json:"attachments"`
}

type Attachment struct {
	AuthorName string   `json:"author_name,omitempty"`
	Title      string   `json:"title,omitempty"`
	Fallback   string   `json:"fallback,omitempty"`
	TitleLink  string   `json:"title_link,omitempty"`
	Text       string   `json:"text,omitempty"` // can contain markup
	PreText    string   `json:"pretext,omitempty"`
	ThumbUrl   string   `json:"thumb_url,omitempty"`
	Footer     string   `json:"footer,omitempty"`
	FooterIcon string   `json:"footer_icon,omitempty"`
	Color      string   `json:"color,omitempty"`
	Markdown   []string `json:"mrkdwn_in,omitempty"`
	Fields     []*Field `json:"fields,omitempty"`
}

type Field struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"` // can contain markup
	Short bool   `json:"short,omitempty"`
}

func slackCaseResponse(res http.ResponseWriter, data *f1Case) {
	resp := &SlackResponse{
		ResponseType: "in_channel",
	}

	// author
	authorSection := Attachment{
		Title:     data.Author.Username,
		TitleLink: caseLinkGen("user", data.Author.Username),
	}
	if data.Author.TopContributor {
		authorSection.Footer = "Top Contributor"
		authorSection.FooterIcon = "http://i.imgur.com/oYpmgwF.jpg"
	} else if data.Author.Verified {
		authorSection.Footer = "Verified"
		authorSection.FooterIcon = "http://i.imgur.com/9eyI61P.jpg"
	}
	resp.Attachments = append(resp.Attachments, &authorSection)

	// case info
	caseInfoSection := Attachment{
		ThumbUrl: caseLinkGen("image", data.Id),
	}
	split := strings.Split(data.Caption, " ")
	const limit int = 36
	var caption string
	if len(split) > limit {
		caption = strings.Join(split[0:limit], " ") + "..."
	} else {
		caption = data.Caption
	}
	authorSection.Fallback = "FIGURE 1 CASE: " + caption
	caseInfoSection.Text = caption
	caseInfoSection.Footer = strings.Join([]string{
		data.ImageViews,
		strconv.Itoa(data.VoteCount) + " stars",
		strconv.Itoa(data.CommentCount) + " comments",
		strconv.Itoa(data.Followers) + " followers",
	}, ", ")
	resp.Attachments = append(resp.Attachments, &caseInfoSection)

	// share links
	shareSection := Attachment{
		Title: "Share case link",
		Text:  caseLinkGen("case", data.Id),
		Color: "#8bcaf1",
	}
	resp.Attachments = append(resp.Attachments, &shareSection)

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(resp)
}

func caseLinkGen(linkType, val string) string {
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

func userLinkGen(username string) string {
	return "https://app.figure1.com/rd/publicprofile?username=" + username
}

func slackUserResponse(res http.ResponseWriter, data *f1User) {
	fmt.Println("Data is", data)
	resp := &SlackResponse{
		ResponseType: "in_channel",
	}

	// main section
	mainSection := Attachment{
		Title:     data.Username,
		TitleLink: userLinkGen(data.Username),
		Text:      data.Category + ", " + data.Specialty,
	}
	if data.TopContributor {
		mainSection.Footer = "Top Contributor"
		mainSection.FooterIcon = "http://i.imgur.com/oYpmgwF.jpg"
	} else if data.Verified {
		mainSection.Footer = "Verified"
		mainSection.FooterIcon = "http://i.imgur.com/9eyI61P.jpg"
	}
	resp.Attachments = append(resp.Attachments, &mainSection)

	// extra content section
	extraContentSection := Attachment{
		Fallback: "FIGURE 1 USER: " + data.Username,
		Title:    data.FullName,
		Text:     data.Bio,
	}

	// link
	if data.Link != "" {
		linkField := &Field{
			Title: "Link",
			Value: data.Link,
			Short: true,
		}
		extraContentSection.Fields = append(extraContentSection.Fields, linkField)
	}

	// institution + location
	loc := ""
	if data.ProfileCountry != "" {
		loc = data.ProfileCountry
	} else {
		loc = data.Country
	}
	if data.Institution != "" {
		loc = data.Institution + ", " + loc
	}
	if loc != "" {
		institutionCountryField := &Field{
			Title: "Institution/Country",
			Value: loc,
			Short: true,
		}
		extraContentSection.Fields = append(extraContentSection.Fields, institutionCountryField)
	}

	// stats
	var stats []string
	for key, count := range map[string]int{
		"comments":  data.CommentsCount,
		"favorites": data.FavoritesCount,
		"followers": data.FollowersCount,
		"following": data.FollowingCount,
		"uploads":   data.UploadsCount,
	} {
		if count != 0 {
			stats = append(stats, strconv.Itoa(count)+" "+key)
		}
	}
	extraContentSection.Footer = strings.Join(stats, ", ")

	resp.Attachments = append(resp.Attachments, &extraContentSection)

	// share links
	shareSection := Attachment{
		Title: "Share profile link",
		Text:  userLinkGen(data.Username),
		Color: "#8bcaf1",
	}
	resp.Attachments = append(resp.Attachments, &shareSection)

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(resp)
}
