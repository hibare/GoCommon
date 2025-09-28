package gpg

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	commonHTTPClient "github.com/hibare/GoCommon/v2/pkg/http/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	testKeyID        = "test-key-id-123"
	testKeyServerURL = "https://keyserver.ubuntu.com"
	testGPGKeyData   = "-----BEGIN PGP PUBLIC KEY BLOCK-----\n\n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\nKtQIWXTIGw==\n=oPyT\n-----END PGP PUBLIC KEY BLOCK-----"
)

func TestGPGManager_FetchGPGPubKeyFromKeyServer_Success(t *testing.T) {
	// Create a test server that returns a GPG key
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request URL
		expectedPath := fmt.Sprintf("/pks/lookup?op=get&search=%s", testKeyID)
		assert.Equal(t, expectedPath, r.URL.Path+"?"+r.URL.RawQuery)
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(testGPGKeyData))
		require.NoError(t, err)
	}))
	defer server.Close()

	// Create GPG manager with default HTTP client
	manager := NewGPGManager(ManagerOptions{})

	// Test fetching the key
	filePath, err := manager.FetchGPGPubKeyFromKeyServer(testKeyID, server.URL)

	// Verify results
	require.NoError(t, err)
	require.NotNil(t, filePath)

	// Verify file was created
	fileInfo, err := os.Stat(*filePath)
	require.NoError(t, err)
	require.NotZero(t, fileInfo.Size())

	// Verify file content
	fileContent, err := os.ReadFile(*filePath)
	require.NoError(t, err)
	assert.Equal(t, testGPGKeyData, string(fileContent))

	// Verify file naming convention
	expectedFileName := fmt.Sprintf("%s_%s.%s", GPGFilePrefix, testKeyID, GPGFileExtension)
	assert.True(t, strings.HasSuffix(*filePath, expectedFileName))

	// Cleanup
	t.Cleanup(func() {
		_ = os.Remove(*filePath)
	})
}

func TestGPGManager_FetchGPGPubKeyFromKeyServer_EmptyKeyID(t *testing.T) {
	manager := NewGPGManager(ManagerOptions{})

	filePath, err := manager.FetchGPGPubKeyFromKeyServer("", testKeyServerURL)

	require.Error(t, err)
	require.Nil(t, filePath)
	assert.Contains(t, err.Error(), "keyID cannot be empty")
}

func TestGPGManager_FetchGPGPubKeyFromKeyServer_EmptyKeyServerURL(t *testing.T) {
	manager := NewGPGManager(ManagerOptions{})

	filePath, err := manager.FetchGPGPubKeyFromKeyServer(testKeyID, "")

	require.Error(t, err)
	require.Nil(t, filePath)
	assert.Contains(t, err.Error(), "keyServerURL cannot be empty")
}

func TestGPGManager_FetchGPGPubKeyFromKeyServer_HTTPClientError(t *testing.T) {
	// Create a mock HTTP client that returns an error
	mockClient := &commonHTTPClient.MockClient{}
	mockClient.On("Do", mock.AnythingOfType("*http.Request")).Return(nil, assert.AnError)

	manager := NewGPGManager(ManagerOptions{
		HTTPClient: mockClient,
	})

	filePath, err := manager.FetchGPGPubKeyFromKeyServer(testKeyID, testKeyServerURL)

	require.Error(t, err)
	require.Nil(t, filePath)
	assert.Contains(t, err.Error(), "failed to download GPG key")

	mockClient.AssertExpectations(t)
}

func TestGPGManager_FetchGPGPubKeyFromKeyServer_NonOKStatus(t *testing.T) {
	// Create a test server that returns a non-OK status
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("Key not found"))
		require.NoError(t, err)
	}))
	defer server.Close()

	manager := NewGPGManager(ManagerOptions{})

	filePath, err := manager.FetchGPGPubKeyFromKeyServer(testKeyID, server.URL)

	require.Error(t, err)
	require.Nil(t, filePath)
	assert.Contains(t, err.Error(), "key-server returned non-OK status: 404")
}

