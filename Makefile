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
# High-functionality modular Makefile for building, testing, and managing the project
# Inspired by vdaas/vald repository structure
# Can be used in Dockerfiles and CI/CD pipelines for standardization

SHELL := bash

# ============================================================================
# Include Variables and Functions
# ============================================================================

include Makefile.d/variables.mk
include Makefile.d/functions.mk

# ============================================================================
# Default Target
# ============================================================================

.DEFAULT_GOAL := help

# ============================================================================
# Help Target
# ============================================================================

.PHONY: help
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
# Check Targets
# ============================================================================

.PHONY: check
check: check-go check-deps check-proto ## Run all checks

.PHONY: check-go
check-go: ## Check Go installation and version
	@echo "$(COLOR_BLUE)Checking Go installation...$(COLOR_RESET)"
	@$(GO) version
	@echo "$(COLOR_GREEN)✓ Go is installed$(COLOR_RESET)"

.PHONY: check-deps
check-deps: ## Check if dependencies are satisfied
	@echo "$(COLOR_BLUE)Checking dependencies...$(COLOR_RESET)"
	@$(GO) mod verify
	@echo "$(COLOR_GREEN)✓ Dependencies are satisfied$(COLOR_RESET)"

.PHONY: check-proto
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
# Include Module Makefiles
# ============================================================================

include Makefile.d/tools.mk
include Makefile.d/build.mk
include Makefile.d/test.mk
include Makefile.d/proto.mk
include Makefile.d/dependencies.mk
include Makefile.d/format.mk
include Makefile.d/lint.mk
include Makefile.d/security.mk
include Makefile.d/docker.mk
include Makefile.d/ci.mk
include Makefile.d/git.mk

# ============================================================================
# Bootstrap Targets (Self-Hosting/Self-Improvement)
# ============================================================================

.PHONY: bootstrap bootstrap-check bootstrap-ci

bootstrap: build ## Run BuildBureau in self-hosting mode to improve itself
	@echo "$(COLOR_BLUE)Starting BuildBureau Bootstrap Mode...$(COLOR_RESET)"
	@./bootstrap/bootstrap.sh

bootstrap-check: ## Verify bootstrap environment is ready
	@echo "Checking bootstrap environment..."
	@test -f bootstrap/config.yaml || (echo "$(COLOR_RED)Error: bootstrap/config.yaml not found$(COLOR_RESET)" && exit 1)
	@test -d bootstrap/agents || (echo "$(COLOR_RED)Error: bootstrap/agents/ not found$(COLOR_RESET)" && exit 1)
	@test -n "$$GEMINI_API_KEY" -o -n "$$OPENAI_API_KEY" -o -n "$$CLAUDE_API_KEY" || \
		(echo "$(COLOR_RED)Error: No LLM API key set$(COLOR_RESET)" && exit 1)
	@echo "$(COLOR_GREEN)✓ Bootstrap environment ready$(COLOR_RESET)"

bootstrap-ci: build ## Run bootstrap mode in CI (non-interactive)
	@echo "Running BuildBureau in bootstrap mode (CI)..."
	@export BUILDBUREAU_CONFIG=bootstrap/config.yaml && \
		./build/buildbureau --non-interactive || true

# ============================================================================
# Example Targets
# ============================================================================

.PHONY: example

example: build ## Run basic example
	@echo "$(COLOR_BLUE)Running basic example...$(COLOR_RESET)"
	$(GO) run examples/test_basic.go

# ============================================================================
# Phony Targets Declaration
# ============================================================================

.PHONY: all

all: clean deps proto build test ## Build everything (clean, deps, proto, build, test)
