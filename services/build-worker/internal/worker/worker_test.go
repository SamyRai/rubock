package worker

import (
	"encoding/json"
	"testing"

	"helios/pkg/events"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
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
	testLogger := zerolog.Nop()
	mockNATS := &MockNatsPublisher{}
	worker := NewWorker(mockNATS, testLogger)

	// Create a sample deployment request event
	request := events.DeploymentRequest{
		AppID:         "app-123",
		GitRepository: "https://github.com/example/app.git",
		GitBranch:     "develop",
	}
	requestData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Create a NATS message
	msg := &nats.Msg{
		Subject: events.SubjectDeploymentRequested,
		Data:    requestData,
	}

	// --- Act ---
	worker.HandleDeploymentRequest(msg)

	// --- Assert ---

	// 1. Check that a message was published
	if mockNATS.PublishedSubject == "" {
		t.Fatal("worker did not publish a NATS message")
	}

	// 2. Check that it was published to the correct subject
	expectedSubject := events.SubjectBuildSucceeded
	if mockNATS.PublishedSubject != expectedSubject {
		t.Errorf("worker published to wrong NATS subject: got %s want %s", mockNATS.PublishedSubject, expectedSubject)
	}

	// 3. Check the payload of the published message
	var publishedEvent events.BuildSucceeded
	if err := json.Unmarshal(mockNATS.PublishedData, &publishedEvent); err != nil {
		t.Fatalf("Could not unmarshal published NATS message payload: %v", err)
	}

	if publishedEvent.AppID != request.AppID {
		t.Errorf("NATS event has wrong AppID: got %s want %s", publishedEvent.AppID, request.AppID)
	}

	if publishedEvent.GitCommitSHA == "" {
		t.Error("NATS event is missing GitCommitSHA")
	}

	if publishedEvent.ImageURI == "" {
		t.Error("NATS event is missing ImageURI")
	}
}