// Package platform holds the central application container and other foundational
// components for the build-worker service.
package platform

import (
	"os"
	"os/signal"
	"syscall"

	"helios/build-worker/internal/worker"
	"helios/pkg/events"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

// App represents the central application container, holding all dependencies.
type App struct {
	Logger zerolog.Logger
	NATS   *nats.Conn
}

// NewApp creates and configures a new application instance.
func NewApp(logger zerolog.Logger, natsConn *nats.Conn) *App {
	return &App{
		Logger: logger,
		NATS:   natsConn,
	}
}

// Run starts the worker, subscribes to NATS, and handles graceful shutdown.
func (a *App) Run() {
	w := worker.NewWorker(a.NATS, a.Logger)

	// Subscribe to deployment requests
	subject := events.SubjectDeploymentRequested
	sub, err := a.NATS.QueueSubscribe(subject, "build-workers", w.HandleDeploymentRequest)
	if err != nil {
		a.Logger.Fatal().Err(err).Str("subject", subject).Msg("FATAL: Could not subscribe to NATS subject")
	}

	a.Logger.Info().Str("subject", subject).Str("queue_group", "build-workers").Msg("Listening for events")

	// Wait for interrupt signal to gracefully shut down the worker
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	a.Logger.Warn().Msg("Shutdown signal received, draining NATS subscription...")

	// Drain the subscription, processing any remaining messages.
	if err := sub.Drain(); err != nil {
		a.Logger.Error().Err(err).Msg("Error draining NATS subscription")
	}

	a.Logger.Info().Msg("Worker exiting")
}