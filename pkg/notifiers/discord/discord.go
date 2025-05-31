// Package discord provides a Discord client for sending messages to Discord webhooks.
package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	commonHTTPClient "github.com/hibare/GoCommon/v2/pkg/http/client"
)

// ErrNoEmbeds is returned when no embeds are available in the message.
var ErrNoEmbeds = errors.New("no embeds available")

// Client is the interface for the Discord client.
type Client interface {
	Send(ctx context.Context, msg *Message) error
}

type client struct {
	webhookURL string
	httpClient commonHTTPClient.Client
}

// Options is the options for the Discord client.
type Options struct {
	WebhookURL string
	HTTPClient commonHTTPClient.Client
}

func (o *Options) validate() error {
	if o.WebhookURL == "" {
		return errors.New("webhook URL is required")
	}
	return nil
}

// Send sends the message to the Discord webhook.
func (c *client) Send(ctx context.Context, msg *Message) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusNoContent { // 204 on success
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

// NewClient creates a new Discord client with the given options.
func NewClient(opts Options) (Client, error) {
	if err := opts.validate(); err != nil {
		return nil, err
	}

	if opts.HTTPClient == nil {
		opts.HTTPClient = commonHTTPClient.NewDefaultClient()
	}

	return &client{
		webhookURL: opts.WebhookURL,
		httpClient: opts.HTTPClient,
	}, nil
}
