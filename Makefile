# KFn FaaS Platform - Makefile
# Provides build, test, and development targets for the KFn platform

# Variables
BINARY_DIR := bin
CMD_DIR := cmd
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Go tools
GOFMT := gofmt
GOLINT := golangci-lint
GOTEST := go test
GOVET := go vet

# Binaries
FAASCTL := $(BINARY_DIR)/faasctl
SCHEDULER := $(BINARY_DIR)/scheduler
API_GATEWAY := $(BINARY_DIR)/api-gateway

# Build targets
.PHONY: all build clean test lint vet fmt check-fmt

all: build

# Create bin directory if it doesn't exist
$(BINARY_DIR):
	mkdir -p $(BINARY_DIR)

# Build all binaries
build: $(BINARY_DIR) faasctl scheduler api-gateway

# Build faasctl CLI
faasctl: $(BINARY_DIR)
	@echo "Building faasctl..."
	@go build $(LDFLAGS) -o $(FAASCTL) $(CMD_DIR)/faasctl

# Build scheduler
scheduler: $(BINARY_DIR)
	@echo "Building scheduler..."
	@go build $(LDFLAGS) -o $(SCHEDULER) $(CMD_DIR)/scheduler

# Build API gateway
api-gateway: $(BINARY_DIR)
	@echo "Building api-gateway..."
	@go build $(LDFLAGS) -o $(API_GATEWAY) $(CMD_DIR)/api-gateway

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BINARY_DIR)
	@go clean

# Run tests
test:
	@echo "Running tests..."
	@$(GOTEST) -v -race -coverprofile=coverage.out ./...

# Run tests with coverage report
test-coverage: test
	@echo "Generating coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out | tail -1

# Run linting
lint:
	@echo "Running linter..."
	@$(GOLINT) run ./...

# Run go vet
vet:
	@echo "Running go vet..."
	@$(GOVET) ./...

# Format code
fmt:
	@echo "Formatting code..."
	@$(GOFMT) -s -w .

# Check if code is formatted
check-fmt:
	@echo "Checking code formatting..."
	@test -z "$$($(GOFMT) -l .)" || (echo "Code is not formatted. Run 'make fmt'." && exit 1)

# Run all checks
check: fmt vet lint test

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Development targets
.PHONY: dev faasctl-run scheduler-run api-gateway-run

# Run faasctl
faasctl-run: faasctl
	@$(FAASCTL) $(ARGS)

# Run scheduler
scheduler-run: scheduler
	@$(SCHEDULER) $(ARGS)

# Run API gateway
api-gateway-run: api-gateway
	@$(API_GATEWAY) $(ARGS)

# Docker targets
.PHONY: docker-build docker-push docker-build-all

# Build Docker image for a specific component
docker-build:
	@echo "Building Docker image for $(COMPONENT)..."
	@docker build -t isat/kfn-$(COMPONENT):$(VERSION) -f deploy/docker/$(COMPONENT).Dockerfile .

# Push Docker image
docker-push: docker-build
	@echo "Pushing Docker image..."
	@docker push isat/kfn-$(COMPONENT):$(VERSION)

# Build all Docker images
docker-build-all:
	@for component in scheduler api-gateway; do \
		$(MAKE) docker-build COMPONENT=$$component; \
	done

# CI/CD targets
.PHONY: ci ci-test ci-lint

# CI test target
ci-test:
	@echo "Running CI tests..."
	@$(GOTEST) -v -race -covermode=atomic ./...

# CI lint target
ci-lint:
	@echo "Running CI linting..."
	@$(GOLINT) run --timeout=5m ./...

# CI target (runs all checks)
ci: fmt vet ci-lint ci-test

# Help target
.PHONY: help
help:
	@echo "KFn FaaS Platform - Makefile"
	@echo ""
	@echo "Build targets:"
	@echo "  make all           - Build all binaries"
	@echo "  make build         - Build all binaries (alias for all)"
	@echo "  make faasctl       - Build faasctl CLI"
	@echo "  make scheduler     - Build scheduler"
	@echo "  make api-gateway   - Build API gateway"
	@echo "  make clean         - Remove build artifacts"
	@echo ""
	@echo "Test targets:"
	@echo "  make test          - Run all tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo ""
	@echo "Linting targets:"
	@echo "  make lint          - Run golangci-lint"
	@echo "  make vet           - Run go vet"
	@echo "  make fmt           - Format code with gofmt"
	@echo "  make check-fmt     - Check if code is formatted"
	@echo "  make check         - Run all checks (fmt, vet, lint, test)"
	@echo ""
	@echo "Development targets:"
	@echo "  make deps          - Install and tidy dependencies"
	@echo "  make faasctl-run   - Run faasctl (use ARGS= for arguments)"
	@echo "  make scheduler-run - Run scheduler (use ARGS= for arguments)"
	@echo "  make api-gateway-run - Run API gateway (use ARGS= for arguments)"
	@echo ""
	@echo "Docker targets:"
	@echo "  make docker-build  - Build Docker image (use COMPONENT=)"
	@echo "  make docker-push   - Push Docker image (use COMPONENT=)"
	@echo "  make docker-build-all - Build all Docker images"
	@echo ""
	@echo "CI/CD targets:"
	@echo "  make ci            - Run all CI checks"
	@echo "  make ci-test       - Run CI tests"
	@echo "  make ci-lint       - Run CI linting"
	@echo ""
	@echo "Example usage:"
	@echo "  make build"
	@echo "  make test"
	@echo "  make docker-build COMPONENT=scheduler"
	@echo "  make faasctl-run ARGS=list"