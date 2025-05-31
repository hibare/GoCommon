// Package logger provides logging utilities for the application.
package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/hibare/GoCommon/v2/pkg/slice"
)

const (
	// LogLevelError is the error log level.
	LogLevelError = "ERROR"

	// LogLevelWarn is the warn log level.
	LogLevelWarn = "WARN"

	// LogLevelInfo is the info log level.
	LogLevelInfo = "INFO"

	// LogLevelDebug is the debug log level.
	LogLevelDebug = "DEBUG"
)

const (
	// LogModePretty is the pretty log mode.
	LogModePretty = "PRETTY"

	// LogModeJSON is the JSON log mode.
	LogModeJSON = "JSON"
)

var (
	LogLevels          = []string{LogLevelError, LogLevelWarn, LogLevelInfo, LogLevelDebug}
	LogModes           = []string{LogModePretty, LogModeJSON}
	DefaultLoggerLevel = LogLevelInfo
	DefaultLoggerMode  = LogModePretty
)

// InitDefaultLogger initializes the default logger.
func InitDefaultLogger() {
	InitLogger(nil, nil)
}

// InitLogger initializes the logger with the given log level and mode.
func InitLogger(logLevel, logMode *string) {
	handler := getHandler(logLevel, logMode)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func getHandler(logLevel, logMode *string) slog.Handler {
	level := getSlogLevelFromString(logLevel)

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}

	if logMode == nil {
		logMode = &DefaultLoggerMode
	}

	switch strings.ToUpper(*logMode) {
	case LogModePretty:
		return slog.NewTextHandler(os.Stdout, opts)
	case LogModeJSON:
		return slog.NewJSONHandler(os.Stdout, opts)
	default:
		return slog.NewJSONHandler(os.Stdout, opts)
	}
}

func getSlogLevelFromString(level *string) slog.Level {
	if level == nil {
		level = &DefaultLoggerLevel
	}

	switch strings.ToUpper(*level) {
	case LogLevelError:
		return slog.LevelError
	case LogLevelWarn:
		return slog.LevelWarn
	case LogLevelInfo:
		return slog.LevelInfo
	case LogLevelDebug:
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}

// IsValidLogLevel checks if the provided log level is valid.
func IsValidLogLevel(level string) bool {
	return slice.SliceContains(strings.ToUpper(level), LogLevels)
}

// IsValidLogMode checks if the provided log mode is valid.
func IsValidLogMode(mode string) bool {
	return slice.SliceContains(strings.ToUpper(mode), LogModes)
}
