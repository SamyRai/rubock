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
	// Setup a valid build succeeded event for reuse
	validEvent := events.BuildSucceeded{
		AppID:        "app-123",
		ImageURI:     "registry.helios.internal/app-123:a1b2c3d4",
		GitCommitSHA: "a1b2c3d4",
	}
	validEventData, err := json.Marshal(validEvent)
	require.NoError(t, err, "Setup failed: could not marshal valid event")

	testCases := []struct {
		name        string
		natsMsgData []byte
	}{
		{
			name:        "Successful Case",
			natsMsgData: validEventData,
		},
		{
			name:        "Failure Case - Invalid JSON",
			natsMsgData: []byte(`{"app_id": "app-123",`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			testLogger := testutil.NewTestLogger()
			worker := NewWorker(testLogger)

			msg := &nats.Msg{
				Subject: events.SubjectBuildSucceeded,
				Data:    tc.natsMsgData,
			}

			// Execute & Assert
			// The handler should be robust and not panic, even with bad input.
			require.NotPanics(t, func() {
				worker.HandleBuildSucceeded(msg)
			}, "HandleBuildSucceeded should not panic")
		})
	}
}