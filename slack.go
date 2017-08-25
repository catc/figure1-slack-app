package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const slackPostMsgLink = "https://slack.com/api/chat.postMessage"
const (
	verifiedBadgeLink       = "http://i.imgur.com/9eyI61P.jpg"
	topContributorBadgeLink = "http://i.imgur.com/oYpmgwF.jpg"
	textCasePlaceholder     = "http://i.imgur.com/9Tpmuwk.png"
)
const (
	colorLightBlue = "#8bcaf1"
	colorRed       = "#fd7f8a"
)

// SlackResponse is the wrapper for responding to slash commands
type SlackResponse struct {
	ResponseType string        `json:"response_type"`
	Text         string        `json:"text,omitempty"`
	Attachments  []*Attachment `json:"attachments"`
}

// Attachment is individual item when posting a message to slack
type Attachment struct {
	AuthorName string   `json:"author_name,omitempty"`
	AuthorLink string   `json:"author_link,omitempty"`
	Title      string   `json:"title,omitempty"`
	Fallback   string   `json:"fallback,omitempty"`
	TitleLink  string   `json:"title_link,omitempty"`
	Text       string   `json:"text,omitempty"` // can contain markup
	PreText    string   `json:"pretext,omitempty"`
	ThumbURL   string   `json:"thumb_url,omitempty"`
	Footer     string   `json:"footer,omitempty"`
	FooterIcon string   `json:"footer_icon,omitempty"`
	Color      string   `json:"color,omitempty"`
	Markdown   []string `json:"mrkdwn_in,omitempty"`
	Fields     []*Field `json:"fields,omitempty"`
}

// Field contains segments of data, part of an attachment
type Field struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"` // can contain markup
	Short bool   `json:"short,omitempty"`
}

func respondToSlashCommand(link string, attachments []*Attachment) {
	body := &SlackResponse{
		ResponseType: "in_channel",
		Attachments:  attachments,
	}

	// marshal body
	reqBody := new(bytes.Buffer)
	if err := json.NewEncoder(reqBody).Encode(body); err != nil {
		msg := "Failed to encode slack JSON"
		(&slackError{msg, msg, err}).handleError(link)
		return
	}

	// create request
	req, err := http.NewRequest("POST", link, reqBody)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		msg := "Failed to connect to slack api"
		(&slackError{msg, msg, err}).handleError(link)
		return
	}
	defer resp.Body.Close()
}

func generateCaseContent(data *f1Case, opUser string) []*Attachment {
	attachments := []*Attachment{}

	// author
	authorSection := Attachment{
		Title:     data.Author.Username,
		TitleLink: userLinkGen(data.Author.Username),
	}
	if data.Author.TopContributor {
		authorSection.Footer = "Top Contributor"
		authorSection.FooterIcon = topContributorBadgeLink
	} else if data.Author.Verified {
		authorSection.Footer = "Verified"
		authorSection.FooterIcon = verifiedBadgeLink
	}
	attachments = append(attachments, &authorSection)

	// case info
	caseInfoSection := Attachment{
		ThumbURL: caseLinkGen("image", data.ID),
	}
	caption := truncateString(data.Caption)
	authorSection.Fallback = "FIGURE 1 CASE: " + caption
	caseInfoSection.Text = caption
	caseInfoSection.Footer = strings.Join([]string{
		data.ImageViews,
		strconv.Itoa(data.VoteCount) + " stars",
		strconv.Itoa(data.CommentCount) + " comments",
		strconv.Itoa(data.Followers) + " followers",
	}, ", ")
	attachments = append(attachments, &caseInfoSection)

	// share links
	shareSection := Attachment{
		Title:  "Share case link",
		Text:   caseLinkGen("case", data.ID),
		Color:  colorLightBlue,
		Footer: fmt.Sprintf("posted by @%v", opUser),
	}
	attachments = append(attachments, &shareSection)

	return attachments
}

func generateUserContent(data *f1User, opUser string) []*Attachment {
	attachments := []*Attachment{}

	// main section
	mainSection := Attachment{
		Title:     data.Username,
		TitleLink: userLinkGen(data.Username),
		Text:      data.Category + ", " + data.Specialty,
	}
	if data.TopContributor {
		mainSection.Footer = "Top Contributor"
		mainSection.FooterIcon = topContributorBadgeLink
	} else if data.Verified {
		mainSection.Footer = "Verified"
		mainSection.FooterIcon = verifiedBadgeLink
	}
	attachments = append(attachments, &mainSection)

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
	attachments = append(attachments, &extraContentSection)

	// share links
	shareSection := Attachment{
		Title:  "Share profile link",
		Text:   userLinkGen(data.Username),
		Color:  colorLightBlue,
		Footer: fmt.Sprintf("posted by @%v", opUser),
	}
	attachments = append(attachments, &shareSection)

	return attachments
}

func generateCollectionContent(data *f1Collection, opUser string) []*Attachment {
	attachments := []*Attachment{}

	// collection info
	author := data.Embedded.Authors[0]
	mainSection := Attachment{
		AuthorLink: userLinkGen(author.Username),
		AuthorName: author.Username,
		Title:      data.Title,
		Fallback:   "FIGURE 1 COLLECTION: " + data.Title,
	}
	if data.Size == 1 {
		mainSection.Footer = "1 case"
	} else {
		mainSection.Footer = fmt.Sprintf("%v cases", data.Size)
	}
	mainSection.Text = truncateString(data.Description)
	attachments = append(attachments, &mainSection)

	// items
	items := data.Embedded.Items
	length := 3
	if len(items) < length {
		length = len(items)
	}
	for _, item := range items[0:length] {
		attachment := Attachment{
			Color:    colorRed,
			Text:     truncateString(item.Caption),
			ThumbURL: genCollectionItemImageLink(item.Links.Image.Href, item.ID),
		}
		attachment.Footer = strings.Join([]string{
			strconv.Itoa(item.VoteCount) + " stars",
			strconv.Itoa(item.CommentCount) + " comments",
			strconv.Itoa(item.Followers) + " followers",
		}, ", ")
		attachments = append(attachments, &attachment)
	}

	// share links
	shareSection := Attachment{
		Title:  "Share collection link",
		Text:   collectionLinkGen(data.ID),
		Color:  colorLightBlue,
		Footer: fmt.Sprintf("posted by @%v", opUser),
	}
	attachments = append(attachments, &shareSection)

	return attachments
}

func caseLinkGen(linkType, val string) string {
	switch linkType {
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

func collectionLinkGen(id string) string {
	return "https://app.figure1.com/rd/collections?id=" + id
}

func genCollectionItemImageLink(link, collectionID string) string {
	if link == "" {
		return textCasePlaceholder
	}
	return caseLinkGen("image", collectionID)
}
