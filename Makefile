.PHONY: help build run test docker-up docker-down migrate-up migrate-down clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

docker-up: ## Start Redis and Postgres with docker-compose
	docker-compose up -d
	@echo "Waiting for services to be healthy..."
	@sleep 3

docker-down: ## Stop and remove containers
	docker-compose down

docker-logs: ## Show logs from docker-compose services
	docker-compose logs -f

build: ## Build the gateway binary
	go build -o bin/gateway ./cmd/gateway

run: ## Run the gateway (requires docker-up first)
	go run ./cmd/gateway

test: ## Run all tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

migrate-up: ## Run database migrations up
	@echo "Migration commands will be added when we set up migrate tool"

migrate-down: ## Run database migrations down
	@echo "Migration commands will be added when we set up migrate tool"

clean: ## Clean build artifacts
	rm -rf bin/ coverage.out

