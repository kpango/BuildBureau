.PHONY: all build clean test run proto install-tools

# Binary name
BINARY_NAME=buildbureau
BINARY_PATH=./bin/$(BINARY_NAME)

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Proto parameters
PROTOC=protoc
PROTO_DIR=./proto
PROTO_GEN_DIR=./proto/gen/go

all: clean build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	$(GOBUILD) -o $(BINARY_PATH) -v ./cmd/buildbureau

clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf bin/
	rm -rf $(PROTO_GEN_DIR)

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

run: build
	@echo "Running $(BINARY_NAME)..."
	$(BINARY_PATH)

deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

proto:
	@echo "Generating protobuf code..."
	@mkdir -p $(PROTO_GEN_DIR)
	$(PROTOC) --go_out=$(PROTO_GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_GEN_DIR) --go-grpc_opt=paths=source_relative \
		-I $(PROTO_DIR) \
		$(PROTO_DIR)/buildbureau/v1/*.proto

install-tools:
	@echo "Installing development tools..."
	$(GOGET) google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(GOGET) google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...

lint: fmt vet
	@echo "Linting complete"
