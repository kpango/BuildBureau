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

# ============================================================================
# Tool Installation and Auto-Install System
# ============================================================================

# Stamp directory for tracking installed tools
STAMP_DIR := .make
$(shell mkdir -p $(STAMP_DIR))

# Stamp files for each tool
GOIMPORTS_STAMP := $(STAMP_DIR)/goimports.stamp
PROTOC_GEN_GO_STAMP := $(STAMP_DIR)/protoc-gen-go.stamp
PROTOC_GEN_GO_GRPC_STAMP := $(STAMP_DIR)/protoc-gen-go-grpc.stamp
PRETTIER_STAMP := $(STAMP_DIR)/prettier.stamp
YAMLFMT_STAMP := $(STAMP_DIR)/yamlfmt.stamp
JQ_STAMP := $(STAMP_DIR)/jq.stamp
GOLANGCI_LINT_STAMP := $(STAMP_DIR)/golangci-lint.stamp
GOSEC_STAMP := $(STAMP_DIR)/gosec.stamp
GOVULNCHECK_STAMP := $(STAMP_DIR)/govulncheck.stamp

# Install protoc-gen-go
$(PROTOC_GEN_GO_STAMP):
	@echo "$(COLOR_BLUE)Installing protoc-gen-go...$(COLOR_RESET)"
	@$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@echo "$(COLOR_GREEN)✓ protoc-gen-go installed$(COLOR_RESET)"
	@touch $@

# Install protoc-gen-go-grpc
$(PROTOC_GEN_GO_GRPC_STAMP):
	@echo "$(COLOR_BLUE)Installing protoc-gen-go-grpc...$(COLOR_RESET)"
	@$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "$(COLOR_GREEN)✓ protoc-gen-go-grpc installed$(COLOR_RESET)"
	@touch $@

# Install goimports
$(GOIMPORTS_STAMP):
	@echo "$(COLOR_BLUE)Installing goimports...$(COLOR_RESET)"
	@$(GO) install golang.org/x/tools/cmd/goimports@latest
	@echo "$(COLOR_GREEN)✓ goimports installed$(COLOR_RESET)"
	@touch $@

# Install prettier (requires npm)
$(PRETTIER_STAMP):
	@if command -v npm >/dev/null 2>&1; then \
echo "$(COLOR_BLUE)Installing prettier...$(COLOR_RESET)"; \
npm install -g prettier 2>/dev/null || (npm install --prefix ~/.npm-global prettier && export PATH=~/.npm-global/bin:$$PATH); \
echo "$(COLOR_GREEN)✓ prettier installed$(COLOR_RESET)"; \
touch $@; \
else \
echo "$(COLOR_YELLOW)⚠ npm not found, skipping prettier installation$(COLOR_RESET)"; \
echo "$(COLOR_YELLOW)  Install Node.js/npm to use prettier$(COLOR_RESET)"; \
touch $@; \
fi

# Install yamlfmt
$(YAMLFMT_STAMP):
	@echo "$(COLOR_BLUE)Installing yamlfmt...$(COLOR_RESET)"
	@$(GO) install github.com/google/yamlfmt/cmd/yamlfmt@latest
	@echo "$(COLOR_GREEN)✓ yamlfmt installed$(COLOR_RESET)"
	@touch $@

# Install jq (system package, just create stamp if exists)
$(JQ_STAMP):
	@if command -v jq >/dev/null 2>&1; then \
echo "$(COLOR_GREEN)✓ jq already installed$(COLOR_RESET)"; \
touch $@; \
else \
echo "$(COLOR_YELLOW)⚠ jq not found$(COLOR_RESET)"; \
echo "$(COLOR_YELLOW)  Install with: apt-get install jq (Ubuntu) or brew install jq (macOS)$(COLOR_RESET)"; \
touch $@; \
fi

# Install golangci-lint
$(GOLANGCI_LINT_STAMP):
	@echo "$(COLOR_BLUE)Installing golangci-lint...$(COLOR_RESET)"
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2
	@echo "$(COLOR_GREEN)✓ golangci-lint installed$(COLOR_RESET)"
	@touch $@

# Install gosec
$(GOSEC_STAMP):
	@echo "$(COLOR_BLUE)Installing gosec...$(COLOR_RESET)"
	@$(GO) install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo "$(COLOR_GREEN)✓ gosec installed$(COLOR_RESET)"
	@touch $@

