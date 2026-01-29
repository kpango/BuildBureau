# BuildBureau Architecture

## Overview

BuildBureau is a multi-agent AI system that simulates a software development company's organizational structure. The system uses Google's Agent Development Kit (ADK) and follows the Agent2Agent (A2A) protocol for inter-agent communication.

## System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    BuildBureau System                        │
│                                                              │
│  ┌────────────┐  ┌──────────┐  ┌──────────┐  ┌─────────┐  │
│  │  Terminal  │  │  Agent   │  │  Slack   │  │  Config │  │
│  │     UI     │  │  System  │  │ Notifier │  │  Loader │  │
│  └────────────┘  └──────────┘  └──────────┘  └─────────┘  │
│        │              │              │              │       │
│        └──────────────┴──────────────┴──────────────┘       │
│                          │                                  │
│                    Main Process                             │
└─────────────────────────────────────────────────────────────┘
```

### Agent Hierarchy

BuildBureau implements a hierarchical multi-agent system:

```
                    ┌─────────────┐
                    │   CEO       │◄──── Client Requests
                    │             │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │ CEO         │
                    │ Secretary   │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │ Department  │
                    │ Head        │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │ DeptHead    │
                    │ Secretary   │
                    └──────┬──────┘
                           │
          ┌────────────────┼────────────────┐
          │                │                │
    ┌─────▼─────┐    ┌─────▼─────┐   ┌─────▼─────┐
    │ Manager   │    │ Manager   │   │ Manager   │
    │ (Frontend)│    │ (Backend) │   │ (QA)      │
    └─────┬─────┘    └─────┬─────┘   └─────┬─────┘
          │                │                │
    ┌─────▼─────┐    ┌─────▼─────┐   ┌─────▼─────┐
    │ Manager   │    │ Manager   │   │ Manager   │
    │ Secretary │    │ Secretary │   │ Secretary │
    └─────┬─────┘    └─────┬─────┘   └─────┬─────┘
          │                │                │
    ┌─────┴─────┐    ┌─────┴─────┐   ┌─────┴─────┐
    │           │    │           │   │           │
┌───▼───┐  ┌───▼───┐│Worker│ ┌──▼───┐│Worker│ ┌──▼───┐
│Worker │  │Worker ││      │ │Worker││      │ │Worker│
└───────┘  └───────┘└──────┘ └──────┘└──────┘ └──────┘
```

## Component Details

### 1. Agent System

**Location:** `internal/agent/`

The agent system manages the entire hierarchy of AI agents:

- **CEO Agent**: Handles client negotiations and requirement clarification
- **CEO Secretary**: Records decisions and manages knowledge base
- **Department Head Agent**: Breaks down projects into major task categories
- **DeptHead Secretary**: Conducts research and details requirements
- **Manager Agents**: Create technical specifications for specific areas (Frontend, Backend, QA)
- **Manager Secretaries**: Research technical details and document specifications
- **Worker Agents**: Execute development tasks

Each agent:
- Uses Google's Generative AI (Gemini) as the LLM backend
- Has a configurable instruction prompt defining its role
- Can use various tools (search, calculator, code executor, etc.)
- Has a configurable temperature setting for response creativity

### 2. Configuration System

**Location:** `internal/config/`

The configuration system provides:

- **YAML-based configuration**: All system behavior is configurable via `configs/config.yaml`
- **Environment variable expansion**: Secrets can be provided via environment variables
- **Validation**: Ensures configuration is valid before system starts
- **Dynamic agent creation**: Agent hierarchy is built based on configuration

Key configuration sections:
- `hierarchy`: Defines organizational structure
- `agents`: Configures each agent type's behavior
- `slack`: Slack integration settings
- `system`: System-level settings (logging, timeouts, UI, etc.)

### 3. Slack Notification System

**Location:** `internal/slack/`

The Slack notifier provides:

- **Event-based notifications**: Configurable triggers for different events
- **Role-based filtering**: Only notify configured roles for specific events
- **Channel mapping**: Route notifications to appropriate channels
- **Message formatting**: Consistent message format with timestamps and agent names

Supported events:
- `task_assigned`: When tasks are delegated
- `task_completed`: When tasks finish
- `project_started`: When new projects begin
- `project_completed`: When projects finish
- `error_occurred`: When errors happen
- `milestone_reached`: When milestones are achieved

### 4. Terminal UI

**Location:** `internal/ui/`

The Terminal UI provides:

- **Interactive input**: Text area for entering client requests
- **Real-time display**: Shows agent conversations as they happen
- **Color coding**: Different colors for different agent roles
- **Hierarchical visualization**: Indentation shows task flow through hierarchy
- **Status monitoring**: Shows when system is processing

Built with Charm's Bubble Tea TUI framework for a modern, responsive experience.

### 5. gRPC Service Definitions

**Location:** `proto/`

Protocol Buffer definitions for agent communication:

- **CEOService**: Client request handling and project outline
- **DeptHeadService**: Task planning and distribution
- **ManagerService**: Technical specification and task assignment
- **WorkerService**: Task execution

These definitions enable:
- **Type-safe communication**: Structured messages between agents
- **Future scalability**: Easy to split into separate processes/services
- **A2A compliance**: Follows Agent2Agent protocol patterns

## Data Flow

### Request Processing Flow

1. **User Input**: User enters client request via Terminal UI
2. **CEO Processing**: CEO agent receives and analyzes request
3. **CEO Secretary**: Records requirements to knowledge base
4. **Delegation to DeptHead**: CEO delegates to Department Head
5. **DeptHead Secretary**: Researches and details requirements
6. **Task Breakdown**: DeptHead breaks into categories
7. **Manager Distribution**: Tasks distributed to specialized managers
8. **Manager Secretaries**: Research technical details
9. **Technical Specs**: Managers create detailed specifications
10. **Worker Assignment**: Tasks assigned to worker agents
11. **Implementation**: Workers execute tasks
12. **Completion**: Results bubble back up the hierarchy

### Notification Flow

```
Event Occurs
    ↓
