package main

import (
	"log"
	"net/http"
	"os"

	"helios/api/internal/handlers"
	"github.com/nats-io/nats.go"
)

func main() {
	// --- Connect to NATS ---
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL // "nats://127.0.0.1:4222"
	}

	natsConn, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to NATS at %s: %v", natsURL, err)
	}
	defer natsConn.Close()
	log.Printf("Successfully connected to NATS at %s", natsConn.ConnectedUrl())

	// --- Setup HTTP Server ---
	// The nats.Conn satisfies the NatsPublisher interface, so we can pass it directly.
	apiHandlers := handlers.NewAPIHandlers(natsConn)

	mux := http.NewServeMux()
	mux.HandleFunc("/projects", apiHandlers.CreateProjectHandler)
	mux.HandleFunc("/applications", apiHandlers.CreateApplicationHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting API server on port %s...", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("FATAL: Could not start server: %v", err)
	}
}