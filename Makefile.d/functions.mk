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

# Color output functions
red    = printf "\x1b[31m\#\# %s\x1b[0m\n" $1
green  = printf "\x1b[32m\#\# %s\x1b[0m\n" $1
yellow = printf "\x1b[33m\#\# %s\x1b[0m\n" $1
blue   = printf "\x1b[34m\#\# %s\x1b[0m\n" $1
pink   = printf "\x1b[35m\#\# %s\x1b[0m\n" $1
cyan   = printf "\x1b[36m\#\# %s\x1b[0m\n" $1

# Tool installation helper
define go-tool-install
	@echo "Installing $(1)..."
	@go install $(1)@$(2)
endef

# Directory creation helper
define mkdir
	@mkdir -p $1
endef

# Go build function with version info
define go-build
	@echo "Building $(1)..."
	CGO_ENABLED=$(CGO_ENABLED) \
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
	-o $(2) \
	$(1)
endef

# Go lint function
define go-lint
	@echo "Running golangci-lint..."
	@golangci-lint run --config .golangci.json $(1)
endef

# Go test function with coverage
define go-test
	@echo "Running tests for $(1)..."
	@go test -race -cover -coverprofile=$(2) $(1)
endef

# Format function
define go-format
	@echo "Formatting Go code..."
	@gofmt -s -w $(1)
	@goimports -w $(1)
endef
