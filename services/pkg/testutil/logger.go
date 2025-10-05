package testutil

import "github.com/rs/zerolog"

// NewTestLogger returns a zerolog.Logger instance that is silenced,
// which is useful for keeping test output clean.
func NewTestLogger() zerolog.Logger {
	return zerolog.Nop()
}