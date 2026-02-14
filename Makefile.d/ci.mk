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

