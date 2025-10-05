package worker

import (
	"encoding/json"
	"log"
	"time"

	"helios/pkg/events"
	"github.com/nats-io/nats.go"
)

// Worker holds dependencies for the message handler.
// It doesn't have any dependencies for now, but is kept for consistency.
type Worker struct{}

// NewWorker creates a new Worker.
func NewWorker() *Worker {
	return &Worker{}
}

// HandleBuildSucceeded processes incoming build succeeded events.
func (w *Worker) HandleBuildSucceeded(m *nats.Msg) {
	var event events.BuildSucceeded
	if err := json.Unmarshal(m.Data, &event); err != nil {
		log.Printf("ERROR: Could not unmarshal build succeeded event: %v", err)
		return
	}

	log.Printf("INFO: Received build succeeded event for App ID: %s, Image: %s", event.AppID, event.ImageURI)

	// Simulate the deployment process
	log.Printf("INFO: Simulating deployment of image %s for App ID %s...", event.ImageURI, event.AppID)
	log.Printf("INFO: > docker-compose up -d (simulation)")
	time.Sleep(1 * time.Second) // Reduced for faster tests
	log.Printf("INFO: Deployment simulation complete for %s.", event.AppID)

	// In a real implementation, we would publish a "deployment.succeeded" event here.
	log.Printf("SUCCESS: End of workflow for App ID %s.", event.AppID)
}