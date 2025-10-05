# Helios Development Roadmap

This document tracks the major development tasks for the Helios project. It is organized by priority, from highest to lowest.

## High Priority

-   [ ] **Implement a Standardized Application Wrapper:**
    -   Create a shared package or struct (e.g., `pkg/app`) that standardizes service initialization and shutdown.
    -   This wrapper should handle:
        -   Loading configuration.
        -   Initializing the logger.
        -   Establishing database and NATS connections.
        -   Implementing graceful shutdown by handling `SIGINT` and `SIGTERM` signals.
    -   Refactor all existing services (`api`, `build-worker`, `oal-worker`) to use this new wrapper, removing redundant boilerplate from their `main.go` files.

-   [ ] **Establish CI/CD Pipeline:**
    -   Create a GitHub Actions workflow file (e.g., `.github/workflows/ci.yml`).
    -   The pipeline should trigger on pushes and pull requests to the `main` branch.
    -   It must run `make tidy` and `make test` to ensure all Go modules are clean and all tests pass before code can be merged.

## Medium Priority

-   [ ] **Increase Test Coverage:**
    -   Write comprehensive unit and integration tests for the core logic of each service.
    -   **`api` service:** Test HTTP handlers, middleware, and request/response validation.
    -   **`build-worker` & `oal-worker`:** Test NATS message handlers and business logic.
    -   Use manual mocks for dependencies like NATS publishers and the database, following the guidelines in `AGENTS.md`.

-   [ ] **Implement Database Migration Tooling:**
    -   Integrate a database migration tool (e.g., `golang-migrate/migrate`).
    -   Add commands to the `Makefile` to create new migrations and apply them.
    -   Create an initial migration that defines the schema required by the existing services.

## Low Priority

-   [ ] **Develop the Helios CLI:**
    -   Begin implementation of the `helios-cli` in the `cmd/helios-cli` directory.
    -   The initial version should support basic commands for user authentication and application deployment.

-   [ ] **Add Service-Level Metrics:**
    -   Integrate a Prometheus client into the standardized application wrapper.
    -   Expose basic metrics for each service (e.g., HTTP request latency, NATS message processing time, error rates).
-   [ ] **Refine Shared `pkg/events`:**
    -   Ensure all NATS subjects are versioned (e.g., `v1.deployment.requested`).
    -   Define and document all event schemas used for inter-service communication.