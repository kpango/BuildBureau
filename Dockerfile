# BuildBureau Dockerfile
# Multi-stage build for optimal image size and security

# Stage 1: Builder
FROM golang:1.25-alpine AS builder

# Install build dependencies
# - make: Required for Makefile build process
# - gcc, musl-dev: Required for CGo compilation (mattn/go-sqlite3)
# - protobuf, protoc-gen-go: Required for gRPC code generation
RUN apk add --no-cache \
    make \
    gcc \
    musl-dev \
    git \
    protobuf \
    protobuf-dev

# Install Go tools for protobuf/gRPC
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Set working directory
WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Use Makefile for standardized build process
# Generate protobuf files
RUN make proto || echo "Proto generation skipped (protoc may not be available)"

# Build the application using Makefile
RUN make build-release

# Stage 2: Runtime
FROM alpine:latest

# Install runtime dependencies
# - ca-certificates: For HTTPS connections to LLM APIs
# - sqlite-libs: Runtime library for SQLite (not needed with static linking but good to have)
RUN apk add --no-cache ca-certificates sqlite-libs

# Create non-root user for security
RUN addgroup -g 1000 buildbureau && \
    adduser -D -u 1000 -G buildbureau buildbureau

# Create directories for data and configuration
RUN mkdir -p /app/data /app/config /app/agents && \
    chown -R buildbureau:buildbureau /app

# Set working directory
WORKDIR /app

# Copy binary from builder (built using Makefile)
COPY --from=builder /build/build/buildbureau /app/buildbureau

# Copy configuration files
COPY --chown=buildbureau:buildbureau config.yaml /app/config/
COPY --chown=buildbureau:buildbureau agents/ /app/agents/

# Switch to non-root user
USER buildbureau

# Expose ports
# 8080: gRPC server
# 8081: Health check/metrics (if implemented)
EXPOSE 8080 8081

# Set default environment variables
ENV BUILDBUREAU_CONFIG=/app/config/config.yaml
ENV BUILDBUREAU_DATA_DIR=/app/data

# Volume for persistent data
VOLUME ["/app/data"]

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/app/buildbureau", "--version"] || exit 1

# Run the application
ENTRYPOINT ["/app/buildbureau"]
CMD []
