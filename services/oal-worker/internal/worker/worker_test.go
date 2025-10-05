package worker

import (
	"encoding/json"
	"testing"

	"helios/pkg/events"
	"helios/pkg/testutil"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
)

func TestHandleBuildSucceeded(t *testing.T) {
	// --- Setup ---
	testLogger := testutil.NewTestLogger()
	worker := NewWorker(testLogger)

	// Create a sample build succeeded event
	event := events.BuildSucceeded{
		AppID:        "app-123",
		ImageURI:     "registry.helios.internal/app-123:a1b2c3d4",
		GitCommitSHA: "a1b2c3d4",
	}
	eventData, err := json.Marshal(event)
	require.NoError(t, err, "Failed to marshal event")

	// Create a NATS message
	msg := &nats.Msg{
		Subject: events.SubjectBuildSucceeded,
		Data:    eventData,
	}

	// --- Act & Assert ---
	// The test will pass if this function executes without panicking.
	// In a real-world scenario, we might check for logs or other side effects.
	require.NotPanics(t, func() {
		worker.HandleBuildSucceeded(msg)
	}, "HandleBuildSucceeded should not panic")
}