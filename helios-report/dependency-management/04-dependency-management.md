# Dependency Management: Inconsistencies and Risks

The project uses Go Modules and `replace` directives, which is a standard and effective strategy for a Go monorepo. However, a closer look at the dependencies reveals several risks related to outdated packages, inconsistent versions, and a lack of automated management.

## Key Issues

### 1. Outdated and Unmaintained PostgreSQL Driver

The shared `database` package relies on `github.com/lib/pq`, which is officially in maintenance mode. The community has largely moved to more actively maintained and performant drivers.

**Code Evidence (`services/pkg/database/database.go`):**
```go
import (
	// ...
	_ "github.com/lib/pq" // PostgreSQL driver
)
```
-   **Impact:** Using a driver that is no longer actively developed means bug fixes and performance improvements will be missed. It also poses a potential security risk if new vulnerabilities are discovered but not patched. The recommended replacement is `github.com/jackc/pgx`.

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