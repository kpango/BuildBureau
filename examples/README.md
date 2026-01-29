# BuildBureau Examples

This directory contains example applications demonstrating various features of the BuildBureau system.

## Demo Application

The `demo` directory contains a complete walkthrough of the BuildBureau multi-agent system.

### Running the Demo

```bash
cd examples/demo
go run main.go
```

## Google ADK Integration Example

The `google-adk` directory demonstrates using Google's Generative AI (Gemini) with BuildBureau.

### Running the Google ADK Example

```bash
# Set your Google AI API key
export GOOGLE_AI_API_KEY="your-api-key-here"

# Run the example
cd examples/google-adk
go run main.go
```

### Getting a Google AI API Key

1. Go to [Google AI Studio](https://aistudio.google.com/app/apikey)
2. Click "Create API Key"
3. Copy the key and set it as an environment variable

### What the Example Shows

- Creating a Google ADK client
- Simple text generation with system instructions
- Streaming responses
- Using temperature and max tokens parameters
- Error handling

## Demo Application (Multi-Agent System)

The `demo` directory contains a complete walkthrough of the BuildBureau multi-agent system.

### What the Demo Shows

The demo application demonstrates:

1. **System Initialization**
   - Creating an agent pool with all agent types
   - Setting up the knowledge base
   - Registering tools
   - Initializing the LLM client

2. **Agent Hierarchy in Action**
   - President Agent: Plans the project and breaks it into tasks
   - Department Manager: Divides tasks into section-level work
   - Section Manager: Creates detailed implementation specifications
   - Employee Agents: Execute the actual work

3. **Knowledge Sharing**
   - All agents store information in the shared knowledge base
   - Information flows through the hierarchy
   - Audit trail with creator tracking

4. **Tool Usage**
   - Tools are available to agents based on permissions
   - Built-in tools for web search, code analysis, document management, etc.

### Expected Output

```
=== BuildBureau System Demo ===

1. Initializing system components...
✓ System initialized with 7 agents

2. Creating gRPC services...
✓ All services created

3. President Agent: Planning project...
✓ Created 1 high-level tasks
  - task-1: Initial Planning

4. Department Manager: Dividing tasks into sections...
✓ Created 1 section plans

5. Section Manager: Preparing implementation specifications...
✓ Created 1 implementation specs

6. Employee Agents: Executing tasks...
✓ Task task-1 completed: success

7. Final System State:
   [Shows all agent statuses and knowledge base entries]
```

## Creating Your Own Example

To create a new example:

1. Create a new directory under `examples/`
2. Add a `main.go` file
3. Import the necessary packages:
   ```go
   import (
       "github.com/kpango/BuildBureau/internal/agent"
       "github.com/kpango/BuildBureau/internal/grpc"
       "github.com/kpango/BuildBureau/internal/knowledge"
       "github.com/kpango/BuildBureau/internal/llm"
       "github.com/kpango/BuildBureau/internal/tools"
   )
   ```

4. Initialize the system components
5. Use the gRPC services to coordinate agents

## Available Components

### Agent Pool
Manages all agents in the system. Supports registration, lookup, and status tracking.

### Knowledge Base
Shared storage for information across agents. Supports CRUD operations and search.

### Tool Registry
Manages available tools that agents can use. Built-in tools include:
- web_search
- code_analyzer
- document_manager
- file_operations
- code_execution

### LLM Client
Abstraction layer for LLM providers. Currently supports mock client for testing.

### gRPC Services
Four main services coordinating the agent hierarchy:
- PresidentService
- DepartmentManagerService
- SectionManagerService
- EmployeeService

## More Examples Coming Soon

- Custom tool implementation
- Real LLM integration
- Streaming responses
- Error recovery patterns
- Advanced agent coordination
