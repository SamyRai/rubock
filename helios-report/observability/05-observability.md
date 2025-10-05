# Observability: A System Flying Blind

The Helios platform has a solid foundation for logging but is critically lacking in the other two pillars of observability: metrics and tracing. This means that while it's possible to see what happened in a single service (via logs), it's nearly impossible to understand the overall health of the system or debug performance issues in a distributed context.

## Strengths

### 1. Standardized, Structured Logging

The shared `pkg/logger` package provides a standardized `zerolog` instance to all services. The ability to switch between a human-readable console format for development and a structured JSON format for production is an excellent feature.

**Code Evidence (`services/pkg/logger/logger.go`):**
```go
// New creates and configures a new zerolog.Logger instance.
func New() zerolog.Logger {
	env := os.Getenv("ENV")
	if env == "development" {
		// Use a pretty, colorized console writer for local development.
		return log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	// Use a structured JSON logger in production.
	return zerolog.New(os.Stderr).With().Timestamp().Logger()
}
```
-   **Impact:** This provides a consistent and machine-readable logging format across the platform, which is the necessary first step for effective observability.

## Critical Weaknesses

### 2. No Metrics Collection

There is no evidence of any metrics collection in the codebase. None of the services expose a metrics endpoint (e.g., `/metrics` for Prometheus), and no metrics libraries are included as dependencies.

-   **Impact:** The platform is a black box from a performance and health perspective. It is impossible to answer critical operational questions like:
    -   What is the average request latency for the `api` service?
    -   What is the error rate for creating applications?
    -   How many messages per second is the `build-worker` processing?
    -   Is CPU or memory usage for a service approaching its limit?
    Without metrics, there is no way to set up proactive alerting, create dashboards to monitor system health, or identify performance bottlenecks. This is a critical gap for any production system.

### 3. No Distributed Tracing

The codebase lacks any form of distributed tracing. There are no tracing libraries (like OpenTelemetry) integrated into the HTTP handlers or NATS message processors.

-   **Impact:** In a distributed system like Helios, a single user request can trigger a chain of events across multiple services (`api` -> NATS -> `build-worker` -> NATS -> `oal-worker`). Without distributed tracing, it is extremely difficult to debug issues that span service boundaries. For example, if a deployment is slow, it's impossible to tell whether the bottleneck is in the `api` service, the `build-worker`, or the messaging queue itself. This dramatically increases the mean time to resolution (MTTR) for production incidents.

### 4. No Centralized Logging System Mentioned

While the services produce structured logs, there is no mention of a system to collect, aggregate, and search these logs from all services in a centralized location.

-   **Impact:** Without a centralized logging platform (like the ELK stack, Grafana Loki, or Datadog), operators are forced to inspect logs on a per-service or per-pod basis (e.g., `kubectl logs <pod-name>`). This is highly inefficient and makes it nearly impossible to correlate events across different services when debugging an issue.