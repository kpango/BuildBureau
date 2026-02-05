# BuildBureau Makefile Documentation

## Overview

The BuildBureau Makefile provides a comprehensive, standardized interface for
building, testing, and managing the project. It's designed to be used
consistently across local development, Docker builds, and CI/CD pipelines.

## Quick Start

```bash
# Display all available targets
make help

# Build the application
make build

# Run tests
make test

# Build for production
make build-release

# Run all CI checks
make ci-all
```

## Target Categories

### Build Targets

#### `make build`

Build the application binary with version information embedded.

**Output**: `./build/buildbureau`

```bash
make build
```

#### `make build-debug`

Build with debug symbols and no optimization for debugging.

**Output**: `./build/buildbureau-debug`

```bash
make build-debug
```

#### `make build-release`

Build optimized release binary with path trimming.

**Output**: `./build/buildbureau`

```bash
make build-release
```

#### `make build-static`

Build a static binary (note: SQLite requires CGO, so this may not be fully
static).

**Output**: `./build/buildbureau-static`

```bash
make build-static
```

#### `make build-all`

Build for multiple platforms (linux/amd64, linux/arm64, darwin/amd64,
darwin/arm64).

**Output**: `./dist/buildbureau-{os}-{arch}`

```bash
make build-all
```

### Test Targets

#### `make test`

Run all tests with race detection.

```bash
make test
```

#### `make test-unit`

Run only unit tests (using `-short` flag).

```bash
make test-unit
```

#### `make test-coverage`

Run tests with coverage reporting.

**Output**: `./coverage/coverage.out`

```bash
make test-coverage
```

#### `make test-coverage-html`

Generate HTML coverage report.

**Output**: `./coverage/coverage.html`

```bash
make test-coverage-html
```

#### `make test-bench`

Run benchmark tests.

```bash
make test-bench
```

#### `make test-race`

Run tests with race detector.

```bash
make test-race
```

#### `make test/llm-integration`

Test with real LLM provider integration (requires API keys).

**Replaced**: `test_real_llm.sh` script

**Requirements**: At least one API key set (GEMINI_API_KEY, OPENAI_API_KEY, or
CLAUDE_API_KEY)

```bash
# Test with Gemini
export GEMINI_API_KEY="your-key"
make test/llm-integration

# Test with multiple providers
export GEMINI_API_KEY="key1"
export OPENAI_API_KEY="key2"
export CLAUDE_API_KEY="key3"
make test/llm-integration
```

This target tests actual LLM provider connectivity and functionality.

### Formatting Targets

BuildBureau includes comprehensive formatting support for multiple file types
with graceful degradation when tools are missing.

#### `make format`

Format all files across all supported types (Go, YAML, JSON, Markdown).

**Formats**:

- Go files: `gofmt` + `goimports` (if available)
- YAML files: `prettier` or `yamlfmt` (whichever is available)
- JSON files: `prettier` or `jq` (whichever is available)
- Markdown files: `prettier` (if available)

```bash
# Format all files
make format
```

**Auto-Install**: Automatically installs missing formatters on first use (except
jq, which requires manual system installation).

#### `make format/go`

Format only Go files using `gofmt` and `goimports`.

**Tools Used**:

- `gofmt`: Built-in Go formatter (always available)
- `goimports`: Import sorting and formatting (auto-installed)

```bash
# Format Go files
make format/go
```

**Behavior**:

- If `goimports` is available: Formats code + sorts imports
- If `goimports` is missing: Uses `gofmt` only + shows info message

#### `make format/yaml`

Format YAML files with graceful degradation.

**Tools Used** (in order of preference):

1. `prettier`: Node.js-based formatter (auto-installed via npm)
2. `yamlfmt`: Go-based YAML formatter (auto-installed)

```bash
# Format YAML files
make format/yaml
```

**Behavior**:

- Tries `prettier` first (better formatting quality)
- Falls back to `yamlfmt` if prettier unavailable
- Shows warning if neither tool is available
- Formats: `*.yaml`, `*.yml`, `.github/workflows/*.yml`

#### `make format/json`

Format JSON files with graceful degradation.

