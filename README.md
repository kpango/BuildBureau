# BuildBureau

BuildBureau is a multi-layer AI Agent system that orchestrates tasks from a President agent down to Workers.

## Architecture

- **President**: Plans the project.
- **Manager**: Breaks down plans into section tasks.
- **Section**: Creates implementation specs.
- **Worker**: Implements the specs.

Technically, it uses:
- **Go 1.25.6** (Generics used for Agent abstraction)
- **ADK (Agent Development Kit)**: Custom abstraction for LLM agents.
- **A2A Protocol**: Internal message bus for agent communication.
- **Bubble Tea**: TUI for interaction.
- **Slack**: Notification integration.

## Configuration

Configuration is managed in `configs/config.yaml`.
You can set API keys via environment variables (e.g. `OPENAI_API_KEY`).

## Running

```bash
go mod tidy
go run cmd/buildbureau/main.go
```

If no API keys are present, the system defaults to **Mock Mode**, where agents return predefined responses for testing the flow.

## Testing

```bash
go test ./...
```
