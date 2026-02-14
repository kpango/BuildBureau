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
# Proto Targets
# ============================================================================

proto: | $(PROTOC_GEN_GO_STAMP) $(PROTOC_GEN_GO_GRPC_STAMP) ## Generate protobuf files from .proto definitions
	@echo "$(COLOR_BLUE)Generating protobuf files...$(COLOR_RESET)"
	@mkdir -p $(PROTO_DIR)
	$(PROTOC) --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/agent.proto
	@echo "$(COLOR_GREEN)✓ Proto generation complete$(COLOR_RESET)"

proto-clean: ## Clean generated proto files
	@echo "$(COLOR_BLUE)Cleaning generated proto files...$(COLOR_RESET)"
	@rm -f $(PROTO_DIR)/*.pb.go
	@echo "$(COLOR_GREEN)✓ Proto files cleaned$(COLOR_RESET)"

