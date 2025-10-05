package testutil

import (
	"io"

	"github.com/rs/zerolog"
)

// NewTestLogger returns a zerolog.Logger instance that is silenced,
// which is useful for keeping test output clean.
func NewTestLogger() zerolog.Logger {
	return zerolog.Nop()
}

// NewTestLoggerWithOutput returns a zerolog.Logger instance that writes to the
// provided io.Writer. This is useful for capturing log output during tests.
func NewTestLoggerWithOutput(w io.Writer) zerolog.Logger {
	return zerolog.New(w)
}