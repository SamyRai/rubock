package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"helios/api/internal/platform"
	"helios/pkg/bootstrap"
	"helios/pkg/database"
	"helios/pkg/logger"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

// Config holds all configuration for the service.
type Config struct {
	DB database.DBConfig
}

// getenv returns the value of an environment variable or a fallback value.
func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// loadConfig loads configuration from environment variables.
func loadConfig(log zerolog.Logger) (*Config, error) {
	dbPort, err := strconv.ParseUint(getenv("DB_PORT", "5432"), 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT value: %w", err)
	}

	maxOpenConns, err := strconv.Atoi(getenv("DB_MAX_OPEN_CONNS", "25"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_OPEN_CONNS value: %w", err)
	}

	maxIdleConns, err := strconv.Atoi(getenv("DB_MAX_IDLE_CONNS", "25"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_IDLE_CONNS value: %w", err)
	}

	maxIdleTime, err := time.ParseDuration(getenv("DB_MAX_IDLE_TIME", "15m"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_IDLE_TIME value: %w", err)
	}

	cfg := &Config{
		DB: database.DBConfig{
			User:         getenv("DB_USER", "user"),
			Password:     getenv("DB_PASSWORD", "password"),
			Host:         getenv("DB_HOST", "localhost"),
			Port:         uint16(dbPort),
			DBName:       getenv("DB_NAME", "helios"),
			SSLMode:      getenv("DB_SSLMODE", "disable"),
			MaxOpenConns: maxOpenConns,
			MaxIdleConns: maxIdleConns,
			MaxIdleTime:  maxIdleTime,
		},
	}

	return cfg, nil
}

func main() {
	// --- Initialize Dependencies ---
	log := logger.New()

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Info().Msg("No .env file found, using environment variables")
	}

	// Load configuration
	cfg, err := loadConfig(log)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Connect to NATS
	natsConn, err := bootstrap.ConnectNATS(log)
	if err != nil {
		log.Fatal().Err(err).Msg("FATAL: Could not connect to NATS")
	}
	defer natsConn.Close()

	// Initialize database connection
	db, err := database.NewDB(cfg.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("FATAL: Could not connect to the database")
	}
	defer db.Close()
	log.Info().Msg("Successfully connected to the database")

	// --- Create and Run Application ---
	app := platform.NewApp(log, natsConn, db)
	app.Run()
}