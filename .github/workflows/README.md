# GitHub Actions Workflows

This directory contains comprehensive CI/CD workflows for the BuildBureau project.

## Workflows

### CI/CD Workflows

#### `ci.yml`
The main continuous integration workflow that runs on every push and pull request.

**Triggers:**
- Push to `main` and `copilot/*` branches
- Pull requests to `main`

**Jobs:**
- Install dependencies
- Generate proto files
- Check formatting
- Run linters
- Build application
- Run tests
- Execute example tests

---

#### `release.yml`
Automated release workflow for creating multi-platform binary releases.

**Triggers:**
- Push of version tags (`v*`)

**Features:**
- Multi-platform builds (Linux/macOS, amd64/arm64)
- CGO support for SQLite
- Automatic changelog generation
- SHA256 checksums
- GitHub release creation
- Installation instructions

**Platforms:**
- `linux/amd64`
- `linux/arm64`
- `darwin/amd64`
- `darwin/arm64`

---

#### `docker-publish.yml`
Docker image build and publishing workflow.

**Triggers:**
- Push to `main`
- Version tags (`v*`)
- Pull requests (build only)
- Manual dispatch

**Features:**
- Multi-architecture Docker images (linux/amd64, linux/arm64)
- Publishes to GitHub Container Registry (GHCR)
- SBOM generation
- Trivy vulnerability scanning
- Automatic image tagging (latest, semver, sha, branch)
- Build caching for faster builds

**Image Tags:**
- `latest` - Latest main branch build
- `v1.2.3` - Semantic version tags
- `v1.2` - Minor version tags
- `v1` - Major version tags
- `main-<sha>` - Branch with commit SHA

---

### Code Quality Workflows

#### `golangci-lint.yml`
Comprehensive Go linting workflow using golangci-lint.

**Triggers:**
- Push to `main`, `develop`
- Pull requests to `main`

**Linters Enabled (30+):**
- errcheck, gosimple, govet, ineffassign
- staticcheck, unused, gofmt, goimports
- misspell, gocritic, revive, stylecheck
- gosec, bodyclose, noctx, gocyclo
- dupl, prealloc, errorlint, and more

**Configuration:** `.golangci.json`

---

#### `coverage.yml`
Code coverage tracking and reporting workflow.

**Triggers:**
- Push to `main`, `develop`
- Pull requests to `main`

**Features:**
- Runs tests with coverage
- Uploads to Codecov
- Comments coverage on PRs
- Generates coverage reports
- Creates coverage artifacts

**Artifacts:**
- `coverage.out` - Coverage data
- `coverage.txt` - Coverage report

---

### Security Workflows

#### `codeql-analysis.yml`
CodeQL security analysis workflow.

**Triggers:**
- Push to `main`, `develop`
- Pull requests to `main`
- Weekly schedule (Monday 6 AM UTC)
- Manual dispatch

**Features:**
- Automated security scanning
- Vulnerability detection
- Security and quality queries
- GitHub Security tab integration

---

#### `dependency-review.yml`
Dependency security review workflow.

**Triggers:**
- Pull requests to `main`

**Features:**
- Dependency vulnerability scanning
- License compliance checking
- govulncheck for Go vulnerabilities
- Nancy security scanning
- Fails on moderate+ severity vulnerabilities
- Denies GPL-3.0, AGPL-3.0 licenses

---

### Automation Workflows

#### `release-drafter.yml`
Automated release notes generation.

**Triggers:**
- Push to `main`
- Pull request events (opened, synchronize, reopened, closed)

**Features:**
- Auto-generates release notes
- Categorizes changes by type
- Version resolution (major/minor/patch)
- Auto-labels PRs

**Categories:**
- 🚀 Features
- 🐛 Bug Fixes
- 🔒 Security
- ⚡ Performance
- 📚 Documentation
- 🧪 Testing
- 🏗️ Build & CI/CD
- 🔧 Maintenance
- 💔 Breaking Changes

**Configuration:** `.github/release-drafter.yml`

---

#### `auto-label.yml`
Automatic labeling for PRs and issues.

**Triggers:**
- Pull requests (opened, synchronize, reopened)
- Issues (opened, edited)

**Features:**
- Labels PRs based on changed files
- Labels PRs by size (xs/s/m/l/xl)
- Auto-labels issues by keywords
- Detects documentation changes
- Detects breaking changes

**Size Thresholds:**
- `size/xs`: ≤10 lines
- `size/s`: ≤100 lines
- `size/m`: ≤500 lines
- `size/l`: ≤1000 lines
- `size/xl`: >1000 lines

