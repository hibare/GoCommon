# File Package Documentation

## Overview

The `file` package provides utilities for file and directory operations, including archiving, reading, hashing, downloading, and comparing files. It is designed to simplify common file-related tasks in Go applications.

---

## Key Types and Functions

- **ArchiveDir(dirPath, exclude) (ArchiveDirResponse, error)**: Creates a zip archive of a directory, excluding files/dirs by regex.
- **ReadFileBytes(path) ([]byte, error)**: Reads the entire content of a file as bytes.
- **ReadFileLines(path) ([]string, error)**: Reads a file and returns its contents as a slice of lines.
- **CalculateFileSHA256(path) (string, error)**: Calculates the SHA-256 checksum of a file.
- **ValidateFileSHA256(path, sha256Str) error**: Validates a file's SHA-256 checksum.
- **DownloadFile(url, destination) error**: Downloads a file from a URL to a destination path.
- **ExtractFileFromTarGz(archivePath, targetFilename) (string, error)**: Extracts a specific file from a .tar.gz archive.
- **ListFilesDirs(root, exclude) ([]string, []string)**: Lists files and directories under a root, excluding by regex.
- **GetHash(filePath) ([]byte, error)**: Computes the SHA-256 hash of a file.
- **IsFilesSameContent(file1, file2) (bool, error)**: Checks if two files have the same content by comparing hashes.

---

## Example Usage

```go
import (
    "github.com/hibare/GoCommon/v2/pkg/file"
)

resp, err := file.ArchiveDir("/path/to/dir", nil)
if err != nil {
    panic(err)
}
fmt.Println("Archive created at:", resp.ArchivePath)
```

---

## Notes

- Useful for backup, validation, and file management tasks.
- Integrates with the `errors` package for error handling.
