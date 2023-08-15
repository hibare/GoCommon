package file

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hibare/GoCommon/pkg/testhelper"
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
	archivePath, totalFiles, totalDirs, successFiles, err := ArchiveDir(tempDir)
	defer os.Remove(archivePath)

	assert.Equal(t, totalFiles, 1)
	assert.Equal(t, totalDirs, 1)
	assert.Equal(t, successFiles, 1)
	assert.NoError(t, err)

	// check archive path exists
	_, err = os.Stat(archivePath)
	assert.NoError(t, err)
}

func TestArchiveDirInvalidDir(t *testing.T) {
	// Create a sample dir in temp
	tempDir := "/tmp/does-not-exists"

	// archive tempDir
	archivePath, totalFiles, totalDirs, successFiles, err := ArchiveDir(tempDir)
	defer os.Remove(archivePath)

	assert.Empty(t, totalDirs)
	assert.Empty(t, totalFiles)
	assert.Empty(t, successFiles)
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

func TestReadFilePass(t *testing.T) {
	_, absPath, err := CreateTestFile()
	defer os.Remove(absPath)

	assert.NoError(t, err)

	lines, err := ReadFileLines(absPath)
	assert.NoError(t, err)
	assert.Len(t, lines, 2)
}

func TestReadFileFail(t *testing.T) {
	lines, err := ReadFileLines("some/random/path")
	assert.Error(t, err)
	assert.Nil(t, lines)
}

func TestCalculateFileSHA256Pass(t *testing.T) {
	_, absPath, err := CreateTestFile()
	defer os.Remove(absPath)

	assert.NoError(t, err)

	expectedSHA256 := "2172154e8979de165445a17dd2bdcba6408df06de67d042a6ae6781a1461e076"

	calculatedSHA256, err := CalculateFileSHA256(absPath)

	assert.NoError(t, err)
	assert.Equal(t, expectedSHA256, calculatedSHA256)
}

func TestCalculateFileSHA256Fail(t *testing.T) {
	_, absPath, err := CreateTestFile()
	defer os.Remove(absPath)

	assert.NoError(t, err)

	expectedSHA256 := "daed58c831385cdebbb45785b1d5e2c5b2d0769a83896affa720bb32a325b5c6"

	calculatedSHA256, err := CalculateFileSHA256(absPath)

	assert.NoError(t, err)
	assert.NotEqual(t, expectedSHA256, calculatedSHA256)
}

func TestValidateFileSHA256Pass(t *testing.T) {
	_, absPath, err := CreateTestFile()
	defer os.Remove(absPath)

	assert.NoError(t, err)

	expectedSHA256 := "2172154e8979de165445a17dd2bdcba6408df06de67d042a6ae6781a1461e076"

	err = ValidateFileSha256(absPath, expectedSHA256)
	assert.NoError(t, err)
}

func TestValidateFileSHA256Fail(t *testing.T) {
	_, absPath, err := CreateTestFile()
	defer os.Remove(absPath)

	assert.NoError(t, err)

	expectedSHA256 := "daed58c831385cdebbb45785b1d5e2c5b2d0769a83896affa720bb32a325b5c6"

	err = ValidateFileSha256(absPath, expectedSHA256)
	assert.Error(t, err)
}

func TestDownloadFilePass(t *testing.T) {
	// Create a test file
	_, absPath, err := CreateTestFile()
	assert.NoError(t, err)

	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve a test file
		http.ServeFile(w, r, absPath)
	}))
	defer server.Close()

	downloadFilePath := filepath.Join(os.TempDir(), "test-file.txt")
	defer os.Remove(downloadFilePath)

	// Download the file using the download function
	err = DownloadFile(server.URL, downloadFilePath)
	assert.NoError(t, err)

	lines, err := ReadFileLines(downloadFilePath)
	assert.NoError(t, err)
	assert.Len(t, lines, 2)

}

func TestDownloadFileFail(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve a test file
		http.ServeFile(w, r, "some/random/path")
	}))
	defer server.Close()

	downloadFilePath := filepath.Join(os.TempDir(), "test-file.txt")
	defer os.Remove(downloadFilePath)

	// Download the file using the download function
	err := DownloadFile(server.URL, downloadFilePath)
	assert.Error(t, err)

	lines, err := ReadFileLines(downloadFilePath)
	assert.Error(t, err)
	assert.Nil(t, lines)
}

func TestExtractFileFromTarGzPass(t *testing.T) {
	archivePath := filepath.Join(testhelper.TestDataDir, "sample.tar.gz")
	targetFilename := "sample.txt"
	extractedFilePath := filepath.Join(os.TempDir(), targetFilename)

	extractedPath, err := ExtractFileFromTarGz(archivePath, targetFilename)
	assert.NoError(t, err)
	assert.Equal(t, extractedFilePath, extractedPath)
}

func TestExtractFileFromTarGzFail(t *testing.T) {
	archivePath := filepath.Join(testhelper.TestDataDir, "sample.tar.gz")
	targetFilename := "sample1.txt"

	_, err := ExtractFileFromTarGz(archivePath, targetFilename)
	assert.Error(t, err)
}

func TestListFilesDirs(t *testing.T) {
	rootDir := "../testhelper"
	expectedFiles := []string{
		"../testhelper/test_data/sample.tar.gz",
		"../testhelper/testhelper.go",
	}
	expectedDirs := []string{
		"../testhelper",
		"../testhelper/test_data",
	}
	files, dirs := ListFilesDirs(rootDir, nil)
	assert.Equal(t, expectedFiles, files)
	assert.Equal(t, expectedDirs, dirs)
}

func TestListFilesDirsExcludeFiles(t *testing.T) {
	rootDir := "../testhelper"
	expectedFiles := []string{
		"../testhelper/test_data/sample.tar.gz",
	}
	expectedDirs := []string{
		"../testhelper",
		"../testhelper/test_data",
	}
	files, dirs := ListFilesDirs(rootDir, []*regexp.Regexp{regexp.MustCompile(".*.go")})
	assert.Equal(t, expectedFiles, files)
	assert.Equal(t, expectedDirs, dirs)
}

func TestListFilesDirsExcludeDirs(t *testing.T) {
	rootDir := "../testhelper"
	expectedFiles := []string{
		"../testhelper/testhelper.go",
	}
	expectedDirs := []string{
		"../testhelper",
	}
	files, dirs := ListFilesDirs(rootDir, []*regexp.Regexp{regexp.MustCompile("test_data")})
	assert.Equal(t, expectedFiles, files)
	assert.Equal(t, expectedDirs, dirs)
}
