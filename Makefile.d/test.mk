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

