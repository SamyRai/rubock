package main

import (
	"helios/build-worker/internal/platform"
	"helios/pkg/bootstrap"
	"helios/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	// --- Initialize Dependencies ---
	log := logger.New()

	// Load environment variables from .env file for local development.
	if err := godotenv.Load(); err != nil {
		log.Info().Msg("No .env file found, using environment variables")
	}

	// Connect to NATS with resilient retry logic.
	natsConn, err := bootstrap.ConnectNATS(log)
	if err != nil {
		log.Fatal().Err(err).Msg("FATAL: Could not connect to NATS")
	}
	defer natsConn.Close()

	// --- Create and Run Application ---
	app := platform.NewApp(log, natsConn)
	app.Run()
}