**Configuration:** `.github/labels.yml`

---

#### `stale.yml`
Stale issue and PR management.

**Triggers:**
- Daily schedule (midnight UTC)
- Manual dispatch

**Configuration:**
- Issues: Stale after 60 days, close after 14 days
- PRs: Stale after 30 days, close after 7 days

**Exempt Labels:**
- Issues: `pinned`, `security`, `help-wanted`, `good-first-issue`
- PRs: `pinned`, `security`, `work-in-progress`

---

## Configuration Files

### `.github/dependabot.yml`
Dependabot configuration for automated dependency updates.

**Update Schedule:** Weekly (Monday 6 AM UTC)

**Ecosystems:**
- Go modules
- GitHub Actions
- Docker

**Features:**
- Grouped updates (production/development)
- Auto-labeling
- Conventional commit messages
- Ignores major version updates

---

### `.github/labels.yml`
Label definitions for auto-labeling based on file patterns.

**Labels:**
- `go` - Go source files
- `documentation` - Markdown and docs
- `ci/cd` - CI/CD files
- `docker` - Docker files
- `configuration` - Config files
- `internal` - Internal packages
- `pkg` - Public API packages
- `agents` - Agent files
- `cli` - CLI files
- `testing` - Test files
- `dependencies` - Dependency files
- `security` - Security files
- `performance` - Performance files

---

### `.github/release-drafter.yml`
Release Drafter configuration for changelog generation.

**Version Resolution:**
- **Major**: `breaking-change`, `major` labels
- **Minor**: `enhancement`, `feature`, `minor` labels
- **Patch**: `bug`, `bugfix`, `fix`, `patch`, `chore`, `documentation`, `dependencies` labels

---

### `.github/CODEOWNERS`
Code ownership definitions for automatic review requests.

**Owners:**
- Default: @kpango
- All sections: @kpango

---

### `.golangci.json`
golangci-lint configuration with comprehensive linter settings.

**Enabled Linters:** 30+
- Code quality: gocritic, revive, stylecheck
- Security: gosec
- Performance: prealloc
- Error handling: errcheck, errorlint, nilerr
- Code complexity: gocyclo, gocognit
- And many more...

---

## Issue Templates

### `ISSUE_TEMPLATE/bug_report.md`
Structured bug report template with:
- Bug description
- Reproduction steps
- Expected vs actual behavior
- Environment details
- Screenshots/logs
- Configuration
- Checklist

### `ISSUE_TEMPLATE/feature_request.md`
Feature request template with:
- Feature description
- Problem statement
- Proposed solution
- Alternatives considered
- Use cases and benefits
- Implementation details
- Checklist

### `ISSUE_TEMPLATE/question.md`
Question template with:
- Question description
- Context and attempts
- Environment details
- Configuration
- Documentation references
- Checklist

---

## Pull Request Template

### `pull_request_template.md`
Comprehensive PR template with:
- Description and linked issues
- Type of change checklist
- Changes made
- Testing information
- Screenshots/recordings
- Review checklist
- Performance impact
- Breaking changes section
- Related issues/PRs

---

## Secrets Required

