# Generic Agent System

## Overview

The Generic Agent System is a flexible, configuration-driven approach to building multi-agent hierarchies in BuildBureau. Instead of hardcoded role-specific agent types (President, Secretary, Director, Manager, Engineer), the system uses a single `GenericAgent` implementation where behavior is defined through configuration files.

## Key Concepts

### Generic Agent

A `GenericAgent` is a flexible agent implementation that derives its behavior from:

1. **System Prompt**: Defines the agent's role, responsibilities, and decision-making approach
2. **Capabilities**: List of skills/abilities the agent possesses
3. **Hierarchical Relationships**: Parent agent and subordinates for delegation
4. **LLM Integration**: Uses language models for intelligent task processing
5. **Memory**: Learns from past interactions and tasks

### Configuration-Driven Behavior

Instead of writing role-specific code, behavior is configured through YAML files:

```yaml
name: CustomRole
role: CustomRole
description: A custom agent role
model: gemini
system_prompt: |
  You are a CustomRole agent in BuildBureau.
  Your responsibilities are:
  - Task A
  - Task B
  - Task C
  
  Always be thorough and strategic.
capabilities:
  - capability_1
  - capability_2
```

## Architecture

### GenericAgent

Located in `internal/agent/generic.go`, the GenericAgent provides:

- **Dynamic Task Processing**: Uses LLM with role-specific prompts
- **Memory-Enhanced Decisions**: Retrieves similar past tasks for context
- **Automatic Delegation**: Determines when to delegate based on LLM response
- **Flexible Hierarchy**: Parent-child relationships for task flow
- **Fallback Mode**: Works without LLM using basic role responses

### GenericOrganization

Located in `internal/agent/generic_organization.go`, the GenericOrganization:

- **Builds Hierarchy from Config**: Creates agents and relationships automatically
- **No Hardcoded Roles**: All roles defined in configuration
- **Flexible Structure**: Supports any organizational structure
- **Agent Management**: Query, start, stop, and monitor agents

## Usage

### Creating a Generic Organization

```go
import (
    "github.com/kpango/BuildBureau/internal/agent"
    "github.com/kpango/BuildBureau/internal/config"
)

// Load configuration
loader := config.NewLoader()
cfg, err := loader.LoadConfig("./config.yaml")
if err != nil {
    log.Fatal(err)
}

// Create generic organization
org, err := agent.NewGenericOrganization(cfg)
if err != nil {
    log.Fatal(err)
}

// Start the organization
ctx := context.Background()
if err := org.Start(ctx); err != nil {
    log.Fatal(err)
}
defer org.Stop(ctx)
```

### Processing Tasks

```go
// Create a task
task := &types.Task{
    ID:          "task-1",
    Title:       "Build Feature",
    Description: "Implement user authentication",
    Priority:    1,
}

// Process through organization
response, err := org.ProcessTask(ctx, task)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Status: %s\n", response.Status)
fmt.Printf("Result: %s\n", response.Result)
```

### Querying Organization

```go
// Get specific agent
agent := org.GetAgent("president-1")

// Get agents by role
engineers := org.GetAgentsByRole(types.RoleEngineer)

// Get all agents
allAgents := org.GetAllAgents()

// Get status
status := org.GetStatus()
for id, stats := range status {
    fmt.Printf("%s: %v\n", id, stats)
}
```

## Configuration

### Organization Structure (config.yaml)

```yaml
organization:
  layers:
    - name: President
      agent: ./agents/president.yaml
    - name: Secretary
      count: 1
      attach_to: [President, Director]
      agent: ./agents/secretary.yaml
    - name: Director
      count: 2
      agent: ./agents/director.yaml
    - name: Manager
      count: 3
      agent: ./agents/manager.yaml
    - name: Engineer
      count: 5
      agent: ./agents/engineer.yaml
```

### Agent Configuration (agents/custom.yaml)

```yaml
name: CustomAgent
role: CustomRole
description: A custom agent role
model: gemini
system_prompt: |
  You are a CustomAgent in BuildBureau.
  Your responsibilities include:
  - Analyzing requirements
  - Making strategic decisions
  - Delegating tasks appropriately
  
  When you receive a task:
  1. Analyze the requirements thoroughly
  2. Break it down into subtasks if needed
  3. Decide whether to delegate or handle yourself
  4. If delegating, indicate which subordinate should handle it
  
  Always provide clear reasoning for your decisions.
capabilities:
  - requirement_analysis
  - strategic_planning
  - task_delegation
```

## Task Processing Flow

1. **Task Received**: Agent receives task from parent or external source
2. **Memory Lookup**: Searches for similar past tasks for context
3. **LLM Processing**: Generates response using system prompt + task details + past context
4. **Delegation Decision**: Analyzes LLM response for delegation keywords
5. **Subordinate Selection**: Chooses appropriate subordinate if delegation needed
6. **Execution**: Processes task or delegates to subordinate
7. **Memory Storage**: Stores task and result for future reference
8. **Response**: Returns result to caller

## Advanced Features

### LLM-Powered Processing

The agent uses LLM to make intelligent decisions:

