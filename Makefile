# Makefile for the Helios monorepo

# Define the services directory
SERVICES_DIR := ./services

# Find all Go services (directories with a go.mod file)
SERVICES := $(shell find $(SERVICES_DIR) -mindepth 1 -maxdepth 1 -type d)

.PHONY: all test tidy

all: test

# Run tests for all services
test:
	@echo "Running tests for all services..."
	@for service in $(SERVICES); do \
		if [ -f "$$service/go.mod" ]; then \
			echo "--> Testing $$service"; \
			go -C $$service test ./...; \
		fi \
	done

# Tidy go.mod files for all services
tidy:
	@echo "Tidying go.mod files for all services..."
	@for service in $(SERVICES); do \
		if [ -f "$$service/go.mod" ]; then \
			echo "--> Tidying $$service"; \
			go -C $$service mod tidy; \
		fi \
	done