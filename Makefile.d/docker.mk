#
# Copyright (C) 2024 BuildBureau team
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

.PHONY: docker/build/all
## build all Docker images
docker/build/all: \
	docker/build/buildbureau

.PHONY: docker/build/buildbureau
## build BuildBureau Docker image
docker/build/buildbureau:
	@$(call green,"Building BuildBureau Docker image...")
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg BUILD_DATE=$(DATETIME) \
		-t $(IMAGE_NAME):$(TAG) \
		-t $(IMAGE_NAME):latest \
		-f Dockerfile \
		.

.PHONY: docker/build/multiarch
## build multi-architecture Docker images
docker/build/multiarch:
	@$(call green,"Building multi-architecture Docker images...")
	@docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--build-arg VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg BUILD_DATE=$(DATETIME) \
		-t $(IMAGE_NAME):$(TAG) \
		-t $(IMAGE_NAME):latest \
		-f Dockerfile \
		--push \
		.

.PHONY: docker/run
## run Docker container
docker/run:
	@$(call green,"Running BuildBureau container...")
	@docker run -d \
		--name buildbureau \
		-e GEMINI_API_KEY=$(GEMINI_API_KEY) \
		-e OPENAI_API_KEY=$(OPENAI_API_KEY) \
		-e CLAUDE_API_KEY=$(CLAUDE_API_KEY) \
		-v $(PWD)/data:/app/data \
		$(IMAGE_NAME):$(TAG)

.PHONY: docker/run/interactive
## run Docker container interactively
docker/run/interactive:
	@$(call green,"Running BuildBureau container interactively...")
	@docker run -it --rm \
		-e GEMINI_API_KEY=$(GEMINI_API_KEY) \
		-e OPENAI_API_KEY=$(OPENAI_API_KEY) \
		-e CLAUDE_API_KEY=$(CLAUDE_API_KEY) \
		-v $(PWD)/data:/app/data \
		$(IMAGE_NAME):$(TAG)

.PHONY: docker/stop
## stop Docker container
docker/stop:
	@$(call green,"Stopping BuildBureau container...")
	@docker stop buildbureau || true
	@docker rm buildbureau || true

.PHONY: docker/push
## push Docker image to registry
docker/push:
	@$(call green,"Pushing Docker image...")
	@docker push $(IMAGE_NAME):$(TAG)
	@docker push $(IMAGE_NAME):latest

.PHONY: docker/scan
## scan Docker image for vulnerabilities
docker/scan:
	@$(call green,"Scanning Docker image for vulnerabilities...")
	@docker scout cves $(IMAGE_NAME):$(TAG) || \
		trivy image --severity HIGH,CRITICAL $(IMAGE_NAME):$(TAG)

.PHONY: docker/compose/up
## start Docker Compose stack
docker/compose/up:
	@$(call green,"Starting Docker Compose stack...")
	@docker-compose up -d

.PHONY: docker/compose/down
## stop Docker Compose stack
docker/compose/down:
	@$(call green,"Stopping Docker Compose stack...")
	@docker-compose down

.PHONY: docker/compose/logs
## view Docker Compose logs
docker/compose/logs:
	@docker-compose logs -f

.PHONY: docker/clean
## clean Docker images and containers
docker/clean:
	@$(call yellow,"Cleaning Docker resources...")
	@docker system prune -f
