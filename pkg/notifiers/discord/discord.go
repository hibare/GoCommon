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
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

type EmbedImage struct {
	URL string `json:"url"`
}

type EmbedFooter struct {
	Text string `json:"text"`
}

type Embed struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Color       int          `json:"color"`
	Footer      EmbedFooter  `json:"footer"`
	Fields      []EmbedField `json:"fields"`
	Image       EmbedImage   `json:"image"`
}

type Component struct {
	// Define struct for  components if needed
}

type Message struct {
	Embeds     []Embed     `json:"embeds"`
	Components []Component `json:"components"`
	Username   string      `json:"username"`
	Content    string      `json:"content"`
	AvatarURL  string      `json:"avatar_url"`
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
