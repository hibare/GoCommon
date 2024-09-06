package file

import (
	"archive/zip"
	"crypto/sha256"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hibare/GoCommon/v2/pkg/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestShouldExclude(t *testing.T) {
	tests := []struct {
		name     string
		patterns []*regexp.Regexp
		input    string
		expected bool
	}{
		{
			name:     "Exact Match",
			patterns: []*regexp.Regexp{regexp.MustCompile(`^file\.txt$`)},
			input:    "file.txt",
			expected: true,
		},
		{
			name:     "Partial Match",
			patterns: []*regexp.Regexp{regexp.MustCompile(`^.*\.tmp$`)},
			input:    "document.tmp",
			expected: true,
		},
		{
			name:     "No Match",
			patterns: []*regexp.Regexp{regexp.MustCompile(`^pattern1$`), regexp.MustCompile(`^pattern2$`)},
			input:    "file.txt",
			expected: false,
		},
		{
			name:     "Invalid Pattern",
			patterns: []*regexp.Regexp{regexp.MustCompile(`invalid`)},
			input:    "file.txt",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := shouldExclude(test.input, test.patterns)
			if actual != test.expected {
				t.Errorf("Expected shouldExclude('%s', patterns) to be %v, but got %v", test.input, test.expected, actual)
			}
		})
	}
}

func TestArchiveDir(t *testing.T) {
	t.Run("Valid Directory", func(t *testing.T) {
		// Create a sample dir in temp
		tempDir, err := os.MkdirTemp("", "test")
		defer os.RemoveAll(tempDir)

		assert.NoError(t, err)

		// create a sample file in temp dir
		_, err = os.CreateTemp(tempDir, "test")
		assert.NoError(t, err)

		// archive tempDir
		archivePath, totalFiles, totalDirs, successFiles, err := ArchiveDir(tempDir, nil)
		defer os.Remove(archivePath)

		assert.Equal(t, totalFiles, 1)
		assert.Equal(t, totalDirs, 1)
		assert.Equal(t, successFiles, 1)
		assert.NoError(t, err)

		// check archive path exists
		_, err = os.Stat(archivePath)
		assert.NoError(t, err)
	})

	t.Run("Invalid Directory", func(t *testing.T) {
		// Create a sample dir in temp
		tempDir := "/tmp/does-not-exists"

		// archive tempDir
		archivePath, totalFiles, totalDirs, successFiles, err := ArchiveDir(tempDir, nil)
		defer os.Remove(archivePath)

		assert.Empty(t, totalDirs)
		assert.Empty(t, totalFiles)
		assert.Empty(t, successFiles)
		assert.Error(t, err)
	})

	t.Run("Exclude Patterns", func(t *testing.T) {
		// Create a temporary directory for testing
		testDir, err := os.MkdirTemp("", "test-archive-dir")
		assert.NoError(t, err)
		defer os.RemoveAll(testDir)

		// Create a few files and directories in the test directory
		err = os.MkdirAll(filepath.Join(testDir, "subdir"), os.ModePerm)
		assert.NoError(t, err)

		_, _, err = testhelper.CreateTestFile(testDir, "")
		assert.NoError(t, err)

		_, _, err = testhelper.CreateTestFile(testDir, "file-*.log")
		assert.NoError(t, err)

		// Define exclusion patterns as regular expressions
		excludePatterns := []*regexp.Regexp{
			regexp.MustCompile(`\.log$`),   // Exclude files ending with .log
			regexp.MustCompile(`^subdir$`), // Exclude the "subdir" directory
		}

		// Archive the directory while excluding files and directories based on patterns
		zipPath, totalFiles, totalDirs, successFiles, err := ArchiveDir(testDir, excludePatterns)
		assert.NoError(t, err)
		defer os.Remove(zipPath)

		// Check the contents of the created ZIP file
		zipReader, err := zip.OpenReader(zipPath)
		assert.NoError(t, err)
		defer zipReader.Close()

		// Verify that the ZIP file only contains the expected files
		expectedFilePattern := regexp.MustCompile("^test-file.*.txt$")
		for _, file := range zipReader.File {
			assert.Regexp(t, expectedFilePattern, file.Name, "Unexpected file found in ZIP archive")
		}

		// Check the counts of total files, total directories, and successfully archived files
		expectedTotalFiles := 1 // Only "file*.txt" is expected to be archived
		expectedTotalDirs := 1  // The root directory itself and the "subdir" directory
		expectedSuccessFiles := 1

		assert.Equal(t, expectedTotalFiles, totalFiles, "Total files mismatch")
		assert.Equal(t, expectedTotalDirs, totalDirs, "Total directories mismatch")
		assert.Equal(t, expectedSuccessFiles, successFiles, "Successfully archived files mismatch")
	})
}

