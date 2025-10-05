# Helios

Helios is a powerful, cloud-native Platform-as-a-Service (PaaS) designed to streamline the deployment and management of modern, containerized applications. It provides a robust and scalable environment, abstracting away the complexities of the underlying infrastructure and allowing developers to focus on writing code.

This repository is a monorepo containing all the services and packages that make up the Helios platform.

## Architecture Overview

Helios is built on a microservices architecture, with each service responsible for a specific domain. The services are written in Go and communicate with each other asynchronously using NATS, a lightweight and high-performance messaging system.

The main components of the Helios platform are:

-   **API Service**: The public-facing entry point for all user interactions. It handles API requests, authentication, and orchestration of the other services.
-   **OAL Worker**: The Open Application Language (OAL) worker is responsible for parsing and interpreting OAL files, which define the structure and configuration of user applications.
-   **Build Worker**: This service is responsible for building container images from user-provided source code, based on the instructions in the OAL files.

## Getting Started

To get started with Helios, you will need to have the following installed:

-   Go (version 1.21 or later)
-   Docker
-   PostgreSQL
-   NATS

Once you have all the prerequisites installed, you can clone the repository and run the services. Each service is a standalone Go application that can be run from its respective directory.

```bash
git clone https://github.com/your-repo/helios.git
cd helios

# To run the API service
go run ./services/api

# To run the OAL worker
go run ./services/oal-worker

# To run the build worker
go run ./services/build-worker
```

## Repository Structure

This repository is a Go monorepo that contains all the services and shared packages for the Helios platform. The repository is organized as follows:

```
.
├── cmd
│   └── helios-cli          # The command-line interface for Helios
├── database
│   └── migrations          # Database schema migrations
└── services
    ├── api                 # The main API service
    ├── build-worker        # The service for building container images
    ├── oal-worker          # The service for parsing OAL files
    └── pkg                 # Shared packages used by multiple services
        ├── events          # Shared NATS event definitions
        ├── logger          # Shared logger implementation
        └── testutil        # Test utilities
```