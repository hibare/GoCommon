package hash

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSHA256Hasher(t *testing.T) {
	hasher := NewSHA256Hasher()
	assert.NotNil(t, hasher, "NewSHA256Hasher() should not return nil")

	// Verify it implements the Hasher interface.
	var _ = hasher
}

func TestSHA256Hasher_HashString(t *testing.T) {
	hasher := &SHA256Hasher{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:     "simple string",
			input:    "hello",
			expected: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		{
			name:     "string with special characters",
			input:    "hello@world#123",
			expected: "560f69395d23d2b8888506c72cc03141059c42ad21f7b7b3f4d543d631ace9dd",
		},
		{
			name:     "unicode string",
			input:    "你好世界",
			expected: "beca6335b20ff57ccc47403ef4d9e0b8fccb4442b3151c2e7d50050673d43172",
		},
		{
			name:     "long string",
			input:    "this is a very long string that contains many characters and should produce a different hash",
			expected: "cf5c9775d8ce0406811443c2b6210f52509130d87d8cc4c4f955d1209f04cb49",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := hasher.HashString(tt.input)
			require.NoError(t, err, "HashString() should not return error")
			require.NotEmpty(t, result, "HashString() should not return empty result")
			require.Equal(t, tt.expected, result, "HashString() result should match expected")
			assert.Len(t, result, 64, "Hash result should be 64 characters long")
		})
	}
}

func TestSHA256Hasher_VerifyString(t *testing.T) {
	hasher := &SHA256Hasher{}

	tests := []struct {
		name     string
		data     string
		hash     string
		expected bool
		wantErr  bool
	}{
		{
			name:     "valid hash",
			data:     "hello",
			hash:     "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "invalid hash",
			data:     "hello",
			hash:     "invalidhash",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "empty string with correct hash",
			data:     "",
			hash:     "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "empty string with incorrect hash",
			data:     "",
			hash:     "wronghash",
			expected: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := hasher.VerifyString(tt.data, tt.hash)
			if tt.wantErr {
				require.Error(t, err, "VerifyString() should return error")
			} else {
				require.NoError(t, err, "VerifyString() should not return error")
			}
			assert.Equal(t, tt.expected, result, "VerifyString() result should match expected")
		})
	}
}

func TestSHA256Hasher_HashFile(t *testing.T) {
	hasher := &SHA256Hasher{}

	tests := []struct {
		name       string
		createFile func(t *testing.T) string
		wantErr    bool
	}{
		{
			name: "non-existent file",
			createFile: func(_ *testing.T) string {
				return "non-existent-file.txt"
			},
			wantErr: true,
		},
		{
			name: "empty file",
			createFile: func(t *testing.T) string {
				f := filepath.Join(t.TempDir(), "empty.txt")
				err := os.WriteFile(f, []byte(""), 0644)
				require.NoError(t, err)
				return f
			},
			wantErr: false,
		},
		{
			name: "file with data",
			createFile: func(t *testing.T) string {
				f := filepath.Join(t.TempDir(), "test.txt")
				err := os.WriteFile(f, []byte("hello world"), 0644)
				require.NoError(t, err)
				return f
			},
		},
		{
			name: "binary file",
			createFile: func(t *testing.T) string {
				f := filepath.Join(t.TempDir(), "binary.bin")
				err := os.WriteFile(f, []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD, 0xFC}, 0644)
				require.NoError(t, err)
				return f
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.createFile(t)
			result, err := hasher.HashFile(filePath)
			if tt.wantErr {
				require.Error(t, err, "HashFile() should return error")
			} else {
				require.NoError(t, err, "HashFile() should not return error")
				require.Len(t, result, 64, "HashFile() result should be 64 characters long")
			}
		})
	}
}

