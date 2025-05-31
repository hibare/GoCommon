// Package file provides utilities for file and directory operations.
package file

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
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

// ArchiveDirResponse represents the result of archiving a directory.
type ArchiveDirResponse struct {
	ArchivePath  string
	TotalFiles   int
	TotalDirs    int
	SuccessFiles int
	FailedFiles  map[string]error
}

// ArchiveDir creates a zip archive of the specified directory, excluding files/dirs matching the exclude patterns.
func ArchiveDir(dirPath string, exclude []*regexp.Regexp) (ArchiveDirResponse, error) {
	dirPath = filepath.Clean(dirPath)
	dirName := filepath.Base(dirPath)
	zipName := fmt.Sprintf("%s.zip", dirName)
	zipPath := filepath.Join(os.TempDir(), zipName)

	zipFile, err := os.Create(zipPath)
	if err != nil {
		return ArchiveDirResponse{}, fmt.Errorf("failed to create zip file: %w", err)
	}
	defer func() {
		_ = zipFile.Close()
	}()

	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		_ = zipWriter.Close()
	}()

	totalFiles, totalDirs, successFiles := 0, 0, 0
	failedFiles := make(map[string]error)

	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk error at %s: %w", path, err)
		}

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

		if !info.Mode().IsRegular() {
			return nil
		}

		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			failedFiles[path] = fmt.Errorf("failed to get relative path: %w", err)
			return nil
		}

		zh, err := zipWriter.CreateHeader(&zip.FileHeader{
			Name:   relPath,
			Method: zip.Deflate,
		})
		if err != nil {
			failedFiles[path] = fmt.Errorf("failed to create zip header: %w", err)
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			failedFiles[path] = fmt.Errorf("failed to open file: %w", err)
			return nil
		}
		defer func() {
			_ = file.Close()
		}()

		_, err = io.Copy(zh, file)
		if err != nil {
			failedFiles[path] = fmt.Errorf("failed to copy file to zip: %w", err)
			return nil
		}

		successFiles++
		return nil
	})

	return ArchiveDirResponse{
		ArchivePath:  zipPath,
		TotalFiles:   totalFiles,
		TotalDirs:    totalDirs,
		SuccessFiles: successFiles,
		FailedFiles:  failedFiles,
	}, err
}

// ReadFileBytes reads the entire content of a file and returns it as a byte slice.
func ReadFileBytes(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

// ReadFileLines reads a file and returns its contents as a slice of strings, one per line.
func ReadFileLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan file: %w", err)
	}
	return lines, nil
}

// CalculateFileSHA256 calculates the SHA-256 checksum of a file and returns it as a hex string.
func CalculateFileSHA256(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// ValidateFileSHA256 checks if the SHA-256 checksum of a file matches the provided checksum string.
func ValidateFileSHA256(path string, sha256Str string) error {
	calculatedSha256, err := CalculateFileSHA256(path)
	if err != nil {
		return fmt.Errorf("failed to calculate file SHA256: %w", err)
	}

	if calculatedSha256 != sha256Str {
		return errors.ErrChecksumMismatch
	}
	return nil
}

// DownloadFile downloads a file from the given URL to the specified destination path.
func DownloadFile(url string, destination string) error {
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed to get url: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", response.Status)
	}

	out, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() {
		_ = out.Close()
	}()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		return fmt.Errorf("failed to copy response body: %w", err)
	}

	return nil
}

// ExtractFileFromTarGz extracts a specific file from a .tar.gz archive and returns the path to the extracted file.
func ExtractFileFromTarGz(archivePath, targetFilename string) (string, error) {
	var targetFilePath string

	file, err := os.Open(archivePath)
	if err != nil {
		return targetFilePath, fmt.Errorf("failed to open archive: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return targetFilePath, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer func() {
		_ = gzipReader.Close()
	}()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			return targetFilePath, os.ErrNotExist
		}
		if err != nil {
			return targetFilePath, fmt.Errorf("failed to read tar header: %w", err)
		}

		if strings.HasSuffix(header.Name, targetFilename) {
			targetFilePath = filepath.Join(os.TempDir(), targetFilename)
			targetFile, err := os.OpenFile(targetFilePath, os.O_CREATE|os.O_WRONLY, header.FileInfo().Mode())
			if err != nil {
				return targetFilePath, fmt.Errorf("failed to create target file: %w", err)
			}
			defer func() {
				_ = targetFile.Close()
			}()

			if _, err := io.CopyN(targetFile, tarReader, header.Size); err != nil {
				_ = os.Remove(targetFilePath)
				return targetFilePath, fmt.Errorf("failed to copy file from archive: %w", err)
			}
			break
		}
	}
	return targetFilePath, nil
}

// ListFilesDirs returns slices of file and directory paths under root, excluding those matching the exclude patterns.
func ListFilesDirs(root string, exclude []*regexp.Regexp) ([]string, []string) {
	var files []string
	var dirs []string

	readDir := func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walkdir error at %s: %w", path, err)
		}

		if d.IsDir() {
			if shouldExclude(d.Name(), exclude) {
				slog.Info("Skipping dir", "name", d.Name(), "path", path)
				return filepath.SkipDir
			}
			dirs = append(dirs, path)
		} else {
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

// GetHash computes the SHA-256 hash of a file and returns the raw bytes.
func GetHash(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, fmt.Errorf("failed to hash file: %w", err)
	}

	return hash.Sum(nil), nil
}

// IsFilesSameContent checks if two files have the same content by comparing their SHA-256 hashes.
func IsFilesSameContent(file1, file2 string) (bool, error) {
	hash1, err := GetHash(file1)
	if err != nil {
		return false, fmt.Errorf("failed to hash file1: %w", err)
	}

	hash2, err := GetHash(file2)
	if err != nil {
		return false, fmt.Errorf("failed to hash file2: %w", err)
	}

	return bytes.Equal(hash1, hash2), nil
}
