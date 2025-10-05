// Package database provides a robust and resilient database connection handler
// with built-in retry logic and standardized configuration.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"net/url"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
	"github.com/rs/zerolog"

	"helios/pkg/config"
)

// init seeds the random number generator for creating jitter in retry delays.
func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// DBConfig holds the configuration for the database connection, loaded from environment variables.
type DBConfig struct {
	User         string
	Password     string
	Host         string
	Port         string
	DBName       string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
	MaxRetries   int
	BaseDelay    time.Duration
	MaxDelay     time.Duration
}

// NewDBConfig creates a database configuration from environment variables.
func NewDBConfig() DBConfig {
	return DBConfig{
		User:         config.Getenv("DB_USER", "postgres"),
		Password:     config.Getenv("DB_PASSWORD", "postgres"),
		Host:         config.Getenv("DB_HOST", "localhost"),
		Port:         config.Getenv("DB_PORT", "5432"),
		DBName:       config.Getenv("DB_NAME", "helios"),
		SSLMode:      config.Getenv("DB_SSLMODE", "disable"),
		MaxOpenConns: config.GetenvInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns: config.GetenvInt("DB_MAX_IDLE_CONNS", 25),
		MaxIdleTime:  config.GetenvDuration("DB_MAX_IDLE_TIME", 15*time.Minute),
		MaxRetries:   config.GetenvInt("DB_MAX_RETRIES", 5),
		BaseDelay:    config.GetenvDuration("DB_BASE_DELAY", 1*time.Second),
		MaxDelay:     config.GetenvDuration("DB_MAX_DELAY", 30*time.Second),
	}
}

// NewDB creates a new database connection pool with a resilient retry mechanism.
func NewDB(log zerolog.Logger) (*sql.DB, error) {
	cfg := NewDBConfig()

	// Use net/url.URL to safely construct the DSN.
	dsn := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Path:   cfg.DBName,
	}
	q := dsn.Query()
	q.Set("sslmode", cfg.SSLMode)
	dsn.RawQuery = q.Encode()

	connStr := dsn.String()
	log.Info().Msgf("Attempting to connect to database at %s:%s", cfg.Host, cfg.Port)

	var db *sql.DB
	var err error

	for i := 0; i < cfg.MaxRetries; i++ {
		db, err = sql.Open("pgx", connStr)
		if err != nil {
			return nil, fmt.Errorf("failed to open database connection: %w", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err = db.PingContext(ctx); err == nil {
			log.Info().Msg("Successfully connected to the database")
			db.SetMaxOpenConns(cfg.MaxOpenConns)
			db.SetMaxIdleConns(cfg.MaxIdleConns)
			db.SetConnMaxIdleTime(cfg.MaxIdleTime)
			return db, nil
		}
		db.Close() // Close the connection if ping fails

		delay := cfg.BaseDelay * time.Duration(1<<i)
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}
		// Add jitter to prevent thundering herd
		jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
		delay += jitter

		log.Warn().Err(err).Msgf("Failed to connect to database, retrying in %s...", delay)
		time.Sleep(delay)
	}

	return nil, fmt.Errorf("failed to connect to the database after %d retries: %w", cfg.MaxRetries, err)
}