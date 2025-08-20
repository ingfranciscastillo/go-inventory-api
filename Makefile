# Inventory API Makefile
.PHONY: help build run clean test migrate seed docker-build docker-run docker-stop

# Variables
BINARY_NAME=inventory-api
DOCKER_IMAGE=inventory-api
DOCKER_TAG=latest

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Development commands
build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) ./cmd/server

run: ## Run the application
	@echo "Starting $(BINARY_NAME)..."
	@go run ./cmd/server/main.go

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@go clean
	@rm -f $(BINARY_NAME)

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Database commands
migrate: ## Run database migrations
	@echo "Running migrations..."
	@go run scripts/migrate/main.go

seed: ## Seed database with example data
	@echo "Seeding database..."
	@go run scripts/seed/main.go

reset-db: migrate seed ## Reset database (migrate + seed)

# Docker commands
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## Run application in Docker
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-stop: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	@docker-compose down

docker-dev: ## Run development environment with Docker Compose
	@echo "Starting development environment..."
	@docker-compose up --build

docker-dev-bg: ## Run development environment in background
	@echo "Starting development environment in background..."
	@docker-compose up -d --build

# Dependencies
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	@go mod tidy
	@go get -u ./...

# Linting and formatting
fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

# Production commands
build-prod: ## Build for production
	@echo "Building for production..."
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BINARY_NAME) ./cmd/server

# Installation commands
install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# API testing
api-test: ## Test API endpoints (requires server to be running)
	@echo "Testing API endpoints..."
	@curl -s http://localhost:8080/health | jq .
	@echo "\nRegistering test user..."
	@curl -s -X POST http://localhost:8080/auth/register \
		-H "Content-Type: application/json" \
		-d '{"email":"test@example.com","password":"test123"}' | jq .

# Cleanup commands
clean-docker: ## Clean Docker images and containers
	@echo "Cleaning Docker..."
	@docker system prune -f
	@docker image prune -f

clean-all: clean clean-docker ## Clean everything

# Development setup
setup: deps migrate seed ## Setup development environment
	@echo "Development environment setup complete!"
	@echo "Run 'make run' to start the server"

# Quick start
start: build run ## Build and run the application

# Documentation
docs: ## Generate API documentation (if using tools like swag)
	@echo "Generating documentation..."
	@echo "Documentation would be generated here..."

# Version
version: ## Show version information
	@echo "Go version: $(shell go version)"
	@echo "Git commit: $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
	@echo "Build date: $(shell date -u +%Y-%m-%dT%H:%M:%SZ)"