# Install govulncheck
$(GOVULNCHECK_STAMP):
	@echo "$(COLOR_BLUE)Installing govulncheck...$(COLOR_RESET)"
	@$(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	@echo "$(COLOR_GREEN)✓ govulncheck installed$(COLOR_RESET)"
	@touch $@

# Master install targets
.PHONY: all build build-debug build-release build-static build-all \
	clean clean-all clean-build clean-coverage clean-cache clean-stamps \
	deps deps-update deps-verify deps-graph deps-tidy \
	docker docker-build docker-build-no-cache docker-push docker-build-multi docker-run docker-run-interactive docker-test docker-compose-up docker-compose-down \
	fmt fmt-check format format/go format/yaml format/json format/md format-check \
	help \
	install install-tools install-formatters install-security-tools install-all \
	lint lint-all lint-go lint-docker lint-fix format-lint format-all \
	proto proto-clean \
	run run-debug \
	test test-unit test-integration test-all test-coverage test-coverage-html test-bench test-race test/llm-integration \
	ci-test ci-build ci-lint ci-all \
	security security-scan security-deps \
	release release-build release-package \
	version \
	check check-go check-deps check-proto

# ============================================================================
# Help Target
# ============================================================================

help: ## Display this help message
	@echo "$(COLOR_BOLD)BuildBureau Makefile$(COLOR_RESET)"
	@echo "$(COLOR_CYAN)Version: $(VERSION)$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BOLD)Usage:$(COLOR_RESET)"
	@echo "  make $(COLOR_GREEN)<target>$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BOLD)Build Targets:$(COLOR_RESET)"
	@awk 'BEGIN {FS = ":.*##"; category="build"} /^[a-zA-Z_-]+:.*?##/ { if ($$0 ~ /Build/) printf "  $(COLOR_GREEN)%-20s$(COLOR_RESET) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(COLOR_BOLD)Test Targets:$(COLOR_RESET)"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { if ($$0 ~ /Test|test|coverage|bench/) printf "  $(COLOR_GREEN)%-20s$(COLOR_RESET) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(COLOR_BOLD)Docker Targets:$(COLOR_RESET)"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { if ($$0 ~ /Docker|docker|container/) printf "  $(COLOR_GREEN)%-20s$(COLOR_RESET) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(COLOR_BOLD)Development Targets:$(COLOR_RESET)"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_\/-]+:.*?##/ { if ($$0 ~ /fmt|format|lint|proto|deps|install/) printf "  $(COLOR_GREEN)%-20s$(COLOR_RESET) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(COLOR_BOLD)CI/CD Targets:$(COLOR_RESET)"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_\/-]+:.*?##/ { if ($$0 ~ /^ci-/) printf "  $(COLOR_GREEN)%-20s$(COLOR_RESET) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(COLOR_BOLD)Other Targets:$(COLOR_RESET)"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_\/-]+:.*?##/ { if ($$0 !~ /Build|Test|test|Docker|docker|fmt|format|lint|proto|deps|install|^ci-|coverage|bench/) printf "  $(COLOR_GREEN)%-20s$(COLOR_RESET) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# ============================================================================
# Build Targets
# ============================================================================

all: clean deps proto build test ## Build everything (clean, deps, proto, build, test)

build: ## Build the application binary
	@echo "$(COLOR_BLUE)Building $(APP_NAME)...$(COLOR_RESET)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build \
		-ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(APP_NAME) \
		$(CMD_DIR)
	@echo "$(COLOR_GREEN)✓ Build complete: $(BUILD_DIR)/$(APP_NAME)$(COLOR_RESET)"

build-debug: ## Build with debug symbols and no optimization
	@echo "$(COLOR_BLUE)Building debug version...$(COLOR_RESET)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build \
		-gcflags="all=-N -l" \
		-ldflags "$(DEBUG_LDFLAGS)" \
		-o $(BUILD_DIR)/$(APP_NAME)-debug \
		$(CMD_DIR)
	@echo "$(COLOR_GREEN)✓ Debug build complete: $(BUILD_DIR)/$(APP_NAME)-debug$(COLOR_RESET)"

build-release: ## Build optimized release binary
	@echo "$(COLOR_BLUE)Building release version...$(COLOR_RESET)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build \
		-ldflags "$(LDFLAGS)" \
		-trimpath \
		-o $(BUILD_DIR)/$(APP_NAME) \
		$(CMD_DIR)
	@echo "$(COLOR_GREEN)✓ Release build complete: $(BUILD_DIR)/$(APP_NAME)$(COLOR_RESET)"

build-static: ## Build static binary (CGO disabled where possible)
	@echo "$(COLOR_BLUE)Building static binary...$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)Note: SQLite requires CGO, so this may not be fully static$(COLOR_RESET)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 $(GO) build \
		-ldflags "$(LDFLAGS) -extldflags '-static'" \
		-tags 'netgo osusergo' \
		-o $(BUILD_DIR)/$(APP_NAME)-static \
		$(CMD_DIR)
	@echo "$(COLOR_GREEN)✓ Static build complete: $(BUILD_DIR)/$(APP_NAME)-static$(COLOR_RESET)"

build-all: ## Build for multiple platforms (linux/amd64, linux/arm64, darwin/amd64, darwin/arm64)
	@echo "$(COLOR_BLUE)Building for multiple platforms...$(COLOR_RESET)"
	@mkdir -p $(DIST_DIR)
	@for os in linux darwin; do \
		for arch in amd64 arm64; do \
			output="$(DIST_DIR)/$(APP_NAME)-$$os-$$arch"; \
			if [ "$$os" = "windows" ]; then output="$$output.exe"; fi; \
			echo "$(COLOR_CYAN)Building $$os/$$arch...$(COLOR_RESET)"; \
			CGO_ENABLED=$(CGO_ENABLED) GOOS=$$os GOARCH=$$arch $(GO) build \
				-ldflags "$(LDFLAGS)" \
				-o $$output \
				$(CMD_DIR) || exit 1; \
		done \
	done
	@echo "$(COLOR_GREEN)✓ Multi-platform builds complete in $(DIST_DIR)/$(COLOR_RESET)"

# ============================================================================
# Test Targets
# ============================================================================

test: ## Run all tests
	@echo "$(COLOR_BLUE)Running tests...$(COLOR_RESET)"
	CGO_ENABLED=1 $(GO) test $(TEST_FLAGS) -timeout $(TEST_TIMEOUT) $$(go list ./... | grep -v /examples)
	@echo "$(COLOR_GREEN)✓ Tests complete$(COLOR_RESET)"

test-unit: ## Run unit tests only
	@echo "$(COLOR_BLUE)Running unit tests...$(COLOR_RESET)"
	CGO_ENABLED=1 $(GO) test $(TEST_FLAGS) -short -timeout $(TEST_TIMEOUT) $$(go list ./... | grep -v /examples)
	@echo "$(COLOR_GREEN)✓ Unit tests complete$(COLOR_RESET)"

test-integration: ## Run integration tests
	@echo "$(COLOR_BLUE)Running integration tests...$(COLOR_RESET)"
	CGO_ENABLED=1 $(GO) test $(TEST_FLAGS) -run Integration -timeout $(TEST_TIMEOUT) $$(go list ./... | grep -v /examples)
	@echo "$(COLOR_GREEN)✓ Integration tests complete$(COLOR_RESET)"

test-all: test ## Run all tests (alias for test)

test-coverage: ## Run tests with coverage
	@echo "$(COLOR_BLUE)Running tests with coverage...$(COLOR_RESET)"
	@mkdir -p $(COVERAGE_DIR)
	CGO_ENABLED=1 $(GO) test -race -coverprofile=$(COVERAGE_OUT) -covermode=atomic ./...
	@echo "$(COLOR_GREEN)✓ Coverage report: $(COVERAGE_OUT)$(COLOR_RESET)"
	@$(GO) tool cover -func=$(COVERAGE_OUT) | tail -n 1

test-coverage-html: test-coverage ## Generate HTML coverage report
	@echo "$(COLOR_BLUE)Generating HTML coverage report...$(COLOR_RESET)"
	$(GO) tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	@echo "$(COLOR_GREEN)✓ HTML report: $(COVERAGE_HTML)$(COLOR_RESET)"

test-bench: ## Run benchmark tests
	@echo "$(COLOR_BLUE)Running benchmarks...$(COLOR_RESET)"
	CGO_ENABLED=1 $(GO) test -bench=. -benchmem -run=^$$ ./...
	@echo "$(COLOR_GREEN)✓ Benchmarks complete$(COLOR_RESET)"

test-race: ## Run tests with race detector
	@echo "$(COLOR_BLUE)Running tests with race detector...$(COLOR_RESET)"
	CGO_ENABLED=1 $(GO) test -race -timeout $(TEST_TIMEOUT) ./...
	@echo "$(COLOR_GREEN)✓ Race detection tests complete$(COLOR_RESET)"

test/llm-integration: ## Test with real LLM integration (replaces test_real_llm.sh)
	@echo "=== BuildBureau Real LLM Integration Test ==="
	@echo ""
	@if [ -z "$(GEMINI_API_KEY)" ] || [ "$(GEMINI_API_KEY)" = "demo-key" ]; then \
		echo "$(COLOR_YELLOW)⚠️  Warning: GEMINI_API_KEY is not set or is using demo value$(COLOR_RESET)"; \
		echo "   The system will work but LLM features will be limited"; \
		echo ""; \
		echo "To test with real LLM:"; \
		echo "  1. Get an API key from: https://aistudio.google.com/app/apikey"; \
		echo "  2. Set it: export GEMINI_API_KEY='your-actual-key'"; \
		echo "  3. Re-run this test"; \
		echo ""; \
	fi
	@echo "Testing with:"
	@echo "  GEMINI_API_KEY: $${GEMINI_API_KEY:0:20}..."
	@echo ""
	@echo "Running BuildBureau test..."
	@echo "=============================="
	@GEMINI_API_KEY="$${GEMINI_API_KEY:-demo-key}" \
		CLAUDE_API_KEY="$${CLAUDE_API_KEY:-demo-key}" \
		CODEX_API_KEY="$${CODEX_API_KEY:-demo-key}" \
		QWEN_API_KEY="$${QWEN_API_KEY:-demo-key}" \
		go run examples/test_basic/main.go
	@echo ""
	@echo "=============================="
	@echo "$(COLOR_GREEN)✓ Test complete!$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BLUE)To test the interactive TUI:$(COLOR_RESET)"
	@echo "  ./buildbureau"

# ============================================================================
# Proto Targets
# ============================================================================

proto: | $(PROTOC_GEN_GO_STAMP) $(PROTOC_GEN_GO_GRPC_STAMP) ## Generate protobuf files from .proto definitions
	@echo "$(COLOR_BLUE)Generating protobuf files...$(COLOR_RESET)"
	@mkdir -p $(PROTO_DIR)
	$(PROTOC) --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/agent.proto
	@echo "$(COLOR_GREEN)✓ Proto generation complete$(COLOR_RESET)"

proto-clean: ## Clean generated proto files
	@echo "$(COLOR_BLUE)Cleaning generated proto files...$(COLOR_RESET)"
	@rm -f $(PROTO_DIR)/*.pb.go
	@echo "$(COLOR_GREEN)✓ Proto files cleaned$(COLOR_RESET)"

# ============================================================================
# Dependencies
# ============================================================================

deps: ## Download and tidy dependencies
	@echo "$(COLOR_BLUE)Downloading dependencies...$(COLOR_RESET)"
	$(GO) mod download
	$(GO) mod tidy
	@echo "$(COLOR_GREEN)✓ Dependencies ready$(COLOR_RESET)"

deps-update: ## Update dependencies to latest versions
	@echo "$(COLOR_BLUE)Updating dependencies...$(COLOR_RESET)"
	$(GO) get -u ./...
	$(GO) mod tidy
	@echo "$(COLOR_GREEN)✓ Dependencies updated$(COLOR_RESET)"

deps-verify: ## Verify dependencies
	@echo "$(COLOR_BLUE)Verifying dependencies...$(COLOR_RESET)"
	$(GO) mod verify
	@echo "$(COLOR_GREEN)✓ Dependencies verified$(COLOR_RESET)"

deps-graph: ## Display dependency graph
	@echo "$(COLOR_BLUE)Dependency graph:$(COLOR_RESET)"
	$(GO) mod graph

deps-tidy: ## Tidy go.mod and go.sum
	@echo "$(COLOR_BLUE)Tidying dependencies...$(COLOR_RESET)"
	$(GO) mod tidy
	@echo "$(COLOR_GREEN)✓ Dependencies tidied$(COLOR_RESET)"

# ============================================================================
# Code Quality
# ============================================================================

fmt: ## Format Go code
	@echo "$(COLOR_BLUE)Formatting code...$(COLOR_RESET)"
	$(GO) fmt ./...
	@echo "$(COLOR_GREEN)✓ Code formatted$(COLOR_RESET)"

fmt-check: ## Check if code is formatted
	@echo "$(COLOR_BLUE)Checking code formatting...$(COLOR_RESET)"
	@UNFORMATTED=$$(gofmt -s -l . | grep -v vendor || true); \
	if [ -n "$$UNFORMATTED" ]; then \
		echo "$(COLOR_YELLOW)The following files need formatting:$(COLOR_RESET)"; \
		echo "$$UNFORMATTED"; \
		exit 1; \
	fi
	@echo "$(COLOR_GREEN)✓ All files are formatted$(COLOR_RESET)"

format: format/go format/yaml format/json format/md ## Format all files (Go, YAML, JSON, Markdown)
	@echo "$(COLOR_GREEN)✓ All files formatted$(COLOR_RESET)"

format/go: | $(GOIMPORTS_STAMP) ## Format Go files with gofmt and goimports
	@echo "$(COLOR_BLUE)Formatting Go files...$(COLOR_RESET)"
	@$(GO) fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		find . -name '*.go' -not -path './vendor/*' -not -path './.git/*' -exec goimports -w {} \; ; \
		echo "$(COLOR_GREEN)✓ Go files formatted (gofmt + goimports)$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_GREEN)✓ Go files formatted (gofmt only - install goimports for import sorting)$(COLOR_RESET)"; \
	fi

format/yaml: | $(PRETTIER_STAMP) $(YAMLFMT_STAMP) ## Format YAML files
	@echo "$(COLOR_BLUE)Formatting YAML files...$(COLOR_RESET)"
	@if command -v prettier >/dev/null 2>&1; then \
		find . -name '*.yaml' -o -name '*.yml' | grep -v vendor | grep -v .git | xargs prettier --write --parser yaml 2>/dev/null || true; \
		echo "$(COLOR_GREEN)✓ YAML files formatted (prettier)$(COLOR_RESET)"; \
	elif command -v yamlfmt >/dev/null 2>&1; then \
		find . -name '*.yaml' -o -name '*.yml' | grep -v vendor | grep -v .git | xargs yamlfmt 2>/dev/null || true; \
		echo "$(COLOR_GREEN)✓ YAML files formatted (yamlfmt)$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)⚠ No YAML formatter found (install prettier or yamlfmt)$(COLOR_RESET)"; \
	fi

format/json: | $(PRETTIER_STAMP) $(JQ_STAMP) ## Format JSON files
	@echo "$(COLOR_BLUE)Formatting JSON files...$(COLOR_RESET)"
	@if command -v prettier >/dev/null 2>&1; then \
		find . -name '*.json' | grep -v vendor | grep -v .git | grep -v node_modules | xargs prettier --write --parser json 2>/dev/null || true; \
		echo "$(COLOR_GREEN)✓ JSON files formatted (prettier)$(COLOR_RESET)"; \
	elif command -v jq >/dev/null 2>&1; then \
		for file in $$(find . -name '*.json' | grep -v vendor | grep -v .git | grep -v node_modules); do \
			jq '.' "$$file" > "$$file.tmp" && mv "$$file.tmp" "$$file" 2>/dev/null || true; \
		done; \
		echo "$(COLOR_GREEN)✓ JSON files formatted (jq)$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)⚠ No JSON formatter found (install prettier or jq)$(COLOR_RESET)"; \
	fi

format/md: | $(PRETTIER_STAMP) ## Format Markdown files
	@echo "$(COLOR_BLUE)Formatting Markdown files...$(COLOR_RESET)"
	@if command -v prettier >/dev/null 2>&1; then \
		find . -name '*.md' | grep -v vendor | grep -v .git | grep -v node_modules | xargs prettier --write --parser markdown --prose-wrap always 2>/dev/null || true; \
		echo "$(COLOR_GREEN)✓ Markdown files formatted (prettier)$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)⚠ No Markdown formatter found (install prettier)$(COLOR_RESET)"; \
	fi

format-check: ## Check if all files are properly formatted
	@echo "$(COLOR_BLUE)Checking all file formatting...$(COLOR_RESET)"
	@EXIT_CODE=0; \
	UNFORMATTED=$$(gofmt -s -l . | grep -v vendor || true); \
	if [ -n "$$UNFORMATTED" ]; then \
		echo "$(COLOR_YELLOW)The following Go files need formatting:$(COLOR_RESET)"; \
		echo "$$UNFORMATTED"; \
		EXIT_CODE=1; \
	fi; \
	if command -v prettier >/dev/null 2>&1; then \
		if ! prettier --check '**/*.{yaml,yml,json,md}' --ignore-path .gitignore 2>/dev/null; then \
			echo "$(COLOR_YELLOW)Some YAML/JSON/Markdown files need formatting$(COLOR_RESET)"; \
			EXIT_CODE=1; \
		fi; \
	fi; \
	if [ $$EXIT_CODE -eq 0 ]; then \
		echo "$(COLOR_GREEN)✓ All files are properly formatted$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)Run 'make format' to fix formatting issues$(COLOR_RESET)"; \
		exit $$EXIT_CODE; \
	fi

lint-fix: | $(GOLANGCI_LINT_STAMP) ## Run golangci-lint with automatic fixes
	@echo "$(COLOR_BLUE)Running golangci-lint with auto-fix...$(COLOR_RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --fix --skip-dirs=examples ./...; \
		echo "$(COLOR_GREEN)✓ Auto-fixes applied$(COLOR_RESET)"; \
		echo "$(COLOR_YELLOW)Run 'git diff' to see changes$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_RED)✗ golangci-lint not installed$(COLOR_RESET)"; \
		exit 1; \
	fi