Check if Slack enabled
    ↓
Check if event type configured
    ↓
Check if role should be notified
    ↓
Get appropriate channel
    ↓
Format message
    ↓
Send to Slack
```

## Technology Stack

- **Language**: Go 1.21+
- **AI Backend**: Google Generative AI (Gemini) via `google/generative-ai-go`
- **TUI Framework**: Charm Bubble Tea
- **API Communication**: Slack Go SDK
- **Configuration**: YAML with gopkg.in/yaml.v3
- **Protocol Definitions**: Protocol Buffers (for future gRPC)

## Design Principles

### 1. Separation of Concerns
Each component has a single, well-defined responsibility:
- Agent system handles AI logic
- Config system handles configuration
- Slack system handles notifications
- UI system handles user interaction

### 2. Configuration over Code
System behavior can be modified without recompilation:
- Agent prompts are in YAML
- Organizational structure is configurable
- Notification settings are adjustable
- UI preferences are customizable

### 3. Hierarchical Task Distribution
Tasks flow from high-level to low-level:
- CEO: Business requirements
- DeptHead: Major task categories
- Managers: Technical specifications
- Workers: Implementation details

### 4. Loose Coupling
Components communicate through well-defined interfaces:
- Agent system is independent of UI
- Notifications are decoupled from agent logic
- Configuration is centralized

### 5. Single Binary Deployment
All components compile into one executable:
- No runtime dependencies
- Easy distribution
- Simple deployment

## Extensibility

### Adding New Agent Types

1. Add agent configuration to `configs/config.yaml`
2. Update `Agent` type in `internal/agent/agent.go`
3. Implement processing logic

### Adding New Event Types

1. Add event type constant in `internal/slack/notifier.go`
2. Add event configuration in YAML
3. Call notifier when event occurs

### Adding New Tools

1. Define tool in agent configuration
2. Implement tool logic in agent processing
3. Update agent instruction prompts

### Scaling to Distributed System

The gRPC service definitions allow future distribution:

1. **Process Separation**: Run agents in separate processes
2. **Service Scaling**: Scale specific agent types independently
3. **Multi-Language**: Implement agents in different languages
4. **Remote Communication**: Agents communicate over network

## Security Considerations

- **API Keys**: Stored in environment variables, never in code
- **Secrets Management**: Use `${VAR}` syntax in YAML
- **Validation**: Configuration is validated before use
- **Logging**: Sensitive data is not logged
- **Token Handling**: Slack tokens are handled securely

## Performance Considerations

- **Parallel Agent Processing**: Agents can process in parallel (via goroutines)
- **Efficient UI Updates**: UI uses buffered updates to avoid flicker
- **Connection Pooling**: Reuse HTTP clients and connections
- **Resource Cleanup**: Proper cleanup of resources on shutdown

## Future Enhancements

1. **Persistence**: Store agent conversations and knowledge in database
2. **Audit Trail**: Full audit log of all agent actions
3. **Web UI**: Add web-based interface alongside TUI
4. **Agent Marketplace**: Pluggable agent implementations
5. **Multi-Project**: Handle multiple projects simultaneously
6. **Real Code Execution**: Integrate actual code generation and testing
7. **Human-in-Loop**: Allow manual intervention at any hierarchy level
8. **Metrics Dashboard**: Real-time metrics and analytics
