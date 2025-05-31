// Package client provides a mock implementation of the Client interface.
package client

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// MockClient is a mock implementation of the Client interface.
type MockClient struct {
	mock.Mock
}

// Do is a mock implementation of the Do method.
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}
