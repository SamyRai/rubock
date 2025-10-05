package main

import (
	"helios/build-worker/internal/platform"
	"helios/pkg/bootstrap"
	"helios/pkg/logger"
)

func main() {
	// --- Initialize Dependencies ---
	log := logger.New()

	// Connect to NATS
	natsConn, err := bootstrap.ConnectNATS(log)
	if err != nil {
		log.Fatal().Err(err).Msg("FATAL: Could not connect to NATS")
	}
	defer natsConn.Close()

	// --- Create and Run Application ---
	app := platform.NewApp(log, natsConn)
	app.Run()
}