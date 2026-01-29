# Development Guide

Developer guide for BuildBureau

## Development Environment Setup

### Required Tools

1. **Go 1.23 or higher**
   ```bash
   go version
   ```

2. **Protocol Buffers compiler (protoc)**
   ```bash
   # macOS
   brew install protobuf
   
   # Linux
   sudo apt install -y protobuf-compiler
   ```

3. **protoc plugins for Go**
   ```bash
   make install-tools
   ```

### Clone the Project

```bash
git clone https://github.com/kpango/BuildBureau.git
cd BuildBureau
```

### Install Dependencies

```bash
make deps
```

## Build

### Standard Build

```bash
make build
```

The built binary will be placed in `./bin/buildbureau`.

### Clean Build

```bash
make clean
make build
```

## Testing

### Run All Tests

```bash
make test
```

### Test Specific Package

```bash
go test -v ./internal/config
go test -v ./internal/agent
```

### Test Coverage

```bash
go test -cover ./...
```

## Code Quality

### Format

```bash
make fmt
```

### Lint

```bash
make vet
```

Or use the integrated command:

```bash
make lint
```

## Project Structure

```
BuildBureau/
├── cmd/
│   └── buildbureau/          # Main application
│       └── main.go
├── internal/                  # Internal packages
│   ├── agent/                # Agent implementation
│   │   ├── agent.go
│   │   └── agent_test.go
│   ├── config/               # Configuration management
│   │   ├── config.go
│   │   └── config_test.go
│   ├── grpc/                 # gRPC service implementation
│   ├── slack/                # Slack notifications
│   │   └── notifier.go
│   └── ui/                   # Terminal UI
│       └── ui.go
├── proto/                     # Protocol Buffers definitions
│   └── buildbureau/v1/
│       └── service.proto
├── pkg/                       # Public packages
│   ├── models/               # Data models
│   └── utils/                # Utilities
├── docs/                      # Documentation
│   ├── ARCHITECTURE.md
│   └── CONFIGURATION.md
├── config.yaml               # Default configuration
├── .env.example              # Environment variable template
├── Makefile                  # Build scripts
├── go.mod                    # Go dependencies
└── README.md
```

## Adding New Features

### 1. Implementing a New Agent

```go
// internal/agent/president.go
package agent

import (
    "context"
    "github.com/kpango/BuildBureau/internal/config"
)

type PresidentAgent struct {
    *BaseAgent
    // Additional fields
}

func NewPresidentAgent(id string, cfg config.AgentConfig) *PresidentAgent {
    return &PresidentAgent{
        BaseAgent: NewBaseAgent(id, AgentTypePresident, cfg),
    }
}

func (a *PresidentAgent) Process(ctx context.Context, input interface{}) (interface{}, error) {
    // Implementation
    return nil, nil
}
```

### 2. Adding a New gRPC Service

1. Add service definition to `proto/buildbureau/v1/service.proto`

```protobuf
service NewService {
    rpc NewMethod(RequestType) returns (ResponseType);
}
```

2. Generate Protocol Buffer code

```bash
make proto
```

3. Add service implementation

```go
// internal/grpc/new_service.go
package grpc

type NewServiceServer struct {
    // Fields
}

func (s *NewServiceServer) NewMethod(ctx context.Context, req *pb.RequestType) (*pb.ResponseType, error) {
    // Implementation
    return &pb.ResponseType{}, nil
}
```

### 3. Adding a New Slack Notification Event

1. Add notification configuration to `internal/config/config.go`

```go
type NotificationsConfig struct {
    // Existing events
    NewEvent NotificationConfig `yaml:"newEvent"`
}
```

2. Add configuration to `config.yaml`

```yaml
slack:
  notifications:
    newEvent:
      enabled: true
      message: "New event: {{.Data}}"
```

3. Add method to `internal/slack/notifier.go`

```go
func (n *Notifier) SendNewEvent(ctx context.Context, data string) error {
    return n.Send(ctx, NotificationNewEvent, NotificationData{
        // Data
    })
}
```

## Debugging

### Setting Log Level

```yaml
# config.yaml
ui:
  logLevel: "debug"
```

Or use environment variable:

```bash
export LOG_LEVEL=debug
```

### Using a Debugger

```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Run with debugger
dlv debug ./cmd/buildbureau
```

## Writing Tests

### Unit Tests

```go
// internal/agent/president_test.go
package agent

import (
    "context"
    "testing"
)

func TestPresidentAgent_Process(t *testing.T) {
    agent := NewPresidentAgent("test-1", config.AgentConfig{})
    
    ctx := context.Background()
    result, err := agent.Process(ctx, "test input")
    
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    
    // Assertions
}
```

### Table-Driven Tests

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"case1", "input1", "output1", false},
        {"case2", "input2", "output2", false},
        {"error case", "bad", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := functionUnderTest(tt.input)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Pull Request Guidelines

### Branch Strategy

- `main`: Production branch
- `develop`: Development branch
- `feature/*`: Feature addition branch
- `fix/*`: Bug fix branch
- `docs/*`: Documentation update branch

### Commit Messages

Follow Conventional Commits:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Code style
- `refactor`: Refactoring
- `test`: Tests
- `chore`: Build/configuration, etc.

Example:
```
feat(agent): Add LLM integration for president agent

Implement LLM API calls using Google ADK.
Add retry logic and error handling.

Closes #123
```

### Pull Request Checklist

- [ ] Code builds successfully
- [ ] All tests pass
- [ ] Added new tests
- [ ] Updated documentation
- [ ] Commit messages are appropriate
- [ ] Received code review

## Common Issues and Solutions

### 1. Build Error: "cannot find package"

```bash
make deps
go mod tidy
```

### 2. Test Error: "no such file or directory"

Check if the path is relative:

```go
// Bad example
os.ReadFile("config.yaml")

// Good example
os.ReadFile("/absolute/path/config.yaml")
```

### 3. Protocol Buffer Generation Error

```bash
make install-tools
make proto
```

## Performance Profiling

### CPU Profile

```bash
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof
```

### Memory Profile

```bash
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

## CI/CD

### GitHub Actions (Planned)

```yaml
# .github/workflows/test.yml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '1.23'
      - run: make deps
      - run: make test
      - run: make build
```

## Release

### Versioning

Follow Semantic Versioning: `MAJOR.MINOR.PATCH`

### Release Process

1. Create version tag
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. Create release notes

3. Build binaries
   ```bash
   make build
   ```

## References

- [Go Documentation](https://golang.org/doc/)
- [gRPC Go Tutorial](https://grpc.io/docs/languages/go/)
- [Protocol Buffers Guide](https://developers.google.com/protocol-buffers)
- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Slack API Documentation](https://api.slack.com/)
