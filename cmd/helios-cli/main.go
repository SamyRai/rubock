package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// --- Define and parse command-line flags ---
	var gitRepo, gitBranch, apiURL, appName string

	flag.StringVar(&gitRepo, "repo", "", "The git repository URL to deploy (required).")
	flag.StringVar(&gitBranch, "branch", "main", "The git branch to deploy.")
	flag.StringVar(&appName, "name", "my-app", "The name of the application.")
	flag.StringVar(&apiURL, "api", "http://localhost:8080", "The URL of the Helios API server.")
	flag.Parse()

	if gitRepo == "" {
		fmt.Println("Error: The --repo flag is required.")
		flag.Usage()
		os.Exit(1)
	}

	// --- Construct the request to the API server ---
	requestBody, err := json.Marshal(map[string]string{
		"name":            appName,
		"git_repository":  gitRepo,
		"git_branch":      gitBranch,
	})
	if err != nil {
		log.Fatalf("FATAL: Failed to marshal request JSON: %v", err)
	}

	// --- Make the HTTP POST request ---
	log.Printf("Sending deployment request for repo %s to %s...", gitRepo, apiURL)

	resp, err := http.Post(apiURL+"/applications", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("FATAL: Failed to send request to API server: %v", err)
	}
	defer resp.Body.Close()

	// --- Process the response ---
	log.Printf("Received response with status code: %d", resp.StatusCode)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("FATAL: Failed to read response body: %v", err)
	}

	if resp.StatusCode >= 400 {
		log.Fatalf("FATAL: API server returned an error:\n%s", string(responseBody))
	}

	// --- Print the successful response ---
	fmt.Println("\nDeployment request accepted by Helios:")
	// Pretty print the JSON response
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, responseBody, "", "  "); err != nil {
		// If indent fails, just print the raw response
		fmt.Println(string(responseBody))
	} else {
		fmt.Println(prettyJSON.String())
	}
}