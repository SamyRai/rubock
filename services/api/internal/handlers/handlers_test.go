package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"helios/pkg/events"
	"github.com/rs/zerolog"
)

// --- Mocks ---

// MockNatsPublisher is a mock implementation of the NatsPublisher interface
// for testing purposes.
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

func TestCreateProjectHandler(t *testing.T) {
	// Use a disabled logger for tests to avoid noisy output.
	testLogger := zerolog.Nop()
	// Create a new set of handlers with a nil NATS publisher since this handler
	// doesn't use it.
	handlers := NewAPIHandlers(nil, testLogger)

	req, err := http.NewRequest("POST", "/projects", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateProjectHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Check the response body
	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Could not parse response body: %v", err)
	}

	if id, ok := response["id"]; !ok || id != "proj_12345" {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestCreateApplicationHandler(t *testing.T) {
	testLogger := zerolog.Nop()
	mockNATS := &MockNatsPublisher{}
	handlers := NewAPIHandlers(mockNATS, testLogger)

	// Create the request body
	reqBody := map[string]string{
		"name":            "my-app",
		"git_repository":  "https://github.com/example/my-app.git",
		"git_branch":      "main",
	}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "/applications", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateApplicationHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusAccepted)
	}

	// Check that a message was published to NATS
	if mockNATS.PublishedSubject == "" {
		t.Errorf("handler did not publish a NATS message")
	}

	if mockNATS.PublishedSubject != events.SubjectDeploymentRequested {
		t.Errorf("handler published to wrong NATS subject: got %s want %s", mockNATS.PublishedSubject, events.SubjectDeploymentRequested)
	}

	// Check the NATS message payload
	var event events.DeploymentRequest
	if err := json.Unmarshal(mockNATS.PublishedData, &event); err != nil {
		t.Fatalf("Could not unmarshal NATS message payload: %v", err)
	}

	if event.AppID != "app_67890" {
		t.Errorf("NATS event has wrong AppID: got %s want %s", event.AppID, "app_67890")
	}
	if event.GitRepository != reqBody["git_repository"] {
		t.Errorf("NATS event has wrong GitRepository: got %s want %s", event.GitRepository, reqBody["git_repository"])
	}
}