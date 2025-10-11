package runtime

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

// MockRuntime is a mock implementation of RuntimeIface for testing.
type MockRuntime struct {
	mock.Mock
}

// GetGOOS provides a mock function with given fields.
func (_m *MockRuntime) GetGOOS() string {
	mockArgs := _m.Called()
	return mockArgs.String(0)
}

// GetPlatform provides a mock function with given fields.
func (_m *MockRuntime) GetPlatform() string {
	mockArgs := _m.Called()
	return mockArgs.String(0)
}

// GetConfigDir provides a mock function with given fields.
func (_m *MockRuntime) GetConfigDir() string {
	mockArgs := _m.Called()
	return mockArgs.String(0)
}

// GetConfigFilePath provides a mock function with given fields.
func (_m *MockRuntime) GetConfigFilePath() string {
	mockArgs := _m.Called()
	return mockArgs.String(0)
}

// SetupMock sets up the mock runtime for testing and ensures expectations are met.
func SetupMock(t *testing.T) RuntimeIface {
	mock := &MockRuntime{}
	New = func() RuntimeIface {
		return mock
	}
	t.Cleanup(func() {
		New = newRuntime
		mock.AssertExpectations(t)
	})
	return mock
}
