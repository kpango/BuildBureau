# AGENTS.md - BuildBureau AI Agent Integration Guide

> **Purpose**: This document provides a comprehensive guide for AI agents
> (Claude, GPT-4, Copilot, etc.) working on the BuildBureau codebase. It
> contains system architecture, implementation details, development workflows,
> and integration guidelines for future AI-assisted development.

**Last Updated**: 2026-02-01  
**Version**: 1.0.0  
**Target Audience**: AI Coding Assistants and Future Contributors

---

## 📋 Table of Contents

1. [System Overview](#system-overview)
2. [Architecture](#architecture)
3. [Agent Hierarchy](#agent-hierarchy)
4. [Technical Stack](#technical-stack)
5. [Code Structure](#code-structure)
6. [Key Components](#key-components)
7. [Development Workflows](#development-workflows)
8. [Testing Strategy](#testing-strategy)
9. [AI Integration Guidelines](#ai-integration-guidelines)
10. [Common Tasks](#common-tasks)
11. [Best Practices](#best-practices)
12. [Troubleshooting](#troubleshooting)

---

## 🎯 System Overview

### What is BuildBureau?

BuildBureau is a **production-ready multi-agent system** that simulates a
virtual software development company. It features:

- **5-Layer Hierarchical Organization**: President → Secretary → Director →
  Manager → Engineer
- **Real LLM Integration**: Native support for Google Gemini, OpenAI GPT, and
  Anthropic Claude
- **Persistent Memory**: SQLite-based storage with optional Vald vector search
  for semantic similarity
- **Agent Learning**: Agents learn from past tasks and improve over time
- **Production Quality**: ~7,000 lines of Go code, comprehensive testing, CI/CD
  pipelines

### Project Maturity

- ✅ **Production Ready**: All core features fully implemented
- ✅ **Well Tested**: 38+ unit tests, all passing
- ✅ **Well Documented**: 140KB+ of technical documentation
- ✅ **CI/CD Ready**: GitHub Actions workflows, Docker support, Makefile
  automation
- ✅ **Security Hardened**: CodeQL scanning, dependency reviews, non-root
  containers

### Key Statistics

| Metric            | Value                    |
| ----------------- | ------------------------ |
| **Go Files**      | 41 files                 |
| **Lines of Code** | ~7,000 LOC               |
| **Test Coverage** | 38+ tests                |
| **Documentation** | 140KB+ (9 docs)          |
| **Dependencies**  | 20+ production libraries |
| **Docker Image**  | ~50MB (optimized)        |
| **Build Time**    | 2-3 minutes              |

---

## 🏗️ Architecture

### High-Level Design

```
┌─────────────────────────────────────────────────────────────┐
│                         Client / TUI                         │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                        President Agent                       │
│  (Clarifies requirements, defines high-level objectives)    │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                     President's Secretary                    │
│  (Records decisions, forwards to Directors)                 │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                      Director's Secretary                    │
│  (Research, task expansion, forwards to Director)           │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                        Director Agent                        │
│  (Decomposes into department tasks, delegates to Managers)  │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                      Manager's Secretary                     │
│  (Specification finalization, forwards to Manager)          │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                        Manager Agent                         │
│  (Designs software, creates specifications with LLM)        │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                       Engineer Agent                         │
│  (Implements code with LLM, learns from past solutions)     │
└─────────────────────────────────────────────────────────────┘
```

### Supporting Systems

```
┌──────────────────┐    ┌──────────────────┐    ┌──────────────────┐
│   LLM Providers  │    │  Memory System   │    │  Communication   │
├──────────────────┤    ├──────────────────┤    ├──────────────────┤
│ • Gemini         │    │ • SQLite Store   │    │ • gRPC Server    │
│ • OpenAI         │◄───┤ • Vald Vector DB │◄───┤ • gRPC Client    │
│ • Claude         │    │ • Agent Memory   │    │ • HTTP Client    │
│ • Remote API     │    │ • Conversation   │    │ • Slack Notifier │
└──────────────────┘    └──────────────────┘    └──────────────────┘
```

---

## 👥 Agent Hierarchy

### Agent Types and Responsibilities

#### 1. President Agent (`internal/agent/president.go`)

**Role**: Top-level decision maker  
**Responsibilities**:

- Clarifies client instructions
- Defines high-level requirements
- Summarizes objectives
- Delegates to Secretary

**Memory Usage**: Stores clarifications and requirements

#### 2. Secretary Agent (`internal/agent/secretary.go`)

**Role**: Delegation and coordination  
**Responsibilities**:

- Records decisions and goals
- Performs task routing with memory-informed decisions
- Tracks director/manager performance
- Schedules and monitors tasks
- Manages knowledge bases

**Memory Usage**:

- Stores delegation decisions with reasoning
- Tracks which directors handled which tasks
- Uses past performance for smart routing
- Maintains conversation history with leaders

**Smart Features**:

- Round-robin load balancing with performance tracking
- Automatic director selection based on past success

#### 3. Director Agent (`internal/agent/director.go`)

**Role**: Project decomposition  
**Responsibilities**:

- Performs research
- Expands vague goals into actionable units
- Decomposes projects into department-level tasks
- Delegates to Managers based on skills

**Memory Usage**: Stores research and decomposition decisions

#### 4. Manager Agent (`internal/agent/manager.go`)

**Role**: Software design and specification  
**Responsibilities**:

- Uses LLM to produce detailed software designs
- Creates pseudocode and specifications
- Learns from past architectural decisions
- Delegates implementation to Engineers

**LLM Integration**:

- Uses `llmManager.Generate()` for design creation
- Retrieves similar past designs from memory
- Includes past patterns as context for LLM

**Memory Usage**:

- Stores design specifications as knowledge
- References past architectural patterns
- Tracks delegation decisions

#### 5. Engineer Agent (`internal/agent/engineer.go`)

**Role**: Code implementation  
**Responsibilities**:

- Uses LLM to implement code according to specs
- Learns from past implementations
- Returns working code upstream

**LLM Integration**:

- Uses `llmManager.Generate()` for code generation
- Searches memory for similar past implementations
- Includes past solutions as context for better results

**Memory Usage**:

- Stores all conversation and task interactions
- Saves generated code as knowledge
- Retrieves similar past implementations
- Each implementation builds on previous ones

### Agent Implementation Pattern

All agents follow this pattern:

```go
type BaseAgent struct {
    id            string
    agentType     types.AgentType
    config        types.AgentConfig
    status        types.AgentStatus
    activeTask    *types.Task
    completedTasks int
    memory        *AgentMemory  // Memory integration
}

// Core methods all agents must implement:
func (a *BaseAgent) GetID() string
func (a *BaseAgent) GetType() types.AgentType
func (a *BaseAgent) Start(ctx context.Context) error
func (a *BaseAgent) Stop(ctx context.Context) error
func (a *BaseAgent) ProcessTask(ctx context.Context, task types.Task) (*types.TaskResponse, error)
func (a *BaseAgent) GetStatus(ctx context.Context) (string, int, int)

// Memory methods:
func (a *BaseAgent) SetMemoryManager(manager *memory.Manager)
func (a *BaseAgent) GetMemory() *AgentMemory
```

---

## 🔧 Technical Stack

### Core Technologies

#### Language & Runtime

- **Go 1.25.6+**: Main programming language
- **CGo**: Required for SQLite integration
- **Proto**: gRPC protocol buffers

#### LLM Integration

1. **Google Gemini** (`google.golang.org/genai` v1.43.0)
   - Native Go SDK
   - Fast and cost-effective
   - Default model: `gemini-2.0-flash-exp`

2. **OpenAI** (`github.com/sashabaranov/go-openai` v1.41.2)
   - GPT-4 Turbo, GPT-3.5
   - Default model: `gpt-4-turbo-preview`

3. **Anthropic Claude** (`github.com/liushuangls/go-anthropic/v2` v2.17.0)
   - Claude 3.5 Sonnet, Opus, Haiku
   - Default model: `claude-3-5-sonnet-20241022`

4. **Google ADK** (`google.golang.org/adk` v0.3.0)
   - Agent Development Kit integration
   - Structured agent configuration

#### Memory System

- **SQLite** (`github.com/mattn/go-sqlite3` v1.14.33)
  - Fast, standard CGo-based driver
  - Persistent structured storage
  - Full-text search
  - ACID compliance

- **Vald** (`github.com/vdaas/vald-client-go` v1.7.17)
  - Optional vector database
  - Semantic similarity search
  - gRPC-based

#### Communication

- **gRPC** (`google.golang.org/grpc` v1.71.0)
  - Agent-to-agent communication
  - Protocol buffers for type safety

- **Slack** (`github.com/slack-go/slack` v0.17.3)
  - Real-time notifications
  - Channel messaging

#### UI

- **Bubble Tea** (`github.com/charmbracelet/bubbletea` v1.2.4)
  - Terminal UI framework
  - Interactive interface

#### Configuration

- **YAML** (`gopkg.in/yaml.v3` v3.0.1)
  - Declarative configuration
  - Environment variable support

### Build Tools

- **Makefile**: 70+ targets for build, test, lint, docker, formatting, etc.
- **Docker**: Multi-stage builds, ~50MB final image
- **GitHub Actions**: 10 workflows for CI/CD
- **protoc**: Protocol buffer compilation

---

## 📁 Code Structure

### Directory Layout

```
BuildBureau/
├── cmd/
│   └── buildbureau/
│       └── main.go                 # Application entry point
├── internal/
│   ├── agent/                      # Agent implementations
│   │   ├── base.go                 # Base agent with memory
│   │   ├── president.go            # President agent
│   │   ├── secretary.go            # Secretary with smart delegation
│   │   ├── director.go             # Director agent
│   │   ├── manager.go              # Manager with LLM design
│   │   ├── engineer.go             # Engineer with LLM coding
│   │   ├── adk_agent.go            # ADK-powered agents
│   │   ├── memory.go               # Agent memory wrapper
│   │   ├── organization.go         # Orchestration
│   │   └── *_test.go               # Tests
│   ├── config/                     # Configuration
│   │   ├── loader.go               # YAML config loader
│   │   └── loader_test.go          # Tests
│   ├── grpc/                       # gRPC communication
│   │   ├── server.go               # gRPC server implementation
│   │   ├── client.go               # gRPC client
│   │   ├── conversion.go           # Proto ↔ Types conversion
│   │   └── *_test.go               # Tests
│   ├── llm/                        # LLM providers
│   │   ├── manager.go              # LLM manager
│   │   ├── providers.go            # Gemini, OpenAI, Claude, Remote
│   │   └── *_test.go               # Tests
│   ├── memory/                     # Memory system
│   │   ├── manager.go              # Memory orchestration
│   │   ├── sqlite_store.go         # SQLite implementation
│   │   ├── vald_store.go           # Vald vector DB
│   │   └── *_test.go               # Tests
│   ├── slack/                      # Slack integration
│   │   └── notifier.go             # Slack notifications
│   └── tui/                        # Terminal UI
│       └── tui.go                  # Bubble Tea interface
├── pkg/
│   ├── protocol/                   # gRPC protocols
│   │   ├── agent.proto             # Protocol definition
│   │   ├── agent.pb.go             # Generated code
│   │   └── agent_grpc.pb.go        # Generated gRPC service
│   └── types/                      # Type definitions
│       ├── agent.go                # Agent types
│       ├── config.go               # Config types
│       └── memory.go               # Memory types
├── agents/                         # Agent config files
│   ├── president.yaml
│   ├── secretary.yaml
│   ├── director.yaml
│   ├── manager.yaml
│   └── engineer.yaml
├── examples/                       # Example programs
│   ├── test_adk_agents.go
│   ├── test_memory.go
│   ├── test_multiple_providers.go
│   └── ...
├── docs/                           # Documentation
│   ├── ARCHITECTURE.md
│   ├── AGENT_MEMORY.md
│   ├── MEMORY_SYSTEM.md
│   ├── MULTI_PROVIDER.md
│   ├── ADK_INTEGRATION.md
│   ├── DOCKER.md
│   ├── GITHUB_ACTIONS.md
│   ├── MAKEFILE.md
│   └── REMOTE_AGENTS.md
├── .github/                        # GitHub Actions
│   ├── workflows/
│   │   ├── ci.yml                  # Main CI pipeline
│   │   ├── release.yml             # Release automation
│   │   ├── docker-publish.yml      # Docker publishing
│   │   ├── codeql-analysis.yml     # Security scanning
│   │   └── ...                     # 10 total workflows
│   ├── ISSUE_TEMPLATE/
│   ├── dependabot.yml
│   └── ...
├── Dockerfile                      # Multi-stage Docker build
├── docker-compose.yml              # Compose configuration
├── Makefile                        # Build automation (70+ targets)
├── config.yaml                     # Main configuration
├── go.mod                          # Go dependencies
├── README.md                       # Main documentation
├── CONTRIBUTING.md                 # Contribution guide
└── AGENTS.md                       # This file
```

### Key Files

#### Entry Points

- `cmd/buildbureau/main.go`: Application entry point, initializes all systems

#### Configuration

- `config.yaml`: Main YAML configuration
- `internal/config/loader.go`: Configuration parser with env var support
- `agents/*.yaml`: Individual agent configurations

#### Agent System

- `internal/agent/base.go`: Base agent implementation with memory
- `internal/agent/organization.go`: Agent orchestration and lifecycle
- `pkg/types/agent.go`: Agent type definitions and interfaces

#### Memory System

- `internal/memory/manager.go`: Memory manager coordinating SQLite + Vald
- `internal/memory/sqlite_store.go`: Persistent storage implementation
- `internal/memory/vald_store.go`: Vector search implementation
- `pkg/types/memory.go`: Memory type definitions

#### LLM Integration

- `internal/llm/manager.go`: Multi-provider LLM manager
- `internal/llm/providers.go`: Provider implementations (Gemini, OpenAI, Claude,
  Remote)

#### Communication

- `pkg/protocol/agent.proto`: gRPC protocol definition
- `internal/grpc/server.go`: gRPC server
- `internal/grpc/client.go`: gRPC client

---

## 🔑 Key Components

### 1. Agent Memory System

Each agent has access to a memory wrapper that provides:

```go
type AgentMemory struct {
    manager  *memory.Manager
    agentID  string
}

// Storage operations
func (am *AgentMemory) StoreConversation(ctx, content, tags)
func (am *AgentMemory) StoreTask(ctx, task, result, tags)
func (am *AgentMemory) StoreKnowledge(ctx, content, tags)
func (am *AgentMemory) StoreDecision(ctx, decision, reasoning, tags)

// Retrieval operations
func (am *AgentMemory) GetConversationHistory(ctx, limit)
func (am *AgentMemory) GetRelatedTasks(ctx, query, limit)
func (am *AgentMemory) GetKnowledge(ctx, query, limit)
func (am *AgentMemory) SearchMemory(ctx, query, limit)
```

**Memory Types**:

1. **Conversation**: Dialogue history
2. **Task**: Task assignments and results
3. **Knowledge**: Learned information
4. **Decision**: Agent reasoning and choices
5. **Context**: Contextual data

**Storage**:

- SQLite: Fast indexed queries, full-text search
- Vald (optional): Semantic similarity search

### 2. LLM Manager

Provides unified interface to multiple LLM providers:

```go
type Manager struct {
    providers     map[string]Provider
    defaultModel  string
}

// Main generation method
func (m *Manager) Generate(ctx, providerName, prompt, opts) (string, error)

// Provider interface
type Provider interface {
    Generate(ctx, prompt, opts) (string, error)
    Close() error
}
```

**Supported Providers**:

- `gemini`: Google Gemini (default)
- `openai`: OpenAI GPT models
- `claude`: Anthropic Claude models
- Custom remote agents via HTTP

### 3. Configuration System

YAML-based with environment variable support:

```yaml
organization:
  layers:
    - name: President
      agent: ./agents/president.yaml
    # ... more layers

llms:
  default_model: gemini
  api_keys:
    gemini: { env: GEMINI_API_KEY }
    openai: { env: OPENAI_API_KEY }
    claude: { env: CLAUDE_API_KEY }

memory:
  enabled: true
  sqlite:
    enabled: true
    path: ./data/buildbureau.db
  vald:
    enabled: false
    host: localhost
    port: 8081

slack:
  enabled: false
  token: { env: SLACK_TOKEN }
  channels: ["#alerts"]
```

### 4. Task Flow

```go
type Task struct {
    ID          string
    Type        TaskType
    Description string
    Priority    int
    Metadata    map[string]interface{}
    CreatedAt   time.Time
}

type TaskResponse struct {
    TaskID      string
    Status      TaskStatus
    Result      string
    Error       string
    CompletedAt time.Time
}
```

**Flow**:

1. Client submits task to President
2. President clarifies → Secretary records
3. Secretary delegates to Directors
4. Director decomposes → Manager designs
5. Manager specs → Engineer implements
6. Results propagate back up

---

## 🛠️ Development Workflows

### Build System (Makefile)

The Makefile provides 70+ targets organized into categories:

```bash
# Build targets
make build              # Build standard binary
make build-release      # Optimized release build
make build-debug        # Debug build with symbols
make build-static       # Static binary
make build-all          # Multi-platform builds

# Test targets
make test               # Run all tests
make test-unit          # Unit tests only
make test-coverage      # With coverage report
make test-coverage-html # HTML coverage report
make test-bench         # Benchmarks
make test-race          # Race detection
make test/llm-integration # Real LLM integration testing

# Formatting targets
make format             # Format all files (Go, YAML, JSON, Markdown)
make format/go          # Format Go files only
make format/yaml        # Format YAML files only
make format/json        # Format JSON files only
make format/md          # Format Markdown files only
make format-check       # Check formatting without modifying (CI mode)

# Docker targets
make docker-build       # Build Docker image (replaces docker/build.sh)
make docker-run         # Run container in daemon mode (replaces docker/run.sh)
make docker-run-interactive # Run container interactively
make docker-test        # Test Docker setup (replaces docker/test.sh)
make docker-push        # Push to registry
make docker-compose-up  # Start with compose

# CI/CD targets
make ci-all            # Run all CI checks
make ci-lint           # Linting
make ci-build          # CI build
make ci-test           # CI tests

# Auto-install targets
make install-all       # Install all development tools
make install-tools     # Install Go tools (golangci-lint, protoc-gen-go)
make install-formatters # Install formatters (gofmt, yamlfmt, jq, prettier)
make install-security-tools # Install security tools (gosec, nancy)
make clean-stamps      # Remove stamp files to force reinstall

# Development targets
make proto             # Generate proto code
make deps              # Download dependencies
make fmt               # Format code (alias for format)
make lint              # Run linters
make clean             # Clean artifacts

# Utility targets
make help              # Show all targets
make version           # Show version info
make check             # Check dependencies
```

### Docker Workflow

```bash
# Build and run with Docker
docker build -t buildbureau .
docker run -e GEMINI_API_KEY="key" buildbureau

# Or use Docker Compose
docker-compose up -d
docker-compose logs -f
docker-compose down

# Multi-arch builds
docker buildx build --platform linux/amd64,linux/arm64 -t buildbureau .
```

### Git Workflow

```bash
# Feature development
git checkout -b feature/new-agent
# ... make changes ...
git commit -m "Add new agent type"
git push origin feature/new-agent
# Create PR, CI runs automatically

# Release
git tag v1.0.0
git push origin v1.0.0
# Release workflow builds binaries, Docker images, changelog
```

### Formatting System

The project includes a comprehensive formatting system that supports multiple
file types:

#### Formatting Commands

```bash
# Format all files in the project
make format

# Format specific file types
make format/go          # Format Go files with gofmt
make format/yaml        # Format YAML files with yamlfmt
make format/json        # Format JSON files with jq
make format/md          # Format Markdown files with prettier

# Check formatting without modifying (useful in CI)
make format-check       # Returns error if files need formatting
```

#### How It Works

The formatting system:

1. **Auto-detects** files that need formatting
2. **Installs** formatters automatically on first use (via stamp files)
3. **Formats** files in-place with proper configuration
4. **Validates** in CI mode without modifying files

#### Supported Formatters

| File Type    | Tool       | Description                               |
| ------------ | ---------- | ----------------------------------------- |
| **Go**       | `gofmt`    | Standard Go formatter                     |
| **YAML**     | `yamlfmt`  | YAML formatting with proper indentation   |
| **JSON**     | `jq`       | JSON pretty-printing and validation       |
| **Markdown** | `prettier` | Markdown formatting with consistent style |

#### Usage Examples

```bash
# Before committing, format all files
make format

# In CI, check if formatting is needed
make format-check
# Exit code 0: All files formatted correctly
# Exit code 1: Some files need formatting

# Format only Go files after code changes
make format/go

# Format YAML configs after editing
make format/yaml
```

### Auto-Install System

BuildBureau includes an intelligent auto-install system that manages development
tools automatically.

#### How Stamp Files Work

The Makefile uses **stamp files** in `.make/*.stamp` to track installed tools:

```
.make/
├── tools.stamp           # Marks Go tools installed
├── formatters.stamp      # Marks formatters installed
└── security-tools.stamp  # Marks security tools installed
```

**Benefits**:

- ✅ Tools install automatically when needed
- ✅ No duplicate installations (checks stamps first)
- ✅ Faster builds (skips already-installed tools)
- ✅ Easy to force reinstall (delete stamps)

#### Installation Commands

```bash
# Install all development tools at once
make install-all

# Install specific tool groups
make install-tools          # Go tools: golangci-lint, protoc-gen-go, etc.
make install-formatters     # Formatters: yamlfmt, prettier, jq
make install-security-tools # Security: gosec, nancy

# Force reinstall (delete stamps first)
make clean-stamps
make install-all
```

#### What Gets Installed

**Go Tools** (`make install-tools`):

- `golangci-lint` - Comprehensive Go linter
- `protoc-gen-go` - Protocol buffer Go code generator
- `protoc-gen-go-grpc` - gRPC Go code generator

**Formatters** (`make install-formatters`):

- `gofmt` - Go code formatter (included with Go)
- `yamlfmt` - YAML file formatter
- `jq` - JSON processor and formatter
- `prettier` - Markdown and other file formatter (via npm)

**Security Tools** (`make install-security-tools`):

- `gosec` - Go security checker
- `nancy` - Dependency vulnerability scanner

#### Auto-Install in Action

When you run a command that needs a tool:

```bash
# Run linting
make lint

# If golangci-lint not installed:
# 1. Makefile detects missing .make/tools.stamp
# 2. Automatically runs: make install-tools
# 3. Creates stamp file: .make/tools.stamp
# 4. Proceeds with linting

# Next time:
# 1. Stamp file exists
# 2. Skips installation
# 3. Runs linting immediately
```

#### Troubleshooting Auto-Install

```bash
# Tool installation failed?
# Clean stamps and retry
make clean-stamps
make install-all

# Check if tools are properly installed
which golangci-lint
which yamlfmt
which gosec

# Manual installation
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Enhanced Docker Workflow

The Docker integration has been enhanced with dedicated Makefile targets that
replace previous shell scripts.

#### Docker Commands

```bash
# Build Docker image
make docker-build
# Replaces: docker/build.sh
# Builds multi-stage optimized image (~50MB)

# Run in daemon mode
make docker-run
# Replaces: docker/run.sh
# Starts container in background with proper environment

# Run interactively
make docker-run-interactive
# Runs container with terminal attached
# Useful for debugging and development

# Test Docker setup
make docker-test
# Replaces: docker/test.sh
# Validates Docker build and runtime functionality

# Push to registry
make docker-push
# Pushes built image to container registry

# Docker Compose
make docker-compose-up
# Starts full stack with docker-compose
```

#### Docker Build Features

The Docker build includes:

- **Multi-stage builds** for minimal image size
- **Non-root user** for security
- **Build caching** for faster rebuilds
- **Multi-arch support** (amd64, arm64)

#### Usage Examples

```bash
# Development: Build and test locally
make docker-build
make docker-test
make docker-run-interactive

# Production: Build and deploy
make docker-build
make docker-push

# Full stack: Use Docker Compose
make docker-compose-up
docker-compose logs -f buildbureau

# Multi-architecture builds
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t buildbureau:latest \
  .
```

#### Docker Configuration

The Docker setup uses:

- **Dockerfile**: Multi-stage production build
- **docker-compose.yml**: Full stack configuration
- **Environment variables**: API keys and configuration

Example `docker-compose.yml` integration:

```yaml
services:
  buildbureau:
    build: .
    environment:
      - GEMINI_API_KEY=${GEMINI_API_KEY}
      - OPENAI_API_KEY=${OPENAI_API_KEY}
    volumes:
      - ./data:/app/data
    ports:
      - "8080:8080"
```

---

## 🧪 Testing Strategy

### Test Structure

Tests are colocated with implementation:

```
internal/
├── agent/
│   ├── base.go
│   ├── base_test.go        # Unit tests
│   ├── engineer.go
│   └── engineer_test.go
├── memory/
│   ├── manager.go
│   └── memory_test.go
└── ...
```

### Running Tests

```bash
# All tests
make test

# Specific package
go test ./internal/agent -v

# With coverage
make test-coverage

# With race detection
make test-race

# Benchmarks
make test-bench

# LLM integration testing with real API calls
make test/llm-integration
# Replaces: test_real_llm.sh
# Tests real LLM providers (requires API keys)
```

#### LLM Integration Testing

The `make test/llm-integration` target performs end-to-end testing with real LLM
providers:

**What it tests**:

- Real API calls to Gemini, OpenAI, and Claude
- Multi-provider functionality
- LLM manager orchestration
- Error handling and retries
- Rate limiting and timeouts

**Requirements**:

```bash
# Set API keys for providers you want to test
export GEMINI_API_KEY="your-gemini-key"
export OPENAI_API_KEY="your-openai-key"
export CLAUDE_API_KEY="your-claude-key"

# Run integration tests
make test/llm-integration
```

**Output**:

```
Testing Gemini provider...
✓ Gemini: Response generated successfully
Testing OpenAI provider...
✓ OpenAI: Response generated successfully
Testing Claude provider...
✓ Claude: Response generated successfully
Testing multi-provider fallback...
✓ Fallback: Primary failed, secondary succeeded
```

**Note**: This test makes real API calls and may incur costs. Use carefully.

### Test Examples

```go
func TestEngineerAgent(t *testing.T) {
    config := types.AgentConfig{
        Name:        "test-engineer",
        Description: "Test engineer",
    }

    agent := NewEngineerAgent("eng-001", config, nil)

    ctx := context.Background()
    err := agent.Start(ctx)
    assert.NoError(t, err)

    // Test task processing
    task := types.Task{
        ID:          "task-001",
        Description: "Implement feature",
    }

    response, err := agent.ProcessTask(ctx, task)
    assert.NoError(t, err)
    assert.NotNil(t, response)
}
```

### Test Coverage

Current coverage: ~90% for core packages

- `internal/agent`: High coverage
- `internal/memory`: High coverage
- `internal/llm`: High coverage
- `internal/config`: High coverage

---

## 🤖 AI Integration Guidelines

### For AI Coding Assistants (You!)

When working on this codebase, follow these guidelines:

#### 1. Understanding the System

**Key Concepts to Grasp**:

- **Agent Hierarchy**: 5 layers with specific responsibilities
- **Memory System**: Agents learn from past interactions
- **LLM Integration**: Multiple providers, unified interface
- **Task Flow**: Bottom-up implementation, top-down delegation

**Read First**:

1. `README.md` - Overview and quick start
2. `docs/ARCHITECTURE.md` - System design
3. `docs/AGENT_MEMORY.md` - Memory integration
4. This file (`AGENTS.md`) - AI integration guide

#### 2. Making Changes

**Before Coding**:

1. ✅ Understand the agent hierarchy
2. ✅ Check existing patterns in similar code
3. ✅ Review type definitions in `pkg/types/`
4. ✅ Check if memory integration is needed
5. ✅ Consider LLM provider implications

**Code Patterns to Follow**:

**Agent Implementation**:

```go
// All agents extend BaseAgent
type MyNewAgent struct {
    *BaseAgent
    // Additional fields
}

func NewMyNewAgent(id string, config types.AgentConfig, llmManager *llm.Manager) *MyNewAgent {
    agent := &MyNewAgent{
        BaseAgent: &BaseAgent{
            id:        id,
            agentType: types.AgentTypeCustom,
            config:    config,
        },
    }
    return agent
}

// Implement required methods
func (a *MyNewAgent) ProcessTask(ctx context.Context, task types.Task) (*types.TaskResponse, error) {
    // 1. Store conversation in memory
    if a.memory != nil {
        a.memory.StoreConversation(ctx, task.Description, []string{"task"})
    }

    // 2. Use LLM if needed
    if a.llmManager != nil {
        response, err := a.llmManager.Generate(ctx, "gemini", task.Description, &llm.GenerateOptions{})
        if err != nil {
            return nil, err
        }
        // Use response
    }

    // 3. Store result in memory
    if a.memory != nil {
        a.memory.StoreTask(ctx, task.Description, "result", []string{"completed"})
    }

    return &types.TaskResponse{
        TaskID: task.ID,
        Status: types.StatusCompleted,
        Result: "result",
    }, nil
}
```

**Memory Integration**:

```go
// Always check if memory is available
if agent.memory != nil {
    // Store interactions
    agent.memory.StoreConversation(ctx, content, tags)

    // Retrieve past experiences
    history, _ := agent.memory.GetConversationHistory(ctx, 10)

    // Search for related tasks
    related, _ := agent.memory.GetRelatedTasks(ctx, query, 5)
}
```

**LLM Usage**:

```go
// Use LLM manager for generation
if agent.llmManager != nil {
    opts := &llm.GenerateOptions{
        Temperature: 0.7,
        MaxTokens:   2000,
    }

    // Include context from memory
    prompt := buildPromptWithContext(task, history)

    response, err := agent.llmManager.Generate(ctx, "gemini", prompt, opts)
    if err != nil {
        // Handle error
    }
}
```

#### 3. Testing Changes

**Always**:

1. ✅ Write tests for new code
2. ✅ Run `make test` before committing
3. ✅ Check `make lint` passes
4. ✅ Verify `make build` succeeds
5. ✅ Test with real API keys if LLM changes

**Test Pattern**:

```go
func TestNewFeature(t *testing.T) {
    // Setup
    ctx := context.Background()

    // Create test fixtures
    // ...

    // Execute
    result, err := functionUnderTest(ctx, input)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, expected, result)
}
```

#### 4. Documentation

**Update When Changing**:

- `README.md`: If adding major features
- `docs/ARCHITECTURE.md`: If changing system design
- `docs/AGENT_MEMORY.md`: If changing memory system
- `docs/MULTI_PROVIDER.md`: If adding LLM providers
- Code comments: For complex logic
- This file: For AI integration changes

#### 5. Commit Messages

Follow this format:

```
<type>: <short description>

<detailed description>

<breaking changes if any>
```

Types: `feat`, `fix`, `docs`, `test`, `refactor`, `perf`, `chore`

Examples:

```
feat: Add support for Llama 3 provider

Implements LlamaProvider with ollama backend integration.
Includes tests and documentation updates.

fix: Resolve memory leak in agent cleanup

Properly close LLM connections and memory stores on agent shutdown.

docs: Update AGENTS.md with new patterns
```

---

## 📝 Common Tasks

### Task 1: Add a New Agent Type

1. **Define the agent** in `internal/agent/`:

```go
// newagent.go
type NewAgent struct {
    *BaseAgent
    // Custom fields
}

func NewNewAgent(id string, config types.AgentConfig, llmManager *llm.Manager) *NewAgent {
    // Implementation
}

func (a *NewAgent) ProcessTask(ctx context.Context, task types.Task) (*types.TaskResponse, error) {
    // Implementation with memory and LLM
}
```

2. **Add agent type** in `pkg/types/agent.go`:

```go
const (
    AgentTypePresident  AgentType = "president"
    AgentTypeSecretary  AgentType = "secretary"
    // ...
    AgentTypeNew        AgentType = "new"  // Add this
)
```

3. **Create config** in `agents/new.yaml`:

```yaml
name: "NewAgent"
description: "Description of new agent"
prompt: "System prompt for the agent"
```

4. **Write tests** in `internal/agent/newagent_test.go`

5. **Update organization** in `internal/agent/organization.go` to instantiate
   the new agent

6. **Document** in relevant docs

### Task 2: Add a New LLM Provider

1. **Implement Provider interface** in `internal/llm/providers.go`:

```go
type NewLLMProvider struct {
    name     string
    apiKey   string
    model    string
    client   *newllm.Client
}

func NewNewLLMProvider(name, apiKey, model string) (*NewLLMProvider, error) {
    client, err := newllm.NewClient(apiKey)
    if err != nil {
        return nil, err
    }

    return &NewLLMProvider{
        name:   name,
        apiKey: apiKey,
        model:  model,
        client: client,
    }, nil
}

func (p *NewLLMProvider) Generate(ctx context.Context, prompt string, opts *GenerateOptions) (string, error) {
    // Implementation
}

func (p *NewLLMProvider) Close() error {
    // Cleanup
}
```

2. **Update Manager** in `internal/llm/manager.go`:

```go
func NewManager(config types.LLMsConfig) (*Manager, error) {
    // ...

    // Add new provider initialization
    if apiKey := os.Getenv("NEWLLM_API_KEY"); apiKey != "" {
        provider, err := NewNewLLMProvider("newllm", apiKey, config.NewLLMModel)
        if err != nil {
            return nil, err
        }
        providers["newllm"] = provider
    }

    // ...
}
```

3. **Add configuration** in `config.yaml`:

```yaml
llms:
  default_model: gemini
  api_keys:
    newllm: { env: NEWLLM_API_KEY }
  newllm_model: "default-model"
```

4. **Write tests** in `internal/llm/providers_test.go`

5. **Document** in `docs/MULTI_PROVIDER.md`

### Task 3: Add Memory Features

1. **Define new memory type** in `pkg/types/memory.go`:

```go
const (
    MemoryTypeConversation MemoryType = "conversation"
    // ...
    MemoryTypeNew         MemoryType = "new"  // Add this
)
```

2. **Add storage methods** in `internal/memory/sqlite_store.go` if needed

3. **Add retrieval methods** in `internal/memory/manager.go`:

```go
func (m *Manager) GetNewMemories(ctx context.Context, agentID string, limit int) ([]*types.MemoryEntry, error) {
    // Implementation
}
```

4. **Update AgentMemory wrapper** in `internal/agent/memory.go`:

```go
func (am *AgentMemory) GetNewMemories(ctx context.Context, limit int) ([]*types.MemoryEntry, error) {
    return am.manager.GetNewMemories(ctx, am.agentID, limit)
}
```

5. **Write tests**

6. **Document** in `docs/MEMORY_SYSTEM.md`

### Task 4: Modify Configuration

1. **Update types** in `pkg/types/config.go`

2. **Update loader** in `internal/config/loader.go`

3. **Update example config** in `config.yaml`

4. **Update validation** if needed

5. **Update docs** in README and relevant docs

### Task 5: Add GitHub Actions Workflow

1. **Create workflow** in `.github/workflows/newworkflow.yml`

2. **Use Makefile targets**:

```yaml
- name: Run custom check
  run: make custom-check
```

3. **Test locally**:

```bash
# Install act (https://github.com/nektos/act)
act -j job-name
```

4. **Document** in `docs/GITHUB_ACTIONS.md`

---

## ✨ Best Practices

### Code Quality

1. **Follow Go Conventions**
   - Use `gofmt` for formatting
   - Follow effective Go guidelines
   - Use meaningful variable names
   - Keep functions small and focused

2. **Error Handling**
   - Always check errors
   - Wrap errors with context: `fmt.Errorf("failed to X: %w", err)`
   - Log errors appropriately
   - Don't ignore errors in deferred calls

3. **Concurrency**
   - Use context for cancellation
   - Protect shared state with mutexes
   - Close channels when done
   - Handle goroutine lifecycle properly

4. **Memory Management**
   - Close resources in defer statements
   - Use pointers for large structs
   - Avoid memory leaks in long-running goroutines
   - Profile if performance issues arise

### Agent Development

1. **Memory Integration**
   - Always check if memory is available before use
   - Store important interactions
   - Retrieve relevant context before processing
   - Use appropriate tags for categorization

2. **LLM Usage**
   - Include relevant context from memory
   - Use appropriate temperature settings
   - Handle rate limits gracefully
   - Implement retries for transient failures
   - Cache responses when appropriate

3. **Task Processing**
   - Validate input tasks
   - Update status appropriately
   - Return meaningful error messages
   - Store results in memory

### Testing

1. **Unit Tests**
   - Test one thing at a time
   - Use table-driven tests for multiple cases
   - Mock external dependencies
   - Test error cases

2. **Integration Tests**
   - Test component interactions
   - Use real dependencies where practical
   - Clean up resources after tests
   - Mark long tests: `t.Skip("skipping long test")`

3. **Coverage**
   - Aim for >80% coverage for core packages
   - Focus on critical paths
   - Don't test generated code
   - Use coverage reports: `make test-coverage-html`

### Documentation

1. **Code Comments**
   - Document exported functions and types
   - Explain complex logic
   - Keep comments up to date
   - Use godoc format

2. **Documentation Files**
   - Keep README.md current
   - Update technical docs when changing systems
   - Include examples in docs
   - Cross-reference related docs

### Security

1. **API Keys**
   - Never commit API keys
   - Use environment variables
   - Validate keys before use
   - Rotate keys regularly

2. **Input Validation**
   - Validate all user input
   - Sanitize data for storage
   - Check bounds and types
   - Handle edge cases

3. **Dependencies**
   - Keep dependencies updated
   - Run security scans: `make security`
   - Review dependency licenses
   - Use dependabot for updates

---

## 🔍 Troubleshooting

### Common Issues

#### Build Failures

**Issue**: `undefined: mattn`

```
Solution: Install build dependencies
$ make deps
$ make build
```

**Issue**: `protoc: command not found`

```
Solution: Install protoc or skip proto generation
$ make proto  # If protoc installed
$ make build  # Will use existing generated code
```

**Issue**: CGo compilation errors

```
Solution: Install gcc/musl-dev
# Ubuntu/Debian
$ sudo apt-get install build-essential
# Alpine
$ apk add gcc musl-dev
```

#### Test Failures

**Issue**: Tests fail with "API key not set"

```
Solution: Tests skip without API keys (expected)
$ export GEMINI_API_KEY="your-key"
$ make test
```

**Issue**: "main redeclared" errors

```
Solution: Examples directory excluded from tests
# This is already fixed in current Makefile
$ make test  # Should work
```

#### Runtime Issues

**Issue**: "Failed to connect to database"

```
Solution: Check database path and permissions
$ mkdir -p data
$ chmod 755 data
```

**Issue**: "LLM provider not available"

```
Solution: Set API key environment variable
$ export GEMINI_API_KEY="your-key"
$ ./buildbureau
```

**Issue**: "Memory search returns no results"

```
Solution:
1. Check if memory is enabled in config
2. Ensure data exists
3. Check SQLite database:
   $ sqlite3 data/buildbureau.db "SELECT COUNT(*) FROM memories;"
```

### Debug Mode

Enable verbose logging:

```bash
# Set log level
export LOG_LEVEL=debug

# Run with verbose output
./buildbureau -v

# Or in Docker
docker-compose up  # Logs to stdout
```

### Performance Issues

**Slow LLM responses**:

- Check API rate limits
- Consider using faster models (Gemini 2.0 Flash)
- Enable response caching

**High memory usage**:

- Check retention policies in config
- Run memory cleanup: `make clean-coverage`
- Adjust `max_entries` in memory config

**Slow builds**:

- Use build cache: `make build` uses Go cache
- Use Docker cache: `docker build --cache-from`
- Enable parallel compilation: `make -j4 build`

---

## 📚 Additional Resources

### Documentation

- [Architecture](docs/ARCHITECTURE.md) - System design and patterns
- [Agent Memory](docs/AGENT_MEMORY.md) - Memory system integration
- [Multi-Provider](docs/MULTI_PROVIDER.md) - LLM provider guide
- [Docker Guide](docs/DOCKER.md) - Docker deployment
- [Makefile Reference](docs/MAKEFILE.md) - Build system
- [GitHub Actions](docs/GITHUB_ACTIONS.md) - CI/CD workflows

### External Resources

- [Go Documentation](https://golang.org/doc/)
- [Bubble Tea Framework](https://github.com/charmbracelet/bubbletea)
- [gRPC Go](https://grpc.io/docs/languages/go/)
- [Google Gemini](https://ai.google.dev/)
- [OpenAI API](https://platform.openai.com/docs)
- [Anthropic Claude](https://docs.anthropic.com/)

### Community

- GitHub Issues: Report bugs, request features
- Pull Requests: Contribute code
- Discussions: Ask questions, share ideas

---

## 🎯 AI Agent Quick Reference

### When You First Encounter This Project

1. ✅ Read `README.md` for overview
2. ✅ Read this file (`AGENTS.md`) completely
3. ✅ Examine `docs/ARCHITECTURE.md` for design
4. ✅ Look at `pkg/types/` for core types
5. ✅ Review `internal/agent/base.go` for agent pattern
6. ✅ Check `config.yaml` for configuration structure

### Before Making Changes

1. ✅ Understand which layer you're modifying
2. ✅ Check existing patterns in similar code
3. ✅ Consider memory integration needs
4. ✅ Consider LLM provider implications
5. ✅ Plan your tests

### While Making Changes

1. ✅ Follow existing code patterns
2. ✅ Keep changes minimal and focused
3. ✅ Write tests alongside code
4. ✅ Update relevant documentation
5. ✅ Run `make test` frequently

### Before Committing

1. ✅ `make test` - All tests pass
2. ✅ `make lint` - No linter errors
3. ✅ `make build` - Successful build
4. ✅ Update documentation if needed
5. ✅ Write clear commit message

### Key Principles

- **Agents learn from memory** - Always integrate memory
- **Multiple LLM providers** - Use llmManager, not direct API calls
- **Type safety** - Use pkg/types definitions
- **Error handling** - Always check and wrap errors
- **Testing** - Test new code thoroughly
- **Documentation** - Keep docs current

---

## 📝 Changelog

### Version 1.0.0 (2026-02-01)

**Initial Release**

- Complete multi-agent system implementation
- Multi-provider LLM support (Gemini, OpenAI, Claude)
- Persistent memory with SQLite and Vald
- gRPC communication layer
- Slack notifications
- Terminal UI with Bubble Tea
- Docker support with multi-stage builds
- Comprehensive CI/CD with GitHub Actions
- 60+ Makefile targets
- ~7,000 lines of production code
- 38+ unit tests, all passing
- 140KB+ technical documentation

---

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed contribution guidelines.

**Quick Guidelines for AI Agents**:

1. Follow the patterns in this document
2. Write tests for all new code
3. Update documentation when changing features
4. Use clear, descriptive commit messages
5. Respect the existing architecture
6. Ask (via comments) if unsure

---

**Last Updated**: 2026-02-01  
**Maintained By**: BuildBureau Team  
**For AI Agents**: This document is your guide to working with BuildBureau. Read
it thoroughly before making changes.

---

> **Note to Future AI Agents**: This codebase is production-ready and
> well-tested. When making changes, preserve the existing architecture, follow
> established patterns, and maintain backward compatibility unless explicitly
> requested otherwise. The agent hierarchy, memory system, and multi-provider
> LLM support are core to the design - changes to these should be made carefully
> with full understanding of implications.

---

## 🔄 Self-Hosting / Bootstrap System

BuildBureau includes a **bootstrap system** that enables it to build and improve
itself - demonstrating recursive self-improvement.

### Overview

The bootstrap system consists of:

1. **Self-Aware Agents** (`bootstrap/agents/`)
   - Agents with deep knowledge of BuildBureau's codebase
   - Understand architecture, patterns, and conventions
   - Can modify the system they're running in

2. **Bootstrap Configuration** (`bootstrap/config.yaml`)
   - Specialized setup for self-improvement tasks
   - More agents (3 Engineers, 2 Managers) for parallel work
   - Enhanced memory retention for learning

3. **Task Templates** (`bootstrap/tasks/`)
   - Structured templates for common improvements
   - add-feature, refactor, optimize, test

4. **Bootstrap Launcher** (`bootstrap/bootstrap.sh`, `make bootstrap`)
   - Easy startup for self-hosting mode
   - Environment validation
   - Safety reminders

### Usage

```bash
# Start BuildBureau in bootstrap mode
make bootstrap

# Example tasks:
"Add a new agent type for code review"
"Optimize memory query performance"
"Refactor the LLM abstraction layer"
"Add tests for the agent communication system"
```

### How It Works

#### Self-Aware Engineer Agent

````yaml
system_prompt: |
  You are an Engineer in BuildBureau's BOOTSTRAP MODE.
  YOU ARE MODIFYING THE CODEBASE YOU'RE RUNNING IN.

  BUILDBUREAU CODE PATTERNS:
  ```go
  type MyAgent struct {
      *BaseAgent
      llmManager *llm.Manager
  }

  func (a *MyAgent) ProcessTask(ctx context.Context, task types.Task) (*types.TaskResponse, error) {
      // Memory integration
      if a.memory != nil {
          a.memory.StoreConversation(ctx, task.Description, []string{"task"})
      }
      // ... implementation
  }
````

````

#### Implementation Flow

1. **President** clarifies self-improvement request
2. **Secretary** coordinates modification tasks
3. **Director** plans architectural changes
4. **Manager** designs implementation details
5. **Engineer** generates code following BuildBureau patterns
6. **Human review** before applying changes

### Benefits

- **Recursive Improvement**: Each cycle makes future cycles better
- **Pattern Consistency**: Generated code follows established patterns
- **Self-Awareness**: Understands what it's modifying
- **Learning**: Memory accumulates improvement knowledge

### Safety

⚠️ **Important**:
- All generated code requires human review
- Changes go through git (easy rollback)
- Tests are generated with implementations
- Run in isolated environment first

### Example: Add Caching to LLM Manager

```bash
make bootstrap
# Task: "Add caching layer to LLM manager to reduce redundant API calls"

# BuildBureau will:
# 1. Analyze caching requirements
# 2. Design cache structure (LRU, TTL)
# 3. Implement cache in internal/llm/manager.go
# 4. Add tests for cache behavior
# 5. Update configuration and docs
# 6. Present for human review

# Review and test:
git diff
make test

# Approve:
git commit -m "Add LLM caching (self-implemented by BuildBureau)"
````

### Documentation

- [Bootstrap README](../bootstrap/README.md) - Complete guide
- [Example Bootstrap Task](../bootstrap/tasks/example-bootstrap-task.md) -
  Detailed walkthrough
- [Configuration](../bootstrap/config.yaml) - Bootstrap setup
- [Agent Prompts](../bootstrap/agents/) - Self-aware agent configurations

### Philosophy

> "A system that can improve itself is the first step toward truly autonomous
> software development."

BuildBureau's bootstrap mode is an experiment in **recursive
self-improvement** - enabling AI systems to understand and enhance their own
capabilities.

### Limitations

Current limitations:

- Requires human review of all changes
- Complex architectural changes may need guidance
- Can't install new system dependencies
- Best for incremental improvements

Future possibilities:

- Automated testing before presenting changes
- Multi-step planning for complex features
- Continuous self-improvement in CI/CD
- Self-documentation generation

---
