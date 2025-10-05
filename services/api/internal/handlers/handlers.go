package handlers

import (
	"encoding/json"
	"net/http"

	"helios/pkg/events"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

// NatsPublisher defines the interface for publishing messages to NATS.
type NatsPublisher interface {
	Publish(subject string, data []byte) error
}

// APIHandlers holds dependencies for the HTTP handlers.
type APIHandlers struct {
	NATS      NatsPublisher
	Logger    zerolog.Logger
	Validator *validator.Validate
}

// NewAPIHandlers creates a new APIHandlers struct.
func NewAPIHandlers(nats NatsPublisher, logger zerolog.Logger) *APIHandlers {
	return &APIHandlers{
		NATS:      nats,
		Logger:    logger,
		Validator: validator.New(),
	}
}

// CreateProjectHandler simulates creating a new project.
func (h *APIHandlers) CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simulate creating a project and returning its ID
	projectID := "proj_12345"
	h.Logger.Info().Str("project_id", projectID).Msg("Simulating project creation")

	response := map[string]string{"id": projectID, "name": "New Project"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// CreateApplicationRequest defines the structure for the application creation request body.
type CreateApplicationRequest struct {
	Name          string `json:"name" validate:"required"`
	GitRepository string `json:"git_repository" validate:"required,url"`
	GitBranch     string `json:"git_branch" validate:"required"`
}

// CreateApplicationHandler simulates creating a new application and triggers a deployment.
func (h *APIHandlers) CreateApplicationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody CreateApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		h.Logger.Warn().Err(err).Msg("Could not decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the request body
	if err := h.Validator.Struct(&reqBody); err != nil {
		h.Logger.Warn().Err(err).Msg("Request body validation failed")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Simulate creating an application and getting an ID
	appID := "app_67890"
	h.Logger.Info().
		Str("app_id", appID).
		Str("app_name", reqBody.Name).
		Msg("Simulating application creation")

	// Create the deployment request event using the shared package
	event := events.DeploymentRequest{
		AppID:         appID,
		GitRepository: reqBody.GitRepository,
		GitBranch:     reqBody.GitBranch,
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Could not marshal event data")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Publish the event to NATS using the interface
	subject := events.SubjectDeploymentRequested
	if err := h.NATS.Publish(subject, eventData); err != nil {
		h.Logger.Error().Err(err).Str("subject", subject).Msg("Failed to publish to NATS")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.Logger.Info().
		Str("subject", subject).
		Str("app_id", appID).
		Msg("Successfully published event to NATS")

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