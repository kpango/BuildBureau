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
# Build Targets
# ============================================================================

all: clean deps proto build test ## Build everything (clean, deps, proto, build, test)

build: ## Build the application binary
	@echo "$(COLOR_BLUE)Building $(APP_NAME)...$(COLOR_RESET)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build \
		-ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(APP_NAME) \
		$(CMD_DIR)
	@echo "$(COLOR_GREEN)✓ Build complete: $(BUILD_DIR)/$(APP_NAME)$(COLOR_RESET)"

build-debug: ## Build with debug symbols and no optimization
	@echo "$(COLOR_BLUE)Building debug version...$(COLOR_RESET)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build \
		-gcflags="all=-N -l" \
		-ldflags "$(DEBUG_LDFLAGS)" \
		-o $(BUILD_DIR)/$(APP_NAME)-debug \
		$(CMD_DIR)
	@echo "$(COLOR_GREEN)✓ Debug build complete: $(BUILD_DIR)/$(APP_NAME)-debug$(COLOR_RESET)"

build-release: ## Build optimized release binary
	@echo "$(COLOR_BLUE)Building release version...$(COLOR_RESET)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build \
		-ldflags "$(LDFLAGS)" \
		-trimpath \
		-o $(BUILD_DIR)/$(APP_NAME) \
		$(CMD_DIR)
	@echo "$(COLOR_GREEN)✓ Release build complete: $(BUILD_DIR)/$(APP_NAME)$(COLOR_RESET)"

build-static: ## Build static binary (CGO disabled where possible)
	@echo "$(COLOR_BLUE)Building static binary...$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)Note: SQLite requires CGO, so this may not be fully static$(COLOR_RESET)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 $(GO) build \
		-ldflags "$(LDFLAGS) -extldflags '-static'" \
		-tags 'netgo osusergo' \
		-o $(BUILD_DIR)/$(APP_NAME)-static \
		$(CMD_DIR)
	@echo "$(COLOR_GREEN)✓ Static build complete: $(BUILD_DIR)/$(APP_NAME)-static$(COLOR_RESET)"

build-all: ## Build for multiple platforms (linux/amd64, linux/arm64, darwin/amd64, darwin/arm64)
	@echo "$(COLOR_BLUE)Building for multiple platforms...$(COLOR_RESET)"
	@mkdir -p $(DIST_DIR)
	@for os in linux darwin; do \
		for arch in amd64 arm64; do \
			output="$(DIST_DIR)/$(APP_NAME)-$$os-$$arch"; \
			if [ "$$os" = "windows" ]; then output="$$output.exe"; fi; \
			echo "$(COLOR_CYAN)Building $$os/$$arch...$(COLOR_RESET)"; \
			CGO_ENABLED=$(CGO_ENABLED) GOOS=$$os GOARCH=$$arch $(GO) build \
				-ldflags "$(LDFLAGS)" \
				-o $$output \
				$(CMD_DIR) || exit 1; \
		done \
	done
	@echo "$(COLOR_GREEN)✓ Multi-platform builds complete in $(DIST_DIR)/$(COLOR_RESET)"

