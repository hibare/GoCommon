package file

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"log"

	"github.com/hibare/GoCommon/v2/pkg/errors"
)

func ArchiveDir(dirPath string) (string, int, int, int, error) {
	dirPath = filepath.Clean(dirPath)
	dirName := filepath.Base(dirPath)
	zipName := fmt.Sprintf("%s.zip", dirName)
	zipPath := filepath.Join(os.TempDir(), zipName)
	totalFiles, totalDirs, successFiles := 0, 0, 0

	// Create a temporary file to hold the zip archive
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return zipPath, totalFiles, totalDirs, successFiles, err
	}
	defer zipFile.Close()

	// Create a new zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			totalDirs++
			return nil
		}

		totalFiles++

		info, err := d.Info()
		if err != nil {
			log.Printf("Failed to get file info: %v", err)
			return nil
		}

		if !info.Mode().IsRegular() {
			log.Printf("%s is not a regular file", path)
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			log.Printf("Failed to create header: %v", err)
			return nil
		}

		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			log.Printf("Failed to get relative path: %v", err)
			return nil
		}
		header.Name = filepath.ToSlash(filepath.Join(relPath))

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			log.Printf("Failed to create header: %v", err)
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			log.Printf("Failed to open file: %v", err)
			return nil
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			file.Close()
			log.Printf("Failed to write file to archive: %v", err)
			return nil
		}
		file.Close()

		successFiles++

		return nil
	})

	log.Printf("Created archive '%s' for directory '%s'", zipPath, dirPath)
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

			if _, err := io.Copy(targetFile, tarReader); err != nil {
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
			for _, e := range exclude {
				if e.MatchString(d.Name()) {
					return filepath.SkipDir
				}
			}

			dirs = append(dirs, path)
		} else {
			// Check if file matches any of the exclude patterns
			for _, e := range exclude {
				if e.MatchString(d.Name()) {
					return nil
				}
			}

			files = append(files, path)
		}

		return nil
	}

	_ = filepath.WalkDir(root, readDir)

	return files, dirs
}