format-lint: lint-fix ## Format code using golangci-lint --fix (alias for lint-fix)

format-all: format lint-fix ## Format using all tools (gofmt + goimports + golangci-lint --fix)
	@echo "$(COLOR_GREEN)✓ All formatting complete (gofmt + goimports + golangci-lint)$(COLOR_RESET)"

lint: ## Run go vet linter
	@echo "$(COLOR_BLUE)Running go vet...$(COLOR_RESET)"
	$(GO) vet $$($(GO) list ./... | grep -v /examples)
	@echo "$(COLOR_GREEN)✓ Linting complete$(COLOR_RESET)"

lint-all: | $(GOLANGCI_LINT_STAMP) lint ## Run all linters (go vet + golangci-lint if available)
	@echo "$(COLOR_BLUE)Running all linters...$(COLOR_RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --skip-dirs=examples ./...; \
		echo "$(COLOR_GREEN)✓ golangci-lint complete$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)golangci-lint not installed, skipping$(COLOR_RESET)"; \
	fi

lint-go: lint ## Run Go linters (alias for lint)

lint-docker: ## Lint Dockerfile
	@echo "$(COLOR_BLUE)Linting Dockerfile...$(COLOR_RESET)"
	@if command -v hadolint >/dev/null 2>&1; then \
		hadolint Dockerfile; \
		echo "$(COLOR_GREEN)✓ Dockerfile lint complete$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)hadolint not installed, skipping$(COLOR_RESET)"; \
	fi

