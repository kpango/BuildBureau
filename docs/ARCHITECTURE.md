# BuildBureau Architecture Documentation

## System Architecture

BuildBureau is a hierarchical multi-agent system designed to mimic the structure of a corporate organization.

### Hierarchical Structure

```
┌─────────────────────────────────────────────┐
│             Client (User)                   │
└────────────────┬────────────────────────────┘
                 │
┌────────────────▼────────────────────────────┐
│          President Agent                    │
│        + President Secretary Agent          │
│  (Overall planning & requirements)          │
└────────────────┬────────────────────────────┘
                 │
┌────────────────▼────────────────────────────┐
│       Department Manager Agent              │
│      + Department Secretary Agent           │
│  (Task division & assignment to sections)   │
└────────────────┬────────────────────────────┘
                 │
┌────────────────▼────────────────────────────┐
│      Section Manager Agent × N              │
│     + Section Secretary Agent × N           │
│  (Implementation planning & instructions)   │
└────────────────┬────────────────────────────┘
                 │
┌────────────────▼────────────────────────────┐
│        Employee Agent × M                   │
│  (Actual implementation & deliverables)     │
└─────────────────────────────────────────────┘
```

## Component Details

### 1. Agent Layer

#### President Agent
- **Role**: Overall project planning and requirements organization
- **Input**: Requirement specifications from client
- **Output**: Task list (instructions to department manager)
- **Functions**:
  - Requirements analysis and understanding
  - Overall project planning
  - Task prioritization
  - Overall progress monitoring

#### President Secretary Agent
- **Role**: Support for president and requirements detailing
- **Functions**:
  - Requirements documentation
  - Knowledge base updates
  - Task detailing
  - Handoff to department secretary

#### Department Manager Agent
- **Role**: Divide tasks into section units
- **Input**: Task list from president
- **Output**: Section task plan
- **Functions**:
  - Technical analysis of tasks
  - Appropriate assignment to section managers
  - Issue coordination and prioritization

#### Department Secretary Agent
- **Role**: Task detailing and research
- **Functions**:
  - Technical research on tasks
  - Collection of related information
  - Coordination with section secretaries

#### Section Manager Agent × N
- **Role**: Implementation planning and specification creation
- **Input**: Section tasks
- **Output**: Implementation specifications
- **Functions**:
  - Creation of detailed implementation plans
  - Work assignment to employees
  - Review of implementation results

#### Section Secretary Agent × N
- **Role**: Draft creation of implementation procedures
- **Functions**:
  - Specification refinement
  - Coordination with other section secretaries
  - Creation of implementation procedure documents

#### Employee Agent × M
- **Role**: Actual implementation work
- **Input**: Implementation specifications
- **Output**: Deliverables (code, documents, etc.)
- **Functions**:
  - Coding
  - Document creation
  - Test execution

### 2. Communication Layer (gRPC)

#### Protocol Buffers Definition

All services are defined in `proto/buildbureau/v1/service.proto`.

Main message types:
- `RequirementSpec`: Project requirements
- `TaskUnit`: Individual task
- `TaskList`: List of tasks
- `SectionTask`: Task per section
- `ImplementationSpec`: Implementation specification
- `ResultArtifact`: Execution result

Main services:
- `PresidentService`: Project planning
- `DepartmentManagerService`: Task division
- `SectionManagerService`: Implementation planning
- `EmployeeService`: Task execution

### 3. Configuration Management Layer

#### YAML Configuration File (config.yaml)

All configuration is managed in a single YAML file:

```yaml
agents:          # Agent configuration
llm:            # LLM provider configuration
grpc:           # gRPC communication configuration
slack:          # Slack notification configuration
ui:             # Terminal UI configuration
system:         # General system configuration
```

#### Environment Variables

Sensitive information is managed via environment variables:
- `SLACK_BOT_TOKEN`: Slack Bot token
- `SLACK_CHANNEL_ID`: Notification channel ID
- `GOOGLE_AI_API_KEY`: Google AI API key

### 4. UI Layer (Bubble Tea)

#### Terminal UI Components

- **Input Area**: Input for project requirements
- **Status Display**: Status of each agent
- **Message Log**: Display of system messages
- **Progress Display**: Visualization of progress

#### Key Operations

- `Alt+Enter`: Submit requirements
- `Esc`: Exit

### 5. Notification Layer (Slack)

#### Event Types

1. **Project Start** (`project_start`)
   - When a project is started

2. **Task Complete** (`task_complete`)
   - When each task is completed

3. **Error Occurred** (`error`)
   - When an error occurs

4. **Project Complete** (`project_complete`)
   - When the entire project is completed

#### Message Templates

The following variables are available in templates:
- `{{.ProjectName}}`: Project name
- `{{.TaskName}}`: Task name
- `{{.Agent}}`: Agent ID
- `{{.ErrorMessage}}`: Error message
- `{{.Timestamp}}`: Timestamp

## Data Flow

### 1. Project Start Flow

```
Client
  ↓ (Requirements input)
President Agent
  ↓ (Requirements analysis)
President Secretary Agent
  ↓ (Requirements detailing)
Department Secretary Agent
  ↓ (Research & coordination)
Department Manager Agent
  ↓ (Task division)
Section Secretary Agent × N
  ↓ (Specification refinement)
Section Manager Agent × N
  ↓ (Implementation planning)
Employee Agent × M
  ↓ (Implementation)
Deliverables
```

### 2. Result Aggregation Flow

```
Employee Agent
  ↓ (Deliverables)
Section Manager Agent
  ↓ (Review & aggregation)
Department Manager Agent
  ↓ (Integration)
President Agent
  ↓ (Final confirmation)
Client
```

## Extensibility

### Scalability

- **Horizontal Scaling**: Number of agents can be configured in YAML
- **Distributed Execution**: Independent processes possible through gRPC interface
- **Microservice Architecture**: Each layer can be separated as different services in the future

### Customization

- **Agent Personality**: Customize prompts in YAML
- **Tool Addition**: Configure available tools for each agent
- **Notification Customization**: Modify Slack message templates

## Security

### Sensitive Information Management

- Separation of sensitive information via environment variables
- Addition of `.env` file to `.gitignore`
- Recommended periodic rotation of API keys

### Communication Security

- gRPC over TLS support (enabled via configuration)
- Authentication and authorization features planned

## Performance

### Concurrent Processing

- Efficient concurrent processing using Go goroutines
- Asynchronous communication between agents
- Parallel task execution (can be limited via configuration)

### Resource Management

- Prevention of infinite waiting through timeout settings
- Limitation on retry attempts
- Maximum limit on concurrent task count

## Troubleshooting

### Common Issues

1. **Build Error**: Reinstall dependencies with `make deps`
2. **Slack Notifications Not Received**: Check token and channel ID
3. **Agent Timeout**: Increase timeout value in `config.yaml`

### Log Checking

```bash
# Log directory
./logs/

# Log level setting
ui.logLevel: "debug"  # in config.yaml
```

## Future Development Plans

1. **Google ADK Integration**: Implementation of actual LLM calls
2. **Knowledge Base**: Knowledge base system shared between agents
3. **Tool System**: Implementation of external tools available to agents
4. **Web Interface**: Addition of browser-based UI
5. **Metrics**: Prometheus-compatible metrics collection
6. **A2A Protocol**: Implementation of Agent-to-Agent communication
