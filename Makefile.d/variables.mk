#
# Copyright (C) 2024-2026 BuildBureau team
#
# Licensed under the Apache License, Version 2.0 (the "License");
# You may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# BuildBureau Makefile
# High-functionality Makefile for building, testing, and managing the project
# Can be used in Dockerfiles and CI/CD pipelines for standardization

# ============================================================================
# Variables and Configuration
# ============================================================================

# Application info
APP_NAME := buildbureau
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")

# Build configuration
BUILD_DIR := ./build
DIST_DIR := ./dist
CMD_DIR := ./cmd/buildbureau
COVERAGE_DIR := ./coverage
PROTO_DIR := ./pkg/protocol

# Go configuration
GO := go
GOPATH ?= $(shell $(GO) env GOPATH)
GOBIN ?= $(GOPATH)/bin
GOOS ?= $(shell $(GO) env GOOS)
GOARCH ?= $(shell $(GO) env GOARCH)

# Build flags
LDFLAGS := -w -s \
	-X main.version=$(VERSION) \
	-X main.buildTime=$(BUILD_TIME) \
	-X main.gitCommit=$(GIT_COMMIT) \
	-X main.gitBranch=$(GIT_BRANCH)

DEBUG_LDFLAGS := -X main.version=$(VERSION) \
	-X main.buildTime=$(BUILD_TIME) \
	-X main.gitCommit=$(GIT_COMMIT) \
	-X main.gitBranch=$(GIT_BRANCH)

# CGO is required for SQLite
CGO_ENABLED ?= 1

# Docker configuration
DOCKER_REGISTRY ?= 
DOCKER_IMAGE ?= $(APP_NAME)
DOCKER_TAG ?= $(VERSION)
DOCKER_FULL_IMAGE := $(if $(DOCKER_REGISTRY),$(DOCKER_REGISTRY)/,)$(DOCKER_IMAGE):$(DOCKER_TAG)
DOCKER_PLATFORMS ?= linux/amd64,linux/arm64

# Test configuration
TEST_TIMEOUT ?= 10m
TEST_FLAGS ?= -v -race -count=1
COVERAGE_OUT := $(COVERAGE_DIR)/coverage.out
COVERAGE_HTML := $(COVERAGE_DIR)/coverage.html

# Tools
GOLANGCI_LINT := $(GOBIN)/golangci-lint
GOSEC := $(GOBIN)/gosec
PROTOC := protoc
PROTOC_GEN_GO := $(GOBIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(GOBIN)/protoc-gen-go-grpc

# Colors for output
COLOR_RESET := \033[0m
COLOR_BOLD := \033[1m
COLOR_GREEN := \033[32m
COLOR_YELLOW := \033[33m
COLOR_BLUE := \033[34m
COLOR_CYAN := \033[36m

