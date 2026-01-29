# Contributing to BuildBureau

Thank you for your interest in contributing to BuildBureau! This document provides guidelines and instructions for contributing.

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code. Please be respectful and considerate in your interactions.

## How to Contribute

### Reporting Bugs

If you find a bug, please open an issue with:

1. **Clear title and description**
2. **Steps to reproduce**
3. **Expected vs actual behavior**
4. **Environment details** (OS, Go version, etc.)
5. **Relevant logs or screenshots**

### Suggesting Features

Feature requests are welcome! Please:

1. **Check existing issues** to avoid duplicates
2. **Describe the problem** you're trying to solve
3. **Propose a solution** with examples
4. **Explain the benefits** to users

### Submitting Pull Requests

1. **Fork the repository**
2. **Create a feature branch** (`git checkout -b feature/amazing-feature`)
3. **Make your changes**
4. **Add tests** for new functionality
5. **Update documentation** as needed
6. **Ensure tests pass** (`make test`)
7. **Format your code** (`make fmt`)
8. **Commit your changes** with clear messages
9. **Push to your fork**
10. **Open a Pull Request**

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional, for convenience)

### Setup Steps

1. Clone your fork:
```bash
git clone https://github.com/YOUR_USERNAME/BuildBureau.git
cd BuildBureau
```

2. Add upstream remote:
```bash
git remote add upstream https://github.com/kpango/BuildBureau.git
```

3. Install dependencies:
```bash
go mod download
```

4. Set up environment:
```bash
cp .env.example .env
# Edit .env with your API keys
```

5. Run tests:
```bash
make test
```

## Project Structure

```
BuildBureau/
├── cmd/
│   └── buildbureau/      # Main application
├── internal/
│   ├── agent/            # Agent system
│   ├── config/           # Configuration
│   ├── slack/            # Slack integration
│   └── ui/               # Terminal UI
├── proto/                # gRPC definitions
├── configs/              # Configuration files
├── docs/                 # Documentation
└── Makefile              # Build tasks
```

## Coding Standards

### Go Style

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` for formatting
- Run `golangci-lint` for linting
- Keep functions focused and small
- Write descriptive variable names

### Code Organization

- **Package by feature**: Group related code together
- **Internal packages**: Use `internal/` for non-exported code
- **Clear interfaces**: Define clear contracts between components
- **Minimal dependencies**: Avoid unnecessary external dependencies

### Commit Messages

Follow conventional commits:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code formatting
- `refactor`: Code restructuring
- `test`: Adding tests
- `chore`: Maintenance tasks

Example:
```
feat(agent): add support for custom tool definitions

- Allow users to define custom tools in YAML
- Add tool validation and registration
- Update documentation

Closes #123
```

### Testing

- **Write tests** for new functionality
- **Maintain coverage**: Aim for >80%
- **Test edge cases**: Consider error conditions
- **Use table-driven tests** where appropriate

Example test:
```go
func TestAgentProcess(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:  "valid input",
            input: "test task",
            want:  "expected output",
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Documentation

- **Update README** for user-facing changes
- **Add godoc comments** for exported functions
- **Update ARCHITECTURE.md** for structural changes
- **Add examples** in documentation

## Areas for Contribution

### High Priority

- [ ] Add more comprehensive tests
- [ ] Improve error handling and recovery
- [ ] Add metrics and monitoring
- [ ] Implement knowledge base persistence
- [ ] Add more agent tool implementations

### Medium Priority

- [ ] Web UI interface
- [ ] Multi-project support
- [ ] Agent marketplace/plugins
- [ ] Performance optimizations
- [ ] Additional notification channels (Discord, Teams, etc.)

### Good First Issues

Look for issues labeled `good-first-issue`:

- Documentation improvements
- Example configurations
- Small bug fixes
- Test coverage improvements

## Release Process

Releases are managed by maintainers:

1. Version bump in code
2. Update CHANGELOG.md
3. Create release tag
4. Build binaries for all platforms
5. Publish GitHub release

## Getting Help

- **Questions**: Open a GitHub Discussion
- **Bugs**: Open an issue
- **Security**: Email maintainers directly
- **General**: Join community channels

## Recognition

Contributors will be:
- Listed in CONTRIBUTORS.md
- Mentioned in release notes
- Given credit in documentation

## License

By contributing, you agree that your contributions will be licensed under the same license as the project (see LICENSE file).

---

Thank you for contributing to BuildBureau! 🙏
