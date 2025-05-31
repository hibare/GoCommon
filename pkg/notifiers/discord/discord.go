// Package discord provides utilities for sending messages to Discord webhooks.
package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// ErrNoEmbeds is returned when no embeds are available in the message.
var ErrNoEmbeds = errors.New("no embeds available")

// EmbedField represents a field in a Discord embed.
type EmbedField struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

// EmbedImage represents an image in a Discord embed.
type EmbedImage struct {
	URL string `json:"url,omitempty"`
}

// EmbedFooter represents a footer in a Discord embed.
type EmbedFooter struct {
	Text    string `json:"text,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

// EmbedThumbnail represents a thumbnail in a Discord embed.
type EmbedThumbnail struct {
	URL string `json:"url,omitempty"`
}

// EmbedAuthor represents an author in a Discord embed.
type EmbedAuthor struct {
	Name    string `json:"name,omitempty"`
	URL     string `json:"url,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

// Embed represents a Discord embed.
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

// Component represents a Discord component (not implemented).
type Component struct {
	// Define struct for  components if needed
}

// Message represents a Discord webhook message.
type Message struct {
	Embeds     []Embed     `json:"embeds,omitempty"`
	Components []Component `json:"components,omitempty"`
	Username   string      `json:"username,omitempty"`
	Content    string      `json:"content,omitempty"`
	AvatarURL  string      `json:"avatar_url,omitempty"`
}

// AddFooter adds a footer to the last embed in the message.
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

// Send sends the message to the specified Discord webhook URL.
func (d *Message) Send(webhook string) error {
	payload, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusNoContent { //  return 204 on success
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