# ============================================================================
# Security
# ============================================================================

security: | $(GOSEC_STAMP) $(GOVULNCHECK_STAMP) security-scan security-deps ## Run all security checks

security-scan: ## Run security scanner (gosec)
	@echo "$(COLOR_BLUE)Running security scan...$(COLOR_RESET)"
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
		echo "$(COLOR_GREEN)✓ Security scan complete$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)gosec not installed, skipping$(COLOR_RESET)"; \
		echo "$(COLOR_YELLOW)Install: go install github.com/securego/gosec/v2/cmd/gosec@latest$(COLOR_RESET)"; \
	fi

security-deps: ## Check for vulnerable dependencies
	@echo "$(COLOR_BLUE)Checking dependencies for vulnerabilities...$(COLOR_RESET)"
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
		echo "$(COLOR_GREEN)✓ Dependency check complete$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)govulncheck not installed, skipping$(COLOR_RESET)"; \
		echo "$(COLOR_YELLOW)Install: go install golang.org/x/vuln/cmd/govulncheck@latest$(COLOR_RESET)"; \
	fi

# ============================================================================
# Docker Targets
# ============================================================================

docker: docker-build ## Build Docker image (alias for docker-build)

docker-build: ## Build Docker image (replaces docker/build.sh)
	@echo "$(COLOR_BLUE)Building Docker image: $(DOCKER_FULL_IMAGE)$(COLOR_RESET)"
	@docker build -t $(DOCKER_FULL_IMAGE) .
	@echo "$(COLOR_GREEN)✓ Build successful!$(COLOR_RESET)"
	@echo "$(COLOR_GREEN)Image: $(DOCKER_FULL_IMAGE)$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BLUE)To run the image:$(COLOR_RESET)"
	@echo "  docker run -e GEMINI_API_KEY=your-key $(DOCKER_FULL_IMAGE)"
	@echo ""
	@echo "$(COLOR_BLUE)Or use docker-compose:$(COLOR_RESET)"
	@echo "  make docker-compose-up"

