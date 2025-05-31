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
	"github.com/stretchr/testify/require"
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
		tempDir := t.TempDir()

		// create a sample file in temp dir
		err := os.WriteFile(filepath.Join(tempDir, "test"), []byte("test"), 0644)
		require.NoError(t, err)

		// archive tempDir
		resp, err := ArchiveDir(tempDir, nil)
		t.Cleanup(func() {
			_ = os.Remove(resp.ArchivePath)
		})

		assert.Equal(t, resp.TotalFiles, 1)
		assert.Equal(t, resp.TotalDirs, 1)
		assert.Equal(t, resp.SuccessFiles, 1)
		assert.NoError(t, err)

		// check archive path exists
		_, err = os.Stat(resp.ArchivePath)
		assert.NoError(t, err)
	})

	t.Run("Invalid Directory", func(t *testing.T) {
		// Create a sample dir in temp
		tempDir := "/tmp/does-not-exists"

		// archive tempDir
		resp, err := ArchiveDir(tempDir, nil)
		t.Cleanup(func() {
			_ = os.Remove(resp.ArchivePath)
		})

		assert.Empty(t, resp.TotalDirs)
		assert.Empty(t, resp.TotalFiles)
		assert.Empty(t, resp.SuccessFiles)
		assert.Error(t, err)
	})

	t.Run("Exclude Patterns", func(t *testing.T) {
		// Create a temporary directory for testing
		testDir := t.TempDir()

		// Create a few files and directories in the test directory
		err := os.MkdirAll(filepath.Join(testDir, "subdir"), os.ModePerm)
		require.NoError(t, err)

		_, _, err = testhelper.CreateTestFile(testDir, "")
		require.NoError(t, err)

		_, _, err = testhelper.CreateTestFile(testDir, "file-*.log")
		require.NoError(t, err)

		// Define exclusion patterns as regular expressions
		excludePatterns := []*regexp.Regexp{
			regexp.MustCompile(`\.log$`),   // Exclude files ending with .log
			regexp.MustCompile(`^subdir$`), // Exclude the "subdir" directory
		}

		// Archive the directory while excluding files and directories based on patterns
		resp, err := ArchiveDir(testDir, excludePatterns)
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = os.Remove(resp.ArchivePath)
		})

		// Check the contents of the created ZIP file
		zipReader, err := zip.OpenReader(resp.ArchivePath)
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = zipReader.Close()
		})

		// Verify that the ZIP file only contains the expected files
		expectedFilePattern := regexp.MustCompile("^test-file.*.txt$")
		for _, file := range zipReader.File {
			assert.Regexp(t, expectedFilePattern, file.Name, "Unexpected file found in ZIP archive")
		}

		// Check the counts of total files, total directories, and successfully archived files
		expectedTotalFiles := 1 // Only "file*.txt" is expected to be archived
		expectedTotalDirs := 1  // The root directory itself and the "subdir" directory
		expectedSuccessFiles := 1

		assert.Equal(t, expectedTotalFiles, resp.TotalFiles, "Total files mismatch")
		assert.Equal(t, expectedTotalDirs, resp.TotalDirs, "Total directories mismatch")
		assert.Equal(t, expectedSuccessFiles, resp.SuccessFiles, "Successfully archived files mismatch")
	})
}

func TestReadFileBytes(t *testing.T) {
	t.Run("Valid File", func(t *testing.T) {
		content, path, err := testhelper.CreateTestFile(t.TempDir(), "")
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
		_, absPath, err := testhelper.CreateTestFile(t.TempDir(), "")

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
		_, absPath, err := testhelper.CreateTestFile(t.TempDir(), "")

		assert.NoError(t, err)

		expectedSHA256 := "2172154e8979de165445a17dd2bdcba6408df06de67d042a6ae6781a1461e076"

		calculatedSHA256, err := CalculateFileSHA256(absPath)

		assert.NoError(t, err)
		assert.Equal(t, expectedSHA256, calculatedSHA256)
	})

	t.Run("Invalid File", func(t *testing.T) {
		_, absPath, err := testhelper.CreateTestFile(t.TempDir(), "")

		assert.NoError(t, err)

		expectedSHA256 := "daed58c831385cdebbb45785b1d5e2c5b2d0769a83896affa720bb32a325b5c6"

		calculatedSHA256, err := CalculateFileSHA256(absPath)

		assert.NoError(t, err)
		assert.NotEqual(t, expectedSHA256, calculatedSHA256)
	})
}

