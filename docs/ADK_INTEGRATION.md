# ADK Integration Guide

This document explains how BuildBureau integrates with Google's Agent
Development Kit (ADK).

## Overview

BuildBureau now includes ADK-powered agents that leverage
`google.golang.org/adk` for enhanced agent capabilities. The ADK (Agent
Development Kit) is Google's official framework for building LLM-based agents
with structured configuration, tool support, and memory management.

## What is ADK?

The Agent Development Kit (google.golang.org/adk) provides:

- **Structured Agent Configuration** - llmagent.Config for consistent agent
  setup
- **Model Abstraction** - Clean interface for different LLM providers
- **Tool Support** - Built-in framework for agent tools and functions
- **Memory Management** - Conversation history and context management
- **Session Handling** - Multi-turn conversation support

## ADK-Powered Agents

BuildBureau provides ADK implementations for all agent roles:

### 1. ADK Engineer Agent

```go
engineer, err := agent.NewADKEngineerAgent("adk-engineer-1", config, apiKey)
if err != nil {
    log.Fatal(err)
}

// ADK agent is now configured with:
// - Gemini 2.0 Flash model
// - Engineer-specific system instructions
// - Tool capabilities (extensible)
// - Memory management (extensible)
```

**System Instruction:**

```
You are a skilled software engineer.
Your responsibilities:
- Implement code according to specifications
- Write clean, maintainable code
- Add appropriate comments and documentation
- Consider edge cases and error handling
- Follow best practices for the target language
```

### 2. ADK Manager Agent

```go
manager, err := agent.NewADKManagerAgent("adk-manager-1", config, apiKey)
```

**System Instruction:**

```
You are a software development manager.
Your responsibilities:
- Create detailed technical specifications
- Design software architecture
- Break down projects into implementable components
- Define interfaces and data structures
- Plan testing strategies
```

### 3. ADK Director Agent

```go
director, err := agent.NewADKDirectorAgent("adk-director-1", config, apiKey)
```

**System Instruction:**

```
You are a technical director.
Your responsibilities:
- Analyze project requirements
- Perform research on technologies and approaches
- Break down large projects into manageable tasks
- Make architectural decisions
- Allocate resources across teams
```

### 4. ADK President Agent

```go
president, err := agent.NewADKPresidentAgent("adk-president-1", config, apiKey)
```

**System Instruction:**

```
You are the president of a software development organization.
Your responsibilities:
- Clarify client requirements
- Define high-level objectives
- Ensure project alignment with goals
- Communicate with stakeholders
- Oversee project success
```

## Implementation Details

### Architecture

```
┌─────────────────────────────────────────┐
│  ADKAgent (BuildBureau wrapper)         │
│  - Implements types.Agent interface     │
│  - Maintains BaseAgent compatibility    │
└───────────────┬─────────────────────────┘
                │
                ▼
┌─────────────────────────────────────────┐
│  llmagent.Config (ADK configuration)    │
│  - Name, Description                    │
│  - System Instructions                  │
│  - Model (gemini.NewModel)              │
│  - Tools (extensible)                   │
│  - Memory (extensible)                  │
└───────────────┬─────────────────────────┘
                │
                ▼
┌─────────────────────────────────────────┐
│  gemini.NewModel (ADK model wrapper)    │
│  - Wraps genai.Client                   │
│  - Implements model.LLM interface       │
└───────────────┬─────────────────────────┘
                │
                ▼
┌─────────────────────────────────────────┐
│  genai.Client (Google AI SDK)           │
│  - Actual Gemini API calls              │
│  - Content generation                   │
└─────────────────────────────────────────┘
```

### Code Structure

**ADKAgent wraps:**

1. `BaseAgent` - For compatibility with existing agent interface
2. `llmagent.Config` - ADK configuration
3. Model name and API key - For dynamic model creation

**Key Methods:**

- `NewADKEngineerAgent()` - Creates engineer with ADK config
- `ProcessTask()` - Executes task using ADK-configured model
- `GetModelName()` - Returns the model being used

### Configuration Options

ADK agents support all standard BuildBureau agent configuration:

```yaml
# Agent configuration
name: "ADK Engineer"
role: "Engineer"
model: "gemini-2.0-flash-exp" # Optional, defaults to this
description: "An ADK-powered engineer agent"
system_prompt: "Custom instructions" # Optional, overrides defaults
```

## Usage Examples

### Basic Usage

```go
package main

import (
    "context"
    "log"
    "os"

    "github.com/kpango/BuildBureau/internal/agent"
    "github.com/kpango/BuildBureau/pkg/types"
)

func main() {
    // Get API key
    apiKey := os.Getenv("GEMINI_API_KEY")

    // Create config
    config := &types.AgentConfig{
        Name: "My ADK Engineer",
        Role: "Engineer",
    }

    // Create ADK agent
    engineer, err := agent.NewADKEngineerAgent("eng-1", config, apiKey)
    if err != nil {
        log.Fatal(err)
    }

    // Start agent
    ctx := context.Background()
    engineer.Start(ctx)
    defer engineer.Stop(ctx)

    // Create task
    task := &types.Task{
        ID:          "task-1",
        Title:       "Implement Function",
        Description: "Create a function",
        Content:     "Write a Python function to sort a list",
        Priority:    1,
    }

    // Process with ADK agent
    response, err := engineer.ProcessTask(ctx, task)
    if err != nil {
        log.Fatal(err)
    }

    // Response contains actual LLM-generated code
    fmt.Println(response.Result)
}
```

