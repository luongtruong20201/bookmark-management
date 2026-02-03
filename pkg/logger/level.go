package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// SetLogLevel configures the global logging level based on the LOG_LEVEL environment variable.
// It parses the environment variable value and sets it as the global zerolog log level.
// If LOG_LEVEL is not set, invalid, or results in NoLevel, it defaults to InfoLevel.
// This function should be called during application initialization to configure logging.
func SetLogLevel() {
	level, err := zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil || level == zerolog.NoLevel {
		level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(level)
}
