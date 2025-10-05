package worker

import (
	"encoding/json"
	"testing"

	"helios/pkg/events"
	"github.com/nats-io/nats.go"
)

func TestHandleBuildSucceeded(t *testing.T) {
	// --- Setup ---
	worker := NewWorker()

	// Create a sample build succeeded event
	event := events.BuildSucceeded{
		AppID:        "app-123",
		ImageURI:     "registry.helios.internal/app-123:a1b2c3d4",
		GitCommitSHA: "a1b2c3d4",
	}
	eventData, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal event: %v", err)
	}

	// Create a NATS message
	msg := &nats.Msg{
		Subject: events.SubjectBuildSucceeded,
		Data:    eventData,
	}

	// --- Act & Assert ---
	// The test will pass if this function executes without panicking.
	// In a real-world scenario, we might check for logs or other side effects.
	worker.HandleBuildSucceeded(msg)
}