docker-build-no-cache: ## Build Docker image without cache
	@echo "$(COLOR_BLUE)Building Docker image (no cache): $(DOCKER_FULL_IMAGE)$(COLOR_RESET)"
	@docker build --no-cache -t $(DOCKER_FULL_IMAGE) .
	@echo "$(COLOR_GREEN)✓ Docker image built: $(DOCKER_FULL_IMAGE)$(COLOR_RESET)"

docker-build-multi: ## Build multi-architecture Docker image
	@echo "$(COLOR_BLUE)Building multi-arch Docker image: $(DOCKER_FULL_IMAGE)$(COLOR_RESET)"
	@docker buildx build --platform $(DOCKER_PLATFORMS) -t $(DOCKER_FULL_IMAGE) .
	@echo "$(COLOR_GREEN)✓ Multi-arch Docker image built$(COLOR_RESET)"

docker-push: ## Push Docker image to registry
	@echo "$(COLOR_BLUE)Pushing Docker image: $(DOCKER_FULL_IMAGE)$(COLOR_RESET)"
	@docker push $(DOCKER_FULL_IMAGE)
	@echo "$(COLOR_GREEN)✓ Docker image pushed$(COLOR_RESET)"

docker-run: ## Run Docker container in daemon mode (replaces docker/run.sh)
	@if [ -z "$(GEMINI_API_KEY)" ] && [ -z "$(OPENAI_API_KEY)" ] && [ -z "$(CLAUDE_API_KEY)" ]; then \
		echo "$(COLOR_YELLOW)Warning: No LLM API key provided!$(COLOR_RESET)"; \
		echo "$(COLOR_YELLOW)Please set at least one of:$(COLOR_RESET)"; \
		echo "  export GEMINI_API_KEY=your-key"; \
		echo "  export OPENAI_API_KEY=your-key"; \
		echo "  export CLAUDE_API_KEY=your-key"; \
		echo ""; \
	fi
	@CONTAINER_NAME=$${CONTAINER_NAME:-buildbureau}; \
	echo "$(COLOR_BLUE)Starting BuildBureau container...$(COLOR_RESET)"; \
	echo "$(COLOR_BLUE)Container name: $(COLOR_GREEN)$$CONTAINER_NAME$(COLOR_RESET)"; \
	echo "$(COLOR_BLUE)Image: $(COLOR_GREEN)$(DOCKER_FULL_IMAGE)$(COLOR_RESET)"; \
	if docker ps -a --format '{{.Names}}' | grep -q "^$$CONTAINER_NAME$$"; then \
		echo "$(COLOR_BLUE)Stopping existing container...$(COLOR_RESET)"; \
		docker stop $$CONTAINER_NAME > /dev/null 2>&1 || true; \
		docker rm $$CONTAINER_NAME > /dev/null 2>&1 || true; \
	fi; \
	docker run -d \
		--name $$CONTAINER_NAME \
		-e GEMINI_API_KEY="$(GEMINI_API_KEY)" \
		-e OPENAI_API_KEY="$(OPENAI_API_KEY)" \
		-e CLAUDE_API_KEY="$(CLAUDE_API_KEY)" \
		-e OPENAI_MODEL="$(OPENAI_MODEL:-gpt-4-turbo-preview)" \
		-e CLAUDE_MODEL="$(CLAUDE_MODEL:-claude-3-5-sonnet-20241022)" \
		-e SLACK_TOKEN="$(SLACK_TOKEN)" \
		-v buildbureau-data:/app/data \
		-p 8080:8080 \
		--restart unless-stopped \
		$(DOCKER_FULL_IMAGE) && \
	echo "$(COLOR_GREEN)✓ Container started successfully!$(COLOR_RESET)" && \
	echo "" && \
	echo "$(COLOR_BLUE)Container details:$(COLOR_RESET)" && \
	docker ps --filter "name=$$CONTAINER_NAME" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" && \
	echo "" && \
	echo "$(COLOR_BLUE)View logs:$(COLOR_RESET)" && \
	echo "  docker logs -f $$CONTAINER_NAME" && \
	echo "" && \
	echo "$(COLOR_BLUE)Execute interactive shell:$(COLOR_RESET)" && \
	echo "  docker exec -it $$CONTAINER_NAME sh" && \
	echo "" && \
	echo "$(COLOR_BLUE)Stop container:$(COLOR_RESET)" && \
	echo "  docker stop $$CONTAINER_NAME"

