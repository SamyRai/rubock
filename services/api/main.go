package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nats-io/nats.go"
)

// Global NATS connection
var natsConn *nats.Conn

// DeploymentRequest is the event payload for a new deployment
type DeploymentRequest struct {
	AppID           string `json:"app_id"`
	GitRepository   string `json:"git_repository"`
	GitBranch       string `json:"git_branch"`
}

// --- Handlers ---

// createProjectHandler simulates creating a new project.
// In a real implementation, this would write to the database.
func createProjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simulate creating a project and returning its ID
	// In a real app, this would come from the database.
	projectID := "proj_12345"
	log.Printf("Simulating project creation. Assigned ID: %s", projectID)

	response := map[string]string{"id": projectID, "name": "New Project"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// createApplicationHandler simulates creating a new application and triggers a deployment.
func createApplicationHandler(w http.ResponseWriter, r *http.Request) {
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

	// Create the deployment request event
	event := DeploymentRequest{
		AppID:           appID,
		GitRepository:   reqBody.GitRepository,
		GitBranch:       reqBody.GitBranch,
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		log.Printf("ERROR: could not marshal event data: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Publish the event to NATS
	subject := "deployment.requested"
	if err := natsConn.Publish(subject, eventData); err != nil {
		log.Printf("ERROR: failed to publish to NATS subject '%s': %v", subject, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("SUCCESS: Published event to NATS subject '%s' for App ID %s", subject, appID)

	// Respond to the client
	response := map[string]string{
		"id":      appID,
		"name":    reqBody.Name,
		"status":  "pending",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}


// --- Main Application ---

func main() {
	// --- Connect to NATS ---
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL // "nats://127.0.0.1:4222"
	}

	var err error
	natsConn, err = nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to NATS at %s: %v", natsURL, err)
	}
	defer natsConn.Close()
	log.Printf("Successfully connected to NATS at %s", natsConn.ConnectedUrl())

	// --- Setup HTTP Server ---
	mux := http.NewServeMux()
	mux.HandleFunc("/projects", createProjectHandler)
	// A real router would be better for path parameters, e.g., /projects/{id}/applications
	// For this MVP, we'll use a simplified path.
	mux.HandleFunc("/applications", createApplicationHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting API server on port %s...", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("FATAL: Could not start server: %v", err)
	}
}