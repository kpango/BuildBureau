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