### Running the Example

```bash
# Set API key
export GEMINI_API_KEY="your-gemini-api-key"

# Run ADK example
go run examples/test_adk_agents/main.go
```

Expected output:

```
╔════════════════════════════════════════════════════════════╗
║       BuildBureau - ADK Agent Integration Example         ║
╚════════════════════════════════════════════════════════════╝

✓ GEMINI_API_KEY found

═══════════════════════════════════════════════════════════
Testing ADK Engineer Agent
═══════════════════════════════════════════════════════════
Creating ADK Engineer Agent...
✓ ADK Engineer created (Model: gemini-2.0-flash-exp)
✓ Engineer started
Sending task to ADK Engineer...
Task: Implement Fibonacci Function

Response received:
─────────────────────────────────────────────────────────
ADK Agent (adk-engineer-1 - gemini-2.0-flash-exp) Response:

[Real LLM-generated Python code with explanations]
─────────────────────────────────────────────────────────
Status: completed
```

## Testing

### Unit Tests

```bash
# Run ADK agent tests (skips without real API key)
go test ./internal/agent -run TestADK -v

# Run with real API key
export GEMINI_API_KEY="your-key"
go test ./internal/agent -run TestADKAgent_ProcessTask -v
```

### Integration Test

```bash
# Full integration test with real LLM
export GEMINI_API_KEY="your-key"
go run examples/test_adk_agents/main.go
```

## Extending ADK Agents

### Adding Tools

ADK supports tools (functions) that agents can call:

```go
// Future enhancement: Add tools to llmagent.Config
config := llmagent.Config{
    Name:        "engineer",
    Model:       model,
    Instruction: "...",
    Tools: []tool.Tool{
        // Add custom tools here
        myCodeAnalysisTool,
        myFilSystemTool,
    },
}
```

### Adding Memory

ADK supports conversation memory:

```go
// Future enhancement: Configure memory
config := llmagent.Config{
    Name:        "engineer",
    Model:       model,
    Instruction: "...",
    // Memory configuration for multi-turn conversations
}
```

### Custom Models

ADK supports different model backends:

```go
// Use different model
model, _ := gemini.NewModel(ctx, "gemini-1.5-pro", &genai.ClientConfig{
    APIKey: apiKey,
})
```

## Comparison: Standard vs ADK Agents

| Feature           | Standard Agent | ADK Agent               |
| ----------------- | -------------- | ----------------------- |
| Configuration     | Custom         | ADK llmagent.Config     |
| Model Integration | Direct genai   | ADK model.LLM interface |
| Tools Support     | Manual         | ADK tool.Tool framework |
| Memory            | Manual         | ADK memory framework    |
| Multi-turn        | Manual session | ADK session management  |
| Extensibility     | Custom code    | ADK plugins             |

## Benefits of ADK Integration

1. **Structured Configuration**
   - Consistent agent setup across roles
   - Standard llmagent.Config format
   - Easy to extend with tools and memory

2. **Google Official Framework**
   - Maintained by Google
   - Best practices built-in
   - Regular updates and improvements

3. **Future-Proof**
   - Ready for advanced ADK features
   - Tool support prepared
   - Memory management prepared

4. **Compatible**
   - Works with existing BuildBureau architecture
   - Implements types.Agent interface
   - Can mix standard and ADK agents

## Limitations

Current implementation:

- Uses genai client directly for generation
- Full session management not yet implemented
- Tool integration not yet enabled
- Memory features not yet configured

These are planned enhancements that ADK's architecture supports.

## Troubleshooting

### "GEMINI_API_KEY is required"

Solution: Set the environment variable:

```bash
export GEMINI_API_KEY="your-api-key"
```

### "failed to create ADK Gemini model"

Solutions:

1. Check API key is valid
2. Verify internet connectivity
3. Check Gemini API status

### Tests Skip

ADK tests skip without real API key (expected behavior):

```
Skipping ADK test: GEMINI_API_KEY not set or using demo value
```

To run with real API: `export GEMINI_API_KEY="real-key"`

## Resources

- **ADK Documentation**: https://google.github.io/adk-docs/
- **ADK Go Package**: https://pkg.go.dev/google.golang.org/adk
- **Gemini API**: https://ai.google.dev/docs
- **Get API Key**: https://aistudio.google.com/app/apikey

## Summary

BuildBureau's ADK integration provides:

- ✅ ADK-powered agents for all roles
- ✅ Structured configuration using llmagent.Config
- ✅ Gemini model integration via ADK
- ✅ Compatibility with existing architecture
- ✅ Ready for advanced ADK features (tools, memory)
- ✅ Comprehensive tests and examples

The integration demonstrates how to use Google's official Agent Development Kit
while maintaining BuildBureau's flexible agent architecture.
