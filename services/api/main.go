package main

import (
	"os"
	"time"

	"helios/api/internal/platform"
	"helios/pkg/database"
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

	// Initialize database connection
	dbCfg := database.DBConfig{
		User:         "user",
		Password:     "password",
		Host:         "localhost",
		Port:         5432,
		DBName:       "helios",
		SSLMode:      "disable",
		MaxOpenConns: 25,
		MaxIdleConns: 25,
		MaxIdleTime:  15 * time.Minute,
	}
	db, err := database.NewDB(dbCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("FATAL: Could not connect to the database")
	}
	defer db.Close()
	log.Info().Msg("Successfully connected to the database")

	// --- Create and Run Application ---
	app := platform.NewApp(log, natsConn, db)
	app.Run()
}