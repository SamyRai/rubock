package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
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

// Publish records the subject and data it was called with, then returns any configured error.
func (m *MockNatsPublisher) Publish(subject string, data []byte) error {
	m.PublishedSubject = subject
	m.PublishedData = data
	return m.PublishError
}

// --- Tests ---

func TestCreateProjectHandler(t *testing.T) {
	testCases := []struct {
		name               string
		method             string
		expectedStatusCode int
	}{
		{
			name:               "Successful Case - POST",
			method:             http.MethodPost,
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Failure Case - GET not allowed",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			testLogger := testutil.NewTestLogger()
			handlers := NewAPIHandlers(nil, testLogger)

			req, err := http.NewRequest(tc.method, "/projects", nil)
			require.NoError(t, err, "Could not create request")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlers.CreateProjectHandler)

			// Execute
			handler.ServeHTTP(rr, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, rr.Code, "handler returned wrong status code")

			if tc.expectedStatusCode == http.StatusCreated {
				var response map[string]string
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err, "Could not parse response body")
				assert.Equal(t, "proj_12345", response["id"], "handler returned unexpected body")
			}
		})
	}
}

func TestCreateApplicationHandler(t *testing.T) {
	testCases := []struct {
		name               string
		method             string
		body               io.Reader
		mockNatsError      error
		expectedStatusCode int
		expectNatsPublish  bool
	}{
		{
			name:   "Successful Case",
			method: http.MethodPost,
			body: bytes.NewBufferString(`{
				"name": "my-app",
				"git_repository": "https://github.com/example/my-app.git",
				"git_branch": "main"
			}`),
			mockNatsError:      nil,
			expectedStatusCode: http.StatusAccepted,
			expectNatsPublish:  true,
		},
		{
			name:               "Failure Case - Invalid JSON",
			method:             http.MethodPost,
			body:               bytes.NewBufferString(`{"name": "my-app",}`),
			mockNatsError:      nil,
			expectedStatusCode: http.StatusBadRequest,
			expectNatsPublish:  false,
		},
		{
			name:   "Failure Case - NATS Publish Error",
			method: http.MethodPost,
			body: bytes.NewBufferString(`{
				"name": "my-app",
				"git_repository": "https://github.com/example/my-app.git",
				"git_branch": "main"
			}`),
			mockNatsError:      errors.New("NATS is down"),
			expectedStatusCode: http.StatusInternalServerError,
			expectNatsPublish:  true, // It will attempt to publish
		},
		{
			name:               "Failure Case - Method Not Allowed",
			method:             http.MethodGet,
			body:               nil,
			mockNatsError:      nil,
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectNatsPublish:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			testLogger := testutil.NewTestLogger()
			mockNATS := &MockNatsPublisher{PublishError: tc.mockNatsError}
			handlers := NewAPIHandlers(mockNATS, testLogger)

			req, err := http.NewRequest(tc.method, "/applications", tc.body)
			require.NoError(t, err, "Could not create request")
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlers.CreateApplicationHandler)

			// Execute
			handler.ServeHTTP(rr, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, rr.Code, "handler returned wrong status code")

			if tc.expectNatsPublish {
				assert.NotEmpty(t, mockNATS.PublishedSubject, "handler should have attempted to publish a NATS message")
				if tc.mockNatsError == nil {
					assert.Equal(t, events.SubjectDeploymentRequested, mockNATS.PublishedSubject, "handler published to wrong NATS subject")
					var event events.DeploymentRequest
					err = json.Unmarshal(mockNATS.PublishedData, &event)
					require.NoError(t, err, "Could not unmarshal NATS message payload")
					assert.Equal(t, "app_67890", event.AppID, "NATS event has wrong AppID")
				}
			} else {
				assert.Empty(t, mockNATS.PublishedSubject, "handler should not have published a NATS message")
			}
		})
	}
}