.PHONY: all build test bench lint clean install help

# Variables
BINARY_NAME=iso8583-lite
CMD_DIR=./cmd/$(BINARY_NAME)
BUILD_DIR=./bin
GO=go
GOTEST=$(GO) test
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOGET=$(GO) get
GOMOD=$(GO) mod

# Build flags
LDFLAGS=-ldflags "-s -w"

all: test build

## help: Display this help message
help:
	@echo "Available targets:"
	@grep -E '^##' $(MAKEFILE_LIST) | sed 's/^## /  /'

## build: Build the iso8583-lite binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "Built: $(BUILD_DIR)/$(BINARY_NAME)"

## install: Install the iso8583-lite binary
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install $(CMD_DIR)

## test: Run all tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.txt -covermode=atomic ./...

## test-short: Run tests without race detector
test-short:
	@echo "Running short tests..."
	$(GOTEST) -v ./...

## bench: Run benchmarks
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem -run=^$$ ./...

## bench-compare: Run benchmarks and save to file for comparison
bench-compare:
	@echo "Running benchmarks..."
	@mkdir -p benchmarks/results
	$(GOTEST) -bench=. -benchmem -run=^$$ ./... | tee benchmarks/results/bench-$$(date +%Y%m%d-%H%M%S).txt

## coverage: Generate test coverage report
coverage: test
	@echo "Generating coverage report..."
	$(GO) tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report: coverage.html"

## lint: Run linter
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, using staticcheck..."; \
		staticcheck ./...; \
	fi

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

## mod-tidy: Tidy go modules
mod-tidy:
	@echo "Tidying modules..."
	$(GOMOD) tidy

## mod-download: Download dependencies
mod-download:
	@echo "Downloading dependencies..."
	$(GOMOD) download

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.txt coverage.html

## check: Run all checks (fmt, vet, lint, test)
check: fmt vet lint test

## ci: Run CI pipeline locally
ci: mod-download check

.DEFAULT_GOAL := help
