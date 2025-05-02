# Variables
APP_NAME := sepay-service
VERSION := 1.0.0
BUILD_DIR := build
MAIN_PKG := ./cmd/server
ENV_FILE := .env

# Go commands
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := $(GOCMD) fmt
GOLINT := golangci-lint
GOTOOL := $(GOCMD) tool
GOVET := $(GOCMD) vet

# Docker commands
DOCKER := docker
DOCKERFILE := Dockerfile
DOCKER_IMAGE := sepay-integration
DOCKER_TAG := latest

# Git information
GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
BUILD_FLAGS := -ldflags "-X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.GitBranch=$(GIT_BRANCH) -X main.BuildTime=$(BUILD_TIME)"

.PHONY: all build clean test coverage lint fmt vet mod-download mod-tidy run docker-build docker-run help vendor

# Default target
all: clean test build

# Build binary
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PKG)
	@echo "Binary built at $(BUILD_DIR)/$(APP_NAME)"

# Build for production (stripped binaries)
build-prod:
	@echo "Building production binary..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) -ldflags="-w -s $(BUILD_FLAGS)" -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PKG)
	@echo "Production binary built at $(BUILD_DIR)/$(APP_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)
	@echo "Done!"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -race -v ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -race -coverprofile=coverage.out ./...
	$(GOTOOL) cover -func=coverage.out
	@echo "Generate HTML coverage report..."
	$(GOTOOL) cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	$(GOLINT) run

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Vet code
vet:
	@echo "Vetting code..."
	$(GOVET) ./...

# Download dependencies
mod-download:
	@echo "Downloading dependencies..."
	$(GOMOD) download

# Tidy dependencies
mod-tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

# Download dependencies into vendor folder
mod-vendor:
	@echo "Downloading dependencies into vendor folder for local development..."
	$(GOMOD) vendor
	@echo "Vendor folder created successfully!"

# Run application locally
run:
	@echo "Starting application..."
	@if [ -f $(ENV_FILE) ]; then \
		set -a; . $(ENV_FILE); set +a; \
		$(GOCMD) run $(MAIN_PKG); \
	else \
		$(GOCMD) run $(MAIN_PKG); \
	fi

# Run application with hot reload
run-dev:
	@echo "Starting application with hot reload..."
	@if command -v air > /dev/null; then \
		air -c .air.toml; \
	else \
		echo "Air is not installed. Installing..."; \
		$(GOGET) -u github.com/cosmtrek/air; \
		air -c .air.toml; \
	fi

# Create database schema
db-schema:
	@echo "Creating database schema..."
	@if [ -f $(ENV_FILE) ]; then \
		export $$(cat $(ENV_FILE) | xargs) && mysql -u$$DB_USER -p$$DB_PASSWORD -h$$DB_HOST $$DB_NAME < ./scripts/schema.sql; \
	else \
		echo "No .env file found. Please create one with DB_USER, DB_PASSWORD, DB_HOST, and DB_NAME variables."; \
		exit 1; \
	fi

# Create test data
db-seed:
	@echo "Seeding test data..."
	@if [ -f $(ENV_FILE) ]; then \
		export $$(cat $(ENV_FILE) | xargs) && mysql -u$$DB_USER -p$$DB_PASSWORD -h$$DB_HOST $$DB_NAME < ./scripts/seed.sql; \
	else \
		echo "No .env file found. Please create one with DB_USER, DB_PASSWORD, DB_HOST, and DB_NAME variables."; \
		exit 1; \
	fi

# Build Docker image
docker-build:
	@echo "Building Docker image $(DOCKER_IMAGE):$(DOCKER_TAG)..."
	$(DOCKER) build -t $(DOCKER_IMAGE):$(DOCKER_TAG) -f $(DOCKERFILE) .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	$(DOCKER) run -p 8080:8080 --env-file $(ENV_FILE) $(DOCKER_IMAGE):$(DOCKER_TAG)

# Create .env file from .env.example
env-setup:
	@if [ ! -f .env ]; then \
		if [ -f .env.example ]; then \
			echo "Creating .env file from .env.example..."; \
			cp .env.example .env; \
			echo ".env file created. Please update it with your settings."; \
		else \
			echo "No .env.example file found. Creating basic .env file..."; \
			echo "# Server configuration\nSERVER_PORT=8080\nSERVER_READ_TIMEOUT=10\nSERVER_WRITE_TIMEOUT=10\nSERVER_SHUTDOWN_TIMEOUT=5\n\n# Database configuration\nDB_DRIVER=mysql\nDB_HOST=localhost\nDB_PORT=3306\nDB_USER=root\nDB_PASSWORD=password\nDB_NAME=sepay\nDB_MAX_OPEN_CONNS=10\nDB_MAX_IDLE_CONNS=5\n\n# Sepay configuration\nSEPAY_API_KEY=your_api_key_here\nSEPAY_BANK_ID=your_bank_id_here\nSEPAY_ACCOUNT_NUMBER=your_account_number_here\nSEPAY_ACCOUNT_NAME=your_account_name_here\nSEPAY_WEBHOOK_SECRET=your_webhook_secret_here\nSEPAY_WEBHOOK_BASE_URL=https://api.example.com" > .env; \
			echo "Basic .env file created. Please update it with your settings."; \
		fi; \
	else \
		echo ".env file already exists."; \
	fi

# Generate mock files for testing
gen-mocks:
	@echo "Generating mocks..."
	@if ! command -v mockgen > /dev/null; then \
		echo "Installing mockgen..."; \
		$(GOGET) -u github.com/golang/mock/mockgen; \
	fi
	@mkdir -p internal/mock
	mockgen -source=internal/domain/repository/order_repository.go -destination=internal/mock/order_repository_mock.go -package=mock
	mockgen -source=internal/domain/repository/transaction_repository.go -destination=internal/mock/transaction_repository_mock.go -package=mock
	@echo "Mocks generated successfully"

# Show help
help:
	@echo "Available targets:"
	@echo "  all            - Clean, test and build the application"
	@echo "  build          - Build the application binary"
	@echo "  build-prod     - Build production-ready binary"
	@echo "  clean          - Remove build artifacts"
	@echo "  test           - Run tests"
	@echo "  coverage       - Run tests with coverage report"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  vet            - Vet code"
	@echo "  mod-download   - Download dependencies"
	@echo "  mod-tidy       - Tidy dependencies"
	@echo "  mod-vendor     - Download dependencies into vendor folder for local development"
	@echo "  run            - Run the application locally"
	@echo "  run-dev        - Run with hot reload (requires air)"
	@echo "  db-schema      - Create database schema"
	@echo "  db-seed        - Seed test data"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  env-setup      - Setup .env file"
	@echo "  gen-mocks      - Generate mock files for testing"
	@echo "  help           - Show this help message"