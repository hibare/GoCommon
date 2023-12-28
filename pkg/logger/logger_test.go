package logger

import (
	"log/slog"
	"os"
	"testing"

	"github.com/hibare/GoCommon/v2/pkg/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestGetHandler(t *testing.T) {
	// Test case 1: logMode is nil
	logLevel := "info"
	logMode := (*string)(nil)
	expectedTextHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})
	assert.Equal(t, expectedTextHandler, getHandler(&logLevel, logMode))

	// Test case 2: logMode is LogModePretty
	logLevel = "debug"
	logMode = testhelper.StringToPtr(LogModePretty)
	expectedTextHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})
	assert.Equal(t, expectedTextHandler, getHandler(&logLevel, logMode))

	// Test case 3: logMode is LogModeJSON
	logLevel = "error"
	logMode = testhelper.StringToPtr(LogModeJSON)
	expectedJSONHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelError,
	})
	assert.Equal(t, expectedJSONHandler, getHandler(&logLevel, logMode))

	// Test case 4: logMode is not LogModePretty or LogModeJSON
	logLevel = "warn"
	logMode = testhelper.StringToPtr("invalid")
	expectedJSONHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelWarn,
	})
	assert.Equal(t, expectedJSONHandler, getHandler(&logLevel, logMode))
}

func TestGetSlogLevelFromString(t *testing.T) {
	// Test case 1: level is nil, should return default logger level
	level := ""
	result := getSlogLevelFromString(&level)
	expected := slog.LevelInfo
	assert.Equal(t, expected, result)

	// Test case 2: level is "ERROR", should return slog.LevelError
	level = "ERROR"
	result = getSlogLevelFromString(&level)
	expected = slog.LevelError
	assert.Equal(t, expected, result)

	// Test case 3: level is "WARN", should return slog.LevelWarn
	level = "WARN"
	result = getSlogLevelFromString(&level)
	expected = slog.LevelWarn
	assert.Equal(t, expected, result)

	// Test case 4: level is "INFO", should return slog.LevelInfo
	level = "INFO"
	result = getSlogLevelFromString(&level)
	expected = slog.LevelInfo
	assert.Equal(t, expected, result)

	// Test case 5: level is "DEBUG", should return slog.LevelDebug
	level = "DEBUG"
	result = getSlogLevelFromString(&level)
	expected = slog.LevelDebug
	assert.Equal(t, expected, result)

	// Test case 6: level is "UNKNOWN", should return slog.LevelInfo
	level = "UNKNOWN"
	result = getSlogLevelFromString(&level)
	expected = slog.LevelInfo
	assert.Equal(t, expected, result)

}

func TestIsValidLogLevel(t *testing.T) {
	// Test valid log levels
	testCases := []struct {
		level    string
		expected bool
	}{
		{LogLevelError, true},
		{LogLevelWarn, true},
		{LogLevelInfo, true},
		{LogLevelDebug, true},
		{"INVALID", false},
	}

	for _, tc := range testCases {
		t.Run(tc.level, func(t *testing.T) {
			actual := IsValidLogLevel(tc.level)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestIsValidLogMode(t *testing.T) {
	// Test valid log modes
	testCases := []struct {
		mode     string
		expected bool
	}{
		{LogModePretty, true},
		{LogModeJSON, true},
		{"INVALID", false},
	}

	for _, tc := range testCases {
		t.Run(tc.mode, func(t *testing.T) {
			actual := IsValidLogMode(tc.mode)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
