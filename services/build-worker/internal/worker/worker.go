package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"helios/pkg/events"
	"github.com/nats-io/nats.go"
)

// NatsPublisher defines the interface for publishing messages to NATS.
type NatsPublisher interface {
	Publish(subject string, data []byte) error
}

// Worker holds dependencies for the message handler.
type Worker struct {
	NATS NatsPublisher
}

// NewWorker creates a new Worker.
func NewWorker(nats NatsPublisher) *Worker {
	return &Worker{NATS: nats}
}

// HandleDeploymentRequest processes incoming deployment request events.
func (w *Worker) HandleDeploymentRequest(m *nats.Msg) {
	var request events.DeploymentRequest
	if err := json.Unmarshal(m.Data, &request); err != nil {
		log.Printf("ERROR: Could not unmarshal deployment request: %v", err)
		return
	}

	log.Printf("INFO: Received deployment request for App ID: %s, Repo: %s", request.AppID, request.GitRepository)

	// Simulate the build process
	log.Printf("INFO: Simulating build process for %s...", request.AppID)
	time.Sleep(1 * time.Second) // Reduced for faster tests
	log.Printf("INFO: Build simulation complete for %s.", request.AppID)

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
		log.Printf("ERROR: could not marshal build succeeded event: %v", err)
		return
	}

	subject := events.SubjectBuildSucceeded
	if err := w.NATS.Publish(subject, eventData); err != nil {
		log.Printf("ERROR: failed to publish to NATS subject '%s': %v", subject, err)
		return
	}

	log.Printf("SUCCESS: Published event to NATS subject '%s' for App ID %s", subject, request.AppID)
}