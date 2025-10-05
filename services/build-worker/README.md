# Build Worker Service

The Build Worker service is a background worker that listens for deployment requests and simulates a build process for an application.

## Functionality

The worker subscribes to a NATS subject for deployment requests. When a message is received, it performs the following actions:

1.  **Receives a `DeploymentRequest` event.** It decodes the message and validates the payload.
2.  **Simulates a build.** It pauses briefly to simulate the time it would take to clone a repository and build a container image.
3.  **Publishes a `BuildSucceeded` event.** Upon successful "build," it publishes a new event to NATS containing the application ID and a simulated container image URI.
4.  **Acknowledges the message.** It uses manual `ack`/`nak`/`term` to ensure reliable message processing.

## NATS Integration

-   **Subscribes to:** `deployment.requested`
-   **Publishes to:** `build.succeeded`

## Running the Service

To run the Build Worker service locally, you will need to have Go installed. You can start the service with the following command from the root of the repository:

```bash
go run ./services/build-worker
```

The service requires a running NATS server to connect to. The connection details are configured via environment variables.