# BuildBureau

BuildBureau is a hierarchical multi-agent system built in Go that simulates an enterprise organization structure. Multiple specialized AI agents collaborate to handle complex tasks through a chain of command, from CEO to employees.

## 🏢 Overview

BuildBureau implements a multi-agent system where agents are organized in a corporate hierarchy:

- **CEO Agent** - Receives client requests and delegates to managers
- **CEO Secretary** - Documents requirements and manages knowledge base
- **Manager Agents** - Break down projects into categories
- **Manager Secretaries** - Conduct research and detailed planning
- **Lead Agents** - Create technical specifications and development plans
- **Lead Secretaries** - Perform technical research and documentation
- **Employee Agents** - Execute implementation tasks

Each level works with their secretary agent to document, research, and refine tasks before delegating down the hierarchy.

## 🏗️ Architecture

### Technologies

- **Go** - Core implementation language
- **A2A Protocol** - Agent-to-agent communication (gRPC-ready)
- **Bubble Tea** - Interactive terminal UI
- **Slack API** - Real-time notifications
- **YAML** - Configuration management

### Key Components

1. **Agent Framework** (`internal/agent/`)
   - Base agent implementation with task handling
   - Specialized agents (CEO, Manager, Lead, Employee, Secretary)
   - Event-driven communication

2. **A2A Protocol** (`pkg/a2a/`)
   - gRPC protocol definitions for agent communication
   - Designed for future inter-process communication

3. **UI System** (`internal/ui/`)
   - Real-time terminal interface using Bubble Tea
   - Color-coded agent activity display
   - Interactive task submission

4. **Slack Integration** (`internal/slack/`)
   - Configurable notifications for agent events
   - Channel mapping by agent level
   - Event filtering

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or later
- (Optional) Slack Bot Token for notifications

### Installation

```bash
# Clone the repository
git clone https://github.com/kpango/BuildBureau.git
cd BuildBureau

# Build the binary
go build -o buildbureau ./cmd/buildbureau

# Run the application
./buildbureau
```

### Configuration

Create or edit `config.yaml`:

```yaml
# LLM Settings (for future ADK integration)
llm:
  provider: "gemini"
  api_key: "${GEMINI_API_KEY}"
  model: "gemini-1.5-pro"

# Slack Configuration (optional)
slack:
  enabled: false
  token: "${SLACK_BOT_TOKEN}"
  channels:
    management: "C123456789"
    updates: "C987654321"
    dev: "C456789123"
  notifications:
    notify_on_task_assigned: ["CEO", "Manager"]
    notify_on_task_completed: ["Employee"]
    notify_on_error: true

# UI Configuration
ui:
  show_agent_logs: true
  color_coding: true
  refresh_rate_ms: 100
```

## 📖 Usage

### Running the Application

```bash
./buildbureau
```

### Using the Interface

1. **Start the application** - The TUI will display with an input area
2. **Enter a client request** - Type your project requirements
3. **Press Enter** - Watch the agents work through the hierarchy
4. **Observe agent activity** - Color-coded logs show each agent's actions
5. **Press Esc** - Exit the application

### Example Request

```
We need to build a new e-commerce website with user authentication, 
product catalog, shopping cart, and payment integration.
```

The system will:
1. CEO receives and clarifies the request with CEO secretary
2. CEO delegates to managers
3. Managers research and break down into technical categories
4. Leads create detailed specifications
5. Employees implement the tasks

## 🎨 Agent Color Coding

- 🔴 **CEO** - Pink/Magenta
- 🟠 **Manager** - Orange
- 🔵 **Lead** - Cyan
- 🟢 **Employee** - Green
- 🟣 **Secretary** - Purple

## 🔔 Slack Notifications

Configure Slack integration to receive real-time updates:

1. Create a Slack Bot and get the token
2. Set `SLACK_BOT_TOKEN` environment variable
3. Configure channel IDs in `config.yaml`
4. Enable notifications: `slack.enabled: true`

Notifications include:
- Task assignments
- Task completions
- Error events
- Agent messages

## 🏗️ Project Structure

```
BuildBureau/
├── cmd/
│   └── buildbureau/        # Main application entry point
│       └── main.go
├── internal/
│   ├── agent/              # Agent implementations
│   │   ├── base.go         # Base agent functionality
│   │   ├── ceo.go          # CEO agent
│   │   ├── manager.go      # Manager agent
│   │   ├── lead.go         # Lead agent
│   │   ├── employee.go     # Employee agent
│   │   └── secretary.go    # Secretary agent
│   ├── config/             # Configuration management
│   │   └── config.go
│   ├── slack/              # Slack integration
│   │   └── notifier.go
│   └── ui/                 # Terminal UI
│       └── tui.go
├── pkg/
│   ├── a2a/                # A2A protocol definitions
│   │   └── agent.proto
│   └── types/              # Shared types
│       └── types.go
├── config.yaml             # Configuration file
├── go.mod                  # Go dependencies
└── README.md               # This file
```

## 🔮 Future Enhancements

### Phase 2 - ADK Integration
- Integrate Google's Agent Development Kit (ADK)
- Connect to Gemini LLM for intelligent agent behavior
- Implement semantic task understanding and planning

### Phase 3 - Advanced Features
- Persistent knowledge base
- Task history and analytics
- Multi-project management
- External tool integrations
- Web-based UI option

### Phase 4 - Distributed System
- Full gRPC implementation for distributed agents
- Agent discovery service
- Load balancing across agent instances
- Cross-platform agent communication

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

See [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the TUI
- Inspired by Google's Agent Development Kit (ADK)
- Implements principles from the Agent2Agent (A2A) protocol