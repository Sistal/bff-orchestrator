.PHONY: help build run test clean docker-build docker-run

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application binary
	@echo "Building application..."
	@go build -o bin/bff-orchestrator ./cmd/api

run: ## Run the application locally
	@echo "Running application..."
	@go run ./cmd/api/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t bff-orchestrator:latest .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run -p 8080:8080 \
		-e PORT=8080 \
		-e ENVIRONMENT=production \
		-e USER_SERVICE_URL=https://jsonplaceholder.typicode.com \
		-e PRODUCT_SERVICE_URL=https://fakestoreapi.com \
		bff-orchestrator:latest

docker-compose-up: ## Start services with docker-compose
	@echo "Starting services with docker-compose..."
	@docker-compose up -d

docker-compose-down: ## Stop services with docker-compose
	@echo "Stopping services with docker-compose..."
	@docker-compose down

fmt: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./...

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@golangci-lint run

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	@go mod tidy
