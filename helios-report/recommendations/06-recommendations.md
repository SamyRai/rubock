# Recommendations: A Prioritized Roadmap for Improvement

This section synthesizes the findings from the deep-dive analysis into a prioritized, actionable roadmap. The recommendations are grouped by severity to guide remediation efforts, focusing first on critical security and reliability flaws before addressing architectural improvements and best practices.

## Critical Priority: Must-Fix Issues

### 1. Remediate Hardcoded Credentials (Security)
-   **Action:** Immediately remove the hardcoded database credentials from `services/api/main.go`.
-   **Implementation:** Load all parts of the `DBConfig` struct (user, password, host, etc.) from environment variables. Use a library like `github.com/joho/godotenv` for local development and a proper secrets management system (e.g., Kubernetes Secrets, Vault) in production.

### 2. Implement Message Acknowledgement in Workers (Reliability)
-   **Action:** Modify the NATS message handlers in both the `build-worker` and `oal-worker` to explicitly acknowledge messages.
-   **Implementation:**
    -   Enable manual acknowledgement when subscribing to the NATS subject.
    -   Call `m.Ack()` after a message has been successfully processed.
    -   Call `m.Nack()` if the message cannot be processed and should be redelivered.
    -   Call `m.Term()` if the message is invalid and should never be processed again (dead-letter queue).

## High Priority: Major Gaps

### 3. Implement Robust Input Validation (Security & Reliability)
-   **Action:** Add validation for all incoming data in API handlers and worker message handlers.
-   **Implementation:** Integrate a library like `go-playground/validator` to define validation rules on request/event structs (e.g., `validate:"required,url"`). Return a `400 Bad Request` for invalid API requests and log and terminate invalid events in workers.

### 4. Introduce Metrics and Tracing (Observability)
-   **Action:** Instrument all services to emit metrics and traces.
-   **Implementation:**
    -   **Metrics:** Integrate a Prometheus client library. Add a `/metrics` endpoint to the `api` service. Expose key metrics like request latency, error rates, and NATS message processing rates.
    -   **Tracing:** Integrate OpenTelemetry. Add middleware to the `chi` router to trace incoming requests. Inject and extract trace context from NATS messages to enable distributed tracing across services.

### 5. Add Integration and End-to-End Tests (Testing)
-   **Action:** Create a new test suite that runs services in Docker Compose and verifies their interactions.
-   **Implementation:** Create a `tests/integration` directory. Write tests that:
    1.  Call the `POST /applications` endpoint on the `api` service.
    2.  Assert that the `build-worker` receives the correct NATS message.
    3.  Assert that the `oal-worker` receives the `BuildSucceeded` message.
    This will validate the entire core workflow of the platform.

### 6. Automate Dependency Management (Security)
-   **Action:** Implement automated dependency updates and vulnerability scanning.
-   **Implementation:**
    -   Configure **Dependabot** (via `.github/dependabot.yml`) to automatically open pull requests for outdated dependencies.
    -   Add a CI step that runs **`govulncheck`** to fail the build if high-severity vulnerabilities are found in the dependency tree.

## Medium Priority: Architectural & Code Quality Improvements

### 7. Refactor Duplicated Startup Logic (Maintainability)
-   **Action:** Create a shared utility package to handle the common NATS connection logic.
-   **Implementation:** Create a `pkg/bootstrap` or similar package with a function like `ConnectNATS()` that includes the robust retry logic (see next point) and is called from each service's `main.go`.

### 8. Implement Exponential Backoff for Retries (Reliability)
-   **Action:** Replace the brittle, linear retry loop in the NATS connection logic with an exponential backoff strategy with jitter.
-   **Implementation:** Use a well-tested library or implement a standard exponential backoff algorithm to prevent "thundering herd" issues.

### 9. Introduce Event Versioning (Architecture)
-   **Action:** Update NATS subjects to include a version number.
-   **Implementation:** Change subjects in `pkg/events` from `"deployment.requested"` to `"v1.deployment.requested"`. This allows for backward-compatible evolution of event schemas.

### 10. Migrate to the `pgx` PostgreSQL Driver (Dependencies)
-   **Action:** Replace the `lib/pq` driver with `jackc/pgx` in the `pkg/database` package.
-   **Implementation:** Update the `go.mod` file and change the driver name in the `sql.Open` call. The `pgx` driver is more actively maintained and offers better performance.

### 11. Centralize and Standardize Configuration (Code Quality)
-   **Action:** Load all configuration for each service at startup in `main.go`.
-   **Implementation:** Create a `Config` struct in each service's `main.go` that holds all environment-dependent values (ports, URLs, etc.). Populate it once and pass it to the `platform.NewApp` constructor.

### 12. Improve Test Quality for `oal-worker` (Testing)
-   **Action:** Refactor the `oal-worker` tests to be more comprehensive.
-   **Implementation:** The tests should assert specific outcomes, such as verifying log messages to ensure the handler logic was executed, not just that it didn't panic.