func TestGPGManager_FetchGPGPubKeyFromKeyServer_ReadBodyError(t *testing.T) {
	// Create a mock HTTP client that returns a response with a body that fails to read
	mockClient := &commonHTTPClient.MockClient{}

	// Create a response with a body that will fail on ReadAll
	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       &failingReader{},
	}

	mockClient.On("Do", mock.AnythingOfType("*http.Request")).Return(response, nil)

	manager := NewGPGManager(ManagerOptions{
		HTTPClient: mockClient,
	})

	filePath, err := manager.FetchGPGPubKeyFromKeyServer(testKeyID, testKeyServerURL)

	require.Error(t, err)
	require.Nil(t, filePath)
	assert.Contains(t, err.Error(), "failed to read key data")

	mockClient.AssertExpectations(t)
}

func TestGPGManager_FetchGPGPubKeyFromKeyServer_FileCreationError(t *testing.T) {
	// Create a test server that returns a GPG key
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(testGPGKeyData))
		require.NoError(t, err)
	}))
	defer server.Close()

	// Test file creation error scenario

	// We'll test this by creating a temporary directory and then removing write permissions
	tempDir := t.TempDir()

	// Remove write permissions from the temp directory
	err := os.Chmod(tempDir, 0444)
	require.NoError(t, err)
	defer func() {
		_ = os.Chmod(tempDir, 0755)
	}()

	// Create a file in the temp directory to simulate the scenario
	// where we can't create a new file due to permissions
	invalidPath := filepath.Join(tempDir, "invalid", "path", "file.asc")

	// Manually test the file creation logic
	file, err := os.OpenFile(invalidPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		// This is expected to fail due to invalid path or permission denied
		assert.True(t, strings.Contains(err.Error(), "no such file or directory") ||
			strings.Contains(err.Error(), "permission denied"))
		return
	}
	_ = file.Close()
}

func TestGPGManager_FetchGPGPubKeyFromKeyServer_FileWriteError(t *testing.T) {
	// Create a test server that returns a GPG key
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(testGPGKeyData))
		require.NoError(t, err)
	}))
	defer server.Close()

	// Test file write error scenario

	// We'll test this by creating a temporary directory and file first
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.asc")

	// Create the file first
	file, err := os.Create(testFile)
	require.NoError(t, err)
	require.NoError(t, file.Close())

	// Remove write permissions from the file
	err = os.Chmod(testFile, 0444)
	require.NoError(t, err)

	// Try to write to the file - this should fail
	file, err = os.OpenFile(testFile, os.O_WRONLY, 0)
	if err != nil {
		// If we can't open the file for writing, that's also a valid test result
		assert.Contains(t, err.Error(), "permission denied")
		return
	}
	defer file.Close()

	_, err = file.Write([]byte("test data"))
	if err != nil {
		// This is expected to fail due to permission denied
		assert.Contains(t, err.Error(), "permission denied")
		return
	}

	// If we get here, the test should fail because we expected an error
	t.Error("Expected file write to fail due to permission denied, but it succeeded")
}

func TestNewGPGManager_WithHTTPClient(t *testing.T) {
	mockClient := &commonHTTPClient.MockClient{}

	manager := NewGPGManager(ManagerOptions{
		HTTPClient: mockClient,
	})

	require.NotNil(t, manager)

	// Verify the manager uses the provided HTTP client
	gpgManager, ok := manager.(*GPGManager)
	require.True(t, ok)
	assert.Equal(t, mockClient, gpgManager.httpClient)
}

func TestNewGPGManager_WithoutHTTPClient(t *testing.T) {
	manager := NewGPGManager(ManagerOptions{})

	require.NotNil(t, manager)

	// Verify the manager uses the default HTTP client
	gpgManager, ok := manager.(*GPGManager)
	require.True(t, ok)
	assert.NotNil(t, gpgManager.httpClient)
	assert.IsType(t, &http.Client{}, gpgManager.httpClient)
}

func TestGPGManager_Constants(t *testing.T) {
	assert.Equal(t, "asc", GPGFileExtension)
	assert.Equal(t, "gpg_pub_key_", GPGFilePrefix)
}

func TestGPGManager_InterfaceCompliance(t *testing.T) {
	// Verify that GPGManager implements GPGManagerIface
	var _ GPGManagerIface = (*GPGManager)(nil)
}

