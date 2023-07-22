package file

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArchiveDir(t *testing.T) {
	// Create a sample dir in temp
	tempDir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)

	// create a sample file in temp dir
	_, err = os.CreateTemp(tempDir, "test")
	assert.NoError(t, err)

	// archive tempDir
	archivePath, err := ArchiveDir(tempDir)
	assert.NoError(t, err)

	// check archive path exists
	_, err = os.Stat(archivePath)
	assert.NoError(t, err)
}

func TestArchiveDirInvalidDir(t *testing.T) {
	// Create a sample dir in temp
	tempDir := "/tmp/does-not-exists"

	// archive tempDir
	_, err := ArchiveDir(tempDir)
	assert.Error(t, err)

}
