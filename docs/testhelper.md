# Testhelper Package Documentation

## Overview

The `testhelper` package provides utilities for test setup and teardown, including functions for creating test files and directories, and pointer helpers. It is designed to simplify writing tests for Go applications.

---

## Key Functions

- **CreateTestFile(dir, pattern) ([]byte, string, error)**: Creates a test file with sample content and returns its content and absolute path.
- **CreateTestDir(dir, pattern) (string, error)**: Creates a test directory and returns its path.
- **StringToPtr(s string) \*string**: Converts a string to a pointer.

---

## Example Usage

```go
import (
    "github.com/hibare/GoCommon/v2/pkg/testhelper"
)

content, path, err := testhelper.CreateTestFile("/tmp", "test-*.txt")
dir, err := testhelper.CreateTestDir("/tmp", "test-dir-")
ptr := testhelper.StringToPtr("hello")
```

---

## Notes

- Useful for writing unit and integration tests that require temporary files or directories.
