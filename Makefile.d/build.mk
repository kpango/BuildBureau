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

.PHONY: binary/build
## build all binaries
binary/build: \
	cmd/buildbureau/buildbureau

.PHONY: binary/build/release
## build release binary
binary/build/release: cmd/buildbureau/buildbureau-release

.PHONY: binary/build/debug
## build debug binary
binary/build/debug: cmd/buildbureau/buildbureau-debug

cmd/buildbureau/buildbureau:
	$(eval CGO_ENABLED = 1)
	$(call go-build,./cmd/buildbureau,buildbureau)

cmd/buildbureau/buildbureau-release:
	$(eval CGO_ENABLED = 1)
	@echo "Building release binary..."
	CGO_ENABLED=1 \
	GO111MODULE=on \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	go build \
	--ldflags="-w -s \
	-X 'main.Version=$(VERSION)' \
	-X 'main.GitCommit=$(GIT_COMMIT)' \
	-X 'main.BuildTime=$(DATETIME)' \
	-X 'main.GoVersion=$(GO_VERSION)'" \
	-trimpath \
	-o buildbureau-release \
	./cmd/buildbureau

cmd/buildbureau/buildbureau-debug:
	$(eval CGO_ENABLED = 1)
	@echo "Building debug binary..."
	CGO_ENABLED=1 \
	GO111MODULE=on \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	go build \
	-gcflags="all=-N -l" \
	--ldflags="-X 'main.Version=$(VERSION)' \
	-X 'main.GitCommit=$(GIT_COMMIT)' \
	-X 'main.BuildTime=$(DATETIME)' \
	-X 'main.GoVersion=$(GO_VERSION)'" \
	-o buildbureau-debug \
	./cmd/buildbureau

.PHONY: binary/install
## install binary to GOBIN
binary/install:
	@echo "Installing buildbureau to $(GOBIN)..."
	@go install ./cmd/buildbureau