func TestValidateFileSHA256(t *testing.T) {
	t.Run("Valid SHA256", func(t *testing.T) {
		_, absPath, err := testhelper.CreateTestFile(t.TempDir(), "")

		assert.NoError(t, err)

		expectedSHA256 := "2172154e8979de165445a17dd2bdcba6408df06de67d042a6ae6781a1461e076"

		err = ValidateFileSHA256(absPath, expectedSHA256)
		assert.NoError(t, err)
	})

	t.Run("Invalid SHA256", func(t *testing.T) {
		_, absPath, err := testhelper.CreateTestFile(t.TempDir(), "")

		assert.NoError(t, err)

		expectedSHA256 := "daed58c831385cdebbb45785b1d5e2c5b2d0769a83896affa720bb32a325b5c6"

		err = ValidateFileSHA256(absPath, expectedSHA256)
		assert.Error(t, err)
	})
}

func TestDownloadFile(t *testing.T) {
	t.Run("Valid Download", func(t *testing.T) {
		// Create a test file
		_, absPath, err := testhelper.CreateTestFile(t.TempDir(), "")
		assert.NoError(t, err)

		// Create a mock HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Serve a test file
			http.ServeFile(w, r, absPath)
		}))
		defer server.Close()

		downloadFilePath := filepath.Join(os.TempDir(), "test-file.txt")
		t.Cleanup(func() {
			_ = os.Remove(downloadFilePath)
		})

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
		t.Cleanup(func() {
			_ = os.Remove(downloadFilePath)
		})

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
		content, filePath, err := testhelper.CreateTestFile(t.TempDir(), "")
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
		t.Cleanup(func() {
			_ = file.Close()
			_ = os.Remove(file.Name())
		})

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
		content, filePath1, err := testhelper.CreateTestFile(t.TempDir(), "")
		assert.NoError(t, err)

		filePath2, err := os.CreateTemp("", "test-file")
		assert.NoError(t, err)
		t.Cleanup(func() {
			_ = filePath2.Close()
			_ = os.Remove(filePath2.Name())
		})

		_, err = filePath2.Write(content)
		assert.NoError(t, err)

		// Call the FilesSameContent function
		same, err := FilesSameContent(filePath1, filePath2.Name())
		assert.NoError(t, err)
		assert.True(t, same)
	})

	t.Run("Different Content", func(t *testing.T) {
		// Create two test files with different content
		_, filePath1, err := testhelper.CreateTestFile(t.TempDir(), "")
		assert.NoError(t, err)

		_, filePath2, err := testhelper.CreateTestFile(t.TempDir(), "different-content")
		assert.NoError(t, err)

		// write random data to the second file
		file, err := os.OpenFile(filePath2, os.O_WRONLY, os.ModePerm)
		assert.NoError(t, err)
		_, err = file.Write([]byte("random data"))
		assert.NoError(t, err)
		_ = file.Close()

		// Call the FilesSameContent function
		same, err := FilesSameContent(filePath1, filePath2)
		assert.NoError(t, err)
		assert.False(t, same)
	})

	t.Run("Non-existent File", func(t *testing.T) {
		// Create a test file
		_, filePath, err := testhelper.CreateTestFile(t.TempDir(), "")
		assert.NoError(t, err)

		// Call the FilesSameContent function with a non-existent file path
		same, err := FilesSameContent(filePath, "/tmp/non-existent-file.txt")
		assert.Error(t, err)
		assert.False(t, same)
	})
}
