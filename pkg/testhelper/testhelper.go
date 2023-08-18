package testhelper

import (
	"os"
	"path/filepath"
)

const (
	TestDataDir = "../testhelper/test_data"
)

func CreateTestFile() ([]byte, string, error) {
	file, err := os.CreateTemp("", "test-file-*.txt")
	if err != nil {
		return []byte{}, "", err
	}
	defer file.Close()

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