```go
// Prompt construction
prompt := system_prompt + "\n\n" +
         "Task: " + task.Title + "\n" +
         "Description: " + task.Description + "\n" +
         "Capabilities: " + capabilities + "\n" +
         "Subordinates: " + subordinate_list + "\n" +
         "Past Experience: " + similar_tasks
         
// LLM generates response
response := llm.Generate(prompt)
```

### Memory Integration

Agents learn from past interactions:

```go
// Store task
agent.memory.StoreTask(ctx, task, result, tags)

// Retrieve similar tasks
similar := agent.memory.GetRelatedTasks(ctx, query, limit)

// Use in decision making
context := buildContextFromMemory(similar)
```

### Delegation Strategies

Delegation can be enhanced with:

1. **Load Balancing**: Consider subordinate workload
2. **Capability Matching**: Match task requirements to subordinate capabilities
3. **Performance History**: Use past performance metrics
4. **LLM-Based Selection**: Let LLM choose best subordinate

## Comparison: Old vs New

### Old System (Role-Specific)

```go
// Hardcoded agent types
president := NewPresidentAgent(id, config)
director := NewDirectorAgent(id, config)
manager := NewManagerAgent(id, config)
engineer := NewEngineerAgent(id, config)

// Hardcoded behavior
func (p *PresidentAgent) ProcessTask(task) {
    // Fixed president logic here
}
```

### New System (Generic)

```go
// Single generic type
president := NewGenericAgent(id, RolePresident, config, llm)
director := NewGenericAgent(id, RoleDirector, config, llm)
manager := NewGenericAgent(id, RoleManager, config, llm)

// Behavior from configuration
// System prompt in YAML defines what agent does
```

## Benefits

1. **Flexibility**: Add new roles without code changes
2. **Maintainability**: Single implementation to maintain
3. **Configurability**: Behavior defined in config files
4. **Extensibility**: Easy to add new capabilities
5. **Testability**: Generic agents are easier to test
6. **LLM-Powered**: Intelligent, adaptive behavior
7. **Memory-Enhanced**: Learns and improves over time

## Migration Guide

### Gradual Migration

The generic system coexists with the old system:

1. **Phase 1**: Create generic versions alongside existing agents
2. **Phase 2**: Test generic system with subset of tasks
3. **Phase 3**: Gradually migrate agents to generic system
4. **Phase 4**: Remove old role-specific implementations

### Converting Existing Agent

**Before (engineer.go):**
```go
type EngineerAgent struct {
    *BaseAgent
    llm *llm.Manager
}

func (e *EngineerAgent) ProcessTask(task) {
    // Hardcoded engineer logic
    result := "Implementing code..."
    return result
}
```

**After (engineer.yaml):**
```yaml
name: Engineer
role: Engineer
system_prompt: |
  You are an Engineer at BuildBureau.
  You implement code according to specifications.
  Focus on:
  - Clean, maintainable code
  - Following best practices
  - Writing tests
```

## Testing

### Unit Tests

```go
func TestGenericAgent(t *testing.T) {
    config := &types.AgentConfig{
        SystemPrompt: "You are a test agent.",
    }
    
    agent := NewGenericAgent("test", RoleEngineer, config, nil)
    
    task := &types.Task{
        Title: "Test",
        Description: "Test task",
    }
    
    response, err := agent.ProcessTask(ctx, task)
    assert.NoError(t, err)
    assert.Equal(t, StatusCompleted, response.Status)
}
```

### Integration Tests

```go
func TestGenericOrganization(t *testing.T) {
    cfg := loadTestConfig()
    org, err := NewGenericOrganization(cfg)
    require.NoError(t, err)
    
    err = org.Start(ctx)
    require.NoError(t, err)
    
    response, err := org.ProcessTask(ctx, testTask)
    assert.NoError(t, err)
}
```

## Examples

See `examples/test_generic_system/main.go` for a complete working example.

## Future Enhancements

1. **Smart Delegation**: ML-based subordinate selection
2. **Dynamic Prompts**: Context-aware prompt modification
3. **Capability Matching**: Automatic task routing based on capabilities
4. **Performance Metrics**: Track and optimize agent performance
5. **Configuration Validation**: Validate YAML configs at load time
6. **Hot Reload**: Update agent behavior without restart
7. **Plugin System**: Load custom behaviors as plugins

## Troubleshooting

### Agent Not Processing Tasks

Check:
- LLM API key is set
- Agent configuration is valid
- Agent is started
- Task structure is correct

### Delegation Not Working

Check:
- Subordinates are properly connected
- LLM response includes delegation keywords
- Subordinate agents are started

### Memory Not Working

Check:
- Memory is enabled in config
- SQLite database is accessible
- Agent memory is initialized

## Best Practices

1. **Clear Prompts**: Write specific, actionable system prompts
2. **Capability Tags**: Use descriptive capability names
3. **Hierarchical Design**: Design clear parent-child relationships
4. **Testing**: Test agents with various task types
5. **Monitoring**: Track agent performance and adjust prompts
6. **Gradual Rollout**: Test thoroughly before production

## Support

For questions or issues:
- Check the tests in `internal/agent/generic_test.go`
- Review example in `examples/test_generic_system/`
- See main documentation in `AGENTS.md`
