# BuildBureau

Multi-layered AI Agent Implementation System - Hierarchical multi-agent configuration from President to Employee

## Overview

BuildBureau is an AI system with a hierarchical multi-agent configuration: President → Department Manager → Section Manager → Employee. Each level has secretary agents that elaborate on instructions from higher-level agents, recording, assisting, and scheduling tasks.

### Architecture

```
Client
    ↓
President Agent + President Secretary
    ↓
Department Manager Agent + Department Secretary  
    ↓
Section Manager Agent + Section Secretary
    ↓
Employee Agent
```

### Key Features

- **Hierarchical Agent Structure**: 4-layer structure with President, Department Manager, Section Manager, and Employee
- **Secretary Agents**: Secretary agents at each level to support task management
- **gRPC Communication**: Loosely coupled communication between agents via gRPC
- **YAML Configuration**: All settings managed through YAML files
- **Slack Notifications**: Automatic notifications to Slack for important events
- **Terminal UI**: Interactive terminal UI using Bubble Tea
- **Single Binary**: Operates as a single binary implemented in Go

## Tech Stack

- **Language**: Go 1.23+
- **AI Agents**: Google ADK (Agent Development Kit) for Go
- **Communication**: gRPC (Protocol Buffers)
- **UI**: Charmbracelet Bubble Tea
- **Notifications**: Slack API (slack-go)
- **Configuration**: YAML (gopkg.in/yaml.v3)

## Installation

### Prerequisites

- Go 1.23 or higher
- protoc (Protocol Buffers compiler)

### Build

```bash
# Install dependencies
make deps

# Generate protocol buffer code (if needed)
make install-tools
make proto

# Build
make build
```

## Configuration

All settings are managed in the `config.yaml` file.

### Main Configuration Items

#### Agent Configuration

For each agent type, configure:

- `count`: Number of agents
- `model`: LLM model to use
- `instruction`: System prompt for the agent
- `allowTools`: Permission to use tools
- `tools`: List of available tools
- `timeout`: Timeout in seconds
- `retryCount`: Number of retries

```yaml
agents:
  president:
    count: 1
    model: "gemini-2.0-flash-exp"
    instruction: |
      You are the President responsible for overseeing the entire project...
    allowTools: true
    tools:
      - web_search
      - knowledge_base
    timeout: 120
    retryCount: 3
```

#### Slack Notification Settings

```yaml
slack:
  enabled: true
  token: "${SLACK_BOT_TOKEN}"
  channelID: "${SLACK_CHANNEL_ID}"
  notifications:
    projectStart:
      enabled: true
      message: "🚀 Project \"{{.ProjectName}}\" has started"
```

Configure token and channel ID via environment variables:

```bash
export SLACK_BOT_TOKEN="xoxb-your-token"
export SLACK_CHANNEL_ID="C01234567"
```

#### UI Settings

```yaml
ui:
  enableTUI: true
  refreshRate: 100  # milliseconds
  theme: "default"
  showProgress: true
  logLevel: "info"
```

## Usage

### Basic Execution

```bash
# Run with default configuration
./bin/buildbureau

# Specify custom configuration file
CONFIG_PATH=/path/to/config.yaml ./bin/buildbureau
```

### Terminal UI

When TUI is enabled, an interactive terminal interface starts:

- Enter project requirements
- `Alt+Enter`: Submit requirements and start project
- `Esc`: Exit

### Agent Operation Flow

1. **President Agent**: Receives requirements from client and develops overall plan
2. **President Secretary**: Records requirements, elaborates details, and passes to department secretary
3. **Department Manager Agent**: Divides tasks into section-level units
4. **Department Secretary**: Elaborates tasks and coordinates with section secretaries
5. **Section Manager Agent**: Develops implementation plan and specifications
6. **Section Secretary**: Creates draft implementation procedures
7. **Employee Agent**: Executes actual implementation

## Development

### Directory Structure

```
BuildBureau/
├── cmd/
│   └── buildbureau/      # Main application
│       └── main.go
├── internal/
│   ├── agent/            # Agent implementation
│   ├── config/           # Configuration management
│   ├── grpc/             # gRPC service implementation
│   ├── slack/            # Slack notifications
│   └── ui/               # Terminal UI
├── proto/
│   └── buildbureau/v1/   # Protocol Buffers definitions
├── pkg/                  # Public packages
├── config.yaml           # Default configuration
├── Makefile             # Build scripts
└── go.mod               # Go dependencies
```

### Testing

```bash
make test
```

### Format and Lint

```bash
make lint
```

## gRPC Services

gRPC services defined for each level:

- **PresidentService**: Project planning
- **DepartmentManagerService**: Task division
- **SectionManagerService**: Implementation planning
- **EmployeeService**: Task execution

See `proto/buildbureau/v1/service.proto` for details.

## Slack Notifications

Slack notifications are sent for the following events:

- Project start
- Task completion
- Error occurrence
- Project completion

Notification enabling/disabling and content can be configured in `config.yaml`.

## License

See the [LICENSE](LICENSE) file for license information.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss the proposed changes.

## TODO

- [ ] Implement Google ADK integration
- [ ] Complete gRPC service implementation
- [ ] Implement agent-to-agent communication
- [ ] Implement knowledge base
- [ ] Implement tool system
- [ ] Support streaming
- [ ] Enhance error handling
- [ ] Improve test coverage
- [ ] Expand documentation