**Tools Used** (in order of preference):

1. `prettier`: Node.js-based formatter (auto-installed via npm)
2. `jq`: Command-line JSON processor (requires manual installation)

```bash
# Format JSON files
make format/json
```

**Behavior**:

- Tries `prettier` first (preserves formatting style)
- Falls back to `jq` if prettier unavailable
- Shows warning if neither tool is available
- Formats: `*.json`, `**/*.json`

**Note**: `jq` requires manual installation:

```bash
# Ubuntu/Debian
sudo apt-get install jq

# macOS
brew install jq

# Alpine
apk add jq
```

#### `make format/md`

Format Markdown files.

**Tools Used**:

- `prettier`: Node.js-based formatter (auto-installed via npm)

```bash
# Format Markdown files
make format/md
```

**Behavior**:

- Formats all `*.md` files
- Shows warning if prettier unavailable
- Preserves code blocks and formatting

#### `make format-check`

Check if all files are properly formatted (CI-friendly).

**Exit Codes**:

- `0`: All files properly formatted
- `1`: Some files need formatting

```bash
# Check formatting
make format-check

# Use in CI pipeline
make format-check || exit 1
```

**Checks**:

- Go files: `gofmt -s -l`
- YAML/JSON/Markdown: `prettier --check` (if available)

**Output Example**:

```
Checking all file formatting...
✓ All files are properly formatted

# Or if issues found:
The following Go files need formatting:
internal/agent/base.go
internal/llm/providers.go
Run 'make format' to fix formatting issues
```

### Docker Targets

**Note**: These Makefile targets replace the old `docker/*.sh` shell scripts for
better maintainability and consistency.

#### `make docker-build`

Build Docker image with current version tag.

**Variables**:

- `DOCKER_REGISTRY`: Docker registry URL (optional)
- `DOCKER_IMAGE`: Image name (default: buildbureau)
- `DOCKER_TAG`: Image tag (default: current git version)

```bash
# Build with default settings
make docker-build

# Build with custom registry and tag
make docker-build DOCKER_REGISTRY=ghcr.io/kpango DOCKER_TAG=latest
```

**Features**:

- Multi-stage build for minimal image size (~50MB)
- Includes version metadata
- Uses build cache for faster builds

#### `make docker-run`

Run Docker container in daemon mode (background).

**Replaced**: `docker/run.sh` script

**Environment Variables**:

- `GEMINI_API_KEY`: Google Gemini API key
- `OPENAI_API_KEY`: OpenAI API key
- `CLAUDE_API_KEY`: Anthropic Claude API key
- `CONTAINER_NAME`: Container name (default: buildbureau)

```bash
# Run with Gemini
export GEMINI_API_KEY="your-key"
make docker-run

# Run with custom container name
CONTAINER_NAME=myapp make docker-run

# Run with multiple API keys
export GEMINI_API_KEY="key1"
export OPENAI_API_KEY="key2"
make docker-run
```

**Features**:

- Runs as daemon (background process)
- Mounts `./data` directory for persistence
- Auto-removes old containers with same name
- Shows container status and logs location
- Displays warning if no API keys provided

**View Logs**:

```bash
docker logs -f buildbureau
```

**Stop Container**:

```bash
docker stop buildbureau
```

#### `make docker-run-interactive`

Run Docker container in interactive mode (foreground).

```bash
# Run interactively
export GEMINI_API_KEY="your-key"
make docker-run-interactive
```

**Features**:

- Interactive terminal (`-it`)
- Removes container on exit (`--rm`)
- Mounts `./data` directory
- Passes all API keys as environment variables

**Use Cases**:

- Debugging
- One-time tasks
- Development testing
- Interactive sessions

#### `make docker-test`

Comprehensive Docker setup testing.

**Replaced**: `docker/test.sh` script

```bash
# Test complete Docker setup
make docker-test
```

**Checks**:

1. ✅ Docker installation and version
2. ✅ Docker Compose installation
3. ✅ API key availability
4. ✅ Docker image build
5. ✅ Container startup
6. ✅ Application health
7. ✅ Container cleanup

**Example Output**:

