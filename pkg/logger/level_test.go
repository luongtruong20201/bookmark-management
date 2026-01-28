package logger

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestSetLogLevel(t *testing.T) {
	originalLevel := zerolog.GlobalLevel()
	originalEnv := os.Getenv("LOG_LEVEL")
	defer func() {
		zerolog.SetGlobalLevel(originalLevel)
		if originalEnv != "" {
			os.Setenv("LOG_LEVEL", originalEnv)
		} else {
			os.Unsetenv("LOG_LEVEL")
		}
	}()

	testCases := []struct {
		name          string
		envValue      string
		expectedLevel zerolog.Level
		shouldSetEnv  bool
	}{
		{
			name:          "success - set to debug level",
			envValue:      "debug",
			expectedLevel: zerolog.DebugLevel,
			shouldSetEnv:  true,
		},
		{
			name:          "success - set to info level",
			envValue:      "info",
			expectedLevel: zerolog.InfoLevel,
			shouldSetEnv:  true,
		},
		{
			name:          "success - set to warn level",
			envValue:      "warn",
			expectedLevel: zerolog.WarnLevel,
			shouldSetEnv:  true,
		},
		{
			name:          "success - set to error level",
			envValue:      "error",
			expectedLevel: zerolog.ErrorLevel,
			shouldSetEnv:  true,
		},
		{
			name:          "success - set to fatal level",
			envValue:      "fatal",
			expectedLevel: zerolog.FatalLevel,
			shouldSetEnv:  true,
		},
		{
			name:          "success - set to panic level",
			envValue:      "panic",
			expectedLevel: zerolog.PanicLevel,
			shouldSetEnv:  true,
		},
		{
			name:          "success - set to trace level",
			envValue:      "trace",
			expectedLevel: zerolog.TraceLevel,
			shouldSetEnv:  true,
		},
		{
			name:          "success - invalid level defaults to info",
			envValue:      "invalid",
			expectedLevel: zerolog.InfoLevel,
			shouldSetEnv:  true,
		},
		{
			name:          "success - empty env defaults to info",
			envValue:      "",
			expectedLevel: zerolog.InfoLevel,
			shouldSetEnv:  true,
		},
		{
			name:          "success - no env variable defaults to info",
			envValue:      "",
			expectedLevel: zerolog.InfoLevel,
			shouldSetEnv:  false,
		},
		{
			name:          "success - NoLevel defaults to info",
			envValue:      "nolevel",
			expectedLevel: zerolog.InfoLevel,
			shouldSetEnv:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.shouldSetEnv {
				if tc.envValue != "" {
					os.Setenv("LOG_LEVEL", tc.envValue)
				} else {
					os.Unsetenv("LOG_LEVEL")
				}
			} else {
				os.Unsetenv("LOG_LEVEL")
			}

			SetLogLevel()
			assert.Equal(t, tc.expectedLevel, zerolog.GlobalLevel())
		})
	}
}

func TestSetLogLevel_CaseInsensitive(t *testing.T) {
	originalLevel := zerolog.GlobalLevel()
	originalEnv := os.Getenv("LOG_LEVEL")
	defer func() {
		zerolog.SetGlobalLevel(originalLevel)
		if originalEnv != "" {
			os.Setenv("LOG_LEVEL", originalEnv)
		} else {
			os.Unsetenv("LOG_LEVEL")
		}
	}()

	testCases := []struct {
		name          string
		envValue      string
		expectedLevel zerolog.Level
	}{
		{
			name:          "success - uppercase DEBUG",
			envValue:      "DEBUG",
			expectedLevel: zerolog.DebugLevel,
		},
		{
			name:          "success - mixed case DeBuG",
			envValue:      "DeBuG",
			expectedLevel: zerolog.DebugLevel,
		},
		{
			name:          "success - uppercase INFO",
			envValue:      "INFO",
			expectedLevel: zerolog.InfoLevel,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("LOG_LEVEL", tc.envValue)
			SetLogLevel()

			assert.Equal(t, tc.expectedLevel, zerolog.GlobalLevel())
		})
	}
}
