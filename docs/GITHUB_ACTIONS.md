# GitHub Actions CI/CD Documentation

This document provides comprehensive documentation for the GitHub Actions
workflows in BuildBureau. Our CI/CD pipeline ensures code quality, security, and
automated releases across multiple platforms.

## Table of Contents

- [Overview](#overview)
- [Workflows](#workflows)
  - [CI Workflow](#ci-workflow)
  - [Release Workflow](#release-workflow)
  - [Docker Publish Workflow](#docker-publish-workflow)
  - [CodeQL Analysis Workflow](#codeql-analysis-workflow)
  - [GolangCI-Lint Workflow](#golangci-lint-workflow)
  - [Coverage Workflow](#coverage-workflow)
  - [Dependency Review Workflow](#dependency-review-workflow)
  - [Release Drafter Workflow](#release-drafter-workflow)
  - [Stale Workflow](#stale-workflow)
  - [Auto Label Workflow](#auto-label-workflow)
- [Setup Guide](#setup-guide)
- [Usage Examples](#usage-examples)
- [Troubleshooting](#troubleshooting)
- [Best Practices](#best-practices)

## Overview

BuildBureau uses a comprehensive GitHub Actions CI/CD pipeline that includes:

- ✅ **Continuous Integration**: Automated builds and tests on every push and PR
- 🚀 **Multi-Platform Releases**: Automated releases for Linux (amd64/arm64) and
  macOS (Intel/Apple Silicon)
- 🐳 **Docker Publishing**: Multi-architecture Docker images to GitHub Container
  Registry
- 🔒 **Security Scanning**: CodeQL analysis, dependency review, and
  vulnerability checks
- 📊 **Code Quality**: 30+ linters via golangci-lint with automated checks
- 📈 **Code Coverage**: Automated coverage reports with Codecov integration
- 🏷️ **Auto-Labeling**: Intelligent PR and issue labeling
- 📝 **Release Notes**: Automatic changelog generation

### Workflow Triggers

| Workflow          | Push (main) | Push (develop) | Push (copilot/\*) | PR  | Tags     | Schedule    | Manual |
| ----------------- | ----------- | -------------- | ----------------- | --- | -------- | ----------- | ------ |
| CI                | ✅          | ❌             | ✅                | ✅  | ❌       | ❌          | ❌     |
| Release           | ❌          | ❌             | ❌                | ❌  | ✅ (v\*) | ❌          | ❌     |
| Docker Publish    | ✅          | ❌             | ❌                | ✅  | ✅ (v\*) | ❌          | ✅     |
| CodeQL Analysis   | ✅          | ✅             | ❌                | ✅  | ❌       | ✅ (Weekly) | ✅     |
| GolangCI-Lint     | ✅          | ✅             | ❌                | ✅  | ❌       | ❌          | ❌     |
| Coverage          | ✅          | ✅             | ❌                | ✅  | ❌       | ❌          | ❌     |
| Dependency Review | ❌          | ❌             | ❌                | ✅  | ❌       | ❌          | ❌     |
| Release Drafter   | ✅          | ❌             | ❌                | ✅  | ❌       | ❌          | ❌     |
| Stale             | ❌          | ❌             | ❌                | ❌  | ❌       | ✅ (Daily)  | ✅     |
| Auto Label        | ❌          | ❌             | ❌                | ✅  | ❌       | ❌          | ❌     |

## Workflows

### CI Workflow

**File**: `.github/workflows/ci.yml`

The main continuous integration workflow that runs on every push and pull
request.

#### Features

- ✅ **Multi-version Go testing**: Tests with Go 1.24
- ✅ **Comprehensive checks**: Formatting, linting, building, and testing
- ✅ **Makefile integration**: Uses standardized Makefile targets
- ✅ **Proto generation**: Generates protobuf code before building
- ✅ **Example tests**: Validates example code works correctly

#### Steps

1. **Checkout code**: Uses `actions/checkout@v4` with full history
2. **Set up Go**: Installs Go 1.24 with caching enabled
3. **Install dependencies**: Runs `make deps` to install all required tools
4. **Generate proto files**: Runs `make proto` (optional, continues on failure)
5. **Check formatting**: Runs `make fmt-check` to verify code is formatted
6. **Run linters**: Executes `make ci-lint` for code quality checks
7. **Build**: Compiles the project with `make ci-build`
8. **Run tests**: Executes `make ci-test` for unit tests
9. **Run example test**: Tests example code with demo API keys

#### Triggers

```yaml
on:
  push:
    branches: [main, copilot/*]
  pull_request:
    branches: [main]
```

#### Permissions

```yaml
permissions:
  contents: read
```

#### Usage

The CI workflow runs automatically. To test locally:

```bash
# Run all CI checks locally
make ci-all

# Or run individual checks
make fmt-check
make ci-lint
make ci-build
make ci-test
```

### Release Workflow

**File**: `.github/workflows/release.yml`

Automated multi-platform binary releases triggered by version tags.

#### Features

- 🚀 **Multi-platform builds**: Linux (amd64/arm64), macOS (Intel/Apple Silicon)
- 📦 **Compressed archives**: tar.gz with SHA256 checksums
- 📝 **Auto changelog**: Generates changelog from git commits
- 🏷️ **Version injection**: Embeds version, build time, and git commit in
  binaries
- 📋 **Installation instructions**: Includes platform-specific install commands
- 🔒 **CGo support**: Enables CGo for SQLite compilation

#### Build Matrix

| Platform              | OS Runner     | GOOS   | GOARCH | Name         |
| --------------------- | ------------- | ------ | ------ | ------------ |
| Linux (amd64)         | ubuntu-latest | linux  | amd64  | linux-amd64  |
| Linux (arm64)         | ubuntu-latest | linux  | arm64  | linux-arm64  |
| macOS (Intel)         | macos-latest  | darwin | amd64  | darwin-amd64 |
| macOS (Apple Silicon) | macos-latest  | darwin | arm64  | darwin-arm64 |

#### Build Flags

```bash
-ldflags "-w -s \
  -X main.version=${VERSION} \
  -X main.buildTime=$(date -u '+%Y-%m-%d_%H:%M:%S') \
  -X main.gitCommit=$(git rev-parse --short HEAD)"
```

- `-w`: Omit DWARF symbol table
- `-s`: Omit symbol table and debug information
- `-X`: Set variable values at link time

#### Triggers

```yaml
on:
  push:
    tags:
      - "v*" # Triggers on any tag starting with 'v'
```

#### Permissions

```yaml
permissions:
  contents: write # Required for creating releases
  packages: write # Required for uploading artifacts
```

#### Usage

**Create a release:**

```bash
# Create and push a version tag
git tag v1.0.0
git push origin v1.0.0

# Or create a tag with annotation
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

**Artifacts produced:**

- `buildbureau-linux-amd64.tar.gz` + `.sha256`
- `buildbureau-linux-arm64.tar.gz` + `.sha256`
- `buildbureau-darwin-amd64.tar.gz` + `.sha256`
- `buildbureau-darwin-arm64.tar.gz` + `.sha256`

### Docker Publish Workflow

**File**: `.github/workflows/docker-publish.yml`

Publishes multi-architecture Docker images to GitHub Container Registry (GHCR).

#### Features

- 🐳 **Multi-architecture**: Builds for linux/amd64 and linux/arm64
- 📦 **GHCR publishing**: Pushes to `ghcr.io/kpango/buildbureau`
- 🏷️ **Smart tagging**: Automatic semantic versioning tags
- 💾 **Build caching**: GitHub Actions cache for faster builds
- 🔒 **Security scanning**: Trivy vulnerability scanner + SBOM generation
- 📊 **SARIF upload**: Security findings uploaded to GitHub Security

#### Image Tags

The workflow generates multiple tags:

- `main` - Latest main branch
- `latest` - Latest stable release (main branch only)
- `pr-123` - Pull request number
- `1.0.0` - Semantic version
- `1.0` - Major.minor version
- `1` - Major version
- `main-abc1234` - Branch with commit SHA

#### Triggers

```yaml
on:
  push:
    branches: [main]
    tags: ["v*"]
  pull_request:
    branches: [main]
  workflow_dispatch:
```

#### Permissions

```yaml
permissions:
  contents: read
  packages: write # Required for GHCR push
  id-token: write # Required for OIDC
```

#### Security Features

1. **Trivy Vulnerability Scanning**
   - Scans built images for CVEs
   - Uploads results to GitHub Security tab
   - SARIF format for integration

2. **SBOM Generation**
   - Software Bill of Materials in SPDX JSON format
   - Lists all dependencies and versions
   - Retained for 30 days

#### Usage

**Pull the image:**

```bash
# Pull latest
docker pull ghcr.io/kpango/buildbureau:latest

# Pull specific version
docker pull ghcr.io/kpango/buildbureau:v1.0.0

# Pull for specific architecture
docker pull --platform linux/arm64 ghcr.io/kpango/buildbureau:latest
```

**Run the container:**

```bash
docker run -d \
  --name buildbureau \
  -e GEMINI_API_KEY="your-key" \
  -v buildbureau-data:/app/data \
  -p 8080:8080 \
  ghcr.io/kpango/buildbureau:latest
```

### CodeQL Analysis Workflow

**File**: `.github/workflows/codeql-analysis.yml`

Security scanning using GitHub's CodeQL for identifying vulnerabilities.

#### Features

- 🔒 **Security scanning**: Detects security vulnerabilities in Go code
- 📊 **Quality analysis**: Identifies code quality issues
- 🔍 **Deep analysis**: Uses `security-and-quality` query suite
- 📅 **Scheduled scans**: Weekly security checks
- 🚨 **Security alerts**: Integrates with GitHub Security tab

#### Query Suite

The workflow uses the `security-and-quality` query suite which includes:

- **Security queries**: Detects common vulnerabilities (SQL injection, XSS,
  etc.)
- **Quality queries**: Identifies code smells and maintainability issues
- **Best practices**: Enforces Go best practices

#### Triggers

```yaml
on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]
  schedule:
    - cron: "0 6 * * 1" # Weekly on Monday at 6 AM UTC
  workflow_dispatch:
```

#### Permissions

```yaml
permissions:
  actions: read
  contents: read
  security-events: write # Required for uploading SARIF
```

#### Languages Analyzed

- Go (primary language)

#### Usage

**View security alerts:**

1. Navigate to: **Security → Code scanning alerts**
2. Filter by severity, state, or rule
3. Click on an alert to see details and recommendations

**Run manually:**

```bash
# Trigger manual workflow run
gh workflow run codeql-analysis.yml
```

**Check results:**

```bash
# List security alerts
gh api repos/:owner/:repo/code-scanning/alerts
```

### GolangCI-Lint Workflow

**File**: `.github/workflows/golangci-lint.yml`

Comprehensive Go code linting with 30+ linters.

#### Features

- 🔍 **30+ linters**: Comprehensive code quality checks
- ⚡ **Fast execution**: Uses caching for quick runs
- 📋 **Custom config**: Uses `.golangci.json` configuration
- 🚨 **PR integration**: Comments on PRs with findings
- ✅ **Format checking**: Verifies code is properly formatted

#### Linters Included

The workflow uses golangci-lint with these categories:

**Bugs**:

- errcheck, govet, staticcheck, typecheck

**Style**:

- gofmt, goimports, gofumpt, whitespace

**Complexity**:

- gocyclo, gocognit, nestif

**Performance**:

- prealloc, bodyclose

**Security**:

- gosec, gas

**And many more** (see `.golangci.json` for full configuration)

#### Triggers

```yaml
on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]
```

#### Permissions

```yaml
permissions:
  contents: read
  pull-requests: read
  checks: write # Required for check runs
```

#### Configuration

The workflow uses a 10-minute timeout and runs against all files:

```yaml
args: --timeout=10m --config=.golangci.json
```

#### Usage

**Run locally:**

```bash
# Run all linters
make lint

# Or use golangci-lint directly
golangci-lint run --timeout=10m

# Auto-fix issues
golangci-lint run --fix
```

**View results in PR:**

- Linting results appear as PR checks
- Annotations show line-by-line issues
- Failed checks block merging (if configured)

### Coverage Workflow

**File**: `.github/workflows/coverage.yml`

Automated code coverage analysis with Codecov integration.

#### Features

- 📊 **Coverage reporting**: Tracks test coverage percentage
- 📈 **Codecov integration**: Uploads coverage to Codecov.io
- 💬 **PR comments**: Posts coverage reports on PRs
- 📋 **Step summaries**: Shows coverage in workflow summary
- 📦 **Artifact uploads**: Saves coverage files for 30 days

#### Coverage Metrics

The workflow generates:

1. **Coverage percentage**: Overall coverage (e.g., 67.3%)
2. **Package coverage**: Per-package breakdown
3. **Function coverage**: Line-by-line coverage data
4. **Trend analysis**: Coverage changes over time (Codecov)

#### Triggers

```yaml
on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]
```

#### Permissions

```yaml
permissions:
  contents: read
  pull-requests: write # Required for PR comments
```

#### Usage

**Run locally:**

```bash
# Run tests with coverage
make test-coverage

# View coverage in browser
make test-coverage
go tool cover -html=./coverage/coverage.out
```

**Codecov setup:**

Add `CODECOV_TOKEN` to repository secrets:

1. Go to [codecov.io](https://codecov.io)
2. Add your repository
3. Copy the upload token
4. Add to GitHub secrets: Settings → Secrets → Actions → `CODECOV_TOKEN`

**Coverage reports include:**

- Total coverage percentage
- Coverage by package
- Coverage changes in PR
- Uncovered lines highlighted

### Dependency Review Workflow

**File**: `.github/workflows/dependency-review.yml`

Security review of dependency changes in pull requests.

#### Features

- 🔒 **Vulnerability scanning**: Detects known CVEs in dependencies
- 📋 **License checking**: Blocks incompatible licenses
- 🚨 **govulncheck**: Go-specific vulnerability database checks
- 📊 **Nancy scanner**: Sonatype OSS Index integration
- 💬 **PR comments**: Automatic security notifications

#### Checks Performed

1. **Dependency Review Action**
   - Compares dependencies between base and PR
   - Checks GitHub Advisory Database
   - Blocks GPL-3.0 and AGPL-3.0 licenses
   - Fails on moderate+ severity vulnerabilities

2. **govulncheck**
   - Official Go vulnerability scanner
   - Checks direct and indirect dependencies
   - Uses Go vulnerability database

3. **Nancy Security Check**
   - Sonatype OSS Index integration
   - Additional vulnerability database
   - Continues on errors (advisory only)

#### Triggers

```yaml
on:
  pull_request:
    branches: [main]
```

#### Permissions

```yaml
permissions:
  contents: read
  pull-requests: write # Required for comments
```

#### Configuration

```yaml
fail-on-severity: moderate # Blocks moderate, high, critical
deny-licenses:
  - GPL-3.0
  - AGPL-3.0
comment-summary-in-pr: always
```

#### Usage

**Run locally:**

```bash
# Run govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# Check dependencies
go list -m all
```

**Review security alerts:**

- Check PR comments for vulnerability warnings
- Review dependency changes in PR files
- Update vulnerable dependencies before merging

### Release Drafter Workflow

**File**: `.github/workflows/release-drafter.yml`

Automatically generates and maintains release notes.

#### Features

- 📝 **Auto changelog**: Generates changelog from PRs
- 🏷️ **Auto categorization**: Groups changes by type
- 🔢 **Semantic versioning**: Auto-suggests next version
- 🏷️ **Auto labeling**: Labels PRs based on content
- 📋 **Draft releases**: Creates draft releases automatically

#### Categories

The workflow groups changes into:

- 🚀 **Features**: New features and enhancements
- 🐛 **Bug Fixes**: Bug fixes and corrections
- 🔒 **Security**: Security fixes and updates
- ⚡ **Performance**: Performance improvements
- 📚 **Documentation**: Documentation updates
- 🧪 **Testing**: Test additions and fixes
- 🏗️ **Build & CI/CD**: Build and CI/CD changes
- 🔧 **Maintenance**: Maintenance and refactoring
- 💔 **Breaking Changes**: Breaking changes

#### Version Resolution

The workflow auto-increments versions based on labels:

- **Major** (1.0.0 → 2.0.0): `breaking-change`, `major`
- **Minor** (1.0.0 → 1.1.0): `enhancement`, `feature`, `minor`
- **Patch** (1.0.0 → 1.0.1): `bug`, `fix`, `patch`, `chore`, `docs`

#### Triggers

```yaml
on:
  push:
    branches: [main]
  pull_request_target:
    types: [opened, synchronize, reopened, closed]
```

#### Permissions

```yaml
permissions:
  contents: write
  pull-requests: write
```

#### Usage

**View draft release:**

1. Navigate to: **Releases**
2. Find the draft release (not published)
3. Review generated changelog
4. Edit if needed and publish

**Configuration:**

Edit `.github/release-drafter.yml` to customize:

- Categories
- Version bumping rules
- Template format
- Auto-labeling rules

### Stale Workflow

**File**: `.github/workflows/stale.yml`

Automatically manages stale issues and pull requests.

#### Features

- 🗂️ **Auto-stale marking**: Labels inactive issues/PRs as stale
- ⏰ **Auto-closing**: Closes stale items after grace period
- 🏷️ **Exemptions**: Protects important labels from stale marking
- 💬 **Notifications**: Comments on stale items with clear messaging
- 🔄 **Reactivation**: Removes stale label when updated

#### Configuration

**Issues:**

- **Stale after**: 60 days of inactivity
- **Close after**: 14 more days
- **Exempted labels**: `pinned`, `security`, `help-wanted`, `good-first-issue`

**Pull Requests:**

- **Stale after**: 30 days of inactivity
- **Close after**: 7 more days
- **Exempted labels**: `pinned`, `security`, `work-in-progress`

#### Triggers

```yaml
on:
  schedule:
    - cron: "0 0 * * *" # Daily at midnight UTC
  workflow_dispatch:
```

#### Permissions

```yaml
permissions:
  issues: write
  pull-requests: write
```

#### Messages

**Stale issue message:**

```
This issue has been automatically marked as stale because it has not had
recent activity. It will be closed if no further activity occurs within
14 days. Thank you for your contributions.
```

**Stale PR message:**

```
This pull request has been automatically marked as stale because it has
not had recent activity. It will be closed if no further activity occurs
within 7 days. Thank you for your contributions.
```

#### Usage

**Prevent stale marking:**

- Add `pinned` label to keep indefinitely
- Add `security` label for security issues
- Add `work-in-progress` label for PRs in development
- Comment on issue/PR to reset timer

**Run manually:**

```bash
# Trigger manual stale check
gh workflow run stale.yml
```

### Auto Label Workflow

**File**: `.github/workflows/auto-label.yml`

Automatically labels pull requests and issues based on content.

#### Features

- 🏷️ **File-based labeling**: Labels based on changed files
- 📏 **Size labeling**: Labels PRs by change size
- 🔍 **Content analysis**: Labels based on title/body keywords
- 📚 **Documentation detection**: Auto-labels documentation changes
- 💥 **Breaking change detection**: Identifies breaking changes

#### PR Labeling

**File-based labels:**

- `go` - Go code changes (\*.go, go.mod, go.sum)
- `documentation` - Documentation changes (\*.md, docs/)
- `ci/cd` - CI/CD changes (.github/, Makefile)
- `docker` - Docker changes (Dockerfile, docker-compose.yml)
- `configuration` - Config changes (_.yaml, _.yml, \*.json)
- `internal` - Internal package changes
- `pkg` - Public API changes
- `agents` - Agent changes
- `cli` - CLI changes
- `testing` - Test changes (\*\_test.go)

**Size labels:**

- `size/xs` - 1-10 lines
- `size/s` - 11-100 lines
- `size/m` - 101-500 lines
- `size/l` - 501-1000 lines
- `size/xl` - 1001+ lines

**Special labels:**

- `documentation` - Title/body contains "doc" or "docs"
- `breaking-change` - Title/body contains "BREAKING"

#### Issue Labeling

**Keyword-based labels:**

- `bug` - Title/body contains "bug"
- `enhancement` - Title/body contains "feature"
- `question` - Title/body contains "question"
- `documentation` - Title/body contains "doc"
- `performance` - Title/body contains "performance"
- `security` - Title/body contains "security"
- `testing` - Title/body contains "test"
- `ci/cd` - Title/body contains "ci" or "cd"

#### Triggers

```yaml
on:
  pull_request:
    types: [opened, synchronize, reopened]
  issues:
    types: [opened, edited]
```

#### Permissions

```yaml
permissions:
  contents: read
  issues: write
  pull-requests: write
```

#### Usage

**Customize labels:**

Edit `.github/labels.yml` to customize file patterns:

```yaml
- label: "custom-label"
  files:
    - "path/to/files/**/*"
```

**Override auto-labels:**

Manual labels always take precedence. Add labels manually in PR/issue to
override automated labeling.

## Setup Guide

### Initial Setup

1. **Enable GitHub Actions**
   - Actions are enabled by default for new repositories
   - Check: Settings → Actions → General → Actions permissions

2. **Configure Secrets**

   Add these secrets in Settings → Secrets → Actions:

   ```
   CODECOV_TOKEN        # Optional: Codecov integration
   ```

   Note: Other secrets like `GITHUB_TOKEN` are automatically provided.

3. **Configure Permissions**

   In Settings → Actions → General → Workflow permissions:
   - Enable "Read and write permissions"
   - Enable "Allow GitHub Actions to create and approve pull requests"

4. **Set Branch Protection**

   In Settings → Branches → Add rule for `main`:
   - ✅ Require status checks to pass before merging
   - ✅ Require branches to be up to date before merging
   - Select required checks:
     - Build and Test (CI)
     - Lint Go Code (golangci-lint)
     - Test Coverage (coverage)
     - Analyze Code (codeql)

### Customization

#### Customize Linting Rules

Edit `.golangci.json` to configure linters:

```yaml
linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
```

#### Customize Release Categories

Edit `.github/release-drafter.yml`:

```yaml
categories:
  - title: "🚀 Features"
    labels:
      - "enhancement"
      - "feature"
```

#### Customize Stale Settings

Edit `.github/workflows/stale.yml`:

```yaml
days-before-issue-stale: 60
days-before-issue-close: 14
exempt-issue-labels: "pinned,security"
```

#### Customize Auto-Labels

Edit `.github/labels.yml`:

```yaml
- label: "custom-area"
  files:
    - "custom/path/**/*"
```

### Advanced Configuration

#### Matrix Builds

Add more Go versions or OS combinations:

```yaml
strategy:
  matrix:
    go-version: [1.23, 1.24, 1.25]
    os: [ubuntu-latest, macos-latest, windows-latest]
```

#### Custom Runners

Use self-hosted runners:

```yaml
runs-on: [self-hosted, linux, x64]
```

#### Conditional Workflows

Add conditions to jobs:

```yaml
if: github.event_name == 'push' && github.ref == 'refs/heads/main'
```

## Usage Examples

### Example 1: Create a Release

```bash
# 1. Create a tag
git tag -a v1.0.0 -m "Release version 1.0.0"

# 2. Push the tag
git push origin v1.0.0

# 3. Wait for workflows to complete:
#    - Release workflow builds binaries
#    - Docker workflow publishes image

# 4. Check release
gh release view v1.0.0
```

### Example 2: Merge PR with Full CI

```bash
# 1. Create feature branch
git checkout -b feature/new-feature

# 2. Make changes and push
git add .
git commit -m "Add new feature"
git push origin feature/new-feature

# 3. Create PR
gh pr create --title "Add new feature" --body "Description"

# 4. Workflows run automatically:
#    - CI: Build and test
#    - GolangCI-Lint: Code quality
#    - Coverage: Test coverage
#    - Dependency Review: Security check
#    - Auto Label: Labels PR
#    - Release Drafter: Updates draft release

# 5. Review and merge
gh pr merge --squash
```

### Example 3: Run Security Scan Manually

```bash
# 1. Trigger CodeQL workflow
gh workflow run codeql-analysis.yml

# 2. Wait for completion
gh run watch

# 3. View results
gh api repos/:owner/:repo/code-scanning/alerts
```

### Example 4: Pull and Run Docker Image

```bash
# 1. Pull latest image
docker pull ghcr.io/kpango/buildbureau:latest

# 2. Run container
docker run -d \
  --name buildbureau \
  -e GEMINI_API_KEY="${GEMINI_API_KEY}" \
  -v buildbureau-data:/app/data \
  -p 8080:8080 \
  ghcr.io/kpango/buildbureau:latest

# 3. Check logs
docker logs -f buildbureau
```

### Example 5: Local Coverage Testing

```bash
# 1. Run tests with coverage
make test-coverage

# 2. View coverage report
go tool cover -html=./coverage/coverage.out

# 3. Check coverage percentage
go tool cover -func=./coverage/coverage.out | grep total
```

## Troubleshooting

### Common Issues

#### Issue: CI Fails with "go: module not found"

**Solution:**

```bash
# Ensure go.mod is up to date
go mod tidy
git add go.mod go.sum
git commit -m "Update dependencies"
```

#### Issue: Docker Build Fails with CGo Errors

**Solution:**

```yaml
# In workflow, ensure CGO is enabled
env:
  CGO_ENABLED: 1
```

#### Issue: Linter Fails on New Code

**Solution:**

```bash
# Run linter locally
make lint

# Auto-fix issues
golangci-lint run --fix
```

#### Issue: Coverage Report Not Appearing

**Solution:**

1. Verify `CODECOV_TOKEN` is set in secrets
2. Check Codecov integration is enabled
3. Review workflow logs for upload errors

#### Issue: Release Workflow Doesn't Trigger

**Solution:**

```bash
# Ensure tag starts with 'v'
git tag v1.0.0  # ✅ Correct
git tag 1.0.0   # ❌ Won't trigger

# Verify tag is pushed
git push origin v1.0.0
```

#### Issue: Stale Bot Closing Active Issues

**Solution:** Add exempt labels to prevent closing:

```yaml
exempt-issue-labels: "pinned,security,in-progress"
```

#### Issue: PR Size Label Incorrect

**Solution:** The size calculation includes all changed lines. To adjust
thresholds:

```yaml
xs_max_size: 10
s_max_size: 100
m_max_size: 500
```

### Debugging Workflows

#### View Workflow Runs

```bash
# List recent workflow runs
gh run list

# View specific run
gh run view <run-id>

# Watch live run
gh run watch
```

#### Download Artifacts

```bash
# List artifacts
gh run view <run-id> --log

# Download artifacts
gh run download <run-id>
```

#### Re-run Failed Workflow

```bash
# Re-run all failed jobs
gh run rerun <run-id> --failed

# Re-run entire workflow
gh run rerun <run-id>
```

### Performance Issues

#### Slow Workflow Execution

**Solutions:**

1. **Enable caching:**

```yaml
- uses: actions/setup-go@v5
  with:
    cache: true
```

2. **Use build caching:**

```yaml
- uses: actions/cache@v4
  with:
    path: |
      ~/.cache/go-build
      ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
```

3. **Parallelize jobs:**

```yaml
strategy:
  matrix:
    task: [lint, test, build]
```

#### Workflow Timeout

**Solutions:**

1. **Increase timeout:**

```yaml
timeout-minutes: 30
```

2. **Optimize test execution:**

```bash
# Run tests in parallel
go test -p 4 ./...

# Skip slow tests in CI
go test -short ./...
```

## Best Practices

### General Best Practices

1. **Always run CI locally before pushing**

   ```bash
   make ci-all
   ```

2. **Keep workflows DRY with composite actions**

   ```yaml
   - uses: ./.github/actions/setup
   ```

3. **Use workflow concurrency to cancel outdated runs**

   ```yaml
   concurrency:
     group: ${{ github.workflow }}-${{ github.ref }}
     cancel-in-progress: true
   ```

4. **Pin action versions for security**
   ```yaml
   uses: actions/checkout@v4  # ✅ Good
   uses: actions/checkout@main  # ❌ Avoid
   ```

### Security Best Practices

1. **Use minimal permissions**

   ```yaml
   permissions:
     contents: read # Only what's needed
   ```

2. **Never commit secrets**
   - Use GitHub Secrets
   - Use environment variables

3. **Review dependency updates**
   - Check Dependabot PRs carefully
   - Review dependency-review reports

4. **Enable branch protection**
   - Require status checks
   - Require code reviews
   - Enable CodeQL scanning

### Performance Best Practices

1. **Use caching aggressively**
   - Go modules cache
   - Build cache
   - Docker layer cache

2. **Parallelize independent jobs**

   ```yaml
   jobs:
     test: ...
     lint: ...
   # These run in parallel
   ```

3. **Use matrix builds efficiently**

   ```yaml
   strategy:
     matrix:
       task: [build, test]
   ```

4. **Skip unnecessary workflows**
   ```yaml
   on:
     push:
       paths-ignore:
         - "**.md"
   ```

### Maintenance Best Practices

1. **Keep workflows updated**
   - Review and update action versions
   - Monitor deprecation warnings

2. **Document custom workflows**
   - Add comments in YAML
   - Update this documentation

3. **Monitor workflow costs**
   - Check Actions usage in Settings
   - Optimize long-running workflows

4. **Test workflow changes in feature branches**
   ```bash
   git checkout -b test-ci-changes
   # Make changes to .github/workflows/
   git push -u origin test-ci-changes
   ```

### Pull Request Best Practices

1. **Keep PRs focused**
   - One feature per PR
   - Smaller PRs are easier to review

2. **Write descriptive PR titles**

   ```
   ✅ "Add user authentication with JWT"
   ❌ "Fix stuff"
   ```

3. **Use conventional commits**

   ```
   feat: add user authentication
   fix: resolve memory leak in cache
   docs: update API documentation
   ```

4. **Address CI failures before review**
   - Fix linting issues
   - Ensure tests pass
   - Resolve security warnings

## Additional Resources

### GitHub Actions Documentation

- [GitHub Actions Official Docs](https://docs.github.com/en/actions)
- [Workflow Syntax](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions)
- [GitHub CLI](https://cli.github.com/)

### Tools Documentation

- [golangci-lint](https://golangci-lint.run/)
- [CodeQL](https://codeql.github.com/)
- [Codecov](https://docs.codecov.com/)
- [Docker Buildx](https://docs.docker.com/buildx/working-with-buildx/)
- [Trivy Scanner](https://aquasecurity.github.io/trivy/)

### BuildBureau Documentation

- [Main README](../README.md)
- [Contributing Guide](../CONTRIBUTING.md)
- [Makefile Documentation](MAKEFILE.md)
- [Docker Documentation](DOCKER.md)
- [Architecture Guide](ARCHITECTURE.md)

### Support

For issues or questions:

1. Check [Troubleshooting](#troubleshooting) section
2. Search [existing issues](https://github.com/kpango/BuildBureau/issues)
3. Create a [new issue](https://github.com/kpango/BuildBureau/issues/new)
4. Join discussions in
   [GitHub Discussions](https://github.com/kpango/BuildBureau/discussions)

---

_Last updated: February 2025_ _BuildBureau CI/CD Pipeline Documentation_
