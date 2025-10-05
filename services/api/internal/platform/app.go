// Package platform holds the central application container and other foundational
// components for the API service.
package platform

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"helios/api/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

// App represents the central application container, holding all dependencies.
type App struct {
	Logger zerolog.Logger
	Router *chi.Mux
	NATS   *nats.Conn
}

// NewApp creates and configures a new application instance.
func NewApp(logger zerolog.Logger, natsConn *nats.Conn) *App {
	app := &App{
		Logger: logger,
		Router: chi.NewRouter(),
		NATS:   natsConn,
	}

	// Register routes
	app.registerRoutes()

	return app
}

// registerRoutes sets up the application's HTTP routes.
func (a *App) registerRoutes() {
	// The handlers now need access to the app's dependencies, which can be
	// passed via methods on the App struct or by passing the app itself.
	// For simplicity, we'll create handlers that have access to the app.
	apiHandlers := handlers.NewAPIHandlers(a.NATS, a.Logger)

	a.Router.Post("/projects", apiHandlers.CreateProjectHandler)
	a.Router.Post("/applications", apiHandlers.CreateApplicationHandler)
}

// Run starts the HTTP server and handles graceful shutdown.
func (a *App) Run() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: a.Router,
	}

	go func() {
		a.Logger.Info().Msgf("Starting API server on port %s...", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Logger.Fatal().Err(err).Msg("FATAL: Could not start server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	a.Logger.Warn().Msg("Shutdown signal received, starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		a.Logger.Fatal().Err(err).Msg("FATAL: Server forced to shutdown")
	}

	a.Logger.Info().Msg("Server exiting")
}