# Contributing to BuildBureau

Thank you for your interest in contributing to BuildBureau! This document
provides guidelines and instructions for contributing.

## Development Setup

### Prerequisites

- Go 1.26.0 or later
- Git
- Make (optional, but recommended)

### Getting Started

1. Fork the repository on GitHub
2. Clone your fork:

   ```bash
   git clone https://github.com/YOUR_USERNAME/BuildBureau.git
   cd BuildBureau
   ```

3. Install dependencies:

   ```bash
   make install
   # or
   go mod download
   ```

4. Set up environment variables:

   ```bash
   cp .env.example .env
   # Edit .env with your API keys (optional for development)
   ```

5. Build the project:
   ```bash
   make build
   # or
   go build -o buildbureau ./cmd/buildbureau
   ```

## Development Workflow

### Building

```bash
make build
```

### Running Tests

```bash
make test
```

### Running the Application

```bash
make run
```

### Code Formatting

Before submitting changes, ensure your code is properly formatted:

```bash
make fmt
```

### Linting

Run the linter to check for common issues:

```bash
make lint
```

## Project Structure

```
BuildBureau/
├── cmd/                  # Main application entry points
├── internal/            # Internal packages (not exported)
│   ├── agent/          # Agent implementations
│   ├── config/         # Configuration loading
│   ├── grpc/           # gRPC communication
│   ├── llm/            # LLM integration
│   ├── slack/          # Slack notifications
│   └── tui/            # Terminal UI
├── pkg/                 # Public packages
│   ├── protocol/       # Protocol definitions
│   └── types/          # Common types
├── agents/              # Agent configuration files
└── examples/            # Example configurations and tests
```

## Adding New Features

### Adding a New Agent Type

1. Create a new file in `internal/agent/` (e.g., `myagent.go`)
2. Implement the `types.Agent` interface
3. Add configuration support in `internal/agent/organization.go`
4. Create a YAML configuration in `agents/`
5. Update documentation

### Adding a New LLM Provider

1. Create a new provider in `internal/llm/providers.go`
2. Implement the `Provider` interface
3. Add initialization in `internal/llm/manager.go`
4. Update configuration types if needed
5. Add tests

### Improving the TUI

1. Edit `internal/tui/tui.go`
2. Follow Bubble Tea patterns
3. Test with various terminal sizes
4. Update help text

## Testing

### Unit Tests

Write unit tests for new functionality:

```go
package agent

import "testing"

func TestMyFeature(t *testing.T) {
    // Test implementation
}
```

Run tests:

```bash
go test ./...
```

### Integration Tests

Add integration tests in the `examples/` directory that demonstrate end-to-end
functionality.

## Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Write clear, descriptive variable and function names
- Add comments for exported functions and types
- Keep functions focused and small

## Commit Guidelines

- Write clear, concise commit messages
- Start with a verb in present tense (e.g., "Add", "Fix", "Update")
- Reference issues when applicable (e.g., "Fix #123")
- Keep commits focused on a single change

Example:

```
Add support for custom LLM endpoints

- Implement endpoint configuration parsing
- Add client connection handling
- Update documentation

Fixes #42
```

## Pull Request Process

1. Create a feature branch:

   ```bash
   git checkout -b feature/my-feature
   ```

2. Make your changes and commit them

3. Push to your fork:

   ```bash
   git push origin feature/my-feature
   ```

4. Open a Pull Request on GitHub

5. Ensure CI passes:
   - All tests pass
   - Code is formatted
   - No linting errors

6. Request review from maintainers

7. Address review feedback

8. Once approved, a maintainer will merge your PR

## Code Review

All submissions require review. We use GitHub pull requests for this purpose.

Reviewers will check:

- Code quality and style
- Test coverage
- Documentation
- Performance implications
- Security considerations

## Reporting Issues

### Bug Reports

When reporting bugs, include:

- BuildBureau version
- Go version
- Operating system
- Steps to reproduce
- Expected vs actual behavior
- Relevant logs or error messages

### Feature Requests

When requesting features:

- Describe the use case
- Explain the expected behavior
- Suggest implementation if possible
- Consider backwards compatibility

## Getting Help

- Check existing documentation
- Search closed issues
- Open a new issue for questions
- Join community discussions

## License

By contributing, you agree that your contributions will be licensed under the
same license as the project.

## Recognition

Contributors will be recognized in:

- The project README
- Release notes
- The GitHub contributors page

Thank you for contributing to BuildBureau! 🏢
