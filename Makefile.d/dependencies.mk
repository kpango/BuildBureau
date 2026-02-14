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
# Dependencies
# ============================================================================

.PHONY: deps
deps: ## Download and tidy dependencies
	@echo "$(COLOR_BLUE)Downloading dependencies...$(COLOR_RESET)"
	$(GO) mod download
	$(GO) mod tidy
	@echo "$(COLOR_GREEN)✓ Dependencies ready$(COLOR_RESET)"

.PHONY: deps-update
deps-update: ## Update dependencies to latest versions
	@echo "$(COLOR_BLUE)Updating dependencies...$(COLOR_RESET)"
	$(GO) get -u ./...
	$(GO) mod tidy
	@echo "$(COLOR_GREEN)✓ Dependencies updated$(COLOR_RESET)"

.PHONY: deps-verify
deps-verify: ## Verify dependencies
	@echo "$(COLOR_BLUE)Verifying dependencies...$(COLOR_RESET)"
	$(GO) mod verify
	@echo "$(COLOR_GREEN)✓ Dependencies verified$(COLOR_RESET)"

.PHONY: deps-graph
deps-graph: ## Display dependency graph
	@echo "$(COLOR_BLUE)Dependency graph:$(COLOR_RESET)"
	$(GO) mod graph

.PHONY: deps-tidy
deps-tidy: ## Tidy go.mod and go.sum
	@echo "$(COLOR_BLUE)Tidying dependencies...$(COLOR_RESET)"
	$(GO) mod tidy
	@echo "$(COLOR_GREEN)✓ Dependencies tidied$(COLOR_RESET)"

.PHONY: deps-reset
deps-reset: ## Reset go.mod from go.mod.default and update all dependencies
	@echo "$(COLOR_BLUE)Resetting dependencies from go.mod.default...$(COLOR_RESET)"
	@if [ ! -f $(ROOTDIR)/go.mod.default ]; then \
		echo "$(COLOR_RED)Error: go.mod.default not found$(COLOR_RESET)"; \
		exit 1; \
	fi
	@echo "$(COLOR_YELLOW)Backing up current go.mod to go.mod.backup$(COLOR_RESET)"
	@cp $(ROOTDIR)/go.mod $(ROOTDIR)/go.mod.backup 2>/dev/null || true
	@echo "$(COLOR_BLUE)Copying go.mod.default to go.mod$(COLOR_RESET)"
	@cp $(ROOTDIR)/go.mod.default $(ROOTDIR)/go.mod
	@echo "$(COLOR_BLUE)Updating Go version to $(shell $(GO) version | awk '{print $$3}' | sed 's/go//')$(COLOR_RESET)"
	@sed -i "s/go [0-9]\+\.[0-9]\+\(\.[0-9]\+\)\?/go $$($(GO) version | awk '{print $$3}' | sed 's/go//')/" $(ROOTDIR)/go.mod
	@echo "$(COLOR_BLUE)Resolving dependencies with go mod tidy$(COLOR_RESET)"
	@$(GO) mod tidy
	@echo "$(COLOR_GREEN)✓ Dependencies reset and updated$(COLOR_RESET)"
	@echo "$(COLOR_CYAN)Note: Original go.mod saved as go.mod.backup$(COLOR_RESET)"

