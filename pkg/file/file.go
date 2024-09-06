package file

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hibare/GoCommon/v2/pkg/errors"
)

func shouldExclude(name string, exclude []*regexp.Regexp) bool {
	for _, e := range exclude {
		if e.MatchString(name) {
			return true
		}
	}
	return false
}

func ArchiveDir(dirPath string, exclude []*regexp.Regexp) (string, int, int, int, error) {
	// Clean the input directory path
	dirPath = filepath.Clean(dirPath)

	// Extract directory name from the path
	dirName := filepath.Base(dirPath)

	// Generate the zip file name
	zipName := fmt.Sprintf("%s.zip", dirName)

	// Create the full path for the zip file in the temporary directory
	zipPath := filepath.Join(os.TempDir(), zipName)

	// Create the zip file
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return zipPath, 0, 0, 0, err
	}
	defer zipFile.Close()

	// Create a zip writer for the zip file
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Initialize counters
	totalFiles, totalDirs, successFiles := 0, 0, 0

	// Recursively add files to the zip archive
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			if shouldExclude(info.Name(), exclude) {
				slog.Info("Skipping dir", "name", info.Name(), "path", path)
				return filepath.SkipDir
			}
			totalDirs++
			return nil
		}

		if shouldExclude(info.Name(), exclude) {
			slog.Info("Skipping file", "name", info.Name(), "path", path)
			return nil
		}

		totalFiles++

		// Check if the path is a regular file
		if !info.Mode().IsRegular() {
			return nil
		}

		// Get the relative path of the file
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		// Create a header for the file in the zip archive
		zh, err := zipWriter.CreateHeader(&zip.FileHeader{
			Name:   relPath,
			Method: zip.Deflate,
		})
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		// Copy the file content to the zip archive
		_, err = io.Copy(zh, file)
		if err != nil {
			return err
		}

		file.Close()
		successFiles++
		return nil
	})

	return zipPath, totalFiles, totalDirs, successFiles, err
}

func ReadFileBytes(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get the file size
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()

	// Read the file content into a buffer
	data := make([]byte, fileSize)
	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func ReadFileLines(path string) ([]string, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func CalculateFileSHA256(path string) (string, error) {
	file, err := os.Open(path)

	if err != nil {
		return "", err
	}

	defer file.Close()

	h := sha256.New()

	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func ValidateFileSha256(path string, sha256Str string) error {
	calculatedSha256, err := CalculateFileSHA256(path)

	if err != nil {
		return err
	}

	if calculatedSha256 != sha256Str {
		return errors.ErrChecksumMismatch
	}
	return nil
}

func DownloadFile(url string, destination string) error {
	response, err := http.Get(url)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", response.Status)
	}

	out, err := os.Create(destination)

	if err != nil {
		return err
	}

	_, err = io.Copy(out, response.Body)

	defer out.Close()

	if err != nil {
		return err
	}

	return nil
}

func ExtractFileFromTarGz(archivePath, targetFilename string) (string, error) {
	var targetFilePath string

	file, err := os.Open(archivePath)
	if err != nil {
		return targetFilePath, err
	}
	defer file.Close()

	// Create a gzip reader for the file
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return targetFilePath, err
	}
	defer gzipReader.Close()

	// Create a tar reader for the gzip reader
	tarReader := tar.NewReader(gzipReader)

	// Find the target file in the archive

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			return targetFilePath, os.ErrNotExist
		}
		if err != nil {
			return targetFilePath, err
		}

		if strings.HasSuffix(header.Name, targetFilename) {

			targetFilePath = filepath.Join(os.TempDir(), targetFilename)

			// Create the target file and copy the content of the file from the archive
			targetFile, err := os.OpenFile(targetFilePath, os.O_CREATE|os.O_WRONLY, header.FileInfo().Mode())
			if err != nil {
				return targetFilePath, err
			}

			if err != nil {
				return targetFilePath, err
			}

			if _, err := io.CopyN(targetFile, tarReader, header.Size); err != nil {
				targetFile.Close()
				os.Remove(targetFilename)
				return targetFilePath, err
			}
			targetFile.Close()
			break
		}
	}
	return targetFilePath, nil
}

func ListFilesDirs(root string, exclude []*regexp.Regexp) ([]string, []string) {
	var files []string
	var dirs []string

	readDir := func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// Check if directory matches any of the exclude patterns
			if shouldExclude(d.Name(), exclude) {
				slog.Info("Skipping dir", "name", d.Name(), "path", path)
				return filepath.SkipDir
			}

			dirs = append(dirs, path)
		} else {
			// Check if file matches any of the exclude patterns
			if shouldExclude(d.Name(), exclude) {
				slog.Info("Skipping file", "name", d.Name(), "path", path)
				return nil
			}

			files = append(files, path)
		}

		return nil
	}

	_ = filepath.WalkDir(root, readDir)

	return files, dirs
}

// FileHash computes the SHA-256 hash of a file
func FileHash(filePath string) ([]byte, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a new SHA-256 hash
	hash := sha256.New()

	// Copy the file contents into the hash
	if _, err := io.Copy(hash, file); err != nil {
		return nil, err
	}

	// Return the computed hash
	return hash.Sum(nil), nil
}

// FilesSameContent checks if two files have the same content by comparing their hashes
func FilesSameContent(file1, file2 string) (bool, error) {
	// Compute the hash of the first file
	hash1, err := FileHash(file1)
	if err != nil {
		return false, err
	}

	// Compute the hash of the second file
	hash2, err := FileHash(file2)
	if err != nil {
		return false, err
	}

	// Compare the hashes
	return bytes.Equal(hash1, hash2), nil
}
