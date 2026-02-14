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

.PHONY: lint/all
## run all linters
lint/all: \
	lint/go \
	lint/yaml \
	lint/actions

.PHONY: lint/go
## run Go linters
lint/go: tools/golangci-lint
	@$(call green,"Running Go linters...")
	@golangci-lint run --config .golangci.json

.PHONY: lint/go/fix
## run Go linters with auto-fix
lint/go/fix: tools/golangci-lint
	@$(call green,"Running Go linters with auto-fix...")
	@golangci-lint run --config .golangci.json --fix

.PHONY: lint/yaml
## lint YAML files
lint/yaml: tools/yamlfmt
	@$(call green,"Linting YAML files...")
	@yamlfmt -lint .

.PHONY: lint/actions
## lint GitHub Actions workflows
lint/actions: tools/actionlint
	@$(call green,"Linting GitHub Actions workflows...")
	@actionlint

.PHONY: lint/security
## run security linters
lint/security: tools/gosec
	@$(call green,"Running security linters...")
	@gosec -quiet ./...

.PHONY: lint/static
## run static analysis
lint/static: tools/staticcheck
	@$(call green,"Running static analysis...")
	@staticcheck ./...
