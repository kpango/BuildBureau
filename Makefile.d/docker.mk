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
# Docker Targets
# ============================================================================

docker: docker-build ## Build Docker image (alias for docker-build)

docker-build: ## Build Docker image (replaces docker/build.sh)
	@echo "$(COLOR_BLUE)Building Docker image: $(DOCKER_FULL_IMAGE)$(COLOR_RESET)"
	@docker build -t $(DOCKER_FULL_IMAGE) .
	@echo "$(COLOR_GREEN)✓ Build successful!$(COLOR_RESET)"
	@echo "$(COLOR_GREEN)Image: $(DOCKER_FULL_IMAGE)$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BLUE)To run the image:$(COLOR_RESET)"
	@echo "  docker run -e GEMINI_API_KEY=your-key $(DOCKER_FULL_IMAGE)"
	@echo ""
	@echo "$(COLOR_BLUE)Or use docker-compose:$(COLOR_RESET)"
	@echo "  make docker-compose-up"

docker-build-no-cache: ## Build Docker image without cache
	@echo "$(COLOR_BLUE)Building Docker image (no cache): $(DOCKER_FULL_IMAGE)$(COLOR_RESET)"
	@docker build --no-cache -t $(DOCKER_FULL_IMAGE) .
	@echo "$(COLOR_GREEN)✓ Docker image built: $(DOCKER_FULL_IMAGE)$(COLOR_RESET)"

docker-build-multi: ## Build multi-architecture Docker image
	@echo "$(COLOR_BLUE)Building multi-arch Docker image: $(DOCKER_FULL_IMAGE)$(COLOR_RESET)"
	@docker buildx build --platform $(DOCKER_PLATFORMS) -t $(DOCKER_FULL_IMAGE) .
	@echo "$(COLOR_GREEN)✓ Multi-arch Docker image built$(COLOR_RESET)"

docker-push: ## Push Docker image to registry
	@echo "$(COLOR_BLUE)Pushing Docker image: $(DOCKER_FULL_IMAGE)$(COLOR_RESET)"
	@docker push $(DOCKER_FULL_IMAGE)
	@echo "$(COLOR_GREEN)✓ Docker image pushed$(COLOR_RESET)"

docker-run: ## Run Docker container in daemon mode (replaces docker/run.sh)
	@if [ -z "$(GEMINI_API_KEY)" ] && [ -z "$(OPENAI_API_KEY)" ] && [ -z "$(CLAUDE_API_KEY)" ]; then \
		echo "$(COLOR_YELLOW)Warning: No LLM API key provided!$(COLOR_RESET)"; \
		echo "$(COLOR_YELLOW)Please set at least one of:$(COLOR_RESET)"; \
		echo "  export GEMINI_API_KEY=your-key"; \
		echo "  export OPENAI_API_KEY=your-key"; \
		echo "  export CLAUDE_API_KEY=your-key"; \
		echo ""; \
	fi
	@CONTAINER_NAME=$${CONTAINER_NAME:-buildbureau}; \
	echo "$(COLOR_BLUE)Starting BuildBureau container...$(COLOR_RESET)"; \
	echo "$(COLOR_BLUE)Container name: $(COLOR_GREEN)$$CONTAINER_NAME$(COLOR_RESET)"; \
	echo "$(COLOR_BLUE)Image: $(COLOR_GREEN)$(DOCKER_FULL_IMAGE)$(COLOR_RESET)"; \
	if docker ps -a --format '{{.Names}}' | grep -q "^$$CONTAINER_NAME$$"; then \
		echo "$(COLOR_BLUE)Stopping existing container...$(COLOR_RESET)"; \
		docker stop $$CONTAINER_NAME > /dev/null 2>&1 || true; \
		docker rm $$CONTAINER_NAME > /dev/null 2>&1 || true; \
	fi; \
	docker run -d \
		--name $$CONTAINER_NAME \
		-e GEMINI_API_KEY="$(GEMINI_API_KEY)" \
		-e OPENAI_API_KEY="$(OPENAI_API_KEY)" \
		-e CLAUDE_API_KEY="$(CLAUDE_API_KEY)" \
		-e OPENAI_MODEL="$(OPENAI_MODEL:-gpt-4-turbo-preview)" \
		-e CLAUDE_MODEL="$(CLAUDE_MODEL:-claude-3-5-sonnet-20241022)" \
		-e SLACK_TOKEN="$(SLACK_TOKEN)" \
		-v buildbureau-data:/app/data \
		-p 8080:8080 \
		--restart unless-stopped \
		$(DOCKER_FULL_IMAGE) && \
	echo "$(COLOR_GREEN)✓ Container started successfully!$(COLOR_RESET)" && \
	echo "" && \
	echo "$(COLOR_BLUE)Container details:$(COLOR_RESET)" && \
	docker ps --filter "name=$$CONTAINER_NAME" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" && \
	echo "" && \
	echo "$(COLOR_BLUE)View logs:$(COLOR_RESET)" && \
	echo "  docker logs -f $$CONTAINER_NAME" && \
	echo "" && \
	echo "$(COLOR_BLUE)Execute interactive shell:$(COLOR_RESET)" && \
	echo "  docker exec -it $$CONTAINER_NAME sh" && \
	echo "" && \
	echo "$(COLOR_BLUE)Stop container:$(COLOR_RESET)" && \
	echo "  docker stop $$CONTAINER_NAME"

