package worker

import (
	"encoding/json"
	"fmt"
	"time"

	"helios/pkg/events"

	"github.com/go-playground/validator/v10"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

// NatsPublisher defines the interface for publishing messages to NATS.
type NatsPublisher interface {
	Publish(subject string, data []byte) error
}

// Worker holds dependencies for the message handler.
type Worker struct {
	NATS      NatsPublisher
	Logger    zerolog.Logger
	Validator *validator.Validate
}

// NewWorker creates a new Worker.
func NewWorker(nats NatsPublisher, logger zerolog.Logger) *Worker {
	return &Worker{
		NATS:      nats,
		Logger:    logger,
		Validator: validator.New(),
	}
}

// HandleDeploymentRequest processes incoming deployment request events.
func (w *Worker) HandleDeploymentRequest(m *nats.Msg) {
	var request events.DeploymentRequest
	if err := json.Unmarshal(m.Data, &request); err != nil {
		w.Logger.Error().Err(err).Msg("Could not unmarshal deployment request, terminating message")
		if err := m.Term(); err != nil {
			w.Logger.Error().Err(err).Msg("Failed to terminate NATS message")
		}
		return
	}

	// Validate the event payload
	if err := w.Validator.Struct(&request); err != nil {
		w.Logger.Error().Err(err).Msg("Invalid deployment request payload, terminating message")
		if err := m.Term(); err != nil {
			w.Logger.Error().Err(err).Msg("Failed to terminate NATS message")
		}
		return
	}

	log := w.Logger.With().Str("app_id", request.AppID).Logger()

	log.Info().Str("repo", request.GitRepository).Msg("Received deployment request")

	// Simulate the build process
	log.Info().Msg("Simulating build process...")
	time.Sleep(1 * time.Second) // Reduced for faster tests
	log.Info().Msg("Build simulation complete")

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
		log.Error().Err(err).Msg("Could not marshal build succeeded event, terminating message")
		if err := m.Term(); err != nil {
			log.Error().Err(err).Msg("Failed to terminate NATS message")
		}
		return
	}

	subject := events.SubjectBuildSucceeded
	if err := w.NATS.Publish(subject, eventData); err != nil {
		log.Error().Err(err).Str("subject", subject).Msg("Failed to publish to NATS, nakking message for redelivery")
		if err := m.Nak(); err != nil {
			log.Error().Err(err).Msg("Failed to nak NATS message")
		}
		return
	}

	log.Info().Str("subject", subject).Msg("Successfully published event to NATS")

	// Acknowledge the message now that processing is complete
	if err := m.Ack(); err != nil {
		log.Error().Err(err).Msg("Failed to acknowledge NATS message")
	}
}