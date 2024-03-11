# Echo Server Makefile
# -----------------------------------------------------
# This Makefile contains targets for building, running,
# and managing the Echo Server project. It simplifies
# common tasks such as compiling the project, running
# the server, and cleaning up build artifacts.
#
# Targets:
#   - build: Compiles the Echo Server into a binary.
#   - run: Starts the Echo Server locally.
#   - run-binary: Builds and runs the compiled binary.
#   - docker-build: Builds a Docker image for the server.
#   - docker-run: Runs the Docker container.
#   - clean: Removes any build artifacts.
#
# Usage:
#   make <target>
#
# Example:
#   make build
# -----------------------------------------------------

# Docker-related variables
DOCKER_IMAGE_NAME=echo-server
DOCKER_CONTAINER_NAME=echo-server-instance

# LOG_LEVEL can be set to debug, info, warn, error, or fatal
LOG_LEVEL ?= info

# PORT can be set to any valid port number
PORT ?= 8080

# Binary output name
BINARY_NAME=echo-server

# Build the application binary
build:
	go build -o ${BINARY_NAME} cmd/$(BINARY_NAME).go


# Makefile target for running the application
run:
	./echo-server -port=${PORT} -logLevel=${LOG_LEVEL}

# Build Docker image
docker-build:
	docker build -t $(DOCKER_IMAGE_NAME) .

# Run Docker container with environment variables passed from Makefile
docker-run: docker-build
	docker run -d --name $(DOCKER_CONTAINER_NAME) -p $(PORT):$(PORT) --env PORT=$(PORT) --env LOG_LEVEL=$(LOG_LEVEL) $(DOCKER_IMAGE_NAME)

# Stop and remove Docker container
docker-clean:
	docker stop $(DOCKER_CONTAINER_NAME)
	docker rm $(DOCKER_CONTAINER_NAME)

.PHONY: build run docker-build docker-run docker-clean