docker-run-interactive: ## Run Docker container interactively
	@echo "$(COLOR_BLUE)Running Docker container interactively...$(COLOR_RESET)"
	@docker run --rm -it \
		-e GEMINI_API_KEY \
		-e OPENAI_API_KEY \
		-e CLAUDE_API_KEY \
		-v $$(pwd)/data:/app/data \
		$(DOCKER_FULL_IMAGE)

docker-test: ## Test Docker setup (replaces docker/test.sh)
	@echo "=== BuildBureau Docker Test ==="
	@echo ""
	@echo "$(COLOR_BLUE)Checking Docker...$(COLOR_RESET)"
	@if ! command -v docker &> /dev/null; then \
		echo "$(COLOR_RED)✗ Docker not found$(COLOR_RESET)"; \
		exit 1; \
	fi
	@echo "$(COLOR_GREEN)✓ Docker found: $$(docker --version)$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BLUE)Checking Docker Compose...$(COLOR_RESET)"
	@if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then \
		echo "$(COLOR_RED)✗ Docker Compose not found$(COLOR_RESET)"; \
		exit 1; \
	fi
	@echo "$(COLOR_GREEN)✓ Docker Compose found$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BLUE)Checking API keys...$(COLOR_RESET)"
	@if [ -z "$(GEMINI_API_KEY)" ] && [ -z "$(OPENAI_API_KEY)" ] && [ -z "$(CLAUDE_API_KEY)" ]; then \
		echo "$(COLOR_RED)✗ No API key found$(COLOR_RESET)"; \
		echo "$(COLOR_BLUE)Please set at least one:$(COLOR_RESET)"; \
		echo "  export GEMINI_API_KEY=your-key"; \
		echo "  export OPENAI_API_KEY=your-key"; \
		echo "  export CLAUDE_API_KEY=your-key"; \
		exit 1; \
	fi
	@echo "$(COLOR_GREEN)✓ API key(s) found$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BLUE)Building Docker image...$(COLOR_RESET)"
	@docker build -t buildbureau:test . || { \
		echo "$(COLOR_RED)✗ Build failed$(COLOR_RESET)"; \
		exit 1; \
	}
	@echo "$(COLOR_GREEN)✓ Build successful$(COLOR_RESET)"
	@echo ""
	@IMAGE_SIZE=$$(docker images buildbureau:test --format "{{.Size}}"); \
	echo "$(COLOR_BLUE)Image size: $(COLOR_GREEN)$$IMAGE_SIZE$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BLUE)Testing container...$(COLOR_RESET)"
	@docker run --rm buildbureau:test --version || { \
		echo "$(COLOR_RED)✗ Container test failed$(COLOR_RESET)"; \
		exit 1; \
	}
	@echo "$(COLOR_GREEN)✓ Container test successful$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_GREEN)=== All checks passed! ===$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BLUE)Next steps:$(COLOR_RESET)"
	@echo "1. Run with Docker Compose:"
	@echo "   make docker-compose-up"
	@echo ""
	@echo "2. Or run directly:"
	@echo "   make docker-run"
	@echo ""
	@echo "3. View logs:"
	@echo "   docker logs -f buildbureau"
	@echo ""
	@echo "$(COLOR_GREEN)BuildBureau is ready to use! 🎉$(COLOR_RESET)"

docker-compose-up: ## Start services with docker-compose
	@echo "$(COLOR_BLUE)Starting Docker Compose services...$(COLOR_RESET)"
	@docker-compose up -d
	@echo "$(COLOR_GREEN)✓ Services started$(COLOR_RESET)"

docker-compose-down: ## Stop services with docker-compose
	@echo "$(COLOR_BLUE)Stopping Docker Compose services...$(COLOR_RESET)"
	@docker-compose down
	@echo "$(COLOR_GREEN)✓ Services stopped$(COLOR_RESET)"

