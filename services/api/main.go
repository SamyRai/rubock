package main

import (
	"helios/api/internal/platform"
	"helios/pkg/bootstrap"
	"helios/pkg/database"
	"helios/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	// --- Initialize Dependencies ---
	log := logger.New()

	// Load environment variables from .env file for local development.
	// In production, environment variables are set directly.
	if err := godotenv.Load(); err != nil {
		log.Info().Msg("No .env file found, using environment variables")
	}

	// Connect to NATS with resilient retry logic.
	natsConn, err := bootstrap.ConnectNATS(log)
	if err != nil {
		log.Fatal().Err(err).Msg("FATAL: Could not connect to NATS")
	}
	defer natsConn.Close()

	// Initialize database connection with resilient retry logic.
	db, err := database.NewDB(log)
	if err != nil {
		log.Fatal().Err(err).Msg("FATAL: Could not connect to the database")
	}
	defer db.Close()

	// --- Create and Run Application ---
	app := platform.NewApp(log, natsConn, db)
	app.Run()
}