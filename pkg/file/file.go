package file

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"log"
)

func ArchiveDir(dirPath string) (string, error) {
	dirPath = filepath.Clean(dirPath)
	dirName := filepath.Base(dirPath)
	zipName := fmt.Sprintf("%s.zip", dirName)
	zipPath := filepath.Join(os.TempDir(), zipName)

	// Create a temporary file to hold the zip archive
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return zipPath, err
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
			return nil
		}

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

		return nil
	})

	log.Printf("Created archive '%s' for directory '%s'", zipPath, dirPath)
	return zipPath, err
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
