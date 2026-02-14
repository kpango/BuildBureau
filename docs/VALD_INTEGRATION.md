# Vald Integration Guide

This document describes the patterns, practices, and components incorporated from the [vdaas/vald](https://github.com/vdaas/vald) repository into BuildBureau.

## Overview

BuildBureau has integrated several enterprise-grade patterns from vdaas/vald, a production-ready vector database system. These patterns improve code quality, maintainability, observability, and developer experience.

## What Was Incorporated

### 1. Modular Makefile Structure

**Location**: `Makefile.d/`

**Inspired by**: Vald's 18-module Makefile system

The monolithic Makefile has been augmented with modular makefiles organized by domain:

- **functions.mk**: Reusable build functions with color-coded output
- **build.mk**: Binary compilation targets with version info
- **test.mk**: Comprehensive test targets (unit, coverage, benchmarks)
- **tools.mk**: Tool installation helpers
- **docker.mk**: Docker build, scan, and management
- **lint.mk**: Linting and code quality checks

**Benefits**:
- Better maintainability - each module handles a single domain
- Improved discoverability - targets are organized logically
- Reusability - common functions can be used across targets
- Consistency - standardized patterns across all builds

**Usage**:
```bash
# Build targets
make binary/build              # Build binary
make binary/build/release      # Optimized release build
make binary/build/debug        # Debug build with symbols

# Test targets
make test/all                  # Run all tests
make test/coverage             # Generate coverage report
make test/coverage/html        # HTML coverage report
make test/bench                # Run benchmarks

# Tool management
make tools/install             # Install all development tools
make tools/golangci-lint       # Install golangci-lint
make tools/reviewdog           # Install reviewdog

# Docker operations
make docker/build/buildbureau  # Build Docker image
make docker/build/multiarch    # Multi-architecture build
make docker/scan               # Scan for vulnerabilities

# Linting
make lint/all                  # Run all linters
make lint/go/fix               # Run linters with auto-fix
make lint/security             # Security scanning
```

### 2. Enhanced Error Handling

**Location**: `internal/errors/errors.go`

**Inspired by**: Vald's domain-specific error factories

A comprehensive error handling package with factory functions for different domains:

**Features**:
- Domain-specific error types (Agent, LLM, Memory, Config, Communication)
- Error wrapping with context preservation
- Function name reflection for better debugging
- Standard library compatibility (errors.Is, errors.As)

**Usage**:
```go
import "github.com/kpango/BuildBureau/internal/errors"

// Agent errors
err := errors.ErrAgentNotFound("agent-123")
err := errors.ErrAgentTimeout("agent-456")

// LLM errors
err := errors.ErrLLMProviderNotFound("gemini")
err := errors.ErrLLMAPIKeyMissing("openai")

// Memory errors
err := errors.ErrMemoryStoreNotInitialized
err := errors.ErrMemoryNotFound("entry-789")

// Wrapping errors
err := errors.Wrap(originalErr, "context message")
err := errors.Wrapf(originalErr, "failed to process %s", taskID)

// Checking errors
if errors.Is(err, errors.ErrTimeout) {
    // Handle timeout
}
```

### 3. Logger Abstraction

**Location**: `internal/observability/logger/`

**Inspired by**: Vald's multi-backend logger system

A structured logging abstraction with multiple levels and formats:

**Features**:
- Multiple log levels (Debug, Info, Warn, Error, Fatal)
- Text and JSON output formats
- Field-based structured logging
- Global logger singleton
- Thread-safe operations

**Usage**:
```go
import "github.com/kpango/BuildBureau/internal/observability/logger"

// Initialize logger
logger.Init(
    logger.WithLevel(logger.LevelInfo),
    logger.WithFormat(logger.FormatJSON),
)

// Basic logging
logger.Info("Agent started")
logger.Error("Task failed", logger.Error(err))

// Structured logging with fields
logger.Info("Processing task",
    logger.String("task_id", "123"),
    logger.String("agent_id", "agent-456"),
    logger.Int("priority", 5),
)

// Create logger with persistent fields
agentLogger := logger.WithFields(
    logger.String("component", "agent"),
    logger.String("agent_id", "agent-123"),
)
agentLogger.Info("Task completed")
```

### 4. Sync Utilities

**Location**: `internal/sync/`

**Inspired by**: Vald's sync package with errgroup, atomic, semaphore

Concurrency utilities for safe parallel operations:

**Features**:
- ErrGroup for goroutine error handling
- Concurrent-safe Map
- Once, Pool, Mutex, RWMutex wrappers

**Usage**:
```go
import (
    "context"
    "github.com/kpango/BuildBureau/internal/sync"
)

// ErrGroup for parallel operations
ctx, g := sync.NewErrGroup(context.Background())

g.Go(func() error {
    // Parallel operation 1
    return nil
})

g.Go(func() error {
    // Parallel operation 2
    return nil
})

if err := g.Wait(); err != nil {
    // Handle first error from any goroutine
}

// Concurrent-safe map
m := sync.NewMap()
m.Store("key", "value")
val, ok := m.Load("key")

// Pool for object reuse
pool := sync.NewPool(func() interface{} {
    return &MyStruct{}
})
obj := pool.Get().(*MyStruct)
// Use obj...
pool.Put(obj)
```

### 5. Enhanced CI/CD Workflows

**Location**: `.github/workflows/`

**Inspired by**: Vald's comprehensive CI/CD system

Three new GitHub Actions workflows:

#### reviewdog.yml
Provides inline PR feedback using reviewdog:
- golangci-lint integration
- staticcheck analysis
- misspell checking
- alex (inclusive language)

**Benefits**:
- Automatic code review comments on PRs
- Catches issues before merge
- Consistent code quality enforcement

#### code-quality.yml
Comprehensive quality checks:
- Format verification (gofmt, goimports)
- Go vet static analysis
- Security scanning (gosec with SARIF)
- Go mod verification
- Test coverage reporting (Codecov)

**Triggers**:
- Pull requests changing Go code
- Pushes to main/release branches
- Prevents redundant runs with concurrency control

#### docker-scan.yml
Container security scanning:
- Trivy vulnerability scanning
- Docker Scout CVE detection
- Hadolint Dockerfile linting
- SARIF results to GitHub Security

**Triggers**:
- Version tag pushes
- Dockerfile changes
- Weekly scheduled scans

## Integration with Existing Code

### Backward Compatibility

All new components are **additive** and maintain full backward compatibility:
- Existing Makefile targets still work
- No changes to existing error handling
- Logger is opt-in, doesn't affect existing code
- New workflows supplement existing CI

### Gradual Migration

Components can be adopted incrementally:

1. **Start with Makefile targets** for improved developer experience
2. **Adopt error handling** in new code first
3. **Integrate logger** when refactoring existing logging
4. **Use sync utilities** for new concurrent code
5. **Observe CI workflows** and adjust as needed

### Migration Examples

#### Error Handling Migration
```go
// Before
if agent == nil {
    return fmt.Errorf("agent %s not found", id)
}

// After
if agent == nil {
    return errors.ErrAgentNotFound(id)
}
```

#### Logger Migration
```go
// Before
log.Printf("Processing task %s", taskID)

// After
logger.Info("Processing task", 
    logger.String("task_id", taskID))
```

## Testing

All new components have comprehensive test coverage:

```bash
# Test error handling
go test ./internal/errors/... -v

# Test logger
go test ./internal/observability/logger/... -v

# Test sync utilities
go test ./internal/sync/... -v

# Test coverage
make test/coverage
```

## Best Practices

### Error Handling
1. Use domain-specific error factories for clarity
2. Wrap errors with context using `errors.Wrap` or `errors.Wrapf`
3. Check error types with `errors.Is` or `errors.As`
4. Include function context in error messages

### Logging
1. Use appropriate log levels (Debug for development, Info for normal ops)
2. Include relevant fields for structured logging
3. Use component-specific loggers with `WithFields`
4. Avoid logging sensitive data (API keys, passwords)

### Concurrency
1. Use ErrGroup for operations that should fail fast
2. Reuse objects with Pool for high-throughput scenarios
3. Use Map for concurrent access to shared state
4. Properly handle context cancellation

### Makefiles
1. Keep targets focused and composable
2. Use color-coded output for better readability
3. Document targets with `## comments`
4. Include tool installation in targets (auto-install)

## Performance Impact

**Minimal overhead**:
- Error factories: No allocation for static errors
- Logger: Buffered writes, lazy evaluation
- Sync utilities: Thin wrappers around stdlib
- Makefile: Same execution, better organization

**Improvements**:
- Faster builds with cached tools
- Parallel test execution
- Reduced CI time with path filtering

## Future Enhancements

### Planned (from Vald patterns)
- [ ] Backoff and retry mechanisms
- [ ] Circuit breaker for LLM calls
- [ ] Utility packages (conv, strings, io)
- [ ] Reusable CI workflow templates
- [ ] Version management system
- [ ] Observability metrics collection

### Under Consideration
- [ ] Distributed tracing integration
- [ ] Prometheus metrics exporter
- [ ] Advanced test utilities (goleak, mocks)
- [ ] Semaphore and singleflight patterns
- [ ] Atomic operations utilities

## References

- **Vald Repository**: https://github.com/vdaas/vald
- **Vald Architecture**: https://vald.vdaas.org/docs/overview/architecture
- **Go Best Practices**: https://golang.org/doc/effective_go
- **Reviewdog**: https://github.com/reviewdog/reviewdog
- **Codecov**: https://codecov.io/

## Contributing

When contributing code that uses these patterns:

1. Follow the established patterns in each package
2. Add tests for new functionality
3. Update documentation when adding features
4. Use the new error types for domain-specific errors
5. Add structured logging to new components
6. Create Makefile targets for new operations

## Support

For questions or issues related to these patterns:

1. Check existing documentation in `docs/` directory
2. Review test files for usage examples
3. Refer to Vald's implementation for advanced patterns
4. Open an issue for bugs or enhancement requests

---

**Last Updated**: 2024-02-14  
**Vald Version Referenced**: Latest (as of 2024-02-14)  
**Integration Status**: Phase 1 Complete (Core Patterns)