```
=== BuildBureau Docker Test ===

Checking Docker...
✓ Docker found: Docker version 24.0.6

Checking Docker Compose...
✓ Docker Compose found

Checking API keys...
✓ API key(s) found

Building Docker image...
✓ Build successful

Starting test container...
✓ Container started

Testing application...
✓ Application healthy

Cleaning up...
✓ Cleanup complete

=== All tests passed! ===
```

**Use Cases**:

- Pre-deployment validation
- CI/CD pipeline checks
- Development environment setup
- Troubleshooting Docker issues

#### `make docker-build-multi`

Build multi-architecture Docker image (linux/amd64, linux/arm64).

```bash
make docker-build-multi
```

#### `make docker-push`

Push Docker image to registry.

```bash
make docker-push
```

#### `make docker-compose-up`

Start services using docker-compose.

```bash
make docker-compose-up
```

#### `make docker-compose-down`

Stop services using docker-compose.

```bash
make docker-compose-down
```

### Install Targets

BuildBureau includes an automated tool installation system that tracks installed
tools and prevents duplicate installations.

#### `make install-all`

Install all development tools at once.

**Installs**:

- All Go development tools (protoc-gen-go, protoc-gen-go-grpc, goimports)
- All formatting tools (prettier, yamlfmt, jq guidance)
- All security tools (gosec, govulncheck)
- golangci-lint

```bash
# Install everything
make install-all
```

**Features**:

- One-command setup for new developers
- Idempotent (safe to run multiple times)
- Uses stamp file tracking (see Auto-Install System)
- Skips already-installed tools
- Perfect for CI/CD and onboarding

**Use Cases**:

- New developer setup
- CI/CD environment initialization
- Container image preparation
- After `make clean-stamps`

#### `make install-tools`

Install Go development tools only.

**Installs**:

- `protoc-gen-go`: Protocol buffer Go code generator
- `protoc-gen-go-grpc`: gRPC Go code generator
- `goimports`: Import formatter and organizer

```bash
# Install Go tools
make install-tools
```

**When to Use**:

- Setting up Go development environment
- After Go version upgrade
- Building proto files for first time

#### `make install-formatters`

Install formatting tools.

**Installs**:

- `prettier`: YAML, JSON, Markdown formatter (via npm)
- `yamlfmt`: YAML formatter (Go-based)
- Provides `jq` installation instructions (system package)

```bash
# Install formatters
make install-formatters
```

**Requirements**:

- Node.js and npm (for prettier)
- Go 1.21+ (for yamlfmt)

**Note**: `jq` requires manual system installation:

```bash
# Ubuntu/Debian
sudo apt-get install jq

# macOS
brew install jq

# Alpine
apk add jq
```

#### `make install-security-tools`

Install security scanning tools.

**Installs**:

- `gosec`: Go security checker (static analysis)
- `govulncheck`: Go vulnerability scanner

```bash
# Install security tools
make install-security-tools
```

**When to Use**:

- Setting up security scanning
- Before running `make security`
- CI/CD security pipeline setup

#### `make clean-stamps`

Remove stamp files to force tool reinstallation.

```bash
# Force reinstall all tools
make clean-stamps
make install-all
```

**Use Cases**:

- Tool upgrade after version change
- Fixing corrupted installations
- Testing installation scripts
- Switching Go versions

**What It Does**:

- Removes `.make/*.stamp` files
- Next `make install-*` will reinstall tools
- Does not uninstall tools (just clears tracking)

### Development Targets

#### `make proto`

Generate protobuf files from .proto definitions.

```bash
make proto
```

#### `make deps`

Download and tidy dependencies.

```bash
make deps
```

#### `make deps-update`

Update all dependencies to latest versions.

```bash
make deps-update
```

#### `make fmt`

Format Go code using `go fmt`.

```bash
make fmt
```

#### `make fmt-check`

Check if code is properly formatted (useful for CI).

```bash
make fmt-check
```

#### `make lint`

Run `go vet` linter.

```bash
make lint
```

#### `make lint-all`

Run all available linters (go vet + golangci-lint if installed).

```bash
make lint-all
```

### CI/CD Targets

