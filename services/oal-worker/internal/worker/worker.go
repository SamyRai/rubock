package worker

import (
	"encoding/json"
	"time"

	"helios/pkg/events"

	"github.com/go-playground/validator/v10"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

// natsMsg defines the interface for a NATS message, allowing for easier testing.
type natsMsg interface {
	GetData() []byte
	Ack() error
	Nak() error
	Term() error
}

// natsMsgAdapter adapts a *nats.Msg to the natsMsg interface.
type natsMsgAdapter struct {
	msg *nats.Msg
}

func (a *natsMsgAdapter) GetData() []byte {
	return a.msg.Data
}

func (a *natsMsgAdapter) Ack() error {
	return a.msg.Ack()
}

func (a *natsMsgAdapter) Nak() error {
	return a.msg.Nak()
}

func (a *natsMsgAdapter) Term() error {
	return a.msg.Term()
}

// Worker holds dependencies for the message handler.
type Worker struct {
	Logger    zerolog.Logger
	Validator *validator.Validate
}

// NewWorker creates a new Worker.
func NewWorker(logger zerolog.Logger) *Worker {
	return &Worker{
		Logger:    logger,
		Validator: validator.New(),
	}
}

// HandleBuildSucceeded is the public handler for NATS messages. It wraps the
// real message and passes it to the testable internal handler.
func (w *Worker) HandleBuildSucceeded(m *nats.Msg) {
	w.handleBuildSucceededInternal(&natsMsgAdapter{msg: m})
}

// handleBuildSucceededInternal contains the core logic for processing events.
func (w *Worker) handleBuildSucceededInternal(m natsMsg) {
	var event events.BuildSucceeded
	if err := json.Unmarshal(m.GetData(), &event); err != nil {
		w.Logger.Error().Err(err).Msg("Could not unmarshal build succeeded event, terminating message")
		if err := m.Term(); err != nil {
			w.Logger.Error().Err(err).Msg("Failed to terminate NATS message")
		}
		return
	}

	// Validate the event payload
	if err := w.Validator.Struct(&event); err != nil {
		w.Logger.Error().Err(err).Msg("Invalid build succeeded event payload, terminating message")
		if err := m.Term(); err != nil {
			w.Logger.Error().Err(err).Msg("Failed to terminate NATS message")
		}
		return
	}

	log := w.Logger.With().
		Str("app_id", event.AppID).
		Str("image_uri", event.ImageURI).
		Logger()

	log.Info().Msg("Received build succeeded event")

	// Simulate the deployment process
	log.Info().Msg("Simulating deployment...")
	time.Sleep(1 * time.Second) // Reduced for faster tests
	log.Info().Msg("Deployment simulation complete")

	// In a real implementation, we would publish a "deployment.succeeded" event here.
	log.Info().Msg("End of workflow")

	// Acknowledge the message now that processing is complete
	if err := m.Ack(); err != nil {
		log.Error().Err(err).Msg("Failed to acknowledge NATS message")
	}
}