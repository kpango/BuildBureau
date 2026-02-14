# Dependency Management Guide

## Overview

BuildBureau uses a template-based dependency management system inspired by the [vdaas/vald](https://github.com/vdaas/vald) repository. This approach provides better control over dependency upgrades and ensures reproducible builds.

## Files

### `go.mod.default`

This is the **template file** that defines the dependency upgrade policy. It contains:

- Module declaration
- Go version
- Replace directives with `upgrade` placeholders

**Example:**
```go
module github.com/kpango/BuildBureau

go 1.26.0

replace (
	github.com/charmbracelet/bubbles => github.com/charmbracelet/bubbles upgrade
	google.golang.org/grpc => google.golang.org/grpc upgrade
	// ...
)
```

The `upgrade` keyword is a placeholder that gets replaced with the latest version during the build process.

### `go.mod`

This is the **actual Go module file** used by the Go toolchain. It's generated from `go.mod.default` and contains resolved versions.

**Note:** `go.mod` should NOT be manually edited for version updates. Instead, update `go.mod.default` and run `make deps-reset`.

## Makefile Targets

### `make deps`
Download and tidy dependencies (standard Go workflow).

```bash
make deps
```

### `make deps-update`
Update all dependencies to their latest versions using `go get -u`.

```bash
make deps-update
```

### `make deps-reset`
**NEW:** Reset `go.mod` from `go.mod.default` template and resolve all dependencies.

```bash
make deps-reset
```

This target:
1. Backs up current `go.mod` to `go.mod.backup`
2. Copies `go.mod.default` to `go.mod`
3. Updates the Go version to match your current Go installation
4. Runs `go mod tidy` to resolve all dependencies

### `make deps-verify`
Verify that dependencies are valid and checksums match.

```bash
make deps-verify
```

### `make deps-graph`
Display the dependency graph.

```bash
make deps-graph
```

### `make deps-tidy`
Tidy `go.mod` and `go.sum` files.

```bash
make deps-tidy
```

## Workflow

### Regular Development

For day-to-day development, use standard targets:

```bash
# Install dependencies
make deps

# Update to latest versions
make deps-update
```

### Systematic Dependency Updates

When you want to systematically update all dependencies:

```bash
# Reset from template and update all
make deps-reset

# Verify everything still works
make build
make test
```

### Adding New Dependencies

1. Add the dependency to your code
2. Run `make deps` to add it to `go.mod`
3. Update `go.mod.default` to include the new dependency:

```go
replace (
	// ... existing dependencies ...
	github.com/new/dependency => github.com/new/dependency upgrade
)
```

4. Commit both `go.mod.default` and `go.mod`

## Benefits

### Version Control
- `go.mod.default` is version controlled and defines upgrade policies
- Easy to see which dependencies should be kept at specific versions vs. upgraded

### Reproducible Builds
- Template ensures consistent dependency resolution across environments
- Clear separation between policy (go.mod.default) and state (go.mod)

### Automated Updates
- CI/CD can use `make deps-reset` to get latest versions
- Makefile handles version resolution automatically

### Inspired by vald
This pattern is used successfully by the vald project, which manages hundreds of dependencies across a large codebase.

## Reference

- [vald Makefile.d/dependencies.mk](https://github.com/vdaas/vald/blob/main/Makefile.d/dependencies.mk)
- [Go Modules Reference](https://go.dev/ref/mod)
