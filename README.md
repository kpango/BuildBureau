# BuildBureau

[![CI](https://github.com/kpango/BuildBureau/actions/workflows/ci.yml/badge.svg)](https://github.com/kpango/BuildBureau/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/kpango/BuildBureau/branch/main/graph/badge.svg)](https://codecov.io/gh/kpango/BuildBureau)
[![CodeQL](https://github.com/kpango/BuildBureau/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/kpango/BuildBureau/actions/workflows/codeql-analysis.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/kpango/BuildBureau)](https://goreportcard.com/report/github.com/kpango/BuildBureau)
[![Release](https://img.shields.io/github/v/release/kpango/BuildBureau)](https://github.com/kpango/BuildBureau/releases/latest)
[![License](https://img.shields.io/github/license/kpango/BuildBureau)](LICENSE)
[![Docker](https://ghcr-badge.deta.dev/kpango/buildbureau/latest_tag?trim=major&label=Docker)](https://github.com/kpango/BuildBureau/pkgs/container/buildbureau)

🏢 **BuildBureau** is a multi-agent system that simulates a virtual software
development company composed of hierarchical AI agents with **real LLM
integration** supporting **multiple providers** and **persistent memory**.

## 🌟 Highlights

- 🤖 **Multi-Provider LLM Support**: Native integration with **Gemini**,
  **OpenAI**, and **Claude**
- 🧠 **Agent Memory with Learning**: Agents remember conversations, learn from
  past tasks, and improve over time
- 💾 **Optimized SQLite Storage**: Fast `mattn/go-sqlite3` driver with smart
  memory usage per agent
- ⚡ **Real AI Integration**: Actual code generation and technical
  specifications
- 🔗 **Live Slack Notifications**: Real-time updates to your Slack workspace
- 🏗️ **Production-Ready Architecture**: Built with industry-standard libraries
- 🎯 **End-to-End Implementation**: Not just stubs - working AI-powered agents

## Overview

BuildBureau implements a five-layer organizational hierarchy:

- **President**: Clarifies client instructions and defines high-level
  requirements
- **Secretary**: Handles delegation with memory-informed decisions, scheduling,
  monitoring, and knowledge management
- **Director**: Performs research and decomposes projects into department-level
  tasks
- **Manager**: Uses LLM + past designs to produce detailed software
  specifications
- **Engineer**: Uses LLM + past implementations to generate actual code

## Features

### ✅ Fully Implemented

- 🎨 **Multi-Provider LLM Support**:
  - ✅ **Google Gemini** (`google.golang.org/genai`) - Fast and cost-effective
  - ✅ **OpenAI** (`github.com/sashabaranov/go-openai`) - GPT-4, GPT-3.5
  - ✅ **Anthropic Claude** (`github.com/liushuangls/go-anthropic/v2`) - Claude
    3.5 Sonnet, Opus, Haiku
  - 🔌 **Remote Agent API** for custom providers
- 🧠 **Persistent Memory System**:
  - ✅ **SQLite Storage** (`github.com/mattn/go-sqlite3`) - Fast, standard
    driver
  - ✅ **Vald Vector DB** - Semantic similarity search (optional)
  - ✅ **Agent Learning** - Engineers learn from past code, Managers learn from
    past designs
  - ✅ **Smart Delegation** - Secretaries track performance and route
    intelligently
  - ✅ **Conversation History** - Track all agent interactions
  - ✅ **Knowledge Base** - Store and retrieve learned information
  - ✅ **Automatic Expiration** - Configurable retention policies
- 🎯 **ADK-Powered Agents**: Google's Agent Development Kit (ADK) integration
  for structured agents
- 💬 **Real Slack Integration**: Live notifications to Slack channels
- 🌐 **Remote Agent HTTP API**: Complete HTTP client for external LLM services
- 📡 **gRPC Server & Client**: Real TCP listeners and connection management
- 🏢 **Hierarchical Agent System**: Five-layer organization with round-robin
  load balancing
- 📝 **YAML-Driven Configuration**: Flexible, declarative system configuration
- 🔒 **Environment Variable Support**: Secure API key management
- 🖥️ **Terminal-Based UI**: Interactive TUI built with Bubble Tea
- 🔄 **Task Delegation**: Automatic routing through agent hierarchy with UUID
  tracking

### 🔧 Optional Enhancements

- 📋 **Proto Code Generation**: Can generate gRPC service code from .proto files
  with protoc
- 🔌 **Extensible Design**: Easy to add new agent types and capabilities

## Requirements

### Option 1: Docker (Recommended - Zero Dependencies!)

- Docker 20.10+ or Docker Desktop
- Docker Compose 2.0+ (optional, for easier orchestration)
- **At least one LLM API Key** (see below)

### Option 2: Local Build

- Go 1.26.0 or later
- gcc (for SQLite CGo compilation)
- protoc (optional, for gRPC code generation)

### API Keys (Required for both options)

- **At least one LLM API Key**:
  - **Gemini API Key** (free tier available):
    https://aistudio.google.com/app/apikey
  - **OpenAI API Key**: https://platform.openai.com/api-keys
  - **Claude API Key**: https://console.anthropic.com/
- Slack Bot Token (optional, for notifications)

## Installation

### 🐳 Option 1: Docker (Zero Dependencies!)

**Easiest way to run BuildBureau - no Go, no SQLite, no dependencies needed!**

```bash
# Clone the repository
git clone https://github.com/kpango/BuildBureau.git
cd BuildBureau

# Set your API key
export GEMINI_API_KEY="your-gemini-key"

# Run with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

**Or with Docker CLI:**

```bash
# Build the image
docker build -t buildbureau:latest .

# Run the container
docker run -d \
  --name buildbureau \
  -e GEMINI_API_KEY="your-key" \
  -v buildbureau-data:/app/data \
  -p 8080:8080 \
  buildbureau:latest
```

📚 **See [Docker Documentation](docs/DOCKER.md) for complete guide**

### 💻 Option 2: Local Build

```bash
# Clone the repository
git clone https://github.com/kpango/BuildBureau.git
cd BuildBureau

# Install dependencies and build
make deps
make build

# Or build everything in one command
make all
```

## 🛠️ Development with Makefile

BuildBureau includes a comprehensive, high-functionality Makefile for
standardized development, testing, and deployment:

```bash
# Show all available commands
make help

# Build targets
make build              # Build the application
make build-release      # Build optimized release binary
make build-all          # Build for multiple platforms

# Test targets
make test               # Run all tests
make test-coverage      # Run tests with coverage report
make test-bench         # Run benchmarks

# Docker targets
make docker-build       # Build Docker image
make docker-compose-up  # Start with docker-compose

# CI/CD targets (used in CI pipelines)
make ci-all            # Run all CI checks
make ci-lint           # Lint code
make ci-build          # Build for CI
make ci-test           # Test with coverage

# Development
make proto             # Generate protobuf files
make fmt               # Format code
make lint              # Run linters
make deps              # Install dependencies

# Utilities
make version           # Show version info
make clean             # Clean build artifacts
```

📚 **See [Makefile Documentation](docs/MAKEFILE.md) for complete guide**

The Makefile is used consistently across:

- ✅ Local development
- ✅ Docker builds (in Dockerfile)
- ✅ CI/CD pipelines (in GitHub Actions)

This ensures all build commands are standardized and reproducible!

## 🚀 CI/CD Pipeline

BuildBureau includes a comprehensive GitHub Actions CI/CD pipeline:

- **✅ Continuous Integration**: Automated builds and tests on every push/PR
- **🚀 Multi-Platform Releases**: Automated releases for Linux (amd64/arm64) and
  macOS (Intel/Apple Silicon)
- **🐳 Docker Publishing**: Multi-architecture images published to GitHub
  Container Registry
- **🔒 Security Scanning**: CodeQL analysis, dependency review, and
  vulnerability checks
- **📊 Code Quality**: 30+ linters via golangci-lint
- **📈 Code Coverage**: Automated coverage reports with Codecov integration
- **🏷️ Auto-Labeling**: Intelligent PR and issue labeling
- **📝 Release Notes**: Automatic changelog generation

### Quick CI Commands

```bash
# Run all CI checks locally
make ci-all

# Individual checks
make fmt-check    # Check code formatting
make ci-lint      # Run linters
make ci-build     # Build project
make ci-test      # Run tests
```

### Workflows

| Workflow          | Purpose                 | Status                                                                                                                                                                             |
| ----------------- | ----------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **CI**            | Build, test, lint       | [![CI](https://github.com/kpango/BuildBureau/actions/workflows/ci.yml/badge.svg)](https://github.com/kpango/BuildBureau/actions/workflows/ci.yml)                                  |
| **Release**       | Multi-platform binaries | [![Release](https://github.com/kpango/BuildBureau/actions/workflows/release.yml/badge.svg)](https://github.com/kpango/BuildBureau/actions/workflows/release.yml)                   |
| **Docker**        | Container publishing    | [![Docker](https://github.com/kpango/BuildBureau/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/kpango/BuildBureau/actions/workflows/docker-publish.yml)      |
| **CodeQL**        | Security scanning       | [![CodeQL](https://github.com/kpango/BuildBureau/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/kpango/BuildBureau/actions/workflows/codeql-analysis.yml)    |
| **GolangCI-Lint** | Code quality            | [![golangci-lint](https://github.com/kpango/BuildBureau/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/kpango/BuildBureau/actions/workflows/golangci-lint.yml) |
| **Coverage**      | Test coverage           | [![Coverage](https://github.com/kpango/BuildBureau/actions/workflows/coverage.yml/badge.svg)](https://github.com/kpango/BuildBureau/actions/workflows/coverage.yml)                |

📚 **See [GitHub Actions Documentation](docs/GITHUB_ACTIONS.md) for complete guide**

## Quick Start

### 1. Get LLM API Keys

Choose one or more providers:

- **Gemini** (free tier): https://aistudio.google.com/app/apikey
- **OpenAI**: https://platform.openai.com/api-keys
- **Claude**: https://console.anthropic.com/

### 2. Set Your API Keys

```bash
# Use one or multiple providers
export GEMINI_API_KEY="your-gemini-key"
export OPENAI_API_KEY="your-openai-key"     # Optional
export CLAUDE_API_KEY="your-claude-key"     # Optional
```

### 3. Run BuildBureau

```bash
./buildbureau
```

### 4. Try a Task

In the TUI, enter a task like:

```
Create a REST API in Go for managing a todo list with CRUD operations
```

Press **Ctrl+S** to submit, and watch as the task flows through the agent
hierarchy with **real AI-generated code**!

### 5. Test Multiple Providers (Optional)

```bash
# Compare responses from different LLM providers
go run examples/test_multiple_providers/main.go
```

### 5. Quick Test (Without TUI)

```bash
make test/llm-integration
```

This runs a quick test showing the agent hierarchy in action.

## Configuration

### Main Configuration (`config.yaml`)

The system is driven by a YAML configuration file that defines:

- Organization structure and layers
- Agent assignments and counts
- LLM API keys (via environment variables)
- Slack integration settings (optional)

Example:

```yaml
organization:
  layers:
    - name: President
      agent: ./agents/president.yaml
    - name: Secretary
      count: 1
      attach_to: [President, Director, Manager]
      agent: ./agents/secretary.yaml
    - name: Director
      count: 1
      agent: ./agents/director.yaml
    - name: Manager
      count: 1
      agent: ./agents/manager.yaml
    - name: Engineer
      count: 2
      agent: ./agents/engineer.yaml

slack:
  enabled: false
  token: { env: SLACK_TOKEN }
  channels: ["#alerts", "#progress"]
  notify_on: ["task_assigned", "task_completed", "error"]

llms:
  default_model: gemini
  api_keys:
    gemini: { env: GEMINI_API_KEY }
    claude: { env: CLAUDE_API_KEY }
    codex: { env: CODEX_API_KEY }
    qwen: { env: QWEN_API_KEY }
```

### Agent Configuration

Each agent type has its own YAML configuration file in the `agents/` directory:

- `president.yaml` - President agent configuration
- `secretary.yaml` - Secretary agent configuration
- `director.yaml` - Director agent configuration
- `manager.yaml` - Manager agent configuration
- `engineer.yaml` - Engineer agent configuration

### Environment Variables

Create a `.env` file or set environment variables for LLM API keys:

```bash
export GEMINI_API_KEY="your-gemini-api-key"
export CLAUDE_API_KEY="your-claude-api-key"
export CODEX_API_KEY="your-codex-api-key"
export QWEN_API_KEY="your-qwen-api-key"
export SLACK_TOKEN="your-slack-token"  # Optional
```

**Note**: The system will validate that required environment variables are set
when loading the configuration. If you want to run in demo mode without actual
LLM integration, you can set placeholder values.

## Usage

### Running the Application

```bash
# Run with default config.yaml
./buildbureau

# Or specify a custom config path
export BUILDBUREAU_CONFIG=/path/to/your/config.yaml
./buildbureau
```

### Using the TUI

1. The application starts with a terminal-based UI
2. Enter your task or instruction in the input box
3. Press `Ctrl+S` to submit the task to the President agent
4. Watch as the task flows through the agent hierarchy
5. Press `Ctrl+C` or `Esc` to quit

### Example Tasks

Try these sample instructions:

- "Create a REST API for a todo list application"
- "Design a microservices architecture for an e-commerce platform"
- "Implement a user authentication system with JWT"

## 🧠 Agent Memory System

BuildBureau agents have persistent memory that enables them to **learn from
experience** and **improve over time**.

### Key Features

- **🎓 Learning**: Agents remember past implementations and use them as context
- **📊 Smart Delegation**: Secretaries track performance and route tasks
  intelligently
- **🔍 Semantic Search**: Find similar past tasks using vector search (optional)
- **💾 Fast Storage**: Optimized `mattn/go-sqlite3` driver for performance
- **♻️ Automatic Cleanup**: Configurable retention policies

### How It Works

#### Engineers Learn from Code

```
Task 1: "Create REST API"
→ Generate fresh solution
→ Store code as knowledge

Task 2: "Add authentication to REST API"
→ Find past REST API implementation
→ Include as context for LLM
→ Generate improved solution with authentication
→ Store as knowledge
```

#### Managers Learn from Designs

```
Task: "Design database schema for users"
→ Search for similar past designs
→ Reference past patterns
→ Create consistent specification
→ Store as design knowledge
```

#### Secretaries Track Performance

```
Task: "Optimize database queries"
→ Check which directors handled similar tasks
→ Director A: 3 successful similar tasks
→ Director B: 1 successful similar task
→ Route to Director A (better track record)
→ Record decision and reasoning
```

### Example

```bash
# Run the memory demo
go run examples/test_agent_memory/main.go

# Output shows agents using memory:
✓ Found 2 related past implementations
✓ Engineer has 5 knowledge entries stored
✓ Manager made 3 delegation decisions
```

**See [Agent Memory Guide](docs/AGENT_MEMORY.md) for detailed documentation.**

## Architecture

### Agent Hierarchy

```
Client
  ↓
President ←→ Secretary
  ↓
Director ←→ Secretary
  ↓
Manager ←→ Secretary
  ↓
Engineer
```

### Communication Flow

1. Client submits task to President via TUI
2. President clarifies requirements and forwards to their Secretary
3. Secretary records the goal and delegates to Director(s)
4. Director performs research and decomposes into tasks for Manager(s)
5. Manager creates specifications and delegates to Engineer(s)
6. Engineer implements code and returns results upstream
7. Results flow back up through the hierarchy to the client

### Technical Stack

- **Language**: Go 1.26.0+
- **UI Framework**: Bubble Tea (TUI)
- **Configuration**: YAML (gopkg.in/yaml.v3)
- **Protocol**: Foundation for gRPC (protobuf definitions included)
- **LLM Integration**: Designed for github.com/google/adk-go

## Project Structure

```
BuildBureau/
├── cmd/
│   └── buildbureau/      # Main application entry point
│       └── main.go
├── internal/
│   ├── agent/            # Agent implementations
│   │   ├── base.go       # Base agent functionality
│   │   ├── president.go  # President agent
│   │   ├── secretary.go  # Secretary agent
│   │   ├── director.go   # Director agent
│   │   ├── manager.go    # Manager agent
│   │   ├── engineer.go   # Engineer agent
│   │   └── organization.go # Organization orchestrator
│   ├── config/           # Configuration loading
│   │   └── loader.go
│   ├── tui/              # Terminal UI
│   │   └── tui.go
│   ├── grpc/             # gRPC server/client (future)
│   ├── llm/              # LLM integration (future)
│   └── slack/            # Slack notifications (future)
├── pkg/
│   ├── protocol/         # gRPC protocol definitions
│   │   └── agent.proto
│   └── types/            # Common types
│       ├── agent.go
│       └── config.go
├── agents/               # Agent YAML configurations
│   ├── president.yaml
│   ├── secretary.yaml
│   ├── director.yaml
│   ├── manager.yaml
│   └── engineer.yaml
├── config.yaml           # Main configuration file
├── docs/                 # Documentation
│   ├── ARCHITECTURE.md          # System architecture
│   ├── REMOTE_AGENTS.md         # Remote agent setup guide
│   └── REAL_IMPLEMENTATION.md   # Real vs stub implementations
├── go.mod
└── README.md
```

## Documentation

### Core Documentation

- **[🎨 Multi-Provider LLM Guide](docs/MULTI_PROVIDER.md)**: Complete guide to
  using Gemini, OpenAI, and Claude
- **[🚀 Real Implementation Guide](docs/REAL_IMPLEMENTATION.md)**: Details on
  real LLM, Slack, and gRPC implementations
- **[🎯 ADK Integration Guide](docs/ADK_INTEGRATION.md)**: Using Google's Agent
  Development Kit with BuildBureau
- **[✅ TODO Implementation Guide](docs/TODO_IMPLEMENTATION.md)**: Complete
  guide to implemented TODO items
- **[Architecture Guide](docs/ARCHITECTURE.md)**: Detailed system architecture,
  data flows, and design patterns
- **[Remote Agents Guide](docs/REMOTE_AGENTS.md)**: Setting up remote LLM
  services for Claude, Codex, and Qwen

### Development & Operations

- **[🚀 GitHub Actions Guide](docs/GITHUB_ACTIONS.md)**: Complete CI/CD workflow
  documentation
- **[🛠️ Makefile Guide](docs/MAKEFILE.md)**: Complete Makefile documentation
- **[🐳 Docker Guide](docs/DOCKER.md)**: Docker setup and deployment
- **[Contributing Guide](CONTRIBUTING.md)**: How to contribute to the project

## Recent Improvements

### 🎨 Multi-Provider LLM Support (NEW!)

BuildBureau now supports multiple LLM providers natively:

- **✅ Google Gemini** - Native SDK integration (`google.golang.org/genai`)
- **✅ OpenAI** - Native SDK integration with GPT-4, GPT-3.5
  (`github.com/sashabaranov/go-openai`)
- **✅ Anthropic Claude** - Native SDK integration with Claude 3.5 Sonnet, Opus,
  Haiku (`github.com/liushuangls/go-anthropic/v2`)
- **🔌 Remote Agent API** - HTTP/gRPC for custom providers

Choose your preferred provider or use multiple simultaneously! See
[Multi-Provider Guide](docs/MULTI_PROVIDER.md) for details.

### ✅ ADK Integration

BuildBureau supports Google's Agent Development Kit (google.golang.org/adk):

- **ADK-Powered Agents**: Engineer, Manager, Director, and President agents
  using ADK framework
- **Structured Configuration**: Uses ADK's llmagent.Config for consistent agent
  setup
- **Gemini Integration**: Via ADK's model abstraction layer
- **Extensible**: Ready for ADK tools, memory, and advanced features

See [ADK Integration Guide](docs/ADK_INTEGRATION.md) for details.

### ✅ All TODO Items Completed

All TODO comments in the codebase have been replaced with real, working
implementations:

- **Remote Agent HTTP API**: Full HTTP client for external LLM services with
  authentication
- **gRPC Server**: Real TCP listener with goroutine-based serving and graceful
  shutdown
- **gRPC Client**: Connection management with pooling, timeouts, and proper
  cleanup

See [TODO_COMPLETED.md](TODO_COMPLETED.md) for detailed information about the
implementations.

## Development

### Building

```bash
go build -o buildbureau ./cmd/buildbureau
```

### Testing

```bash
go test ./...
```

### Adding New Agent Types

1. Create a new agent struct in `internal/agent/`
2. Implement the `types.Agent` interface
3. Add configuration in `config.yaml`
4. Create agent YAML config in `agents/`

## Persistent Memory System

BuildBureau includes a comprehensive memory system for agent knowledge
retention:

### Features

- **SQLite Storage**: Persistent structured data
- **Vald Vector DB**: Semantic similarity search (optional)
- **Memory Types**: Conversation, Task, Knowledge, Decision, Context
- **Full-Text Search**: Query memories by content
- **Tag Organization**: Categorize memories with tags
- **Automatic Expiration**: Configurable retention policies

### Configuration

```yaml
memory:
  enabled: true
  sqlite:
    enabled: true
    path: ./data/buildbureau.db
    in_memory: false
  vald:
    enabled: false # Optional vector search
    host: localhost
    port: 8081
    dimension: 768
  retention:
    conversation_days: 30
    task_days: 60
    knowledge_days: 0 # Forever
```

### Test Memory System

```bash
go run examples/test_memory/main.go
```

See [docs/MEMORY_SYSTEM.md](docs/MEMORY_SYSTEM.md) for complete documentation.

## Future Enhancements

- [x] ~~Full gRPC implementation for agent communication~~ ✅
- [x] ~~ADK-go integration for actual LLM interactions~~ ✅
- [x] ~~Remote Agent API support for Claude, Codex, and Qwen~~ ✅
- [x] ~~Slack notification implementation~~ ✅
- [x] ~~Persistent task storage and history~~ ✅ (Memory System)
- [ ] Web UI in addition to TUI
- [ ] Multi-threaded task processing
- [ ] Agent performance metrics and monitoring
- [ ] Automatic embedding generation for semantic search
- [ ] Memory importance scoring and summarization

## License

See LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 🔄 Self-Hosting / Bootstrap Mode

BuildBureau can **build itself** using its own multi-agent system - a form of
recursive self-improvement.

### Quick Start

```bash
# Run BuildBureau in bootstrap mode
make bootstrap

# Give it a self-improvement task:
"Add a caching layer to reduce redundant LLM API calls"
```

BuildBureau will:

1. Analyze the request with deep understanding of its own architecture
2. Design the implementation following its own patterns
3. Generate code, tests, and documentation
4. Present changes for human review

### What Can Bootstrap Do?

- **Add Features**: New agent types, LLM providers, memory backends
- **Refactor Code**: Improve structure while maintaining functionality
- **Optimize Performance**: Identify and fix bottlenecks
- **Add Tests**: Generate test coverage for existing code
- **Fix Bugs**: Understand and resolve issues

### See Also

- [Bootstrap README](bootstrap/README.md) - Complete guide
- [Example Bootstrap Task](bootstrap/tasks/example-bootstrap-task.md) - Detailed
  walkthrough
- [Task Templates](bootstrap/tasks/) - Templates for common improvements

**Status**: Experimental - Demonstrates recursive self-improvement capabilities