These targets are specifically designed for CI/CD pipelines.

#### `make ci-all`

Run all CI checks (lint, build, test).

```bash
make ci-all
```

#### `make ci-lint`

Run linters for CI (includes format checking).

```bash
make ci-lint
```

#### `make ci-build`

Build for CI (proto + release build).

```bash
make ci-build
```

#### `make ci-test`

Run tests with coverage for CI.

```bash
make ci-test
```

### Security Targets

#### `make security`

Run all security checks.

```bash
make security
```

#### `make security-scan`

Run security scanner (gosec) if installed.

```bash
make security-scan
```

#### `make security-deps`

Check for vulnerable dependencies (govulncheck) if installed.

```bash
make security-deps
```

### Clean Targets

#### `make clean`

Remove build artifacts.

```bash
make clean
```

#### `make clean-all`

Remove all generated files (build, coverage, dist).

```bash
make clean-all
```

#### `make clean-coverage`

Remove coverage reports.

```bash
make clean-coverage
```

#### `make clean-cache`

Clean Go build cache.

```bash
make clean-cache
```

## Auto-Install System

BuildBureau features an intelligent auto-install system that automatically
manages development tool dependencies. This system eliminates manual tool
installation while preventing duplicate installations.

### How It Works

The auto-install system uses **stamp files** to track installed tools:

```
.make/
├── goimports.stamp
├── protoc-gen-go.stamp
├── protoc-gen-go-grpc.stamp
├── prettier.stamp
├── yamlfmt.stamp
├── jq.stamp
├── golangci-lint.stamp
├── gosec.stamp
└── govulncheck.stamp
```

**Mechanism**:

1. When a target requires a tool (e.g., `format/go` needs `goimports`)
2. Makefile checks if stamp file exists (`.make/goimports.stamp`)
3. If missing: Installs tool + creates stamp file
4. If exists: Skips installation, uses existing tool

**Benefits**:

- ✅ **Zero manual setup** - Tools install automatically when needed
- ✅ **No duplicate installs** - Stamp files prevent redundant installations
- ✅ **Fast rebuilds** - Skips installation if already done
- ✅ **CI/CD friendly** - Deterministic, cacheable installations
- ✅ **Developer friendly** - Just run `make format`, tools auto-install
- ✅ **Transparent** - Shows what's being installed and why

### Auto-Install in Action

#### First Run (Tools Missing)

```bash
$ make format/go

Installing goimports...
go install golang.org/x/tools/cmd/goimports@latest
✓ goimports installed

Formatting Go files...
✓ Go files formatted (gofmt + goimports)
```

#### Second Run (Tools Present)

```bash
$ make format/go

Formatting Go files...
✓ Go files formatted (gofmt + goimports)
```

_Notice: No installation, just formatting!_

### Stamp File Tracking

Stamp files are created in `.make/` directory:

```bash
# View installed tools
$ ls -la .make/
drwxr-xr-x  2 user user 4096 Feb  1 12:00 .
drwxr-xr-x 12 user user 4096 Feb  1 11:59 ..
-rw-r--r--  1 user user    0 Feb  1 12:00 goimports.stamp
-rw-r--r--  1 user user    0 Feb  1 12:00 prettier.stamp
-rw-r--r--  1 user user    0 Feb  1 12:00 yamlfmt.stamp

# Check if tool is tracked
$ test -f .make/goimports.stamp && echo "Installed" || echo "Not installed"
Installed
```

**Note**: Stamp files are empty marker files. Their presence indicates
installation, not contents.

### Workflow Examples

#### Developer Onboarding

New developer clones repo and runs:

```bash
# Clone repository
git clone https://github.com/yourusername/BuildBureau.git
cd BuildBureau

# First build - auto-installs all needed tools
make build

# Output:
Installing protoc-gen-go...
✓ protoc-gen-go installed

Installing protoc-gen-go-grpc...
✓ protoc-gen-go-grpc installed

Generating proto files...
Building application...
✓ Build complete
```

No separate installation step needed!

#### CI/CD Pipeline

