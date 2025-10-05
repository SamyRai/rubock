# API Service

The API service is the main entry point for the Helios platform. It exposes an HTTP API for managing projects and applications.

## Running the Service

To run the API service locally, you will need to have Go installed. You can start the service with the following command from the root of the repository:

```bash
go run ./services/api
```

The service will start on port `8080` by default. You can change the port by setting the `PORT` environment variable.

## API Endpoints

### Create a New Project

*   **Endpoint:** `POST /projects`
*   **Description:** Simulates the creation of a new project.
*   **Response:**
    *   `201 Created` with a JSON body containing the new project's ID and name.

### Create a New Application

*   **Endpoint:** `POST /applications`
*   **Description:** Creates a new application and triggers a deployment by publishing a message to NATS.
*   **Request Body:**
    ```json
    {
      "name": "my-cool-app",
      "git_repository": "https://github.com/user/repo.git",
      "git_branch": "main"
    }
    ```
*   **Response:**
    *   `202 Accepted` with a JSON body containing the new application's ID, name, and a "pending" status.
    *   `400 Bad Request` if the request body is invalid.
    *   `500 Internal Server Error` if the service fails to publish the deployment event to NATS.