func TestReadFileBytes(t *testing.T) {
	t.Run("Valid File", func(t *testing.T) {
		content, path, err := testhelper.CreateTestFile("", "")
		defer os.Remove(path)
		assert.NoError(t, err)

		readBytes, err := ReadFileBytes(path)
		assert.NoError(t, err)
		assert.Equal(t, content, readBytes)
	})

	t.Run("Non-existent File", func(t *testing.T) {
		_, err := ReadFileBytes("/tmp/non-exists-file.txt")
		assert.Error(t, err)
	})
}

func TestReadFileLines(t *testing.T) {
	t.Run("Valid File", func(t *testing.T) {
		_, absPath, err := testhelper.CreateTestFile("", "")
		defer os.Remove(absPath)

		assert.NoError(t, err)

		lines, err := ReadFileLines(absPath)
		assert.NoError(t, err)
		assert.Len(t, lines, 2)
	})

	t.Run("Non-existent File", func(t *testing.T) {
		lines, err := ReadFileLines("some/random/path")
		assert.Error(t, err)
		assert.Nil(t, lines)
	})
}

func TestCalculateFileSHA256(t *testing.T) {
	t.Run("Valid File", func(t *testing.T) {
		_, absPath, err := testhelper.CreateTestFile("", "")
		defer os.Remove(absPath)

		assert.NoError(t, err)

		expectedSHA256 := "2172154e8979de165445a17dd2bdcba6408df06de67d042a6ae6781a1461e076"

		calculatedSHA256, err := CalculateFileSHA256(absPath)

		assert.NoError(t, err)
		assert.Equal(t, expectedSHA256, calculatedSHA256)
	})

	t.Run("Invalid File", func(t *testing.T) {
		_, absPath, err := testhelper.CreateTestFile("", "")
		defer os.Remove(absPath)

		assert.NoError(t, err)

		expectedSHA256 := "daed58c831385cdebbb45785b1d5e2c5b2d0769a83896affa720bb32a325b5c6"

		calculatedSHA256, err := CalculateFileSHA256(absPath)

		assert.NoError(t, err)
		assert.NotEqual(t, expectedSHA256, calculatedSHA256)
	})
}

func TestValidateFileSHA256(t *testing.T) {
	t.Run("Valid SHA256", func(t *testing.T) {
		_, absPath, err := testhelper.CreateTestFile("", "")
		defer os.Remove(absPath)

		assert.NoError(t, err)

		expectedSHA256 := "2172154e8979de165445a17dd2bdcba6408df06de67d042a6ae6781a1461e076"

		err = ValidateFileSha256(absPath, expectedSHA256)
		assert.NoError(t, err)
	})

	t.Run("Invalid SHA256", func(t *testing.T) {
		_, absPath, err := testhelper.CreateTestFile("", "")
		defer os.Remove(absPath)

		assert.NoError(t, err)

		expectedSHA256 := "daed58c831385cdebbb45785b1d5e2c5b2d0769a83896affa720bb32a325b5c6"

		err = ValidateFileSha256(absPath, expectedSHA256)
		assert.Error(t, err)
	})
}

func TestDownloadFile(t *testing.T) {
	t.Run("Valid Download", func(t *testing.T) {
		// Create a test file
		_, absPath, err := testhelper.CreateTestFile("", "")
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
	})

	t.Run("Invalid Download", func(t *testing.T) {
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
	})
}

func TestExtractFileFromTarGz(t *testing.T) {
	t.Run("Valid Extraction", func(t *testing.T) {
		archivePath := filepath.Join(testhelper.TestDataDir, "sample.tar.gz")
		targetFilename := "sample.txt"
		extractedFilePath := filepath.Join(os.TempDir(), targetFilename)

		extractedPath, err := ExtractFileFromTarGz(archivePath, targetFilename)
		assert.NoError(t, err)
		assert.Equal(t, extractedFilePath, extractedPath)
	})

	t.Run("Invalid Extraction", func(t *testing.T) {
		archivePath := filepath.Join(testhelper.TestDataDir, "sample.tar.gz")
		targetFilename := "sample1.txt"

		_, err := ExtractFileFromTarGz(archivePath, targetFilename)
		assert.Error(t, err)
	})
}

