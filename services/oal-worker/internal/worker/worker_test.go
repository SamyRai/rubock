package worker

import (
	"bytes"
	"encoding/json"
	"testing"

	"helios/pkg/events"
	"helios/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockNatsMsg is a mock implementation of the natsMsg interface for testing.
type mockNatsMsg struct {
	data   []byte
	acked  bool
	nakked bool
	termed bool
}

func (m *mockNatsMsg) GetData() []byte {
	return m.data
}

func (m *mockNatsMsg) Ack() error {
	m.acked = true
	return nil
}

func (m *mockNatsMsg) Nak() error {
	m.nakked = true
	return nil
}

func (m *mockNatsMsg) Term() error {
	m.termed = true
	return nil
}

func TestHandleBuildSucceededInternal(t *testing.T) {
	// Setup a valid build succeeded event for reuse
	validEvent := events.BuildSucceeded{
		AppID:        "app-123",
		ImageURI:     "registry.helios.internal/app-123:a1b2c3d4",
		GitCommitSHA: "a1b2c3d4",
	}
	validEventData, err := json.Marshal(validEvent)
	require.NoError(t, err, "Setup failed: could not marshal valid event")

	// Setup an invalid event (missing required field)
	invalidEvent := events.BuildSucceeded{
		AppID: "app-123",
		// ImageURI is missing
		GitCommitSHA: "a1b2c3d4",
	}
	invalidEventData, err := json.Marshal(invalidEvent)
	require.NoError(t, err, "Setup failed: could not marshal invalid event")

	testCases := []struct {
		name                  string
		natsMsgData           []byte
		expectAck             bool
		expectTerm            bool
		expectedLogContains   []string
		unexpectedLogContains []string
	}{
		{
			name:        "Successful Case",
			natsMsgData: validEventData,
			expectAck:   true,
			expectedLogContains: []string{
				"Received build succeeded event",
				"Simulating deployment...",
				"Deployment simulation complete",
				"End of workflow",
			},
		},
		{
			name:        "Failure Case - Invalid JSON",
			natsMsgData: []byte(`{"app_id": "app-123",`),
			expectTerm:  true,
			expectedLogContains: []string{
				"Could not unmarshal build succeeded event, terminating message",
			},
			unexpectedLogContains: []string{
				"Received build succeeded event",
			},
		},
		{
			name:        "Failure Case - Validation Error",
			natsMsgData: invalidEventData,
			expectTerm:  true,
			expectedLogContains: []string{
				"Invalid build succeeded event payload, terminating message",
			},
			unexpectedLogContains: []string{
				"Received build succeeded event",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			var logBuffer bytes.Buffer
			testLogger := testutil.NewTestLoggerWithOutput(&logBuffer)
			worker := NewWorker(testLogger)

			// Use the mock message
			mockMsg := &mockNatsMsg{
				data: tc.natsMsgData,
			}

			// Execute the internal handler
			worker.handleBuildSucceededInternal(mockMsg)

			// Assert
			logOutput := logBuffer.String()

			for _, expected := range tc.expectedLogContains {
				assert.Contains(t, logOutput, expected, "Log output should contain expected message")
			}
			for _, unexpected := range tc.unexpectedLogContains {
				assert.NotContains(t, logOutput, unexpected, "Log output should not contain unexpected message")
			}

			assert.Equal(t, tc.expectAck, mockMsg.acked, "Message acknowledgement state does not match expectation")
			assert.Equal(t, tc.expectTerm, mockMsg.termed, "Message termination state does not match expectation")
		})
	}
}