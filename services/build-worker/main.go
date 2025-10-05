package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

// --- Event Payloads ---

// DeploymentRequest is the incoming event from the API
type DeploymentRequest struct {
	AppID         string `json:"app_id"`
	GitRepository string `json:"git_repository"`
	GitBranch     string `json:"git_branch"`
}

// BuildSucceeded is the event this worker publishes on success
type BuildSucceeded struct {
	AppID        string `json:"app_id"`
	ImageURI     string `json:"image_uri"`
	GitCommitSHA string `json:"git_commit_sha"`
}

// --- NATS Message Handler ---

func handleDeploymentRequest(m *nats.Msg, nc *nats.Conn) {
	var request DeploymentRequest
	if err := json.Unmarshal(m.Data, &request); err != nil {
		log.Printf("ERROR: Could not unmarshal deployment request: %v", err)
		return // Don't negatively acknowledge, just log the error
	}

	log.Printf("INFO: Received deployment request for App ID: %s, Repo: %s", request.AppID, request.GitRepository)

	// --- Simulate the Build Process ---
	// In a real implementation, this is where we would:
	// 1. Clone the git repository.
	// 2. Run Cloud Native Buildpacks (e.g., `pack build ...`).
	// 3. Push the resulting image to an internal container registry.
	log.Printf("INFO: Simulating build process for %s...", request.AppID)
	time.Sleep(5 * time.Second) // Simulate work being done
	log.Printf("INFO: Build simulation complete for %s.", request.AppID)

	// --- Publish Build Succeeded Event ---
	// For the MVP, we'll use a placeholder image URI and commit SHA.
	commitSHA := "a1b2c3d4e5f6"
	imageURI := "registry.helios.internal/" + request.AppID + ":" + commitSHA

	event := BuildSucceeded{
		AppID:        request.AppID,
		ImageURI:     imageURI,
		GitCommitSHA: commitSHA,
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		log.Printf("ERROR: could not marshal build succeeded event: %v", err)
		return
	}

	subject := "build.succeeded"
	if err := nc.Publish(subject, eventData); err != nil {
		log.Printf("ERROR: failed to publish to NATS subject '%s': %v", subject, err)
		return
	}

	log.Printf("SUCCESS: Published event to NATS subject '%s' for App ID %s", subject, request.AppID)
}

// --- Main Application ---

func main() {
	// --- Connect to NATS ---
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL // "nats://127.0.0.1:4222"
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to NATS at %s: %v", natsURL, err)
	}
	defer nc.Close()
	log.Printf("Successfully connected to NATS at %s", nc.ConnectedUrl())

	// --- Subscribe to Deployment Requests ---
	subject := "deployment.requested"
	// Use a queue subscription to ensure only one worker picks up a message.
	_, err = nc.QueueSubscribe(subject, "build-workers", func(m *nats.Msg) {
		handleDeploymentRequest(m, nc)
	})
	if err != nil {
		log.Fatalf("FATAL: Could not subscribe to NATS subject '%s': %v", subject, err)
	}

	log.Printf("Listening on subject '%s' with queue group 'build-workers'", subject)

	// Keep the process running
	select {}
}