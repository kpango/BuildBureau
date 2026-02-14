# BuildBureau Examples

This directory contains example programs demonstrating BuildBureau features.

## Structure

Each example is in its own subdirectory to avoid `main()` function conflicts:

```
examples/
├── test_basic/              # Basic agent demonstration
├── test_adk_agents/         # ADK-powered agents
├── test_agent_memory/       # Agent memory integration
├── test_implementations/    # Implementation demonstrations
├── test_memory/             # Memory system features
├── test_multiple_providers/ # Multi-provider LLM usage
└── remote_agent_server/     # Remote agent HTTP server
```

## Running Examples

Each example can be run independently:

```bash
# Basic example
go run examples/test_basic/main.go

# ADK agents
go run examples/test_adk_agents/main.go

# Agent memory
go run examples/test_agent_memory/main.go

# Memory system
go run examples/test_memory/main.go

# Multiple providers
go run examples/test_multiple_providers/main.go

# Remote agent server
go run examples/remote_agent_server/main.go
```

## Environment Variables

Most examples require API keys:

```bash
export GEMINI_API_KEY="your-gemini-key"
export OPENAI_API_KEY="your-openai-key"
export CLAUDE_API_KEY="your-claude-key"
```

See the main [README.md](../README.md) for detailed setup instructions.
