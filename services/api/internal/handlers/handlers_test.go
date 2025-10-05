package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"helios/pkg/events"
	"helios/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	// Use the test logger from the shared package.
	testLogger := testutil.NewTestLogger()
	handlers := NewAPIHandlers(nil, testLogger)

	req, err := http.NewRequest("POST", "/projects", nil)
	require.NoError(t, err, "Could not create request")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateProjectHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusCreated, rr.Code, "handler returned wrong status code")

	// Check the response body
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err, "Could not parse response body")

	assert.Equal(t, "proj_12345", response["id"], "handler returned unexpected body")
}

func TestCreateApplicationHandler(t *testing.T) {
	testLogger := testutil.NewTestLogger()
	mockNATS := &MockNatsPublisher{}
	handlers := NewAPIHandlers(mockNATS, testLogger)

	// Create the request body
	reqBody := map[string]string{
		"name":           "my-app",
		"git_repository": "https://github.com/example/my-app.git",
		"git_branch":     "main",
	}
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/applications", bytes.NewBuffer(body))
	require.NoError(t, err, "Could not create request")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateApplicationHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusAccepted, rr.Code, "handler returned wrong status code")

	// Check that a message was published to NATS
	assert.NotEmpty(t, mockNATS.PublishedSubject, "handler did not publish a NATS message")
	assert.Equal(t, events.SubjectDeploymentRequested, mockNATS.PublishedSubject, "handler published to wrong NATS subject")

	// Check the NATS message payload
	var event events.DeploymentRequest
	err = json.Unmarshal(mockNATS.PublishedData, &event)
	require.NoError(t, err, "Could not unmarshal NATS message payload")

	assert.Equal(t, "app_67890", event.AppID, "NATS event has wrong AppID")
	assert.Equal(t, reqBody["git_repository"], event.GitRepository, "NATS event has wrong GitRepository")
}