docker-run-interactive: ## Run Docker container interactively
	@echo "$(COLOR_BLUE)Running Docker container interactively...$(COLOR_RESET)"
	@docker run --rm -it \
		-e GEMINI_API_KEY \
		-e OPENAI_API_KEY \
		-e CLAUDE_API_KEY \
		-v $$(pwd)/data:/app/data \
		$(DOCKER_FULL_IMAGE)

docker-test: ## Test Docker setup (replaces docker/test.sh)
	@echo "=== BuildBureau Docker Test ==="
	@echo ""
	@echo "$(COLOR_BLUE)Checking Docker...$(COLOR_RESET)"
	@if ! command -v docker &> /dev/null; then \
		echo "$(COLOR_RED)✗ Docker not found$(COLOR_RESET)"; \
		exit 1; \
	fi
	@echo "$(COLOR_GREEN)✓ Docker found: $$(docker --version)$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BLUE)Checking Docker Compose...$(COLOR_RESET)"
	@if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then \
		echo "$(COLOR_RED)✗ Docker Compose not found$(COLOR_RESET)"; \
		exit 1; \
	fi
	@echo "$(COLOR_GREEN)✓ Docker Compose found$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BLUE)Checking API keys...$(COLOR_RESET)"
	@if [ -z "$(GEMINI_API_KEY)" ] && [ -z "$(OPENAI_API_KEY)" ] && [ -z "$(CLAUDE_API_KEY)" ]; then \
		echo "$(COLOR_RED)✗ No API key found$(COLOR_RESET)"; \
		echo "$(COLOR_BLUE)Please set at least one:$(COLOR_RESET)"; \
		echo "  export GEMINI_API_KEY=your-key"; \
		echo "  export OPENAI_API_KEY=your-key"; \
		echo "  export CLAUDE_API_KEY=your-key"; \
		exit 1; \
	fi
	@echo "$(COLOR_GREEN)✓ API key(s) found$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BLUE)Building Docker image...$(COLOR_RESET)"
	@docker build -t buildbureau:test . || { \
		echo "$(COLOR_RED)✗ Build failed$(COLOR_RESET)"; \
		exit 1; \
	}
	@echo "$(COLOR_GREEN)✓ Build successful$(COLOR_RESET)"
	@echo ""
	@IMAGE_SIZE=$$(docker images buildbureau:test --format "{{.Size}}"); \
	echo "$(COLOR_BLUE)Image size: $(COLOR_GREEN)$$IMAGE_SIZE$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BLUE)Testing container...$(COLOR_RESET)"
	@docker run --rm buildbureau:test --version || { \
		echo "$(COLOR_RED)✗ Container test failed$(COLOR_RESET)"; \
		exit 1; \
	}
	@echo "$(COLOR_GREEN)✓ Container test successful$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_GREEN)=== All checks passed! ===$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BLUE)Next steps:$(COLOR_RESET)"
	@echo "1. Run with Docker Compose:"
	@echo "   make docker-compose-up"
	@echo ""
	@echo "2. Or run directly:"
	@echo "   make docker-run"
	@echo ""
	@echo "3. View logs:"
	@echo "   docker logs -f buildbureau"
	@echo ""
	@echo "$(COLOR_GREEN)BuildBureau is ready to use! 🎉$(COLOR_RESET)"

