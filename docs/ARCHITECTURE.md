# BuildBureau Architecture

## System Overview

BuildBureau is a hierarchical multi-agent system that simulates a virtual
software development company. The system consists of five organizational layers
with specialized agents at each level.

## High-Level Architecture

```
                    ┌─────────────┐
                    │   Client    │
                    │    (You)    │
                    └──────┬──────┘
                           │
                           ▼
                    ┌─────────────┐
                    │  President  │◄───────┐
                    │   Agent     │        │
                    └──────┬──────┘        │
                           │               │
                           ▼               │
                    ┌─────────────┐        │
                    │  Secretary  │────────┘
                    │   Agent     │
                    └──────┬──────┘
                           │
                    ┌──────┴───────┐
                    ▼              ▼
            ┌─────────────┐ ┌─────────────┐
            │  Director   │ │  Director   │
            │   Agent     │ │   Agent     │
            └──────┬──────┘ └──────┬──────┘
                   │                │
                   ▼                ▼
            ┌─────────────┐ ┌─────────────┐
            │  Manager    │ │  Manager    │
            │   Agent     │ │   Agent     │
            └──────┬──────┘ └──────┬──────┘
                   │                │
            ┌──────┴─────┬──────────┴─────┐
            ▼            ▼                 ▼
    ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
    │  Engineer   │ │  Engineer   │ │  Engineer   │
    │   Agent     │ │   Agent     │ │   Agent     │
    └─────────────┘ └─────────────┘ └─────────────┘
```

## Agent Responsibilities

### President

- **Role**: Top-level client interface
- **Responsibilities**:
  - Clarifies client instructions
  - Summarizes objectives
  - Defines high-level requirements
  - Delegates to Secretary
- **Reports to**: Client
- **Delegates to**: Secretary

### Secretary

- **Role**: Administrative assistant to leadership
- **Responsibilities**:
  - Records goals and decisions
  - Schedules and monitors tasks
  - Delegates tasks downward
  - Maintains knowledge base
- **Reports to**: President, Directors, Managers
- **Delegates to**: Directors, Managers, Engineers

### Director

- **Role**: Project decomposition and research
- **Responsibilities**:
  - Performs research
  - Expands vague goals into actionable units
  - Decomposes projects into department-level tasks
  - Assigns work to managers
- **Reports to**: Secretary
- **Delegates to**: Managers

### Manager

- **Role**: Software design and specification
- **Responsibilities**:
  - Finalizes implementation specifications
  - Produces software designs
  - Creates pseudocode and technical specs
  - Delegates implementation to engineers
- **Reports to**: Directors
- **Delegates to**: Engineers

### Engineer

- **Role**: Code implementation
- **Responsibilities**:
  - Implements code per specifications
  - Writes tests
  - Debugs issues
  - Returns completed work
- **Reports to**: Managers
- **Delegates to**: None (leaf nodes)

## Data Flow

### Task Delegation Flow (Top-Down)

```
Client Request
      ↓
┌──────────────────────┐
│ 1. President         │ ← Clarifies requirements
│    - Parse request   │
│    - Define scope    │
└──────┬───────────────┘
       ↓
┌──────────────────────┐
│ 2. Secretary         │ ← Records and routes
│    - Log task        │
│    - Select director │
└──────┬───────────────┘
       ↓
┌──────────────────────┐
│ 3. Director          │ ← Research and decompose
│    - Research        │
│    - Break down      │
└──────┬───────────────┘
       ↓
┌──────────────────────┐
│ 4. Manager           │ ← Design specifications
│    - Create specs    │
│    - Design solution │
└──────┬───────────────┘
       ↓
┌──────────────────────┐
│ 5. Engineer          │ ← Implement code
│    - Write code      │
│    - Test            │
└──────────────────────┘
```

### Result Consolidation Flow (Bottom-Up)

