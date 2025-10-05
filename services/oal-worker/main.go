package main

import (
	"os"
	"time"

	"helios/oal-worker/internal/platform"
	"helios/pkg/logger"

	"github.com/nats-io/nats.go"
)

func main() {
	// --- Initialize Dependencies ---
	log := logger.New()

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	var natsConn *nats.Conn
	var err error
	for i := 0; i < 5; i++ {
		natsConn, err = nats.Connect(natsURL)
		if err == nil {
			break
		}
		log.Warn().Err(err).Msgf("Failed to connect to NATS, retrying in %d seconds...", i+1)
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	if err != nil {
		log.Fatal().Err(err).Msgf("FATAL: Could not connect to NATS at %s", natsURL)
	}
	defer natsConn.Close()
	log.Info().Msgf("Successfully connected to NATS at %s", natsConn.ConnectedUrl())

	// --- Create and Run Application ---
	app := platform.NewApp(log, natsConn)
	app.Run()
}