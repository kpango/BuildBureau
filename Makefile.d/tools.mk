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

.PHONY: tools/install
## install all tools
tools/install: \
	tools/golangci-lint \
	tools/goimports \
	tools/gofumpt \
	tools/yamlfmt \
	tools/reviewdog

.PHONY: tools/golangci-lint
## install golangci-lint
tools/golangci-lint:
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) $(GOLANGCILINT_VERSION); \
	else \
		echo "golangci-lint already installed"; \
	fi

.PHONY: tools/goimports
## install goimports
tools/goimports:
	@if ! command -v goimports >/dev/null 2>&1; then \
		echo "Installing goimports..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
	else \
		echo "goimports already installed"; \
	fi

.PHONY: tools/gofumpt
## install gofumpt
tools/gofumpt:
	@if ! command -v gofumpt >/dev/null 2>&1; then \
		echo "Installing gofumpt..."; \
		go install mvdan.cc/gofumpt@latest; \
	else \
		echo "gofumpt already installed"; \
	fi

.PHONY: tools/yamlfmt
## install yamlfmt
tools/yamlfmt:
	@if ! command -v yamlfmt >/dev/null 2>&1; then \
		echo "Installing yamlfmt..."; \
		go install github.com/google/yamlfmt/cmd/yamlfmt@latest; \
	else \
		echo "yamlfmt already installed"; \
	fi

.PHONY: tools/reviewdog
## install reviewdog
tools/reviewdog:
	@if ! command -v reviewdog >/dev/null 2>&1; then \
		echo "Installing reviewdog..."; \
		go install github.com/reviewdog/reviewdog/cmd/reviewdog@latest; \
	else \
		echo "reviewdog already installed"; \
	fi

.PHONY: tools/buf
## install buf
tools/buf:
	@if ! command -v buf >/dev/null 2>&1; then \
		echo "Installing buf..."; \
		go install github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION); \
	else \
		echo "buf already installed"; \
	fi

.PHONY: tools/actionlint
## install actionlint
tools/actionlint:
	@if ! command -v actionlint >/dev/null 2>&1; then \
		echo "Installing actionlint..."; \
		go install github.com/rhysd/actionlint/cmd/actionlint@latest; \
	else \
		echo "actionlint already installed"; \
	fi

.PHONY: tools/gosec
## install gosec
tools/gosec:
	@if ! command -v gosec >/dev/null 2>&1; then \
		echo "Installing gosec..."; \
		go install github.com/securego/gosec/v2/cmd/gosec@latest; \
	else \
		echo "gosec already installed"; \
	fi

.PHONY: tools/staticcheck
## install staticcheck
tools/staticcheck:
	@if ! command -v staticcheck >/dev/null 2>&1; then \
		echo "Installing staticcheck..."; \
		go install honnef.co/go/tools/cmd/staticcheck@latest; \
	else \
		echo "staticcheck already installed"; \
	fi
