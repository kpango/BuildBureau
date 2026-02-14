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

