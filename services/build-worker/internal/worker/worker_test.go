package worker

import (
	"encoding/json"
	"errors"
	"testing"

	"helios/pkg/events"
	"helios/pkg/testutil"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Mocks ---

// MockNatsPublisher is a mock implementation of the NatsPublisher interface.
type MockNatsPublisher struct {
	PublishedSubject string
	PublishedData    []byte
	PublishError     error
}

// Publish records the subject and data it was called with, then returns any configured error.
func (m *MockNatsPublisher) Publish(subject string, data []byte) error {
	m.PublishedSubject = subject
	m.PublishedData = data
	return m.PublishError
}

// --- Tests ---

func TestHandleDeploymentRequest(t *testing.T) {
	// Setup a valid deployment request for reuse
	validRequest := events.DeploymentRequest{
		AppID:         "app-123",
		GitRepository: "https://github.com/example/app.git",
		GitBranch:     "develop",
	}
	validRequestData, err := json.Marshal(validRequest)
	require.NoError(t, err, "Setup failed: could not marshal valid request")

	testCases := []struct {
		name              string
		natsMsgData       []byte
		mockNatsError     error
		expectNatsPublish bool
	}{
		{
			name:              "Successful Case",
			natsMsgData:       validRequestData,
			mockNatsError:     nil,
			expectNatsPublish: true,
		},
		{
			name:              "Failure Case - Invalid JSON",
			natsMsgData:       []byte(`{"app_id": "app-123",`),
			mockNatsError:     nil,
			expectNatsPublish: false,
		},
		{
			name:              "Failure Case - NATS Publish Error",
			natsMsgData:       validRequestData,
			mockNatsError:     errors.New("NATS is down"),
			expectNatsPublish: true, // It will still attempt to publish
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			testLogger := testutil.NewTestLogger()
			mockNATS := &MockNatsPublisher{PublishError: tc.mockNatsError}
			worker := NewWorker(mockNATS, testLogger)

			msg := &nats.Msg{
				Subject: events.SubjectDeploymentRequested,
				Data:    tc.natsMsgData,
			}

			// Execute
			worker.HandleDeploymentRequest(msg)

			// Assert
			if tc.expectNatsPublish {
				assert.NotEmpty(t, mockNATS.PublishedSubject, "worker should have attempted to publish a NATS message")
				if tc.mockNatsError == nil {
					assert.Equal(t, events.SubjectBuildSucceeded, mockNATS.PublishedSubject, "worker published to wrong NATS subject")
					var publishedEvent events.BuildSucceeded
					err := json.Unmarshal(mockNATS.PublishedData, &publishedEvent)
					require.NoError(t, err, "Could not unmarshal published NATS message payload")
					assert.Equal(t, validRequest.AppID, publishedEvent.AppID, "NATS event has wrong AppID")
					assert.NotEmpty(t, publishedEvent.GitCommitSHA, "NATS event is missing GitCommitSHA")
					assert.NotEmpty(t, publishedEvent.ImageURI, "NATS event is missing ImageURI")
				}
			} else {
				assert.Empty(t, mockNATS.PublishedSubject, "worker should not have published a NATS message")
			}
		})
	}
}