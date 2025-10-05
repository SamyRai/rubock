package main

import (
	"log"
	"os"

	"helios/oal-worker/internal/worker"
	"helios/pkg/events"
	"github.com/nats-io/nats.go"
)

func main() {
	// --- Connect to NATS ---
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to NATS at %s: %v", natsURL, err)
	}
	defer nc.Close()
	log.Printf("Successfully connected to NATS at %s", nc.ConnectedUrl())

	// --- Setup Worker ---
	w := worker.NewWorker()

	// --- Subscribe to Build Succeeded Events ---
	subject := events.SubjectBuildSucceeded
	_, err = nc.QueueSubscribe(subject, "oal-workers", w.HandleBuildSucceeded)
	if err != nil {
		log.Fatalf("FATAL: Could not subscribe to NATS subject '%s': %v", subject, err)
	}

	log.Printf("Listening on subject '%s' with queue group 'oal-workers'", subject)

	// Keep the process running
	select {}
}