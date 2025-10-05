# Dockerfile for the build-worker service

# --- Build Stage ---
# Use the official Go image as a builder.
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container.
WORKDIR /app

# Copy all go.mod and go.sum files to leverage Docker layer caching.
# This is a common pattern for monorepos. We copy all of them to ensure
# that `go mod download` can resolve the `replace` directives.
COPY services/go.mod services/go.sum ./services/
COPY services/api/go.mod services/api/go.sum ./services/api/
COPY services/build-worker/go.mod services/build-worker/go.sum ./services/build-worker/
COPY services/oal-worker/go.mod services/oal-worker/go.sum ./services/oal-worker/

# Change to a subdirectory to run `go mod download`
# This is necessary because the parent `services` directory does not have a go.mod
# that includes all the children. We run it in `build-worker` which will download its own
# dependencies and the shared ones from `pkg`.
RUN cd services/build-worker && go mod download

# Copy the rest of the source code.
COPY . .

# Build the Go application binary.
# CGO_ENABLED=0 is used to build a statically linked binary.
# -o specifies the output file name.
RUN CGO_ENABLED=0 go build -o /app/build-worker-server ./services/build-worker

# --- Final Stage ---
# Use a minimal, non-root base image for the final container.
FROM alpine:latest

# Set the working directory.
WORKDIR /app

# Copy the compiled binary from the builder stage.
COPY --from=builder /app/build-worker-server .

# Set the command to run the application.
CMD ["./build-worker-server"]