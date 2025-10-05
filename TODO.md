# Helios: Strategic Project & Development Plan

## 1. Vision & Mission

*   **Vision:** To empower development teams of all sizes to deploy and manage their applications effortlessly, allowing them to focus on innovation rather than infrastructure.
*   **Mission:** To build a robust, intuitive, and scalable Platform-as-a-Service (PaaS) that automates the entire application lifecycle, from build to deployment and operations, with a relentless focus on developer experience and system reliability.

## 2. Core Development Principles

*   **Developer-Centric:** The developer is our primary customer. Every feature, CLI command, and API endpoint will be designed for clarity, simplicity, and efficiency.
*   **Secure by Design:** Security is not an afterthought. It is a foundational requirement, integrated into every stage of the development lifecycle, from dependency scanning in CI to role-based access control (RBAC) in the API.
*   **Resilient by Default:** Services must be fault-tolerant. We will engineer for failure, with robust error handling, graceful degradation, and automated recovery mechanisms built into the core architecture.
*   **Operationally Excellent:** A platform we build must be a platform we can operate. We will prioritize comprehensive logging, metrics, and tracing to ensure deep visibility into system health and performance.

---

## 3. Strategic Roadmap

### Phase 1: Foundational Hardening & Automation (Timeline: Current - Next 3 Sprints)

**Primary Goal:** To forge a stable, secure, and automated foundation, eliminating architectural inconsistencies and enabling rapid, safe development in subsequent phases.

*   **Initiative 1.1: Standardized Service Architecture**
    *   `[ ]` **Epic:** Implement a shared `pkg/app` wrapper for a consistent service lifecycle.
        *   `[ ]` Story: The wrapper must handle initialization of config, logging, and connections.
        *   `[ ]` Story: The wrapper must implement graceful shutdown on `SIGINT`/`SIGTERM` signals.
    *   `[ ]` **Epic:** Refactor all existing services (`api`, `build-worker`, `oal-worker`) to use the `pkg/app` wrapper, removing all boilerplate from their `main.go` files.

*   **Initiative 1.2: Continuous Integration & Quality Gates**
    *   `[ ]` **Epic:** Establish a CI pipeline in GitHub Actions.
        *   `[ ]` Story: The pipeline must run `make tidy` and `make test` automatically on all pushes and pull requests to `main`.
        *   `[ ]` Story: The pipeline must block merging if any checks fail.

*   **Initiative 1.3: Deterministic Database Management**
    *   `[ ]` **Epic:** Integrate a database migration tool (`golang-migrate/migrate`).
        *   `[ ]` Story: Create `Makefile` targets for creating and applying migrations (`make db/migrate-up`, `make db/new-migration NAME=...`).
        *   `[ ]` Story: Define the initial, complete database schema in the first migration script.

*   **Success Metrics for Phase 1:**
    *   All services use the `pkg/app` wrapper.
    *   CI pipeline is 100% mandatory for all merges to `main`.
    *   The database can be fully provisioned from zero using a single migration command.

---

### Phase 2: Core Product MVP & Observability (Timeline: Next Quarter)

**Primary Goal:** To deliver the core user value proposition—deploying an application—and to build the observability stack required to operate the platform reliably.

*   **Initiative 2.1: The Developer's Toolkit (CLI v1.0)**
    *   `[ ]` **Epic:** Implement the core user journey in the `helios-cli`.
        *   `[ ]` Story: `helios login` - Authenticate the user.
        *   `[ ]` Story: `helios deploy` - Deploy an application from an OAL file.
        *   `[ ]` Story: `helios status <app-name>` - View deployment status and logs.

*   **Initiative 2.2: Foundational Observability**
    *   `[ ]` **Epic:** Integrate a full observability stack.
        *   `[ ]` Story: Integrate a Prometheus client into `pkg/app` to expose standard service metrics (HTTP latencies, error rates, queue depths).
        *   `[ ]` Story: Establish a centralized logging sink (e.g., ELK stack or Grafana Loki).
        *   `[ ]` Story: Set up a Grafana instance with initial dashboards for monitoring service health.

*   **Initiative 2.3: Bulletproof Business Logic**
    *   `[ ]` **Epic:** Achieve comprehensive test coverage.
        *   `[ ]` Story: Write unit tests for all business-critical logic in each service, aiming for >90% coverage.
        *   `[ ]` Story: Write integration tests for the full deployment flow (API -> NATS -> Worker).

*   **Success Metrics for Phase 2:**
    *   A user can deploy a simple "hello-world" application using only the CLI.
    *   Core service health is visible on a production-ready Grafana dashboard.
    *   Test coverage metrics are met and enforced by the CI pipeline.

---

### Phase 3: Enterprise Readiness & Scale (Timeline: Next 6-12 Months)

**Primary Goal:** To evolve Helios into a secure, multi-tenant, and highly available platform capable of supporting enterprise workloads and advanced deployment scenarios.

*   **Initiative 3.1: Security & Identity Management**
    *   `[ ]` **Epic:** Implement a robust security model.
        *   `[ ]` Story: Design and implement a full Role-Based Access Control (RBAC) system.
        *   `[ ]` Story: Integrate a centralized secrets management solution (e.g., HashiCorp Vault) for handling sensitive application data.

*   **Initiative 3.2: Advanced Deployment Workflows**
    *   `[ ]` **Epic:** Support modern deployment strategies.
        *   `[ ]` Story: Implement a blue-green deployment controller.
        *   `[ ]` Story: Investigate and prototype support for canary releases.

*   **Initiative 3.3: User & Management Interface**
    *   `[ ]` **Epic:** Build a web-based user dashboard.
        *   `[ ]` Story: Design a UI for visualizing application status, logs, and metrics.
        *   `[ ]` Story: Allow users to manage their applications (deploy, restart, delete) from the UI.

*   **Success Metrics for Phase 3:**
    *   Administrators can define and enforce user permissions via RBAC.
    *   Users can perform zero-downtime deployments using the blue-green strategy.
    *   Users can manage the entire application lifecycle through the web dashboard.