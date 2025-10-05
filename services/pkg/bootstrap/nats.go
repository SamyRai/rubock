// Package bootstrap provides shared startup and initialization logic for Helios services.
package bootstrap

import (
	"math/rand"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

// getenv returns the value of an environment variable or a fallback value if the variable is not set.
func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// ConnectNATS establishes a connection to the NATS server with a robust,
// exponential backoff retry mechanism.
func ConnectNATS(log zerolog.Logger) (*nats.Conn, error) {
	natsURL := getenv("NATS_URL", nats.DefaultURL)
	log.Info().Msgf("Attempting to connect to NATS at %s", natsURL)

	var natsConn *nats.Conn
	var err error

	// Exponential backoff settings
	maxRetries := 10
	baseDelay := 1 * time.Second
	maxDelay := 30 * time.Second

	for i := 0; i < maxRetries; i++ {
		natsConn, err = nats.Connect(natsURL)
		if err == nil {
			log.Info().Msgf("Successfully connected to NATS at %s", natsConn.ConnectedUrl())
			return natsConn, nil
		}

		delay := baseDelay * time.Duration(1<<i)
		if delay > maxDelay {
			delay = maxDelay
		}
		// Add jitter to prevent thundering herd
		jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
		delay = delay + jitter

		log.Warn().Err(err).Msgf("Failed to connect to NATS, retrying in %s...", delay)
		time.Sleep(delay)
	}

	return nil, err
}