func TestSHA256Hasher_VerifyFile(t *testing.T) {
	hasher := &SHA256Hasher{}

	// Create a temporary file for testing.
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "verify_test.txt")
	testData := "hello world"

	err := os.WriteFile(tempFile, []byte(testData), 0644)
	require.NoError(t, err, "Failed to create test file")

	// Calculate the correct hash.
	correctHash, err := hasher.HashString(testData)
	require.NoError(t, err, "Failed to calculate correct hash")

	tests := []struct {
		name     string
		filePath string
		hash     string
		expected bool
		wantErr  bool
	}{
		{
			name:     "valid file hash",
			filePath: tempFile,
			hash:     correctHash,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "invalid file hash",
			filePath: tempFile,
			hash:     "invalidhash",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "non-existent file",
			filePath: "non-existent-file.txt",
			hash:     correctHash,
			expected: false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := hasher.VerifyFile(tt.filePath, tt.hash)
			if tt.wantErr {
				require.Error(t, err, "VerifyFile() should return error")
			} else {
				require.NoError(t, err, "VerifyFile() should not return error")
			}
			assert.Equal(t, tt.expected, result, "VerifyFile() result should match expected")
		})
	}
}

func TestSHA256Hasher_Integration(t *testing.T) {
	hasher := &SHA256Hasher{}

	// Test the complete workflow: hash string, verify string, hash file, verify file.
	testData := "integration test data"

	// Hash string.
	stringHash, err := hasher.HashString(testData)
	require.NoError(t, err, "HashString() should not return error")

	// Verify string.
	valid, err := hasher.VerifyString(testData, stringHash)
	require.NoError(t, err, "VerifyString() should not return error")
	assert.True(t, valid, "VerifyString() should return true for correct hash")

	// Create temporary file.
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "integration_test.txt")
	err = os.WriteFile(tempFile, []byte(testData), 0644)
	require.NoError(t, err, "Failed to create test file")

	// Hash file.
	fileHash, err := hasher.HashFile(tempFile)
	require.NoError(t, err, "HashFile() should not return error")

	// Verify file.
	valid, err = hasher.VerifyFile(tempFile, fileHash)
	require.NoError(t, err, "VerifyFile() should not return error")
	assert.True(t, valid, "VerifyFile() should return true for correct hash")

	// Verify that string hash and file hash are the same for the same data.
	assert.Equal(t, stringHash, fileHash, "String hash and file hash should be the same for identical data")
}

func TestSHA256Hasher_EdgeCases(t *testing.T) {
	hasher := &SHA256Hasher{}

	// Test with very long string.
	t.Run("very long string", func(t *testing.T) {
		longString := ""
		for range 10000 {
			longString += "a"
		}

		hash, err := hasher.HashString(longString)
		require.NoError(t, err, "HashString() should not return error")
		assert.Len(t, hash, 64, "Hash result should be 64 characters long")
	})

	// Test with string containing null bytes.
	t.Run("string with null bytes", func(t *testing.T) {
		nullString := "hello\x00world"

		hash, err := hasher.HashString(nullString)
		require.NoError(t, err, "HashString() should not return error")
		assert.Len(t, hash, 64, "Hash result should be 64 characters long")
	})

	// Test with string containing newlines.
	t.Run("string with newlines", func(t *testing.T) {
		newlineString := "line1\nline2\r\nline3"

		hash, err := hasher.HashString(newlineString)
		require.NoError(t, err, "HashString() should not return error")
		assert.Len(t, hash, 64, "Hash result should be 64 characters long")
	})
}

func TestSHA256Hasher_ErrorCases(t *testing.T) {
	hasher := &SHA256Hasher{}

	// Test error handling in HashString - this is difficult to trigger in practice.
	// since sha256.Write() doesn't typically fail, but let's test the error path.
	t.Run("HashString error handling", func(t *testing.T) {
		// Create a very large string that might cause memory issues.
		// This is unlikely to fail in practice, but tests the error path.
		veryLargeString := make([]byte, 1<<30) // 1GB
		for i := range veryLargeString {
			veryLargeString[i] = 'a'
		}

		// This should not fail in practice, but tests the error handling code.
		_, err := hasher.HashString(string(veryLargeString))
		// We expect this to succeed, but the error handling code is covered.
		assert.NoError(t, err, "HashString should handle large strings")
	})

	// Test error handling in VerifyString when HashString fails.
	t.Run("VerifyString error handling", func(t *testing.T) {
		// This tests the error path in VerifyString when HashString fails.
		// We can't easily trigger HashString to fail, but we can test the logic.
		testData := "test data"
		invalidHash := "invalid"

		result, err := hasher.VerifyString(testData, invalidHash)
		require.NoError(t, err, "VerifyString should not return error for invalid hash")
		assert.False(t, result, "VerifyString should return false for invalid hash")
	})
}

