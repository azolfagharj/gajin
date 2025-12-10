.PHONY: build test lint clean release help verify

# Variables
BINARY_NAME=gajin
BINARY_PATH=bin/$(BINARY_NAME)
MAIN_PATH=./cmd/gajin

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	@go build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Binary created at $(BINARY_PATH)"

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linters
	@echo "Running linters..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install it from https://golangci-lint.run/"; \
	fi

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@go clean

release: ## Build release version
	@echo "Building release version..."
	@mkdir -p bin
	@go build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Release binary created at $(BINARY_PATH)"

install: build ## Install the binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	@go install $(MAIN_PATH)
	@echo "Installed to $$(go env GOPATH)/bin/$(BINARY_NAME)"

run: build ## Build and run the binary
	@$(BINARY_PATH)

verify: ## Verify dependencies and run tests
	@echo "Verifying project..."
	@go mod verify
	@go vet ./...
	@go test -race ./...
	@echo "Verification passed!"

.DEFAULT_GOAL := help