docker-compose-up: ## Start services with docker-compose
	@echo "$(COLOR_BLUE)Starting Docker Compose services...$(COLOR_RESET)"
	@docker-compose up -d
	@echo "$(COLOR_GREEN)✓ Services started$(COLOR_RESET)"

docker-compose-down: ## Stop services with docker-compose
	@echo "$(COLOR_BLUE)Stopping Docker Compose services...$(COLOR_RESET)"
	@docker-compose down
	@echo "$(COLOR_GREEN)✓ Services stopped$(COLOR_RESET)"

# ============================================================================
# CI/CD Targets
# ============================================================================

ci-all: ci-lint ci-build ci-test ## Run all CI checks (lint, build, test)

ci-lint: fmt-check lint ## Run linters for CI
	@echo "$(COLOR_GREEN)✓ CI linting complete$(COLOR_RESET)"

ci-build: ## Build for CI (proto generation attempted but not required)
	@echo "$(COLOR_BLUE)Building for CI...$(COLOR_RESET)"
	@$(MAKE) proto 2>/dev/null || echo "$(COLOR_YELLOW)Proto generation skipped (protoc not available)$(COLOR_RESET)"
	@$(MAKE) build-release
	@echo "$(COLOR_GREEN)✓ CI build complete$(COLOR_RESET)"

ci-test: test ## Run tests for CI (without coverage due to Go 1.26.0 compatibility)
	@echo "$(COLOR_GREEN)✓ CI tests complete$(COLOR_RESET)"

# ============================================================================
# Development Targets
# ============================================================================

run: build ## Build and run the application
	@echo "$(COLOR_BLUE)Starting $(APP_NAME)...$(COLOR_RESET)"
	$(BUILD_DIR)/$(APP_NAME)

run-debug: build-debug ## Build and run in debug mode
	@echo "$(COLOR_BLUE)Starting $(APP_NAME) in debug mode...$(COLOR_RESET)"
	$(BUILD_DIR)/$(APP_NAME)-debug


# Master install targets
.PHONY: install-all install-tools install-formatters install-security-tools

install-tools: $(PROTOC_GEN_GO_STAMP) $(PROTOC_GEN_GO_GRPC_STAMP) $(GOIMPORTS_STAMP) ## Install all Go development tools
	@echo "$(COLOR_GREEN)✓ All Go development tools installed$(COLOR_RESET)"

install-formatters: $(PRETTIER_STAMP) $(YAMLFMT_STAMP) $(JQ_STAMP) ## Install all formatting tools
	@echo "$(COLOR_GREEN)✓ All formatting tools installed$(COLOR_RESET)"

install-security-tools: $(GOSEC_STAMP) $(GOVULNCHECK_STAMP) ## Install security scanning tools
	@echo "$(COLOR_GREEN)✓ All security tools installed$(COLOR_RESET)"

install-all: install-tools install-formatters install-security-tools $(GOLANGCI_LINT_STAMP) ## Install all tools
	@echo "$(COLOR_GREEN)✓ All tools installed$(COLOR_RESET)"

# Clean stamp files to force reinstall
.PHONY: clean-stamps
clean-stamps: ## Remove all tool installation stamps (force reinstall)
	@echo "$(COLOR_BLUE)Removing tool installation stamps...$(COLOR_RESET)"
	@rm -rf $(STAMP_DIR)
	@echo "$(COLOR_GREEN)✓ Stamps removed$(COLOR_RESET)"



# ============================================================================
# Release Targets
# ============================================================================

release-build: clean deps proto build-all ## Build release for all platforms
	@echo "$(COLOR_GREEN)✓ Release builds complete$(COLOR_RESET)"

