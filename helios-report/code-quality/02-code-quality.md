# Code Quality: Critical Flaws and Inconsistencies

While the codebase adheres to Go formatting standards and leverages good patterns like structured logging, a detailed review reveals several critical code-level issues related to security, configuration management, and robustness.

## Critical Issues

### 1. Hardcoded Database Credentials (High-Severity Security Risk)

The `api` service contains hardcoded database credentials directly in its source code. This is a major security vulnerability that must be remediated immediately.

**Code Evidence (`services/api/main.go`):**
```go
// Initialize database connection
dbCfg := database.DBConfig{
    User:         "user",
    Password:     "password", // <-- Hardcoded password
    Host:         "localhost",
    Port:         5432,
    DBName:       "helios",
    SSLMode:      "disable",
    MaxOpenConns: 25,
    MaxIdleConns: 25,
    MaxIdleTime:  15 * time.Minute,
}
```
-   **Impact:** Anyone with access to the source code can retrieve the database password. Secrets must be externalized and supplied via environment variables or a dedicated secrets management system (e.g., HashiCorp Vault).

## Medium-Severity Issues

### 2. Missing Input Validation

Handlers across the `api` and worker services decode incoming data but fail to validate it. This can lead to malformed data propagating through the system, causing unexpected errors or behavior.

**Code Evidence (`services/api/internal/handlers/handlers.go`):**
```go
func (h *APIHandlers) CreateApplicationHandler(w http.ResponseWriter, r *http.Request) {
    // ...
	var reqBody struct {
		Name          string `json:"name"`
		GitRepository string `json:"git_repository"`
		GitBranch     string `json:"git_branch"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		// ...
		return
	}

    // No validation is performed here.
    // What if GitRepository is not a valid URL?
    // What if Name is an empty string?

	event := events.DeploymentRequest{
		AppID:         appID,
		GitRepository: reqBody.GitRepository, // <-- Potentially invalid data is used
		GitBranch:     reqBody.GitBranch,
	}
    // ...
}
```
-   **Impact:** Lack of validation can lead to panics, incorrect processing, and security vulnerabilities (e.g., if a field is used to construct a shell command). A validation library (e.g., `go-playground/validator`) should be used.

### 3. Brittle Connection Retry Logic

All services use a simple, fixed-interval retry loop to connect to NATS. This is not a robust strategy for a production system.

**Code Evidence (Duplicated across all `main.go` files):**
```go
for i := 0; i < 5; i++ {
    natsConn, err = nats.Connect(natsURL)
    if err == nil {
        break
    }
    log.Warn().Err(err).Msgf("Failed to connect to NATS, retrying in %d seconds...", i+1)
    time.Sleep(time.Duration(i+1) * time.Second) // <-- Linear backoff
}
```
-   **Impact:** During a NATS outage, all services will retry in a synchronized, aggressive manner, which can overwhelm the message broker when it comes back online (a "thundering herd" problem). An exponential backoff with jitter should be implemented.

## Minor Issues & Inconsistencies

### 4. Redundant HTTP Method Checks

The API handlers in `services/api/internal/handlers/handlers.go` manually check the HTTP request method. This is redundant, as the `chi` router is already configured to handle method-based routing.

**Code Evidence (`services/api/internal/handlers/handlers.go`):**
```go
func (h *APIHandlers) CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // <-- Unnecessary check
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
    // ...
}
```
-   **Impact:** This adds unnecessary boilerplate to the handlers. The routing definition in `platform/app.go` (`a.Router.Post(...)`) already ensures that this handler only receives POST requests. Removing this check would make the code cleaner.

### 5. Inconsistent Configuration Loading

Configuration is loaded in an ad-hoc manner. For example, the `PORT` is read from the environment deep inside the `Run` method, while the NATS URL is read in `main`.

**Code Evidence (`services/api/internal/platform/app.go`):**
```go
func (a *App) Run() {
	port := os.Getenv("PORT") // <-- Loaded late in the lifecycle
	if port == "" {
		port = "8080"
	}
    //...
}
```
-   **Impact:** This makes it difficult to understand a service's full configuration dependency at a glance. All configuration should be loaded and validated once at startup in `main.go` and passed to the application via a config struct.