# Logger Package Documentation

## Overview

The `logger` package provides logging utilities for Go applications, built on top of Go's `slog` package. It supports multiple log levels and output modes, and is designed for easy initialization and integration.

---

## Key Types and Constants

- **LogLevelError, LogLevelWarn, LogLevelInfo, LogLevelDebug**: Log level constants.
- **LogModePretty, LogModeJSON**: Log output mode constants.
- **InitDefaultLogger()**: Initializes the logger with default settings.
- **InitLogger(logLevel, logMode)**: Initializes the logger with the specified log level and mode.
- **IsValidLogLevel(level string) bool**: Checks if a log level is valid.
- **IsValidLogMode(mode string) bool**: Checks if a log mode is valid.

---

## Example Usage

```go
import (
    "github.com/hibare/GoCommon/v2/pkg/logger"
)

logger.InitDefaultLogger() // Pretty, info-level by default
// or
level := "DEBUG"
mode := "JSON"
logger.InitLogger(&level, &mode)
```

---

## Notes

- Uses Go's `slog` for structured logging.
- Supports both pretty (text) and JSON output modes.
