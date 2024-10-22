package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrNoEmbeds = errors.New("no embeds available")
)

type EmbedField struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

type EmbedImage struct {
	URL string `json:"url,omitempty"`
}

type EmbedFooter struct {
	Text    string `json:"text,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

type EmbedThumbnail struct {
	URL string `json:"url,omitempty"`
}

type EmbedAuthor struct {
	Name    string `json:"name,omitempty"`
	URL     string `json:"url,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

type Embed struct {
	Title       string         `json:"title,omitempty"`
	URL         string         `json:"url,omitempty"`
	Description string         `json:"description,omitempty"`
	Color       int            `json:"color,omitempty"`
	Footer      EmbedFooter    `json:"footer,omitempty"`
	Fields      []EmbedField   `json:"fields,omitempty"`
	Image       EmbedImage     `json:"image,omitempty"`
	Thumbnail   EmbedThumbnail `json:"thumbnail,omitempty"`
	Author      EmbedAuthor    `json:"author,omitempty"`
}

type Component struct {
	// Define struct for  components if needed
}

type Message struct {
	Embeds     []Embed     `json:"embeds,omitempty"`
	Components []Component `json:"components,omitempty"`
	Username   string      `json:"username,omitempty"`
	Content    string      `json:"content,omitempty"`
	AvatarURL  string      `json:"avatar_url,omitempty"`
}

func (d *Message) AddFooter(footerStr string) error {
	if len(d.Embeds) == 0 {
		return ErrNoEmbeds
	}

	lastEmbedIndex := len(d.Embeds) - 1
	d.Embeds[lastEmbedIndex].Footer = EmbedFooter{
		Text: footerStr,
	}

	return nil
}

func (d *Message) Send(webhook string) error {
	payload, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent { //  return 204 on success
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
