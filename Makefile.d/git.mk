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

