// Package logger provides a standardized zerolog instance for consistent logging
// across all services in the Helios platform.
package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// New creates and configures a new zerolog.Logger instance.
// It reads the `ENV` environment variable to determine whether to use a
// pretty, human-readable console logger (for "development") or a structured
// JSON logger (for "production" or any other value).
func New() zerolog.Logger {
	env := os.Getenv("ENV")
	if env == "development" {
		// Use a pretty, colorized console writer for local development.
		return log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	// Use a structured JSON logger in production.
	return zerolog.New(os.Stderr).With().Timestamp().Logger()
}