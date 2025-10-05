package worker

import (
	"encoding/json"
	"fmt"
	"time"

	"helios/pkg/events"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

// NatsPublisher defines the interface for publishing messages to NATS.
type NatsPublisher interface {
	Publish(subject string, data []byte) error
}

// Worker holds dependencies for the message handler.
type Worker struct {
	NATS   NatsPublisher
	Logger zerolog.Logger
}

// NewWorker creates a new Worker.
func NewWorker(nats NatsPublisher, logger zerolog.Logger) *Worker {
	return &Worker{
		NATS:   nats,
		Logger: logger,
	}
}

// HandleDeploymentRequest processes incoming deployment request events.
func (w *Worker) HandleDeploymentRequest(m *nats.Msg) {
	var request events.DeploymentRequest
	if err := json.Unmarshal(m.Data, &request); err != nil {
		w.Logger.Error().Err(err).Msg("Could not unmarshal deployment request")
		return
	}

	w.Logger.Info().
		Str("app_id", request.AppID).
		Str("repo", request.GitRepository).
		Msg("Received deployment request")

	// Simulate the build process
	w.Logger.Info().Str("app_id", request.AppID).Msg("Simulating build process...")
	time.Sleep(1 * time.Second) // Reduced for faster tests
	w.Logger.Info().Str("app_id", request.AppID).Msg("Build simulation complete")

	// Publish Build Succeeded Event
	commitSHA := "a1b2c3d4e5f6" // Placeholder
	imageURI := fmt.Sprintf("registry.helios.internal/%s:%s", request.AppID, commitSHA)

	event := events.BuildSucceeded{
		AppID:        request.AppID,
		ImageURI:     imageURI,
		GitCommitSHA: commitSHA,
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		w.Logger.Error().Err(err).Str("app_id", request.AppID).Msg("Could not marshal build succeeded event")
		return
	}

	subject := events.SubjectBuildSucceeded
	if err := w.NATS.Publish(subject, eventData); err != nil {
		w.Logger.Error().Err(err).Str("subject", subject).Str("app_id", request.AppID).Msg("Failed to publish to NATS")
		return
	}

	w.Logger.Info().
		Str("subject", subject).
		Str("app_id", request.AppID).
		Msg("Successfully published event to NATS")
}