release-package: release-build ## Package release artifacts
	@echo "$(COLOR_BLUE)Packaging release artifacts...$(COLOR_RESET)"
	@mkdir -p $(DIST_DIR)/packages
	@cd $(DIST_DIR) && \
	for binary in $(APP_NAME)-*; do \
		if [ -f "$$binary" ]; then \
			tar czf packages/$${binary}.tar.gz $$binary; \
			echo "Packaged: packages/$${binary}.tar.gz"; \
		fi \
	done
	@echo "$(COLOR_GREEN)✓ Release artifacts packaged in $(DIST_DIR)/packages/$(COLOR_RESET)"

# ============================================================================
# Clean Targets
# ============================================================================

clean: ## Remove build artifacts
	@echo "$(COLOR_BLUE)Cleaning build artifacts...$(COLOR_RESET)"
	@rm -rf $(BUILD_DIR)
	@rm -f $(APP_NAME)
	$(GO) clean
	@echo "$(COLOR_GREEN)✓ Build artifacts cleaned$(COLOR_RESET)"

clean-all: clean clean-coverage clean-cache ## Remove all generated files
	@echo "$(COLOR_BLUE)Cleaning all generated files...$(COLOR_RESET)"
	@rm -rf $(DIST_DIR)
	@echo "$(COLOR_GREEN)✓ All artifacts cleaned$(COLOR_RESET)"

clean-build: ## Remove build directory
	@echo "$(COLOR_BLUE)Cleaning build directory...$(COLOR_RESET)"
	@rm -rf $(BUILD_DIR)
	@echo "$(COLOR_GREEN)✓ Build directory cleaned$(COLOR_RESET)"

clean-coverage: ## Remove coverage reports
	@echo "$(COLOR_BLUE)Cleaning coverage reports...$(COLOR_RESET)"
	@rm -rf $(COVERAGE_DIR)
	@echo "$(COLOR_GREEN)✓ Coverage reports cleaned$(COLOR_RESET)"

clean-cache: ## Clean Go build cache
	@echo "$(COLOR_BLUE)Cleaning Go cache...$(COLOR_RESET)"
	$(GO) clean -cache -testcache -modcache
	@echo "$(COLOR_GREEN)✓ Go cache cleaned$(COLOR_RESET)"

# ============================================================================
# Utility Targets
# ============================================================================

version: ## Display version information
	@echo "$(COLOR_BOLD)BuildBureau Version Information$(COLOR_RESET)"
	@echo "  Version:    $(VERSION)"
	@echo "  Git Commit: $(GIT_COMMIT)"
	@echo "  Git Branch: $(GIT_BRANCH)"
	@echo "  Build Time: $(BUILD_TIME)"
	@echo "  Go Version: $(shell $(GO) version)"
	@echo "  Platform:   $(GOOS)/$(GOARCH)"

check: check-go check-deps check-proto ## Run all checks

check-go: ## Check Go installation and version
	@echo "$(COLOR_BLUE)Checking Go installation...$(COLOR_RESET)"
	@$(GO) version
	@echo "$(COLOR_GREEN)✓ Go is installed$(COLOR_RESET)"

check-deps: ## Check if dependencies are satisfied
	@echo "$(COLOR_BLUE)Checking dependencies...$(COLOR_RESET)"
	@$(GO) mod verify
	@echo "$(COLOR_GREEN)✓ Dependencies are satisfied$(COLOR_RESET)"

check-proto: ## Check if protoc is installed
	@echo "$(COLOR_BLUE)Checking protoc installation...$(COLOR_RESET)"
	@if command -v $(PROTOC) >/dev/null 2>&1; then \
		$(PROTOC) --version; \
		echo "$(COLOR_GREEN)✓ protoc is installed$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)protoc is not installed$(COLOR_RESET)"; \
		echo "$(COLOR_YELLOW)Install: https://grpc.io/docs/protoc-installation/$(COLOR_RESET)"; \
	fi

# ============================================================================
# Example Targets
# ============================================================================

example: build ## Run basic example
	@echo "$(COLOR_BLUE)Running basic example...$(COLOR_RESET)"
	$(GO) run examples/test_basic.go

# ============================================================================
# End of Makefile
# ============================================================================

# ============================================================================
# Bootstrap Targets (Self-Hosting/Self-Improvement)
# ============================================================================

.PHONY: bootstrap bootstrap-check bootstrap-ci

## bootstrap: Run BuildBureau in self-hosting mode to improve itself
bootstrap: build
	@echo "$(BLUE)Starting BuildBureau Bootstrap Mode...$(NC)"
	@./bootstrap/bootstrap.sh

## bootstrap-check: Verify bootstrap environment is ready
bootstrap-check:
	@echo "Checking bootstrap environment..."
	@test -f bootstrap/config.yaml || (echo "$(RED)Error: bootstrap/config.yaml not found$(NC)" && exit 1)
	@test -d bootstrap/agents || (echo "$(RED)Error: bootstrap/agents/ not found$(NC)" && exit 1)
	@test -n "$$GEMINI_API_KEY" -o -n "$$OPENAI_API_KEY" -o -n "$$CLAUDE_API_KEY" || \
(echo "$(RED)Error: No LLM API key set$(NC)" && exit 1)
	@echo "$(GREEN)✓ Bootstrap environment ready$(NC)"

## bootstrap-ci: Run bootstrap mode in CI (non-interactive)
bootstrap-ci: build
	@echo "Running BuildBureau in bootstrap mode (CI)..."
	@export BUILDBUREAU_CONFIG=bootstrap/config.yaml && \
./build/buildbureau --non-interactive || true