func TestListFilesDirs(t *testing.T) {
	t.Run("No Exclusions", func(t *testing.T) {
		rootDir := "../testhelper"
		expectedFiles := []string{
			"../testhelper/test_data/sample.tar.gz",
			"../testhelper/testhelper.go",
			"../testhelper/testhelper_test.go",
		}
		expectedDirs := []string{
			"../testhelper",
			"../testhelper/test_data",
		}
		files, dirs := ListFilesDirs(rootDir, nil)
		assert.Equal(t, expectedFiles, files)
		assert.Equal(t, expectedDirs, dirs)
	})

	t.Run("Exclude Files", func(t *testing.T) {
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
	})

	t.Run("Exclude Dirs", func(t *testing.T) {
		rootDir := "../testhelper"
		expectedFiles := []string{
			"../testhelper/testhelper.go",
			"../testhelper/testhelper_test.go",
		}
		expectedDirs := []string{
			"../testhelper",
		}
		files, dirs := ListFilesDirs(rootDir, []*regexp.Regexp{regexp.MustCompile("test_data")})
		assert.Equal(t, expectedFiles, files)
		assert.Equal(t, expectedDirs, dirs)
	})
}

func TestFileHash(t *testing.T) {
	t.Run("Valid File", func(t *testing.T) {
		// Create a test file with known content
		content, filePath, err := testhelper.CreateTestFile("", "")
		defer os.Remove(filePath)
		assert.NoError(t, err)

		// Calculate the expected hash using the same content
		expectedHash := sha256.Sum256(content)

		// Call the FileHash function
		actualHash, err := FileHash(filePath)
		assert.NoError(t, err)

		// Compare the actual hash with the expected hash
		assert.Equal(t, expectedHash[:], actualHash)
	})

	t.Run("Non-existent File", func(t *testing.T) {
		// Call the FileHash function with a non-existent file path
		_, err := FileHash("/tmp/non-existent-file.txt")
		assert.Error(t, err)
	})

	t.Run("Empty File", func(t *testing.T) {
		// Create an empty test file
		file, err := os.CreateTemp("", "empty-file")
		assert.NoError(t, err)
		defer os.Remove(file.Name())
		defer file.Close()

		// Call the FileHash function
		actualHash, err := FileHash(file.Name())
		assert.NoError(t, err)

		// Calculate the expected hash for an empty file
		expectedHash := sha256.Sum256(nil)

		// Compare the actual hash with the expected hash
		assert.Equal(t, expectedHash[:], actualHash)
	})
}

func TestFilesSameContent(t *testing.T) {
	t.Run("Same Content", func(t *testing.T) {
		// Create two test files with the same content
		content, filePath1, err := testhelper.CreateTestFile("", "")
		defer os.Remove(filePath1)
		assert.NoError(t, err)

		filePath2, err := os.CreateTemp("", "test-file")
		assert.NoError(t, err)
		defer os.Remove(filePath2.Name())
		defer filePath2.Close()

		_, err = filePath2.Write(content)
		assert.NoError(t, err)

		// Call the FilesSameContent function
		same, err := FilesSameContent(filePath1, filePath2.Name())
		assert.NoError(t, err)
		assert.True(t, same)
	})

	t.Run("Different Content", func(t *testing.T) {
		// Create two test files with different content
		_, filePath1, err := testhelper.CreateTestFile("", "")
		defer os.Remove(filePath1)
		assert.NoError(t, err)

		_, filePath2, err := testhelper.CreateTestFile("", "different-content")
		defer os.Remove(filePath2)
		assert.NoError(t, err)

		// write random data to the second file
		file, err := os.OpenFile(filePath2, os.O_WRONLY, os.ModePerm)
		assert.NoError(t, err)
		file.Write([]byte("random data"))
		file.Close()

		// Call the FilesSameContent function
		same, err := FilesSameContent(filePath1, filePath2)
		assert.NoError(t, err)
		assert.False(t, same)
	})

	t.Run("Non-existent File", func(t *testing.T) {
		// Create a test file
		_, filePath, err := testhelper.CreateTestFile("", "")
		defer os.Remove(filePath)
		assert.NoError(t, err)

		// Call the FilesSameContent function with a non-existent file path
		same, err := FilesSameContent(filePath, "/tmp/non-existent-file.txt")
		assert.Error(t, err)
		assert.False(t, same)
	})
}