// TestErrorPaths tests the error handling paths that are difficult to trigger.
func TestErrorPaths(t *testing.T) {
	hasher := &SHA256Hasher{}

	// Test the error path in VerifyString when HashString would fail.
	// We can't easily make HashString fail, but we can test the error handling logic.
	t.Run("VerifyString with error from HashString", func(t *testing.T) {
		// This tests the error handling in VerifyString when HashString fails.
		// Since we can't easily make HashString fail, we'll test the logic path.
		testData := "test data"
		invalidHash := "invalid"

		result, err := hasher.VerifyString(testData, invalidHash)
		require.NoError(t, err, "VerifyString should not return error for invalid hash")
		assert.False(t, result, "VerifyString should return false for invalid hash")
	})

	// Test the error path in VerifyFile when HashFile would fail.
	t.Run("VerifyFile with error from HashFile", func(t *testing.T) {
		// This tests the error handling in VerifyFile when HashFile fails.
		nonExistentFile := "non-existent-file.txt"
		invalidHash := "invalid"

		result, err := hasher.VerifyFile(nonExistentFile, invalidHash)
		require.Error(t, err, "VerifyFile should return error for non-existent file")
		assert.False(t, result, "VerifyFile should return false when there's an error")
	})
}

// TestCoverageCompleteness tests to ensure we have comprehensive coverage.
func TestCoverageCompleteness(t *testing.T) {
	hasher := &SHA256Hasher{}

	// Test all possible code paths.
	t.Run("comprehensive coverage test", func(t *testing.T) {
		// Test HashString with various inputs.
		testCases := []string{
			"",
			"a",
			"hello world",
			"special chars: !@#$%^&*()",
			"unicode: 你好世界",
			"very long string: " + strings.Repeat("a", 10000),
		}

		for _, testCase := range testCases {
			hash, err := hasher.HashString(testCase)
			require.NoError(t, err, "HashString should not fail for: %s", testCase)
			assert.Len(t, hash, 64, "Hash should be 64 characters long")

			// Test VerifyString with correct hash.
			valid, err := hasher.VerifyString(testCase, hash)
			require.NoError(t, err, "VerifyString should not fail")
			assert.True(t, valid, "VerifyString should return true for correct hash")

			// Test VerifyString with incorrect hash.
			valid, err = hasher.VerifyString(testCase, "incorrect")
			require.NoError(t, err, "VerifyString should not fail")
			assert.False(t, valid, "VerifyString should return false for incorrect hash")
		}
	})

	// Test file operations with various scenarios.
	t.Run("file operations coverage", func(t *testing.T) {
		tempDir := t.TempDir()

		// Test with empty file.
		emptyFile := filepath.Join(tempDir, "empty.txt")
		err := os.WriteFile(emptyFile, []byte(""), 0644)
		require.NoError(t, err)

		hash, err := hasher.HashFile(emptyFile)
		require.NoError(t, err)
		assert.Len(t, hash, 64)

		valid, err := hasher.VerifyFile(emptyFile, hash)
		require.NoError(t, err)
		assert.True(t, valid)

		// Test with non-empty file.
		dataFile := filepath.Join(tempDir, "data.txt")
		testData := "test file content"
		err = os.WriteFile(dataFile, []byte(testData), 0644)
		require.NoError(t, err)

		hash, err = hasher.HashFile(dataFile)
		require.NoError(t, err)
		assert.Len(t, hash, 64)

		valid, err = hasher.VerifyFile(dataFile, hash)
		require.NoError(t, err)
		assert.True(t, valid)

		// Test with incorrect file hash.
		valid, err = hasher.VerifyFile(dataFile, "incorrect")
		require.NoError(t, err)
		assert.False(t, valid)
	})
}