```
┌──────────────────────┐
│ Engineer             │ → Code implementation
│ - Completed code     │
│ - Test results       │
└──────┬───────────────┘
       ↓
┌──────────────────────┐
│ Manager              │ → Code review
│ - Validates code     │
│ - Integrates         │
└──────┬───────────────┘
       ↓
┌──────────────────────┐
│ Director             │ → Project assembly
│ - Combines modules   │
│ - Final review       │
└──────┬───────────────┘
       ↓
┌──────────────────────┐
│ Secretary            │ → Documentation
│ - Creates summary    │
│ - Updates knowledge  │
└──────┬───────────────┘
       ↓
┌──────────────────────┐
│ President            │ → Client delivery
│ - Formats response   │
│ - Delivers to client │
└──────────────────────┘
       ↓
   Client Result
```

## Technical Architecture

### Component Diagram

```
┌─────────────────────────────────────────────────┐
│                   BuildBureau                    │
│                                                  │
│  ┌────────────┐  ┌────────────┐  ┌───────────┐ │
│  │    TUI     │  │   Config   │  │   Types   │ │
│  │  (Bubble   │  │  (YAML)    │  │ (Shared)  │ │
│  │   Tea)     │  │            │  │           │ │
│  └─────┬──────┘  └──────┬─────┘  └─────┬─────┘ │
│        │                │               │        │
│        └────────┬───────┴───────┬───────┘        │
│                 │               │                │
│           ┌─────▼───────────────▼─────┐          │
│           │   Organization Manager     │          │
│           │  - Agent lifecycle         │          │
│           │  - Hierarchy management    │          │
│           └─────┬──────────────────────┘          │
│                 │                                 │
│     ┌───────────┼───────────┐                    │
│     ▼           ▼           ▼                    │
│  ┌──────┐  ┌──────┐    ┌──────┐                 │
│  │Agents│  │ gRPC │    │ LLM  │                 │
│  │      │  │      │    │      │                 │
│  └───┬──┘  └──┬───┘    └──┬───┘                 │
│      │        │           │                      │
└──────┼────────┼───────────┼──────────────────────┘
       │        │           │
       │        │           │
       ▼        ▼           ▼
   ┌───────┐ ┌────────┐ ┌──────────┐
   │In-Mem │ │ Network│ │   APIs   │
   │Agents │ │ Agents │ │ (Gemini, │
   └───────┘ └────────┘ │  Claude) │
                        └──────────┘
```

### Communication Patterns

#### 1. In-Memory Communication (Current)

```
Agent A → Direct Method Call → Agent B
```

#### 2. gRPC Communication (Future)

```
Agent A → gRPC Client → Network → gRPC Server → Agent B
```

#### 3. LLM Integration

```
Agent → LLM Manager → Provider → LLM API
                    ↓
              Remote Agent API (for Claude, Codex, Qwen)
```

## Configuration Architecture

### YAML Configuration Hierarchy

```
config.yaml
    ├── organization
    │   └── layers[]
    │       ├── name
    │       ├── agent (→ agents/*.yaml)
    │       ├── count
    │       └── attach_to[]
    │
    ├── slack
    │   ├── enabled
    │   ├── token {env}
    │   ├── channels[]
    │   └── notify_on[]
    │
    └── llms
        ├── default_model
        └── api_keys
            ├── gemini {env}
            ├── claude {env}
            ├── codex {env}
            └── qwen {env}

agents/*.yaml (per agent)
    ├── name
    ├── role
    ├── description
    ├── model
    ├── system_prompt
    ├── capabilities[]
    └── sub_agents[] (optional)
        ├── name
        └── remote
            ├── endpoint
            └── capabilities[]
```

## Deployment Architecture

### Single Binary Deployment

