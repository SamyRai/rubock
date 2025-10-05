// Package config provides standardized helper functions for reading configuration
// from the environment.
package config

import (
	"os"
	"strconv"
	"time"
)

// Getenv reads an environment variable with a fallback value.
func Getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetenvInt reads an integer environment variable with a fallback value.
func GetenvInt(key string, fallback int) int {
	s := Getenv(key, "")
	if s == "" {
		return fallback
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return v
}

// GetenvDuration reads a time.Duration environment variable with a fallback value.
func GetenvDuration(key string, fallback time.Duration) time.Duration {
	s := Getenv(key, "")
	if s == "" {
		return fallback
	}
	v, err := time.ParseDuration(s)
	if err != nil {
		return fallback
	}
	return v
}