// Package testhelper provides utilities for test setup and teardown.
package testhelper

import (
	"os"
	"path/filepath"
)

const (
	// TestDataDir is the directory containing test data.
	TestDataDir = "../testhelper/test_data"
)

// CreateTestFile creates a test file with the given content and returns its path.
func CreateTestFile(dir, pattern string) ([]byte, string, error) {
	if pattern == "" {
		pattern = "test-file-*.txt"
	}

	file, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return []byte{}, "", err
	}
	defer func() {
		_ = file.Close()
	}()

	content := []byte("This is a test file.\nIt contains some sample content.")
	_, err = file.Write(content)
	if err != nil {
		return []byte{}, "", err
	}

	absPath, err := filepath.Abs(file.Name())
	if err != nil {
		return []byte{}, "", err
	}

	return content, absPath, err
}

// CreateTestDir creates a test directory and returns its path.
func CreateTestDir(dir, pattern string) (string, error) {
	if dir == "" {
		dir = os.TempDir()
	}

	if pattern == "" {
		pattern = "test-dir-"
	}

	randomDirName, err := os.MkdirTemp(dir, pattern)
	if err != nil {
		return "", err
	}

	_, _, err = CreateTestFile(randomDirName, "")
	return randomDirName, err
}

// StringToPtr converts a string to a pointer.
func StringToPtr(s string) *string {
	return &s
}
