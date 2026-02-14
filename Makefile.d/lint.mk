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

