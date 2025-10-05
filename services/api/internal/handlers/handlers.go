package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"helios/pkg/events"
)

// NatsPublisher defines the interface for publishing messages to NATS.
// This allows for mocking in tests.
type NatsPublisher interface {
	Publish(subject string, data []byte) error
}

// APIHandlers holds dependencies for the HTTP handlers.
type APIHandlers struct {
	NATS NatsPublisher
}

// NewAPIHandlers creates a new APIHandlers struct.
func NewAPIHandlers(nats NatsPublisher) *APIHandlers {
	return &APIHandlers{NATS: nats}
}

// CreateProjectHandler simulates creating a new project.
func (h *APIHandlers) CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simulate creating a project and returning its ID
	projectID := "proj_12345"
	log.Printf("Simulating project creation. Assigned ID: %s", projectID)

	response := map[string]string{"id": projectID, "name": "New Project"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// CreateApplicationHandler simulates creating a new application and triggers a deployment.
func (h *APIHandlers) CreateApplicationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody struct {
		Name          string `json:"name"`
		GitRepository string `json:"git_repository"`
		GitBranch     string `json:"git_branch"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Simulate creating an application and getting an ID
	appID := "app_67890"
	log.Printf("Simulating application creation '%s'. Assigned ID: %s", reqBody.Name, appID)

	// Create the deployment request event using the shared package
	event := events.DeploymentRequest{
		AppID:         appID,
		GitRepository: reqBody.GitRepository,
		GitBranch:     reqBody.GitBranch,
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		log.Printf("ERROR: could not marshal event data: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Publish the event to NATS using the interface
	subject := events.SubjectDeploymentRequested
	if err := h.NATS.Publish(subject, eventData); err != nil {
		log.Printf("ERROR: failed to publish to NATS subject '%s': %v", subject, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("SUCCESS: Published event to NATS subject '%s' for App ID %s", subject, appID)

	// Respond to the client
	response := map[string]string{
		"id":     appID,
		"name":   reqBody.Name,
		"status": "pending",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}