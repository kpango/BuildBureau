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