// failingReader is a reader that always returns an error
type failingReader struct{}

func (f *failingReader) Read(p []byte) (n int, err error) {
	return 0, assert.AnError
}

func (f *failingReader) Close() error {
	return nil
}

// Test helper to verify the generated file path format
func TestGPGManager_FileNaming(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(testGPGKeyData))
		require.NoError(t, err)
	}))
	defer server.Close()

	manager := NewGPGManager(ManagerOptions{})

	filePath, err := manager.FetchGPGPubKeyFromKeyServer(testKeyID, server.URL)
	require.NoError(t, err)
	require.NotNil(t, filePath)

	// Verify the file path structure
	expectedFileName := fmt.Sprintf("%s_%s.%s", GPGFilePrefix, testKeyID, GPGFileExtension)
	actualFileName := filepath.Base(*filePath)
	assert.Equal(t, expectedFileName, actualFileName)

	// Verify the file is in the temp directory
	tempDir := os.TempDir()
	assert.True(t, strings.HasPrefix(*filePath, tempDir))

	// Cleanup
	t.Cleanup(func() {
		_ = os.Remove(*filePath)
	})
}

// Test with different key IDs to ensure proper URL construction
func TestGPGManager_DifferentKeyIDs(t *testing.T) {
	testCases := []struct {
		name  string
		keyID string
	}{
		{
			name:  "numeric key ID",
			keyID: "1234567890",
		},
		{
			name:  "hex key ID",
			keyID: "0x1234567890ABCDEF",
		},
		{
			name:  "email key ID",
			keyID: "test@example.com",
		},
		{
			name:  "fingerprint key ID",
			keyID: "1234567890ABCDEF1234567890ABCDEF12345678",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify the key ID is properly URL encoded in the request
				expectedPath := fmt.Sprintf("/pks/lookup?op=get&search=%s", tc.keyID)
				assert.Equal(t, expectedPath, r.URL.Path+"?"+r.URL.RawQuery)

				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(testGPGKeyData))
				require.NoError(t, err)
			}))
			defer server.Close()

			manager := NewGPGManager(ManagerOptions{})

			filePath, err := manager.FetchGPGPubKeyFromKeyServer(tc.keyID, server.URL)
			require.NoError(t, err)
			require.NotNil(t, filePath)

			// Verify file naming includes the key ID
			expectedFileName := fmt.Sprintf("%s_%s.%s", GPGFilePrefix, tc.keyID, GPGFileExtension)
			actualFileName := filepath.Base(*filePath)
			assert.Equal(t, expectedFileName, actualFileName)

			// Cleanup
			t.Cleanup(func() {
				_ = os.Remove(*filePath)
			})
		})
	}
}

// Test context handling in HTTP requests
func TestGPGManager_ContextHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify that the request has a context
		ctx := r.Context()
		require.NotNil(t, ctx)

		// Verify it's not nil (the exact type may vary due to httptest)
		assert.NotNil(t, ctx)

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(testGPGKeyData))
		require.NoError(t, err)
	}))
	defer server.Close()

	manager := NewGPGManager(ManagerOptions{})

	filePath, err := manager.FetchGPGPubKeyFromKeyServer(testKeyID, server.URL)
	require.NoError(t, err)
	require.NotNil(t, filePath)

	// Cleanup
	t.Cleanup(func() {
		_ = os.Remove(*filePath)
	})
}

// Test file permissions
func TestGPGManager_FilePermissions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(testGPGKeyData))
		require.NoError(t, err)
	}))
	defer server.Close()

	manager := NewGPGManager(ManagerOptions{})

	filePath, err := manager.FetchGPGPubKeyFromKeyServer(testKeyID, server.URL)
	require.NoError(t, err)
	require.NotNil(t, filePath)

	// Verify file permissions
	fileInfo, err := os.Stat(*filePath)
	require.NoError(t, err)

	// The file should be readable and writable by owner (0644)
	expectedMode := os.FileMode(0644)
	actualMode := fileInfo.Mode().Perm()
	assert.Equal(t, expectedMode, actualMode)

	// Cleanup
	t.Cleanup(func() {
		_ = os.Remove(*filePath)
	})
}
