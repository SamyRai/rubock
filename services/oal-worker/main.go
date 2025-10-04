package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

// --- Event Payloads ---

// BuildSucceeded is the incoming event from the build-worker
type BuildSucceeded struct {
	AppID        string `json:"app_id"`
	ImageURI     string `json:"image_uri"`
	GitCommitSHA string `json:"git_commit_sha"`
}

// --- NATS Message Handler ---

func handleBuildSucceeded(m *nats.Msg, nc *nats.Conn) {
	var event BuildSucceeded
	if err := json.Unmarshal(m.Data, &event); err != nil {
		log.Printf("ERROR: Could not unmarshal build succeeded event: %v", err)
		return
	}

	log.Printf("INFO: Received build succeeded event for App ID: %s, Image: %s", event.AppID, event.ImageURI)

	// --- Simulate the Deployment Process ---
	// In a real implementation, this is where we would:
	// 1. Fetch the `Heliosfile.yml` from the git repo at the specified commit SHA.
	// 2. Parse the file.
	// 3. Translate the Heliosfile into a `docker-compose.yml`.
	// 4. SSH into the target server and run `docker-compose up -d`.
	log.Printf("INFO: Simulating deployment of image %s for App ID %s...", event.ImageURI, event.AppID)
	log.Printf("INFO: > docker-compose up -d (simulation)")
	time.Sleep(3 * time.Second) // Simulate work
	log.Printf("INFO: Deployment simulation complete for %s.", event.AppID)

	// --- Publish Deployment Succeeded Event ---
	// In a real implementation, we would publish a "deployment.succeeded" event here.
	// For the MVP, we'll just log to signify the end of the workflow.
	log.Printf("SUCCESS: End of workflow for App ID %s.", event.AppID)
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

	// --- Subscribe to Build Succeeded Events ---
	subject := "build.succeeded"
	// Use a queue subscription to ensure only one worker picks up a message.
	_, err = nc.QueueSubscribe(subject, "oal-workers", func(m *nats.Msg) {
		handleBuildSucceeded(m, nc)
	})
	if err != nil {
		log.Fatalf("FATAL: Could not subscribe to NATS subject '%s': %v", subject, err)
	}

	log.Printf("Listening on subject '%s' with queue group 'oal-workers'", subject)

	// Keep the process running
	select {}
}