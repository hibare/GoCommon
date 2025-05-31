package discord

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/mock"
)

// MockClient is a mock implementation of the Client interface.
type MockClient struct {
	mock.Mock
}

// Send is a mock implementation of the Send method.
func (m *MockClient) Send(ctx context.Context, msg *Message) error {
	args := m.Called(ctx, msg)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return args.Error(0)
}

// MockHTTPClient is a mock implementation of the HTTPClient interface.
type MockHTTPClient struct {
	mock.Mock
}

// Do is a mock implementation of the Do method.
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}
