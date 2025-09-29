// Package client provides a HTTP client.
package client

import (
	"net/http"

	commonHTTP "github.com/hibare/GoCommon/v2/pkg/http"
)

// ClientIface is the interface for the HTTP client.
type ClientIface interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewDefaultClient creates a new default HTTP client.
func NewDefaultClient() ClientIface {
	return &http.Client{
		Timeout: commonHTTP.DefaultHTTPClientTimeout,
	}
}
