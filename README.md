# BuildBureau

BuildBureau is a sophisticated multi-agent AI system built in Go that simulates a software development company's organizational structure. Using Google's Agent Development Kit (ADK) and the Agent2Agent (A2A) protocol, multiple specialized AI agents collaborate hierarchically to handle client requests, from negotiation to final implementation.

## 🏗️ Architecture

BuildBureau implements a hierarchical multi-agent system that mirrors a real company structure:

```
CEO + CEO Secretary
    ↓
Department Head + Department Head Secretary
    ↓
Managers (by specialty) + Manager Secretaries
    ↓
Workers (developers)
```

### Agent Roles

- **CEO & CEO Secretary**: Handle client negotiations, clarify requirements, record decisions to knowledge base, and notify stakeholders via Slack
- **Department Head & Secretary**: Break down projects into major tasks, conduct research, and distribute work to specialized managers
- **Managers & Manager Secretaries**: Create detailed technical specifications, select technology stacks, and assign tasks to development teams
- **Workers**: Execute actual development tasks according to specifications

Each leader-level agent has a dedicated secretary agent that handles:
- Information gathering and research
- Documentation and knowledge base management
- Task scheduling and coordination
- Detailed requirement analysis

## ✨ Features

- **Multi-Agent Collaboration**: Agents communicate using gRPC-based interfaces with A2A protocol patterns
- **Hierarchical Task Distribution**: Tasks flow from CEO → Department Head → Manager → Worker with progressive refinement
- **Slack Integration**: Configurable notifications for task assignments, completions, milestones, and errors
- **Interactive Terminal UI**: Real-time visualization of agent conversations using Bubble Tea TUI framework
- **YAML Configuration**: Fully configurable agent hierarchy, prompts, tools, and notification settings
- **Single Binary Deployment**: Compile once, run anywhere with no additional dependencies
- **Knowledge Base**: Persistent storage of project requirements, specifications, and decisions

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or later
- Google Gemini API key
- (Optional) Slack Bot Token for notifications

### Installation

1. Clone the repository:
```bash
git clone https://github.com/kpango/BuildBureau.git
cd BuildBureau
```

2. Set up environment variables:
```bash
export GOOGLE_API_KEY="your-gemini-api-key"
export SLACK_BOT_TOKEN="xoxb-your-slack-bot-token"  # Optional
export SLACK_CHANNEL_MAIN="C01234567"  # Optional
export SLACK_CHANNEL_MANAGEMENT="C01234568"  # Optional
export SLACK_CHANNEL_DEV="C01234569"  # Optional
export SLACK_CHANNEL_ERRORS="C01234570"  # Optional
```

3. Install dependencies and build:
```bash
make build
```

### Running BuildBureau

Run with default configuration:
```bash
make run
```

Run with custom configuration:
```bash
make run-config CONFIG=path/to/your/config.yaml
```

Or run the binary directly:
```bash
./bin/buildbureau -config configs/config.yaml
```

## 📝 Configuration

BuildBureau is configured through YAML files. See `configs/config.yaml` for a complete example.

### Key Configuration Sections

#### Agent Hierarchy
```yaml
hierarchy:
  departments: 1
  managers_per_department: 3
  manager_specialties:
    - Frontend Development
    - Backend Development
    - Quality Assurance
  workers_per_manager: 2
```

#### Agent Configuration
Each agent type can be customized with its own model, instruction prompt, tools, and temperature:
```yaml
agents:
  ceo:
    model: "gemini-2.0-flash-exp"
    instruction: "You are an experienced CEO..."
    tools: ["search", "calculator"]
    temperature: 0.7
```

#### Slack Notifications
Configure which events trigger notifications and which roles receive them:
```yaml
slack:
  enabled: true
  token: "${SLACK_BOT_TOKEN}"
  notify_on:
    task_assigned:
      enabled: true
      roles: ["CEO", "Manager"]
      channel: "main"
```

## 🎨 Terminal UI

BuildBureau features an interactive Terminal User Interface (TUI) that shows:

- Real-time agent conversations with color-coded roles
- Hierarchical indentation showing task flow
- Timestamp for each agent action
- Input prompt for entering client requests
- Status indicators for ongoing processing

### UI Controls

- **Enter**: Submit client request
- **Ctrl+C** or **Esc**: Quit application

## 🔧 Development

### Project Structure

```
BuildBureau/
├── cmd/
│   └── buildbureau/          # Main application entry point
├── internal/
│   ├── agent/                # Agent implementation and orchestration
│   ├── config/               # Configuration loading and validation
│   ├── slack/                # Slack notification system
│   ├── ui/                   # Terminal UI implementation
│   └── grpc/                 # gRPC service implementations (future)
├── proto/                    # Protocol buffer definitions
├── configs/                  # Configuration files
├── docs/                     # Additional documentation
└── Makefile                  # Build automation
```

### Available Make Targets

```bash
make build          # Build the binary
make run            # Build and run
make test           # Run tests
make test-coverage  # Run tests with coverage
make fmt            # Format code
make lint           # Run linter
make clean          # Clean build artifacts
make help           # Show all available targets
```

## 🔐 Security

- API keys should be provided via environment variables
- Slack tokens are never logged or displayed
- All credentials use the `${VAR}` syntax in YAML for environment variable expansion

## 🏗️ Technical Details

### Technologies Used

- **Go 1.21+**: Main programming language
- **Google Gemini API**: LLM backend for agents via genai SDK
- **Agent Development Kit (ADK)**: Framework for agent creation and orchestration
- **Agent2Agent (A2A) Protocol**: Standard for inter-agent communication
- **gRPC**: Service definitions for agent communication interfaces
- **Bubble Tea**: Terminal UI framework
- **Slack API**: Notification and collaboration integration
- **YAML**: Configuration management

### Design Principles

1. **Modular Architecture**: Each component (agent, config, UI, Slack) is independently testable
2. **Configuration over Code**: Behavior can be modified without recompilation
3. **Hierarchical Task Flow**: Progressive task refinement from high-level goals to specific implementations
4. **Loose Coupling**: gRPC interfaces allow future distribution across processes or machines
5. **Single Binary**: Easy deployment with no runtime dependencies

## 📚 Documentation

For more detailed documentation, see the `docs/` directory:

- Architecture diagrams
- Agent communication flow
- Configuration guide
- Deployment instructions

## 🤝 Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## 📄 License

See [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Google Agent Development Kit (ADK) team
- Agent2Agent (A2A) protocol contributors
- Go community and open source contributors

## 📧 Contact

For questions or support, please open an issue on GitHub.

---

Built with ❤️ using Go and AI