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

.PHONY: test/all
## run all tests
test/all:
	@$(call green,"Running all tests...")
	@go test -race -cover -timeout=$(GOTEST_TIMEOUT) ./cmd/... ./internal/... ./pkg/...

.PHONY: test/cmd
## run tests for cmd
test/cmd:
	@$(call green,"Running cmd tests...")
	@go test -race -cover -timeout=$(GOTEST_TIMEOUT) ./cmd/...

.PHONY: test/internal
## run tests for internal
test/internal:
	@$(call green,"Running internal tests...")
	@go test -race -cover -timeout=$(GOTEST_TIMEOUT) ./internal/...

.PHONY: test/pkg
## run tests for pkg
test/pkg:
	@$(call green,"Running pkg tests...")
	@go test -race -cover -timeout=$(GOTEST_TIMEOUT) ./pkg/...

.PHONY: test/coverage
## run tests with coverage
test/coverage:
	@$(call green,"Running tests with coverage...")
	@mkdir -p coverage
	@go test -race -covermode=atomic -coverprofile=coverage/coverage.out ./cmd/... ./internal/... ./pkg/...

.PHONY: test/coverage/html
## generate HTML coverage report
test/coverage/html: test/coverage
	@$(call green,"Generating HTML coverage report...")
	@go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@echo "Coverage report generated at coverage/coverage.html"

.PHONY: test/coverage/func
## show function coverage
test/coverage/func: test/coverage
	@go tool cover -func=coverage/coverage.out

.PHONY: test/bench
## run benchmarks
test/bench:
	@$(call green,"Running benchmarks...")
	@go test -bench=. -benchmem -run=^$ ./...

.PHONY: test/race
## run tests with race detector
test/race:
	@$(call green,"Running tests with race detector...")
	@go test -race -timeout=$(GOTEST_TIMEOUT) ./...

.PHONY: test/short
## run short tests
test/short:
	@$(call green,"Running short tests...")
	@go test -short -timeout=60s ./...

.PHONY: test/verbose
## run tests with verbose output
test/verbose:
	@$(call green,"Running tests with verbose output...")
	@go test -v -race -cover -timeout=$(GOTEST_TIMEOUT) ./...
