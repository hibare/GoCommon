package logger

import (
	stdlog "log"
	"os"
	"strings"

	"github.com/hibare/GoCommon/v2/pkg/slice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	LogLevelError = "ERROR"
	LogLevelWarn  = "WARN"
	LogLevelInfo  = "INFO"
	LogLevelDebug = "DEBUG"
)

const (
	LogModePretty = "PRETTY"
	LogModeJSON   = "JSON"
)

var (
	LogLevels = []string{LogLevelError, LogLevelWarn, LogLevelInfo, LogLevelDebug}
	LogModes  = []string{LogModePretty, LogModeJSON}
)

const (
	DefaultLoggerLevel = LogLevelInfo
	DefaultLoggerMode  = LogModePretty
)

func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	stdlog.SetFlags(0)
	stdlog.SetOutput(log.Logger)
	SetLoggingLevel(DefaultLoggerLevel)
	SetLoggingMode(DefaultLoggerMode)
}

func SetLoggingLevel(level string) {
	level = strings.ToUpper(level)

	switch level {
	case LogLevelError:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case LogLevelWarn:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case LogLevelInfo:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case LogLevelDebug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func GetLoggingLevel() string {
	return strings.ToUpper(zerolog.GlobalLevel().String())
}

func SetLoggingMode(mode string) {
	mode = strings.ToUpper(mode)

	switch mode {
	case LogModePretty:
		log.Logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			NoColor:    true,
			TimeFormat: "2006/01/02 03:04PM",
		}).With().Timestamp().Caller().Logger()
	case LogModeJSON:
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Caller().Logger()
	}
}

func IsValidLogLevel(level string) bool {
	return slice.SliceContains(strings.ToUpper(level), LogLevels)
}

func IsValidLogMode(mode string) bool {
	return slice.SliceContains(strings.ToUpper(mode), LogModes)
}