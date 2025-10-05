package worker

import (
	"encoding/json"
	"time"

	"helios/pkg/events"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

// Worker holds dependencies for the message handler.
type Worker struct {
	Logger zerolog.Logger
}

// NewWorker creates a new Worker.
func NewWorker(logger zerolog.Logger) *Worker {
	return &Worker{Logger: logger}
}

// HandleBuildSucceeded processes incoming build succeeded events.
func (w *Worker) HandleBuildSucceeded(m *nats.Msg) {
	var event events.BuildSucceeded
	if err := json.Unmarshal(m.Data, &event); err != nil {
		w.Logger.Error().Err(err).Msg("Could not unmarshal build succeeded event")
		return
	}

	w.Logger.Info().
		Str("app_id", event.AppID).
		Str("image_uri", event.ImageURI).
		Msg("Received build succeeded event")

	// Simulate the deployment process
	w.Logger.Info().
		Str("app_id", event.AppID).
		Str("image_uri", event.ImageURI).
		Msg("Simulating deployment...")
	time.Sleep(1 * time.Second) // Reduced for faster tests
	w.Logger.Info().Str("app_id", event.AppID).Msg("Deployment simulation complete")

	// In a real implementation, we would publish a "deployment.succeeded" event here.
	w.Logger.Info().Str("app_id", event.AppID).Msg("End of workflow")
}