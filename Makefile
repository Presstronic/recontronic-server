.PHONY: help build test lint clean docker-build docker-up docker-down k8s-deploy k8s-delete dev run migrate

# Variables
APP_NAME=recontronic-server
API_BINARY=bin/api
WORKER_BINARY=bin/worker
CLI_BINARY=bin/cli
DOCKER_REGISTRY=ghcr.io
IMAGE_TAG?=latest
GO_VERSION=1.25.3

# Colors for output
BLUE=\033[0;34m
GREEN=\033[0;32m
RED=\033[0;31m
NC=\033[0m # No Color

## help: Display this help message
help:
	@echo "$(BLUE)Recontronic Server - Available Commands:$(NC)"
	@echo ""
	@grep -E '^## ' Makefile | sed 's/## /  /'
	@echo ""

## build: Build all binaries
build: build-api build-worker
	@echo "$(GREEN)✓ All binaries built successfully$(NC)"

## build-api: Build API server binary
build-api:
	@echo "$(BLUE)Building API server...$(NC)"
	@go build -o $(API_BINARY) ./cmd/api
	@echo "$(GREEN)✓ API server built: $(API_BINARY)$(NC)"

## build-worker: Build worker binary
build-worker:
	@echo "$(BLUE)Building worker...$(NC)"
	@go build -o $(WORKER_BINARY) ./cmd/worker
	@echo "$(GREEN)✓ Worker built: $(WORKER_BINARY)$(NC)"

## build-cli: Build CLI binary
build-cli:
	@echo "$(BLUE)Building CLI...$(NC)"
	@go build -o $(CLI_BINARY) ./cmd/cli
	@echo "$(GREEN)✓ CLI built: $(CLI_BINARY)$(NC)"

## test: Run all tests
test:
	@echo "$(BLUE)Running tests...$(NC)"
	@go test -v -race -coverprofile=coverage.out ./...
	@echo "$(GREEN)✓ Tests completed$(NC)"

## test-coverage: Run tests with coverage report
test-coverage: test
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(NC)"

## lint: Run linter
lint:
	@echo "$(BLUE)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(BLUE)Running staticcheck...$(NC)"
	@which staticcheck > /dev/null 2>&1 || (echo "$(BLUE)Installing staticcheck...$(NC)" && go install honnef.co/go/tools/cmd/staticcheck@latest)
	@$(shell go env GOPATH)/bin/staticcheck ./...
	@echo "$(GREEN)✓ Linting completed$(NC)"

## fmt: Format code
fmt:
	@echo "$(BLUE)Formatting code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

## vet: Run go vet
vet:
	@echo "$(BLUE)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✓ Vet completed$(NC)"

## clean: Clean build artifacts
clean:
	@echo "$(BLUE)Cleaning...$(NC)"
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)✓ Cleaned$(NC)"

## deps: Download dependencies
deps:
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✓ Dependencies downloaded$(NC)"

## docker-build: Build Docker images
docker-build:
	@echo "$(BLUE)Building Docker images...$(NC)"
	@docker build -f deployments/docker/Dockerfile -t $(APP_NAME)-api:$(IMAGE_TAG) .
	@docker build -f deployments/docker/Dockerfile.worker -t $(APP_NAME)-worker:$(IMAGE_TAG) .
	@echo "$(GREEN)✓ Docker images built$(NC)"

## docker-up: Start Docker Compose environment
docker-up:
	@echo "$(BLUE)Starting Docker Compose...$(NC)"
	@cd deployments/docker && docker-compose up -d
	@echo "$(GREEN)✓ Docker Compose started$(NC)"

## docker-down: Stop Docker Compose environment
docker-down:
	@echo "$(BLUE)Stopping Docker Compose...$(NC)"
	@cd deployments/docker && docker-compose down
	@echo "$(GREEN)✓ Docker Compose stopped$(NC)"

## docker-logs: View Docker Compose logs
docker-logs:
	@cd deployments/docker && docker-compose logs -f

## k8s-deploy: Deploy to Kubernetes
k8s-deploy:
	@echo "$(BLUE)Deploying to Kubernetes...$(NC)"
	@kubectl apply -f deployments/k8s/base/
	@echo "$(GREEN)✓ Deployed to Kubernetes$(NC)"

## k8s-delete: Delete from Kubernetes
k8s-delete:
	@echo "$(BLUE)Deleting from Kubernetes...$(NC)"
	@kubectl delete -f deployments/k8s/base/
	@echo "$(GREEN)✓ Deleted from Kubernetes$(NC)"

## k8s-status: Check Kubernetes deployment status
k8s-status:
	@echo "$(BLUE)Kubernetes Status:$(NC)"
	@kubectl get all -n recontronic

## dev: Run in development mode
dev: build-api
	@echo "$(BLUE)Starting development server...$(NC)"
	@./$(API_BINARY)

## run-api: Run API server
run-api: build-api
	@./$(API_BINARY)

## run-worker: Run worker
run-worker: build-worker
	@./$(WORKER_BINARY)

## migrate-up: Run database migrations up
migrate-up:
	@echo "$(BLUE)Running migrations...$(NC)"
	@go run cmd/migrate/main.go up
	@echo "$(GREEN)✓ Migrations completed$(NC)"

## migrate-down: Rollback database migrations
migrate-down:
	@echo "$(BLUE)Rolling back migrations...$(NC)"
	@go run cmd/migrate/main.go down
	@echo "$(GREEN)✓ Rollback completed$(NC)"

## migrate-create: Create new migration (usage: make migrate-create name=create_users_table)
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "$(RED)Error: name parameter required. Usage: make migrate-create name=your_migration_name$(NC)"; \
		exit 1; \
	fi
	@echo "$(BLUE)Creating migration: $(name)$(NC)"
	@mkdir -p migrations
	@timestamp=$$(date +%Y%m%d%H%M%S); \
	touch migrations/$${timestamp}_$(name).up.sql migrations/$${timestamp}_$(name).down.sql
	@echo "$(GREEN)✓ Migration files created$(NC)"

## proto-gen: Generate code from proto files
proto-gen:
	@echo "$(BLUE)Generating protobuf code...$(NC)"
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/*.proto
	@echo "$(GREEN)✓ Protobuf code generated$(NC)"

## install-tools: Install development tools
install-tools:
	@echo "$(BLUE)Installing development tools...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "$(GREEN)✓ Tools installed$(NC)"

## ci: Run CI checks (lint, test, build)
ci: lint test build
	@echo "$(GREEN)✓ CI checks passed$(NC)"

## version: Display version information
version:
	@echo "$(BLUE)Recontronic Server$(NC)"
	@echo "Go version: $(GO_VERSION)"
	@go version