### Optional Secrets
- `CODECOV_TOKEN` - For Codecov integration (optional, public repos don't need it)

### Automatic Secrets
- `GITHUB_TOKEN` - Automatically provided by GitHub Actions

---

## Usage Examples

### Creating a Release

1. Tag your commit with a version:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. The `release.yml` workflow will:
   - Build binaries for all platforms
   - Generate checksums
   - Create a GitHub release
   - Upload all artifacts

### Building Docker Images

Docker images are automatically built and pushed on:
- Every push to `main` (tagged as `latest`)
- Every version tag (tagged with semver)
- Pull requests (build only, not pushed)

Access images:
```bash
docker pull ghcr.io/kpango/buildbureau:latest
docker pull ghcr.io/kpango/buildbureau:v1.0.0
```

### Running Linters Locally

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linters
golangci-lint run ./...

# Or use Make
make lint-all
```

### Checking Security Locally

```bash
# Install govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest

# Check vulnerabilities
govulncheck ./...

# Or use Make
make security
```

### Running Tests with Coverage

```bash
# Using Make
make test-coverage

# View HTML report
make test-coverage-html
open coverage/coverage.html
```

---

## Branch Protection Recommendations

Configure branch protection for `main`:

1. **Require status checks:**
   - CI / Build and Test
   - golangci-lint / Lint Go Code
   - Code Coverage / Test Coverage
   - CodeQL Analysis / Analyze Code

2. **Require reviews:**
   - At least 1 approval
   - Dismiss stale reviews on push
   - Require review from code owners

3. **Additional rules:**
   - Require branches to be up to date
   - Require conversation resolution
   - Include administrators

---

## Troubleshooting

### Release Workflow Fails

**Issue:** Cross-compilation fails for ARM64
**Solution:** Ensure CGO is properly configured for cross-compilation

### Docker Build Fails

**Issue:** Multi-arch build fails
**Solution:** Ensure Docker Buildx is properly set up

### Coverage Upload Fails

**Issue:** Codecov token missing
**Solution:** Add `CODECOV_TOKEN` secret (optional for public repos)

### Linter Fails

**Issue:** Too many linter errors
**Solution:** Run `make fmt` and fix reported issues

---

## Maintenance

### Updating Dependencies

Dependabot automatically creates PRs for:
- Go module updates (weekly)
- GitHub Actions updates (weekly)
- Docker base image updates (weekly)

### Updating Workflows

When updating workflows:
1. Test in a feature branch
2. Verify workflow syntax
3. Check all triggers
4. Update this README if needed

### Adding New Labels

1. Edit `.github/labels.yml`
2. Add file patterns
3. Test with a PR

---

## Best Practices

1. **Always run tests before pushing:**
   ```bash
   make ci-all
   ```

2. **Use conventional commits:**
   - `feat:` - New features
   - `fix:` - Bug fixes
   - `docs:` - Documentation
   - `chore:` - Maintenance
   - `test:` - Tests
   - `ci:` - CI/CD changes

3. **Tag releases properly:**
   - Use semantic versioning: `v1.2.3`
   - Create annotated tags: `git tag -a v1.0.0 -m "Release 1.0.0"`

4. **Review security alerts:**
   - Check GitHub Security tab regularly
   - Fix vulnerabilities promptly
   - Keep dependencies updated

---

## New Advanced Features

### Custom Composite Actions

BuildBureau now includes reusable composite actions in `.github/actions/`:

#### `setup-go`
- Automatically detects Go version from `go.mod`
- Supports version override via input
- Includes caching support
- Verifies installation

#### `dump-context`
- Dumps all GitHub context for debugging
- Shows environment variables, job context, runner info
- Useful for troubleshooting workflows

#### `notify-slack`
- Sends formatted notifications to Slack
- Supports success/failure/info status colors
- Includes workflow details and links
- Customizable message and workflow name

### ChatOps Workflow (`chatops.yaml`)

Interact with PRs using commands in comments:

```bash
# Add labels to a PR
/label bug priority/high

# Rebase PR on base branch
/rebase

# Auto-format code
/format

# Show help
/help
```

**Permissions:** Defined in `.github/chatops_permissions.yaml`
- **Owner**: Full access to all commands
- **Maintainer**: label, rebase, format, approve
- **Contributor**: label, format
- **Member**: label, format

### Enhanced Labeling (`labeler.yaml`, `pr-size-labeler.yaml`)

Automatic labeling based on:
- Changed file patterns (area/*, type/*, language/*)
- PR size (size/XS through size/XL)
- Content detection

### Benchmark Tracking (`benchmark.yaml`)

- Runs Go benchmarks on code changes
- Tracks performance over time
- Comments results on PRs
- Alerts on performance regressions

### Merge Conflict Detection (`check-conflict.yaml`)

- Automatically detects merge conflicts
- Comments on PR with conflict details
- Adds/removes `merge-conflict` label

### Go Version Consistency (`check-go-version.yaml`)

- Ensures all workflows use the same Go version as go.mod
- Comments on inconsistencies
- Fails if mismatches found

### Coverage Reporting (`coverage.yaml`)

- Comprehensive test coverage tracking
- Uploads to Codecov
- Generates HTML reports
- Artifacts for 30 days

---

## Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [golangci-lint Documentation](https://golangci-lint.run/)
- [CodeQL Documentation](https://codeql.github.com/docs/)
- [Dependabot Documentation](https://docs.github.com/en/code-security/dependabot)
- [Semantic Versioning](https://semver.org/)
- [vdaas/vald GitHub Actions Reference](https://github.com/vdaas/vald/tree/main/.github)

---

## Support

For issues or questions about the CI/CD workflows:
1. Check this README
2. Review workflow logs in GitHub Actions tab
3. Open an issue using the question template
4. Contact @kpango (code owner)
