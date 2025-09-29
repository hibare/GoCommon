package hash

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

// MockHasher is a mock implementation of the Hasher interface.
type MockHasher struct {
	mock.Mock
}

// HashString is a mock implementation of the HashString method.
func (m *MockHasher) HashString(data string) (string, error) {
	args := m.Called(data)
	return args.String(0), args.Error(1)
}

// VerifyString is a mock implementation of the VerifyString method.
func (m *MockHasher) VerifyString(data string, hash string) (bool, error) {
	args := m.Called(data, hash)
	return args.Bool(0), args.Error(1)
}

// HashFile is a mock implementation of the HashFile method.
func (m *MockHasher) HashFile(filePath string) (string, error) {
	args := m.Called(filePath)
	return args.String(0), args.Error(1)
}

// VerifyFile is a mock implementation of the VerifyFile method.
func (m *MockHasher) VerifyFile(filePath string, hash string) (bool, error) {
	args := m.Called(filePath, hash)
	return args.Bool(0), args.Error(1)
}

// SetupMockSHA256HasherWithT is a helper function to setup a mock SHA256 hasher.
func SetupMockSHA256HasherWithT(t *testing.T) *MockHasher {
	mock := &MockHasher{}
	NewSHA256Hasher = func() Hasher {
		return mock
	}

	t.Cleanup(func() {
		NewSHA256Hasher = newSHA256Hasher
	})

	return mock
}
