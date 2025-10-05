// Package platform holds the central application container and other foundational
// components for the build-worker service.
package platform

import (
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Create a synchronous queue subscriber for manual acknowledgement.
	subject := events.SubjectDeploymentRequested
	sub, err := a.NATS.QueueSubscribeSync(subject, "build-workers")
	if err != nil {
		a.Logger.Fatal().Err(err).Str("subject", subject).Msg("FATAL: Could not create queue subscription")
	}

	a.Logger.Info().Str("subject", subject).Str("queue_group", "build-workers").Msg("Listening for events")

	// Start a goroutine for message processing.
	processingDone := make(chan struct{})
	go func() {
		defer close(processingDone)
		for {
			msg, err := sub.NextMsg(10 * time.Second)
			if err != nil {
				// ErrTimeout is expected when no messages are pending.
				if err == nats.ErrTimeout {
					continue
				}
				// ErrConnectionClosed or ErrBadSubscription indicate the subscription is done.
				if err == nats.ErrConnectionClosed || err == nats.ErrBadSubscription {
					a.Logger.Info().Msg("Subscription closed, stopping message processing.")
					return
				}
				a.Logger.Error().Err(err).Msg("Error receiving message from NATS")
				continue
			}
			w.HandleDeploymentRequest(msg)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the worker.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	a.Logger.Warn().Msg("Shutdown signal received, draining NATS subscription...")

	// Drain the subscription. This will wait for the processing goroutine to finish.
	if err := sub.Drain(); err != nil {
		a.Logger.Error().Err(err).Msg("Error draining NATS subscription")
	}

	// Wait for the processing goroutine to exit cleanly.
	<-processingDone

	a.Logger.Info().Msg("Worker exiting")
}