```
┌────────────────────────────────────┐
│         Host Machine                │
│                                     │
│  ┌──────────────────────────────┐  │
│  │      buildbureau binary       │  │
│  │                               │  │
│  │  ┌─────────┐  ┌─────────┐    │  │
│  │  │President│  │Secretary│    │  │
│  │  └────┬────┘  └────┬────┘    │  │
│  │       │            │         │  │
│  │  ┌────▼────┐  ┌────▼────┐    │  │
│  │  │Director │  │ Manager │    │  │
│  │  └────┬────┘  └────┬────┘    │  │
│  │       │            │         │  │
│  │  ┌────▼────────────▼────┐    │  │
│  │  │    Engineers          │    │  │
│  │  └───────────────────────┘    │  │
│  │                               │  │
│  │  ┌────────────────────────┐   │  │
│  │  │     TUI Interface      │   │  │
│  │  └────────────────────────┘   │  │
│  └──────────────────────────────┘  │
│                                     │
└────────────────────────────────────┘
```

### Distributed Deployment (Future)

```
┌──────────────┐      ┌──────────────┐
│   Frontend   │      │   Backend    │
│    (TUI)     │─────▶│  President   │
└──────────────┘      │  Secretary   │
                      └──────┬───────┘
                             │
                    ┌────────┴────────┐
                    │                 │
              ┌─────▼─────┐    ┌──────▼──────┐
              │ Director  │    │  Director   │
              │  Node 1   │    │   Node 2    │
              └─────┬─────┘    └──────┬──────┘
                    │                 │
          ┌─────────┴─────┐    ┌──────┴──────┐
     ┌────▼────┐    ┌─────▼────▼────┐  ┌─────▼────┐
     │ Manager │    │   Manager     │  │ Manager  │
     │ Pod 1   │    │    Pod 2      │  │  Pod 3   │
     └────┬────┘    └─────┬─────────┘  └─────┬────┘
          │               │                   │
          └───────┬───────┴───────┬───────────┘
                  │               │
         ┌────────▼──┐    ┌───────▼────┐
         │ Engineer  │    │  Engineer  │
         │  Pool     │    │   Pool     │
         └───────────┘    └────────────┘
```

## Security Architecture

### Secrets Management

```
Environment Variables
        ↓
  ┌──────────────┐
  │   .env file  │
  └──────┬───────┘
         │
         ▼
  ┌──────────────┐
  │  Config      │ ← Validates env vars
  │  Loader      │
  └──────┬───────┘
         │
         ▼
  ┌──────────────┐
  │  LLM Manager │ ← Uses API keys
  └──────────────┘
```

### Authentication Flow (Remote Agents)

```
BuildBureau Agent
      ↓
   API Key
      ↓
   HTTP/gRPC Request
      ↓
Remote Agent Service
      ↓
   Validate Key
      ↓
   LLM Provider API
```

## Performance Considerations

- **Concurrency**: Multiple agents can process tasks simultaneously
- **Caching**: Agent responses can be cached for similar tasks
- **Load Balancing**: Tasks distributed across multiple agents of same type
- **Resource Limits**: Each agent has configurable resource constraints

## Extensibility Points

1. **New Agent Types**: Add by implementing `Agent` interface
2. **New LLM Providers**: Add by implementing `Provider` interface
3. **New Communication Protocols**: Extend transport layer
4. **Custom Business Logic**: Override agent methods
5. **Integration Points**: Webhooks, APIs, message queues

## Monitoring and Observability

```
Agents
  ↓
Metrics Collection
  ├─ Active Tasks
  ├─ Completed Tasks
  ├─ Error Rates
  └─ Response Times
  ↓
Monitoring Dashboard
  ↓
Alerts (Slack/Email)
```

## Future Architecture Enhancements

- [ ] Distributed agent deployment
- [ ] Event sourcing for task history
- [ ] GraphQL API for external integrations
- [ ] WebSocket support for real-time updates
- [ ] Multi-tenancy support
- [ ] Agent hot-reloading
- [ ] Kubernetes operator
- [ ] Service mesh integration
