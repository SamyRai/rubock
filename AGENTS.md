This document provides guidelines for AI agents working on the Helios codebase. By following these guidelines, you will help to ensure that the codebase remains clean, consistent, and maintainable.

## Go Development

All Go services in this repository are part of a monorepo. The shared packages are located in the `services/pkg` directory and are included in the `helios` module.

### Dependency Management

Each service has its own `go.mod` file. When adding or updating dependencies for a service, make sure to run `go mod tidy` from within the service's directory to ensure that the `go.sum` file is updated correctly.

To add a new dependency to a Go service, use the following command from the project root:

```bash
go -C <service_path> mod edit -require <package_name>@<version>
```

After modifying the `go.mod` file, run the following command to synchronize the dependencies:

```bash
go -C <service_path> mod tidy
```

### Testing

All new features and bug fixes must be accompanied by tests. The tests should be placed in a file with the `_test.go` suffix in the same package as the code being tested.

To run the tests for a specific service, use the following command from the `services` directory:

```bash
go -C ./<service-name> test ./...
```

When writing tests, use the `testify/assert` package for assertions. This will make the tests more readable and provide more detailed error messages.

### Mocks

When testing code that interacts with external services or dependencies, use mocks to isolate the code under test. Manual mocks are preferred over generated mocks. When creating a manual mock, ensure that the mock implementation correctly records interactions (e.g., function calls, arguments) before returning any simulated error. This allows tests to assert that an action was attempted, even if it failed.

## General Guidelines

-   All code must be formatted with `gofmt`.
-   All new features must be documented.
-   All new services must have a `README.md` file that explains what the service does and how to run it.
-   All new services must have a `Dockerfile` that can be used to build a container image for the service.
-   All new services must have a CI/CD pipeline that runs tests and builds a container image for the service.