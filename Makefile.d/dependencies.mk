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

