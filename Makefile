# Echo Server Makefile
# -----------------------------------------------------
# This Makefile contains targets for building, running,
# and managing the Echo Server project.
#
# Usage:
#   make <target>
#
# Run 'make' or 'make help' to see available targets.
# -----------------------------------------------------

# Default target
.DEFAULT_GOAL := help

# Docker-related variables
DOCKER_IMAGE_NAME=echo-server
DOCKER_CONTAINER_NAME=echo-server-instance

# LOG_LEVEL can be set to debug, info, warn, error, or fatal
LOG_LEVEL ?= info

# PORT can be set to any valid port number
PORT ?= 8080

# Binary output name
BINARY_NAME=echo-server

##@ General

.PHONY: help
help: ## Show this help message
	@echo "Echo Server - Available targets:"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n\n"} \
		/^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } \
		/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) }' $(MAKEFILE_LIST)
	@echo ""
	@echo "Configuration (override with environment variables):"
	@echo "  PORT=$(PORT)  LOG_LEVEL=$(LOG_LEVEL)"

##@ Development

.PHONY: build
build: ## Build the application binary
	go build -o ${BINARY_NAME} cmd/$(BINARY_NAME).go

.PHONY: run
run: ## Run the application locally
	./echo-server -port=${PORT} -logLevel=${LOG_LEVEL}

.PHONY: test
test: ## Run unit tests
	go test ./...

.PHONY: test-verbose
test-verbose: ## Run unit tests with verbose output
	go test -v ./...

.PHONY: lint
lint: ## Run golangci-lint (installs if not found)
	@which golangci-lint > /dev/null 2>&1 || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	$(shell go env GOPATH)/bin/golangci-lint run ./...

.PHONY: fmt
fmt: ## Format Go source files
	go fmt ./...

.PHONY: coverage
coverage: ## Run tests with coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	@rm -f coverage.out

.PHONY: coverage-html
coverage-html: ## Generate HTML coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

##@ Docker

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE_NAME) .

.PHONY: docker-run
docker-run: docker-build ## Build and run Docker container
	docker run -d --name $(DOCKER_CONTAINER_NAME) -p $(PORT):$(PORT) --env PORT=$(PORT) --env LOG_LEVEL=$(LOG_LEVEL) $(DOCKER_IMAGE_NAME)

.PHONY: docker-clean
docker-clean: ## Stop and remove Docker container
	docker stop $(DOCKER_CONTAINER_NAME)
	docker rm $(DOCKER_CONTAINER_NAME)
