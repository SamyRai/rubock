# OAL Worker Service

The OAL (Observe, Analyze, and Log) Worker service is a background worker that listens for successful build events and simulates the final deployment step of an application.

## Functionality

The worker subscribes to a NATS subject for successful build events. When a message is received, it performs the following actions:

1.  **Receives a `BuildSucceeded` event.** It decodes the message and validates the payload.
2.  **Simulates a deployment.** It pauses briefly to simulate the time it would take to deploy a container image to the platform.
3.  **Logs the end of the workflow.** In the current implementation, this worker represents the end of the deployment pipeline and logs a final message.
4.  **Acknowledges the message.** It uses manual `ack`/`term` to ensure reliable message processing.

## NATS Integration

-   **Subscribes to:** `build.succeeded`
-   **Publishes to:** None. This worker is the current end of the event chain.

## Running the Service

To run the OAL Worker service locally, you will need to have Go installed. You can start the service with the following command from the root of the repository:

```bash
go run ./services/oal-worker
```

The service requires a running NATS server to connect to. The connection details are configured via environment variables.