```yaml
# .github/workflows/ci.yml
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5

      # Option 1: Let auto-install handle it
      - name: Run tests
        run: make test

      # Option 2: Explicit install for caching
      - name: Install tools
        run: make install-all

      - name: Cache tools
        uses: actions/cache@v3
        with:
          path: .make
          key: tools-${{ hashFiles('.make/*.stamp') }}

      - name: Run tests
        run: make test
```

**Cache Efficiency**: Stamp files enable effective tool caching in CI.

#### Local Development

```bash
# Developer wants to format code
$ make format

# Auto-install happens transparently:
Installing prettier...
✓ prettier installed

Installing yamlfmt...
✓ yamlfmt installed

Formatting Go files...
✓ Go files formatted

Formatting YAML files...
✓ YAML files formatted

Formatting JSON files...
✓ JSON files formatted

Formatting Markdown files...
✓ Markdown files formatted

✓ All files formatted

# Next time - instant:
$ make format
✓ All files formatted
```

#### Upgrading Tools

```bash
# Force reinstall all tools (e.g., after Go upgrade)
make clean-stamps
make install-all

# Or just specific tool
rm .make/goimports.stamp
make format/go  # Will reinstall goimports
```

### Graceful Degradation

Some formatters degrade gracefully when tools are unavailable:

**Example: YAML Formatting**

```bash
# Scenario 1: Both prettier and yamlfmt available
$ make format/yaml
✓ YAML files formatted (prettier)

# Scenario 2: Only yamlfmt available
$ make format/yaml
✓ YAML files formatted (yamlfmt)

# Scenario 3: Neither available
$ make format/yaml
⚠ No YAML formatter found (install prettier or yamlfmt)
```

**Tools with Graceful Degradation**:

- `format/go`: Falls back from goimports to gofmt
- `format/yaml`: Tries prettier → yamlfmt → warning
- `format/json`: Tries prettier → jq → warning
- `format/md`: Uses prettier or shows warning

### Manual Tool Installation

For system tools like `jq`:

```bash
# Ubuntu/Debian
sudo apt-get install jq

# macOS
brew install jq

# Alpine Linux
apk add jq

# Then create stamp manually if desired
touch .make/jq.stamp
```

### Troubleshooting Auto-Install

#### Installation Failed

```bash
# Check if tool is in PATH
which goimports

# Verify Go environment
go env GOPATH

# Ensure GOPATH/bin is in PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Retry installation
make clean-stamps
make install-tools
```

#### Stamp File Out of Sync

```bash
# Tool installed but stamp missing
which goimports  # Tool exists
test -f .make/goimports.stamp  # Stamp missing

# Recreate stamp
touch .make/goimports.stamp

# Or let auto-install handle it
make clean-stamps
make install-all
```

#### CI Cache Issues

```yaml
# GitHub Actions - robust caching
- name: Cache Go modules
  uses: actions/cache@v3
  with:
    path: ~/go/pkg/mod
    key: go-${{ hashFiles('go.sum') }}

- name: Cache installed tools
  uses: actions/cache@v3
  with:
    path: |
      ~/go/bin
      .make
    key: tools-${{ runner.os }}-${{ hashFiles('Makefile') }}

- name: Install tools
  run: make install-all
```

### Best Practices

1. **Let Auto-Install Work**

   ```bash
   # Don't: manually install then run
   go install golang.org/x/tools/cmd/goimports@latest
   make format/go

   # Do: let auto-install handle it
   make format/go
   ```

2. **Cache Stamp Directory in CI**

   ```yaml
   - uses: actions/cache@v3
     with:
       path: .make
       key: stamps-${{ runner.os }}
   ```

3. **Clean Stamps After Tool Upgrades**

   ```bash
   # After upgrading Go or tools
   make clean-stamps
   make install-all
   ```

4. **Check Installation Status**

   ```bash
   # See what's installed
   ls .make/*.stamp

   # Count installed tools
   ls .make/*.stamp 2>/dev/null | wc -l
   ```

