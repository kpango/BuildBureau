# Contributing to BuildBureau

Thank you for contributing to BuildBureau!

## Code of Conduct

All participants in the project are expected to act with respect and courtesy.

## How to Contribute

### Bug Reports

If you find a bug:

1. Check existing issues
2. Create a new issue including:
   - Detailed description of the bug
   - Steps to reproduce
   - Expected behavior
   - Actual behavior
   - Environment information (OS, Go version, etc.)

### Feature Proposals

When proposing a new feature:

1. Discuss the proposal in an issue
2. Get approval from maintainers before implementation
3. Include implementation details

### Pull Requests

1. **Fork and Clone**
   ```bash
   git clone https://github.com/YOUR_USERNAME/BuildBureau.git
   cd BuildBureau
   ```

2. **Create Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Implement Changes**
   - Write code
   - Add tests
   - Update documentation

4. **Test**
   ```bash
   make test
   make lint
   ```

5. **Commit**
   ```bash
   git add .
   git commit -m "feat: Add your feature"
   ```

6. **Push**
   ```bash
   git push origin feature/your-feature-name
   ```

7. **Create Pull Request**
   - Create a PR on GitHub
   - Describe the changes
   - Link related issues

## Coding Conventions

### Go Style Guide

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Format with `gofmt`
- Check with `go vet`

### Naming Conventions

```go
// Package: lowercase, single word
package agent

// Exported types: PascalCase
type AgentPool struct {}

// Unexported types: camelCase
type agentInternal struct {}

// Functions: PascalCase (exported), camelCase (unexported)
func NewAgent() {}
func processTask() {}

// Constants: PascalCase (exported), camelCase (unexported)
const MaxRetryCount = 3
const defaultTimeout = 60
```

### Comments

```go
// Package agent provides AI agent implementations.
package agent

// Agent represents an AI agent interface.
// All agent types must implement this interface.
type Agent interface {
    // Process handles the given input and returns output.
    Process(ctx context.Context, input interface{}) (interface{}, error)
}
```

### Error Handling

```go
// Good example
if err != nil {
    return fmt.Errorf("failed to process task: %w", err)
}

// Wrap errors appropriately
if err := doSomething(); err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

## Testing

### Testing Requirements

- Write tests for all public functions
- Aim for 80%+ test coverage
- Use table-driven tests

### Test Example

```go
func TestNewAgent(t *testing.T) {
    tests := []struct {
        name    string
        id      string
        want    string
        wantErr bool
    }{
        {"valid", "agent-1", "agent-1", false},
        {"empty id", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            agent, err := NewAgent(tt.id)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewAgent() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if agent != nil && agent.ID() != tt.want {
                t.Errorf("NewAgent() ID = %v, want %v", agent.ID(), tt.want)
            }
        })
    }
}
```

## Documentation

### Code Documentation

- Add GoDoc comments to all public APIs
- Add explanatory comments to complex logic

### README Updates

When adding features, update the following:

- README.md
- docs/ARCHITECTURE.md (as needed)
- docs/CONFIGURATION.md (when adding configuration)

## Commit Messages

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code (whitespace, formatting, etc.)
- `refactor`: Code changes that neither fix bugs nor add features
- `perf`: Performance improvements
- `test`: Adding or modifying tests
- `chore`: Changes to build process or tools

### Example

```
feat(agent): Add support for streaming responses

Implement streaming support for long-running agent tasks.
This allows clients to receive progress updates in real-time.

Closes #123
```

## Pull Request Review

### For Reviewers

- Verify code quality
- Verify test adequacy
- Verify documentation updates
- Provide constructive feedback

### For Authors

- Respond to feedback
- Explain in comments when discussion is needed
- Incorporate changes promptly

## Release Process

### Versioning

Uses Semantic Versioning (SemVer):

- `MAJOR`: Incompatible changes
- `MINOR`: Backward-compatible feature additions
- `PATCH`: Backward-compatible bug fixes

### Release Steps

1. Update CHANGELOG
2. Create version tag
3. Create release on GitHub
4. Write release notes

## Community

### Questions

- Use GitHub Discussions
- Report bugs in issues

### Communication

- Both Japanese and English are acceptable
- Strive for courteous and constructive communication

## License

Contributed code is provided under the same license as the project (see LICENSE).

## Acknowledgments

Thank you to all contributors!
