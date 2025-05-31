// Package client provides a HTTP client.
package client

import (
	"net/http"
)

// Client is the interface for the HTTP client.
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}
