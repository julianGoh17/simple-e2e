package util

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestConfigureGlobalLogLevel(t *testing.T) {
	values := []struct {
		loggerLevel      string
		expectedLogLevel zerolog.Level
	}{
		{
			"trace",
			zerolog.TraceLevel,
		},
		{
			"debug",
			zerolog.DebugLevel,
		},
		{
			"info",
			zerolog.InfoLevel,
		},
		{
			"random",
			zerolog.InfoLevel,
		},
	}

	for _, value := range values {
		ConfigureGlobalLogLevel(value.loggerLevel)
		assert.Equal(t, value.expectedLogLevel.String(), zerolog.GlobalLevel().String())
	}
}

func TestConfigureLoggerLogLevel(t *testing.T) {
	values := []struct {
		logLevel zerolog.Level
	}{
		{
			zerolog.TraceLevel,
		},
		{
			zerolog.DebugLevel,
		},
		{
			zerolog.InfoLevel,
		},
		{
			zerolog.InfoLevel,
		},
	}

	for _, value := range values {
		logger := GetStandardLogger().With().Logger().Level(zerolog.ErrorLevel)
		newLogger := ConfigureLoggerLogLevel(value.logLevel, &logger)
		assert.NotEqual(t, logger, newLogger)
		assert.Equal(t, newLogger.GetLevel(), value.logLevel)
	}
}
