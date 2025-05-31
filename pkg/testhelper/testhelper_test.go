package testhelper

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTestFile(t *testing.T) {
	dir := os.TempDir()
	pattern := "test-file-*.txt"

	// Test case 1: When pattern is empty
	t.Run("PatternEmpty", func(t *testing.T) {
		content, absPath, err := CreateTestFile(dir, "")
		t.Cleanup(func() { _ = os.Remove(absPath) })

		assert.NoError(t, err)
		assert.Equal(t, []byte("This is a test file.\nIt contains some sample content."), content)
		assert.FileExists(t, absPath)
	})

	// Test case 2: When pattern is not empty
	t.Run("PatternNotEmpty", func(t *testing.T) {
		content, absPath, err := CreateTestFile(dir, pattern)
		t.Cleanup(func() { _ = os.Remove(absPath) })

		assert.NoError(t, err)
		assert.Equal(t, []byte("This is a test file.\nIt contains some sample content."), content)
		assert.FileExists(t, absPath)
	})

	// Test case 3: When base dir does not exists
	t.Run("InvalidBaseDir", func(t *testing.T) {
		_, _, err := CreateTestFile("/invalid/dir", pattern)
		assert.Error(t, err)
	})
}

func TestCreateTestDir(t *testing.T) {
	tempDir := os.TempDir()

	// Test case 1: Testing when dir and pattern are empty
	// Expect the function to create a temporary directory with the default pattern and return the directory path without error
	dir1, err1 := CreateTestDir("", "")
	assert.NoError(t, err1)
	assert.NotEmpty(t, dir1)
	t.Cleanup(func() { _ = os.RemoveAll(dir1) })

	// Test case 2: Testing when dir is empty and pattern is not empty
	// Expect the function to create a temporary directory with the provided pattern and return the directory path without error
	dir2, err2 := CreateTestDir("", "custom-pattern-")
	assert.NoError(t, err2)
	assert.NotEmpty(t, dir2)
	assert.Contains(t, dir2, "custom-pattern-")
	t.Cleanup(func() { _ = os.RemoveAll(dir2) })

	// Test case 3: Testing when dir is not empty and pattern is empty
	// Expect the function to create a temporary directory with the default pattern inside the provided directory and return the directory path without error
	dir3, err3 := CreateTestDir(tempDir, "")
	assert.NoError(t, err3)
	assert.NotEmpty(t, dir3)
	t.Cleanup(func() { _ = os.RemoveAll(dir3) })

	// Test case 4: Testing when dir and pattern are not empty
	// Expect the function to create a temporary directory with the provided pattern inside the provided directory and return the directory path without error
	dir4, err4 := CreateTestDir(tempDir, "custom-pattern-")
	assert.NoError(t, err4)
	assert.NotEmpty(t, dir4)
	assert.Contains(t, dir4, "custom-pattern-")
	t.Cleanup(func() { _ = os.RemoveAll(dir4) })
}

func TestStringToPtr(t *testing.T) {
	// Test case 1: Testing with non-empty string
	str1 := "Hello"
	ptr1 := StringToPtr(str1)
	assert.Equal(t, str1, *ptr1, "Expected %s, but got %s", str1, *ptr1)

	// Test case 2: Testing with empty string
	str2 := ""
	ptr2 := StringToPtr(str2)
	assert.Empty(t, *ptr2, "Expected empty string, but got %s", *ptr2)
}
