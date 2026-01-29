.PHONY: all build run clean test proto install-deps help

# Variables
BINARY_NAME=buildbureau
BUILD_DIR=bin
PROTO_DIR=proto
MAIN_PATH=cmd/buildbureau
GO=go

# Default target
all: build

# Install dependencies
install-deps:
	@echo "Installing Go dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Build the binary
build: install-deps
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run with custom config
run-config: build
	@echo "Running $(BINARY_NAME) with custom config..."
	./$(BUILD_DIR)/$(BINARY_NAME) -config=$(CONFIG)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f logs/*.log
	$(GO) clean

# Run tests
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Generate protocol buffer code (if needed in future)
proto:
	@echo "Generating protobuf code..."
	protoc --go_out=. --go-grpc_out=. $(PROTO_DIR)/*.proto

# Display help
help:
	@echo "BuildBureau - Multi-Agent AI System"
	@echo ""
	@echo "Available targets:"
	@echo "  all            - Build the project (default)"
	@echo "  build          - Build the binary"
	@echo "  run            - Build and run the application"
	@echo "  run-config     - Run with custom config (use CONFIG=path/to/config.yaml)"
	@echo "  clean          - Remove build artifacts"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run linter"
	@echo "  proto          - Generate protobuf code"
	@echo "  install-deps   - Install Go dependencies"
	@echo "  help           - Display this help message"
