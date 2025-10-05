# Dependency Management: Inconsistencies and Risks

The project uses Go Modules and `replace` directives, which is a standard and effective strategy for a Go monorepo. However, a closer look at the dependencies reveals several risks related to outdated packages, inconsistent versions, and a lack of automated management.

## Key Issues

### 1. COMPLETED: PostgreSQL Driver Migration

The project has successfully migrated from the unmaintained `github.com/lib/pq` driver to the modern, performant `github.com/jackc/pgx`. This action has resolved the previously identified security and maintenance risks.

**Code Evidence (`services/pkg/database/database.go`):**
```go
import (
	// ...
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
)
```
-   **Impact:** This migration to an actively maintained driver ensures the project benefits from ongoing bug fixes, performance improvements, and security patches, aligning with current best practices.

### 2. Inconsistent Go and Toolchain Versions

The `go.mod` file for the `api` service specifies conflicting Go language and toolchain versions. This can lead to subtle build inconsistencies and makes the development environment less predictable.

**Code Evidence (`services/api/go.mod`):**
```
module helios/api

go 1.23.0 // <-- Go language version

toolchain go1.24.3 // <-- Go toolchain version
```
-   **Impact:** Different developers or CI/CD systems might resolve these versions differently, leading to "works on my machine" issues. The `go` and `toolchain` versions should be aligned across all modules in the repository to ensure reproducible builds.

### 3. No Automated Dependency Updates or Vulnerability Scanning

The `AGENTS.md` file implies a manual process for dependency updates. The repository lacks automated tools for two critical functions:
1.  **Keeping dependencies up-to-date.**
2.  **Scanning dependencies for known vulnerabilities.**

-   **Impact:**
    -   A manual update process is slow and error-prone, often leading to dependencies becoming stale and exposing the project to security risks. A tool like **Dependabot** should be configured to automatically create pull requests for dependency updates.
    -   Without vulnerability scanning, the project is blind to potential security holes in its third-party code. A tool like **`govulncheck`** (from the official Go team) or **Snyk** should be integrated into the CI/CD pipeline to fail builds when critical vulnerabilities are detected.