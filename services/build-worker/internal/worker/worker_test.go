package worker

import (
	"encoding/json"
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

// Publish records the subject and data it was called with.
func (m *MockNatsPublisher) Publish(subject string, data []byte) error {
	if m.PublishError != nil {
		return m.PublishError
	}
	m.PublishedSubject = subject
	m.PublishedData = data
	return nil
}

// --- Tests ---

func TestHandleDeploymentRequest(t *testing.T) {
	// --- Setup ---
	testLogger := testutil.NewTestLogger()
	mockNATS := &MockNatsPublisher{}
	worker := NewWorker(mockNATS, testLogger)

	// Create a sample deployment request event
	request := events.DeploymentRequest{
		AppID:         "app-123",
		GitRepository: "https://github.com/example/app.git",
		GitBranch:     "develop",
	}
	requestData, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal request")

	// Create a NATS message
	msg := &nats.Msg{
		Subject: events.SubjectDeploymentRequested,
		Data:    requestData,
	}

	// --- Act ---
	worker.HandleDeploymentRequest(msg)

	// --- Assert ---

	// 1. Check that a message was published
	assert.NotEmpty(t, mockNATS.PublishedSubject, "worker did not publish a NATS message")

	// 2. Check that it was published to the correct subject
	assert.Equal(t, events.SubjectBuildSucceeded, mockNATS.PublishedSubject, "worker published to wrong NATS subject")

	// 3. Check the payload of the published message
	var publishedEvent events.BuildSucceeded
	err = json.Unmarshal(mockNATS.PublishedData, &publishedEvent)
	require.NoError(t, err, "Could not unmarshal published NATS message payload")

	assert.Equal(t, request.AppID, publishedEvent.AppID, "NATS event has wrong AppID")
	assert.NotEmpty(t, publishedEvent.GitCommitSHA, "NATS event is missing GitCommitSHA")
	assert.NotEmpty(t, publishedEvent.ImageURI, "NATS event is missing ImageURI")
}