// Package logger provides a standardized zerolog instance for consistent,
// configurable logging across all services in the Helios platform.
package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"helios/pkg/config"
)

// New creates and configures a new zerolog.Logger instance.
// It reads the `ENV` environment variable to determine the output format:
// - "development": A pretty, human-readable console logger.
// - "production" (or any other value): A structured JSON logger.
//
// It also reads the `LOG_LEVEL` environment variable to set the logging level
// (e.g., "debug", "info", "warn", "error"). Defaults to "info".
func New() zerolog.Logger {
	// Set the logging level.
	logLevelStr := config.Getenv("LOG_LEVEL", "info")
	logLevel, err := zerolog.ParseLevel(strings.ToLower(logLevelStr))
	if err != nil {
		// Default to info level if parsing fails.
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	// Set the output format based on the environment.
	env := config.Getenv("ENV", "production")
	if env == "development" {
		// Use a pretty, colorized console writer for local development.
		return log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Caller().Logger()
	}

	// Use a structured JSON logger in production.
	// Add caller information to help trace log origins.
	return zerolog.New(os.Stderr).With().Timestamp().Caller().Logger()
}