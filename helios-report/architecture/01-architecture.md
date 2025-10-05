# Architecture: A Solid Foundation with Key Flaws

The Helios platform is built on a strong architectural foundation, utilizing a Go-based microservices architecture within a monorepo. This structure promotes code sharing and consistency. However, a detailed analysis reveals several critical architectural flaws that undermine the platform's reliability and scalability.

## Key Architectural Patterns

-   **Monorepo:** The use of a single repository for all services (`api`, `build-worker`, `oal-worker`) and shared packages (`pkg/*`) is a sound choice, simplifying dependency management via Go workspaces and `replace` directives.
-   **Centralized Dependencies:** A central `App` struct is used in each service (e.g., `services/api/internal/platform/app.go`) to manage dependencies like the logger, NATS connection, and database handle. This is an excellent pattern for clean, organized code.
-   **Asynchronous Communication:** The services communicate asynchronously via NATS, using well-defined events from the shared `pkg/events` package. This decouples services and enhances resilience.
-   **Graceful Shutdown:** All services correctly implement graceful shutdown, handling `SIGINT` and `SIGTERM` signals to ensure clean termination. The workers properly drain NATS subscriptions, preventing message loss during deployments.

## Critical Architectural Flaws

### 1. Missing Message Acknowledgement in Workers

This is the most severe architectural issue. The message handlers in both the `build-worker` and `oal-worker` do not acknowledge messages after processing them.

**Code Evidence (`services/build-worker/internal/worker/worker.go`):**
```go
// HandleDeploymentRequest processes incoming deployment request events.
func (w *Worker) HandleDeploymentRequest(m *nats.Msg) {
	var request events.DeploymentRequest
	if err := json.Unmarshal(m.Data, &request); err != nil {
		w.Logger.Error().Err(err).Msg("Could not unmarshal deployment request")
		return // <-- Message is not acknowledged here
	}
    // ... process message
    // <-- Message is never acknowledged
}
```

-   **Impact:** Without an `m.Ack()`, if the worker restarts, the message broker will assume the message was not processed and will redeliver it. This leads to **duplicate processing**, causing the same build or deployment to be triggered multiple times, which is a major reliability failure.

### 2. Duplicated Bootstrap Logic

The `main.go` file for all three services contains nearly identical, copy-pasted code for initializing the logger and connecting to NATS.

**Code Evidence (`services/api/main.go`, `services/build-worker/main.go`, etc.):**
```go
// This entire block is duplicated across all services.
var natsConn *nats.Conn
var err error
for i := 0; i < 5; i++ {
    natsConn, err = nats.Connect(natsURL)
    if err == nil {
        break
    }
    log.Warn().Err(err).Msgf("Failed to connect to NATS, retrying in %d seconds...", i+1)
    time.Sleep(time.Duration(i+1) * time.Second)
}
```

-   **Impact:** This duplication makes the system harder to maintain. A change to the connection logic (e.g., to implement exponential backoff) must be manually applied to every service, increasing the risk of inconsistencies.

### 3. Lack of Event Versioning

The shared `pkg/events/events.go` package defines event structures and subjects, which is good. However, the subjects are not versioned.

**Code Evidence (`services/pkg/events/events.go`):**
```go
const (
	SubjectDeploymentRequested = "deployment.requested" // <-- No version
	SubjectBuildSucceeded      = "build.succeeded"      // <-- No version
)
```

-   **Impact:** As the system evolves, event payloads will inevitably change. Without versioning (e.g., `deployment.requested.v1`), it's impossible to introduce breaking changes to an event structure without coordinating a simultaneous deployment of all consuming and producing services, which is brittle and operationally complex.