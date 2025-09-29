package file

import (
	"archive/zip"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hibare/GoCommon/v2/pkg/testhelper"
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

		require.Equal(t, 1, resp.TotalFiles)
		require.Equal(t, 1, resp.TotalDirs)
		require.Equal(t, 1, resp.SuccessFiles)
		require.NoError(t, err)

		// check archive path exists
		_, err = os.Stat(resp.ArchivePath)
		require.NoError(t, err)
	})

	t.Run("Invalid Directory", func(t *testing.T) {
		// Create a sample dir in temp
		tempDir := "/tmp/does-not-exists"

		// archive tempDir
		resp, err := ArchiveDir(tempDir, nil)
		t.Cleanup(func() {
			_ = os.Remove(resp.ArchivePath)
		})

		require.Empty(t, resp.TotalDirs)
		require.Empty(t, resp.TotalFiles)
		require.Empty(t, resp.SuccessFiles)
		require.Error(t, err)
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
			require.Regexp(t, expectedFilePattern, file.Name, "Unexpected file found in ZIP archive")
		}

		// Check the counts of total files, total directories, and successfully archived files
		expectedTotalFiles := 1 // Only "file*.txt" is expected to be archived
		expectedTotalDirs := 1  // The root directory itself and the "subdir" directory
		expectedSuccessFiles := 1

		require.Equal(t, expectedTotalFiles, resp.TotalFiles, "Total files mismatch")
		require.Equal(t, expectedTotalDirs, resp.TotalDirs, "Total directories mismatch")
		require.Equal(t, expectedSuccessFiles, resp.SuccessFiles, "Successfully archived files mismatch")
	})
}

func TestReadFileBytes(t *testing.T) {
	t.Run("Valid File", func(t *testing.T) {
		content, path, err := testhelper.CreateTestFile(t.TempDir(), "")
		require.NoError(t, err)

		readBytes, err := ReadFileBytes(path)
		require.NoError(t, err)
		require.Equal(t, content, readBytes)
	})

	t.Run("Non-existent File", func(t *testing.T) {
		_, err := ReadFileBytes("/tmp/non-exists-file.txt")
		require.Error(t, err)
	})
}

func TestReadFileLines(t *testing.T) {
	t.Run("Valid File", func(t *testing.T) {
		_, absPath, err := testhelper.CreateTestFile(t.TempDir(), "")

		require.NoError(t, err)

		lines, err := ReadFileLines(absPath)
		require.NoError(t, err)
		require.Len(t, lines, 2)
	})

	t.Run("Non-existent File", func(t *testing.T) {
		lines, err := ReadFileLines("some/random/path")
		require.Error(t, err)
		require.Nil(t, lines)
	})
}

func TestDownloadFile(t *testing.T) {
	t.Run("Valid Download", func(t *testing.T) {
		// Create a test file
		_, absPath, err := testhelper.CreateTestFile(t.TempDir(), "")
		require.NoError(t, err)

		// Create a mock HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Serve a test file
			http.ServeFile(w, r, absPath)
		}))
		defer server.Close()

		downloadFilePath := filepath.Join(t.TempDir(), "test-file.txt")
		t.Cleanup(func() {
			_ = os.Remove(downloadFilePath)
		})

		// Download the file using the download function
		err = DownloadFile(t.Context(), server.URL, downloadFilePath)
		require.NoError(t, err)

		lines, err := ReadFileLines(downloadFilePath)
		require.NoError(t, err)
		require.Len(t, lines, 2)
	})

	t.Run("Invalid Download", func(t *testing.T) {
		// Create a mock HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Serve a test file
			http.ServeFile(w, r, "some/random/path")
		}))
		defer server.Close()

		downloadFilePath := filepath.Join(t.TempDir(), "test-file.txt")
		t.Cleanup(func() {
			_ = os.Remove(downloadFilePath)
		})

		// Download the file using the download function
		err := DownloadFile(t.Context(), server.URL, downloadFilePath)
		require.Error(t, err)

		lines, err := ReadFileLines(downloadFilePath)
		require.Error(t, err)
		require.Nil(t, lines)
	})
}

func TestExtractFileFromTarGz(t *testing.T) {
	t.Run("Valid Extraction", func(t *testing.T) {
		archivePath := filepath.Join(testhelper.TestDataDir, "sample.tar.gz")
		targetFilename := "sample.txt"

		extractedPath, err := ExtractFileFromTarGz(archivePath, targetFilename)
		require.NoError(t, err)
		require.Contains(t, extractedPath, targetFilename)
	})

	t.Run("Invalid Extraction", func(t *testing.T) {
		archivePath := filepath.Join(testhelper.TestDataDir, "sample.tar.gz")
		targetFilename := "sample1.txt"

		_, err := ExtractFileFromTarGz(archivePath, targetFilename)
		require.Error(t, err)
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
		require.Equal(t, expectedFiles, files)
		require.Equal(t, expectedDirs, dirs)
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
		require.Equal(t, expectedFiles, files)
		require.Equal(t, expectedDirs, dirs)
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
		require.Equal(t, expectedFiles, files)
		require.Equal(t, expectedDirs, dirs)
	})
}

func TestFilesSameContent(t *testing.T) {
	t.Run("Same Content", func(t *testing.T) {
		// Create two test files with the same content
		content, filePath1, err := testhelper.CreateTestFile(t.TempDir(), "")
		require.NoError(t, err)

		filePath2, err := os.CreateTemp(t.TempDir(), "test-file")
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = filePath2.Close()
			_ = os.Remove(filePath2.Name())
		})

		_, err = filePath2.Write(content)
		require.NoError(t, err)

		// Call the FilesSameContent function
		same, err := IsFilesSameContent(filePath1, filePath2.Name())
		require.NoError(t, err)
		require.True(t, same)
	})

	t.Run("Different Content", func(t *testing.T) {
		// Create two test files with different content
		_, filePath1, err := testhelper.CreateTestFile(t.TempDir(), "")
		require.NoError(t, err)

		_, filePath2, err := testhelper.CreateTestFile(t.TempDir(), "different-content")
		require.NoError(t, err)

		// write random data to the second file
		file, err := os.OpenFile(filePath2, os.O_WRONLY, os.ModePerm)
		require.NoError(t, err)
		_, err = file.WriteString("random data")
		require.NoError(t, err)
		_ = file.Close()

		// Call the FilesSameContent function
		same, err := IsFilesSameContent(filePath1, filePath2)
		require.NoError(t, err)
		require.False(t, same)
	})

	t.Run("Non-existent File", func(t *testing.T) {
		// Create a test file
		_, filePath, err := testhelper.CreateTestFile(t.TempDir(), "")
		require.NoError(t, err)

		// Call the FilesSameContent function with a non-existent file path
		same, err := IsFilesSameContent(filePath, "/tmp/non-existent-file.txt")
		require.Error(t, err)
		require.False(t, same)
	})
}
