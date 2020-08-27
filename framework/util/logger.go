package util

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

const (
	// DebugLevel is a the string that sets the projects log level to 'debug'
	DebugLevel = "debug"
	// LogLevel is a the string that sets the projects log level to 'log'
	LogLevel = "log"
	// TraceLevel is a the string that sets the projects log level to 'trace'
	TraceLevel = "trace"
)

// ConfigureGlobalLogLevel will configure the global logger with the log level that will be used by the project
func ConfigureGlobalLogLevel(level string) {
	logLevel := ""
	switch level {
	case DebugLevel:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		logLevel = DebugLevel
	case TraceLevel:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		logLevel = TraceLevel
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		logLevel = LogLevel
	}

	logger := GetStandardLogger()
	logger.Info().
		Str("logLevel", logLevel).
		Msg(fmt.Sprintf("Global logger has been configured to log level '%s'.", logLevel))
}

// GetStandardLogger will retrieve the standard logger for the project
func GetStandardLogger() *zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()
	return &log
}

// ConfigureLoggerLogLevel will configure the logger level on the specified logger
func ConfigureLoggerLogLevel(logLevel zerolog.Level, logger *zerolog.Logger) *zerolog.Logger {
	previousLogLevel := logger.GetLevel().String()
	newLogger := logger.With().Logger().Level(logLevel)

	newLogger.Debug().
		Str("logLevel", logLevel.String()).
		Str("oldLogLevel", previousLogLevel).
		Msg(fmt.Sprintf("Logger has been configured to log level '%s'.", logLevel.String()))

	return &newLogger
}
