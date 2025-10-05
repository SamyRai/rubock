// Package bootstrap provides shared startup and initialization logic for Helios services.
package bootstrap

import (
	"math/rand"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"

	"helios/pkg/config"
)

// init seeds the random number generator for creating jitter in retry delays.
func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// NATSConfig holds the configuration for the NATS connection.
type NATSConfig struct {
	URL        string
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
}

// NewNATSConfig creates a NATS configuration from environment variables.
func NewNATSConfig() NATSConfig {
	return NATSConfig{
		URL:        config.Getenv("NATS_URL", nats.DefaultURL),
		MaxRetries: config.GetenvInt("NATS_MAX_RETRIES", 10),
		BaseDelay:  config.GetenvDuration("NATS_BASE_DELAY", 1*time.Second),
		MaxDelay:   config.GetenvDuration("NATS_MAX_DELAY", 30*time.Second),
	}
}

// ConnectNATS establishes a connection to the NATS server with a robust,
// exponential backoff retry mechanism.
func ConnectNATS(log zerolog.Logger) (*nats.Conn, error) {
	cfg := NewNATSConfig()
	log.Info().Msgf("Attempting to connect to NATS at %s", cfg.URL)

	var natsConn *nats.Conn
	var err error

	for i := 0; i < cfg.MaxRetries; i++ {
		natsConn, err = nats.Connect(cfg.URL)
		if err == nil {
			log.Info().Msgf("Successfully connected to NATS at %s", natsConn.ConnectedUrl())
			return natsConn, nil
		}

		delay := cfg.BaseDelay * time.Duration(1<<i)
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}
		// Add jitter to prevent thundering herd
		jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
		delay += jitter

		log.Warn().Err(err).Msgf("Failed to connect to NATS, retrying in %s...", delay)
		time.Sleep(delay)
	}

	return nil, err
}