5. **Commit .gitignore for .make/**
   ```gitignore
   # .gitignore
   .make/
   ```

### Summary

The auto-install system provides:

- ✅ Zero-configuration tool management
- ✅ Automatic installation on demand
- ✅ Duplicate prevention via stamp files
- ✅ CI/CD caching support
- ✅ Graceful degradation
- ✅ Manual override capability
- ✅ Transparent operation

**Philosophy**: "Just run `make <target>`, we'll handle the tools."

### Utility Targets

#### `make version`

Display version information.

```bash
make version
```

#### `make help`

Display help message with all available targets.

```bash
make help
```

#### `make check`

Run all checks (Go installation, dependencies, protoc).

```bash
make check
```

## Variables

### Build Variables

- `VERSION`: Version string (default: git describe)
- `BUILD_DIR`: Build output directory (default: ./build)
- `DIST_DIR`: Distribution directory (default: ./dist)
- `CGO_ENABLED`: Enable CGO (default: 1, required for SQLite)

### Test Variables

- `TEST_TIMEOUT`: Test timeout duration (default: 10m)
- `TEST_FLAGS`: Additional test flags (default: -v -race -count=1)

### Docker Variables

- `DOCKER_REGISTRY`: Docker registry URL
- `DOCKER_IMAGE`: Image name (default: buildbureau)
- `DOCKER_TAG`: Image tag (default: VERSION)
- `DOCKER_PLATFORMS`: Build platforms (default: linux/amd64,linux/arm64)

## Usage in Different Contexts

### Local Development

```bash
# Daily development workflow
make deps          # Install dependencies
make proto         # Generate proto files
make build         # Build binary
make test          # Run tests
make run           # Run application

# Before committing
make format        # Format all files (Go, YAML, JSON, Markdown)
make lint          # Run linters
make test-coverage # Check coverage

# Alternative: format specific file types
make format/go     # Format only Go files
make format/yaml   # Format only YAML files

# Check formatting without changes (CI-friendly)
make format-check
```

### Docker Build

The Dockerfile uses Makefile targets for standardization:

```dockerfile
# In Dockerfile
RUN make proto
RUN make build-release
```

Build and run Docker:

```bash
# Build image
make docker-build

# Test Docker setup
make docker-test

# Run in background
export GEMINI_API_KEY="your-key"
make docker-run

# Or run interactively
make docker-run-interactive

# Push to registry
make docker-push
```

### CI/CD Pipeline

The CI workflow uses Makefile targets:

```yaml
# In .github/workflows/ci.yml
- name: Install dependencies
  run: make deps

- name: Run linters
  run: make ci-lint

- name: Build
  run: make ci-build

- name: Run tests
  run: make ci-test
```

### Release Process

```bash
# Build for all platforms
make release-build

# Package releases
make release-package

# Build and push Docker images
make docker-build-multi
make docker-push
```

## Advanced Usage

### Custom Build Configuration

```bash
# Build for specific platform
GOOS=linux GOARCH=arm64 make build

# Build with custom version
VERSION=v1.0.0 make build-release

# Build with custom registry
DOCKER_REGISTRY=ghcr.io/myorg make docker-build
```

### Parallel Builds

```bash
# Use make's parallel execution
make -j4 build-all
```

### Debug Mode

```bash
# Build and run in debug mode
make build-debug
make run-debug
```

## Tool Requirements

### Required Tools

- Go 1.24+
- make
- git

### Optional Tools

**Most tools auto-install when needed!** The auto-install system handles:

- `protoc-gen-go`: Auto-installed on first `make proto`
- `protoc-gen-go-grpc`: Auto-installed on first `make proto`
- `goimports`: Auto-installed on first `make format/go`
- `prettier`: Auto-installed on first `make format` (requires npm)
- `yamlfmt`: Auto-installed on first `make format/yaml`
- `golangci-lint`: Auto-installed on first `make lint-all`
- `gosec`: Auto-installed on first `make security`
- `govulncheck`: Auto-installed on first `make security-deps`

**Manual installation required**:

- `protoc`: Protocol buffer compiler
  ([installation guide](https://grpc.io/docs/protoc-installation/))
- `jq`: JSON processor (system package: `apt install jq` or `brew install jq`)
- `hadolint`: Dockerfile linter
- `npm`: Required for prettier (formatting)

Install all tools at once:

```bash
# Install everything auto-installable
make install-all

# Or install by category
make install-tools           # Go development tools
make install-formatters      # Formatting tools
make install-security-tools  # Security scanners
```

## Troubleshooting

### "protoc: No such file or directory"

Install protoc: https://grpc.io/docs/protoc-installation/

Or skip proto generation if not needed:

```bash
make build  # Will use existing proto files
```

### "golangci-lint not installed, skipping"

This is a warning, not an error. The basic `go vet` linter still runs.

Install golangci-lint: https://golangci-lint.run/usage/install/

### CGO Errors

SQLite requires CGO. Ensure you have gcc/build tools installed:

```bash
# Ubuntu/Debian
sudo apt-get install build-essential

# macOS (via Xcode)
xcode-select --install

# Alpine
apk add gcc musl-dev
```

### Test Failures in CI

Run the same tests locally:

```bash
make ci-test
```

## Best Practices

### 1. Format All Files Before Committing

```bash
# Format everything (Go, YAML, JSON, Markdown)
make format

# Or check formatting in CI
make format-check
```

**Why**: Ensures consistent code style across all file types, not just Go.

### 2. Use Auto-Install System

```bash
# Don't manually install tools
# ❌ go install golang.org/x/tools/cmd/goimports@latest

# Let Makefile handle it
# ✅ make format/go
```

**Why**: Auto-install prevents version mismatches and simplifies onboarding.

### 3. Run Tests Before Committing

```bash
make test
```

### 4. Check Coverage Regularly

```bash
make test-coverage-html
# Open ./coverage/coverage.html
```

### 5. Use CI Targets in CI

Use `ci-*` targets in CI pipelines for consistency:

```bash
make ci-all
```

### 6. Test Docker Setup Early

```bash
# Validate Docker configuration
make docker-test

# Then run
make docker-run
```

**Why**: Catches Docker issues before deployment.

### 7. Build Release Binaries for Production

```bash
make build-release
```

### 8. Keep Tools Updated

```bash
# After Go version upgrade or tool updates
make clean-stamps
make install-all
```

### 6. Keep Dependencies Updated

```bash
make deps-update
make test  # Verify everything still works
```

## Integration Examples

### GitHub Actions

```yaml
name: CI
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.24"

      - name: CI Pipeline
        run: make ci-all
```

### GitLab CI

```yaml
test:
  image: golang:1.24
  script:
    - make ci-all
```

### Jenkins

```groovy
pipeline {
    agent any
    stages {
        stage('Test') {
            steps {
                sh 'make ci-all'
            }
        }
        stage('Build') {
            steps {
                sh 'make build-all'
            }
        }
    }
}
```

## Contributing

When adding new Makefile targets:

1. Add clear comments
2. Use consistent naming (verb-noun format)
3. Add to appropriate .PHONY declaration
4. Include ## comment for help text
5. Use color output for better UX
6. Test in CI context

## Summary

The BuildBureau Makefile provides:

- ✅ **Standardized commands** across environments
- ✅ **Comprehensive build options** (debug, release, static, multi-platform)
- ✅ **Multiple test configurations** (unit, coverage, race, LLM integration)
- ✅ **Multi-format formatting** (Go, YAML, JSON, Markdown)
- ✅ **Auto-install system** for zero-configuration tool management
- ✅ **Enhanced Docker support** (build, run, test, interactive mode)
- ✅ **CI/CD friendly targets** with format checking and caching
- ✅ **Security scanning** (gosec, govulncheck)
- ✅ **Development tools** with stamp-based tracking
- ✅ **Graceful degradation** when optional tools are missing
- ✅ **Clear documentation** and colored output

**New in Latest Version**:

- 🎉 Auto-install system with stamp file tracking
- 🎉 Multi-format support (format/go, format/yaml, format/json, format/md)
- 🎉 Enhanced Docker targets (docker-run, docker-test, docker-run-interactive)
- 🎉 LLM integration testing (test/llm-integration)
- 🎉 Tool installation targets (install-all, install-tools, install-formatters)
- 🎉 Replaced shell scripts with Makefile targets for consistency

Use `make help` to see all available targets at any time.
