# CRONEX Makefile

.PHONY: help build test test-verbose test-coverage test-short clean run benchmark lint fmt vet tidy install

# Default target
help:
	@echo "CRONEX - Cron Job Scheduler"
	@echo ""
	@echo "Available targets:"
	@echo "  make build         - Build the application"
	@echo "  make test          - Run all tests"
	@echo "  make test-verbose  - Run tests with verbose output"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make test-short    - Run short tests only"
	@echo "  make benchmark     - Run benchmarks"
	@echo "  make lint          - Run linter (requires golangci-lint)"
	@echo "  make fmt           - Format code"
	@echo "  make vet           - Run go vet"
	@echo "  make tidy          - Tidy go modules"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make run           - Run the application"
	@echo "  make run-test      - Run with test configuration"
	@echo "  make install       - Install the application"
	@echo "  make all           - Format, vet, test, and build"

# Build the application
build:
	@echo "Building CRONEX..."
	go build -o cronex .

# Build for multiple platforms
build-all:
	@echo "Building for all platforms..."
	GOOS=windows GOARCH=amd64 go build -o cronex-windows-x64.exe .
	GOOS=linux GOARCH=amd64 go build -o cronex-linux-x64 .
	GOOS=linux GOARCH=arm64 go build -o cronex-linux-arm64 .
	@echo "Build complete!"

# Run all tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run short tests
test-short:
	@echo "Running short tests..."
	go test -short ./...

# Run benchmarks
benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem -run=^$$ ./...

# Run linter
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Tidy go modules
tidy:
	@echo "Tidying go modules..."
	go mod tidy

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f cronex cronex.exe cronex-windows-x64.exe cronex-linux-x64
	rm -f coverage.out coverage.html
	@echo "Clean complete!"

# Run the application with default config
run: build
	@echo "Running CRONEX..."
	./cronex

# Run the application with test config
run-test: build
	@echo "Running CRONEX with test configuration..."
	./cronex -config cron.test.json

# Install the application
install:
	@echo "Installing CRONEX..."
	go install .

# Run all checks and build
all: fmt vet test build
	@echo "All checks passed and build complete!"
