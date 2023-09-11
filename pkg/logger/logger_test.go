package logger_test

import (
	"testing"

	"github.com/hibare/GoCommon/v2/pkg/logger"
)

func TestInitLogger(t *testing.T) {
	// Test that InitLogger doesn't produce any errors
	logger.InitLogger()
}

func TestSetLoggingLevel(t *testing.T) {
	// Test setting each log level individually
	testCases := []struct {
		level string
	}{
		{logger.LogLevelError},
		{logger.LogLevelWarn},
		{logger.LogLevelInfo},
		{logger.LogLevelDebug},
	}

	for _, tc := range testCases {
		t.Run(tc.level, func(t *testing.T) {
			logger.SetLoggingLevel(tc.level)
			// Verify that the global log level is set as expected
			actual := logger.GetLoggingLevel()
			if actual != tc.level {
				t.Errorf("Expected log level '%s', but got '%s'", tc.level, actual)
			}
		})
	}
}

func TestIsValidLogLevel(t *testing.T) {
	// Test valid log levels
	testCases := []struct {
		level    string
		expected bool
	}{
		{logger.LogLevelError, true},
		{logger.LogLevelWarn, true},
		{logger.LogLevelInfo, true},
		{logger.LogLevelDebug, true},
		{"INVALID", false},
	}

	for _, tc := range testCases {
		t.Run(tc.level, func(t *testing.T) {
			actual := logger.IsValidLogLevel(tc.level)
			if actual != tc.expected {
				t.Errorf("Expected IsValidLogLevel('%s') to be %v, but got %v", tc.level, tc.expected, actual)
			}
		})
	}
}

func TestIsValidLogMode(t *testing.T) {
	// Test valid log modes
	testCases := []struct {
		mode     string
		expected bool
	}{
		{logger.LogModePretty, true},
		{logger.LogModeJSON, true},
		{"INVALID", false},
	}

	for _, tc := range testCases {
		t.Run(tc.mode, func(t *testing.T) {
			actual := logger.IsValidLogMode(tc.mode)
			if actual != tc.expected {
				t.Errorf("Expected IsValidLogMode('%s') to be %v, but got %v", tc.mode, tc.expected, actual)
			}
		})
	}
}
