package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestArchiveDir(t *testing.T) {
	// Create a sample dir in temp
	tempDir, err := os.MkdirTemp("", "test")
	defer os.RemoveAll(tempDir)

	assert.NoError(t, err)

	// create a sample file in temp dir
	_, err = os.CreateTemp(tempDir, "test")
	assert.NoError(t, err)

	// archive tempDir
	archivePath, err := ArchiveDir(tempDir)
	defer os.Remove(archivePath)
	assert.NoError(t, err)

	// check archive path exists
	_, err = os.Stat(archivePath)
	assert.NoError(t, err)
}

func TestArchiveDirInvalidDir(t *testing.T) {
	// Create a sample dir in temp
	tempDir := "/tmp/does-not-exists"

	// archive tempDir
	archivePath, err := ArchiveDir(tempDir)
	defer os.Remove(archivePath)
	assert.Error(t, err)

}

func TestReadFileBytes(t *testing.T) {
	content, path, err := CreateTestFile()
	defer os.Remove(path)
	assert.NoError(t, err)

	readBytes, err := ReadFileBytes(path)
	assert.NoError(t, err)
	assert.Equal(t, content, readBytes)
}

func TestReadFileBytesNoFile(t *testing.T) {
	_, err := ReadFileBytes("/tmp/non-exists-file.txt")
	assert.Error(t, err)
}
