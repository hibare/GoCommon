package hash

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) HashString(data string) (string, error) {
	args := m.Called(data)
	return args.String(0), args.Error(1)
}

func (m *MockHasher) VerifyString(data string, hash string) (bool, error) {
	args := m.Called(data, hash)
	return args.Bool(0), args.Error(1)
}

func (m *MockHasher) HashFile(filePath string) (string, error) {
	args := m.Called(filePath)
	return args.String(0), args.Error(1)
}

func (m *MockHasher) VerifyFile(filePath string, hash string) (bool, error) {
	args := m.Called(filePath, hash)
	return args.Bool(0), args.Error(1)
}

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
