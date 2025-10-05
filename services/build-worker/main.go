package main

import (
	"log"
	"os"

	"helios/build-worker/internal/worker"
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
	// The nats.Conn satisfies the NatsPublisher interface.
	w := worker.NewWorker(nc)

	// --- Subscribe to Deployment Requests ---
	subject := events.SubjectDeploymentRequested
	_, err = nc.QueueSubscribe(subject, "build-workers", w.HandleDeploymentRequest)
	if err != nil {
		log.Fatalf("FATAL: Could not subscribe to NATS subject '%s': %v", subject, err)
	}

	log.Printf("Listening on subject '%s' with queue group 'build-workers'", subject)

	// Keep the process running
	select {}
}