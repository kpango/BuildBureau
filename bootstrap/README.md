# BuildBureau Bootstrap System

## Overview

This directory contains the **self-hosting/bootstrapping system** for
BuildBureau - enabling BuildBureau to build and improve itself using its own
multi-agent capabilities.

## Concept

**Self-hosting** (or **bootstrapping**) is when a system uses its own
capabilities to develop, improve, or build itself. For BuildBureau, this means:

- BuildBureau's agents can understand BuildBureau's codebase
- Agents can design and implement improvements to BuildBureau
- The system can refactor, optimize, and extend its own code
- BuildBureau learns from building itself, improving future self-modifications

This is a proof-of-concept showing **recursive self-improvement** - a key
capability for advanced AI systems.

## Architecture

### Bootstrap Configuration

The bootstrap system uses a specialized configuration (`config.yaml`) with:

- **Self-aware agents**: Agents with deep knowledge of BuildBureau's
  architecture
- **Enhanced memory**: Longer retention for self-improvement learnings
- **Multiple agents**: 3 Engineers, 2 Managers for parallel self-development
- **Dedicated database**: Separate `bootstrap.db` for self-improvement context

### Specialized Agents

Located in `agents/`, these agents have prompts that include:

1. **President**: Understands BuildBureau's high-level architecture and
   requirements
2. **Secretary**: Coordinates self-modification tasks
3. **Director**: Makes architectural decisions for BuildBureau
4. **Manager**: Designs implementations following BuildBureau patterns
5. **Engineer**: Implements code changes using BuildBureau conventions

Each agent is "self-aware" - they know they're modifying the system they're
running in.

### Task Templates

Located in `tasks/`, these templates help structure self-improvement work:

- `add-feature.yaml` - Add new capabilities to BuildBureau
- `refactor.yaml` - Improve code structure
- `optimize.yaml` - Enhance performance
- `test.yaml` - Add test coverage

## Usage

### Quick Start

```bash
# Run Bootstrap Mode
make bootstrap

# Or manually
export BUILDBUREAU_CONFIG=bootstrap/config.yaml
./build/buildbureau
```

### Example: Add a New Feature

```bash
# 1. Start BuildBureau in bootstrap mode
make bootstrap

# 2. In the TUI, provide a task:
"Add a new agent type called 'Reviewer' that can review code changes"

# 3. BuildBureau will:
#    - Analyze the request (President)
#    - Coordinate implementation (Secretary)
#    - Design the solution (Director → Manager)
#    - Implement the code (Engineer)
#    - Generate tests
#    - Update documentation

# 4. Review the changes in your git working directory
git diff

# 5. Test the changes
make test

# 6. Commit if satisfied
git commit -m "Add Reviewer agent (self-implemented by BuildBureau)"
```

### Example: Refactor Code

```bash
make bootstrap

# Task: "Refactor the memory system to use a more efficient caching strategy"

# BuildBureau will analyze the memory system and propose improvements
```

### Example: Optimize Performance

```bash
make bootstrap

# Task: "Optimize agent task delegation to reduce latency"

# BuildBureau will benchmark, identify bottlenecks, and implement optimizations
```

## How It Works

### 1. Self-Awareness

Each bootstrap agent's prompt includes:

- BuildBureau's architecture and design patterns
- Code structure and conventions
- Technology stack (Go, SQLite, LLMs, etc.)
- Development best practices

### 2. Contextual Code Generation

When implementing changes, Engineer agents:

- Reference existing patterns in BuildBureau
- Follow established conventions (BaseAgent, memory integration, etc.)
- Generate code that fits the existing architecture
- Include appropriate tests

### 3. Memory-Enhanced Learning

Bootstrap agents remember:

- Past self-modifications and their outcomes
- Architectural decisions made
- Patterns that work well
- Mistakes to avoid

Each self-improvement cycle improves future cycles.

### 4. Safety Mechanisms

- **Human review**: Changes are generated but not automatically applied
- **Git integration**: All changes go through version control
- **Testing**: Generated code includes tests
- **Rollback**: Easy to revert if something breaks

## Configuration

### Environment Variables

```bash
# Required: LLM API key (at least one)
export GEMINI_API_KEY="your-key"
# or
export OPENAI_API_KEY="your-key"
# or
export CLAUDE_API_KEY="your-key"

# Optional: Custom config path
export BUILDBUREAU_CONFIG="bootstrap/config.yaml"
```

### Customization

Edit `bootstrap/config.yaml` to adjust:

- Number of agents per layer
- Memory retention periods
- Default LLM model
- Database path

Edit agent prompts in `bootstrap/agents/` to change behavior.

## Advanced Usage

### Batch Processing

Create a task file and process multiple improvements:

```yaml
# improvements.yaml
tasks:
  - type: feature
    description: "Add streaming response support for LLM providers"
  - type: optimization
    description: "Optimize memory query performance"
  - type: refactor
    description: "Improve error handling in agent communication"
```

Process with a custom script:

```bash
./bootstrap/batch-process.sh improvements.yaml
```

### CI/CD Integration

Run BuildBureau in bootstrap mode as part of CI:

```yaml
# .github/workflows/self-improve.yml
name: Self-Improvement

on:
  schedule:
    - cron: "0 0 * * 0" # Weekly

jobs:
  bootstrap:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run Bootstrap Mode
        env:
          GEMINI_API_KEY: ${{ secrets.GEMINI_API_KEY }}
        run: make bootstrap-ci
      - name: Create PR with improvements
        # ... create PR with generated changes
```

## Benefits

### 1. Recursive Self-Improvement

Each improvement cycle makes BuildBureau better at improving itself:

- Better self-awareness over time
- More efficient implementation patterns
- Accumulated knowledge in memory

### 2. Consistency

Self-generated code follows BuildBureau's patterns because agents learned them.

### 3. Rapid Iteration

Add features faster by letting BuildBureau implement them.

### 4. Learning Platform

Demonstrates advanced AI capabilities:

- Self-modification
- Meta-reasoning
- Recursive improvement

## Limitations

### Current Limitations

1. **Human review needed**: Generated code should be reviewed before merging
2. **Complex changes**: May need human guidance for major architectural changes
3. **External dependencies**: Can't install new system-level dependencies
4. **Testing required**: Generated tests need validation

### Future Improvements

- Automated code review and safety checks
- Ability to run tests before presenting changes
- Integration with CI/CD for continuous self-improvement
- Multi-step planning for complex features
- Automatic documentation generation

## Safety Notes

⚠️ **Important Considerations:**

- Always review generated code before running it
- Test thoroughly in a separate environment first
- Keep git history clean for easy rollback
- Don't auto-apply changes in production
- Monitor for infinite improvement loops

## Examples

See `examples/bootstrap/` for:

- Example self-improvement tasks
- Sample outputs
- Success stories
- Common patterns

## Troubleshooting

### Agent Not Understanding Codebase

- Enhance agent prompts with more context
- Use memory system to store architectural decisions
- Break task into smaller pieces

### Generated Code Doesn't Follow Patterns

- Review agent prompts for pattern examples
- Add more examples to agent configurations
- Use Manager agent to specify patterns explicitly

### Tests Failing

- Ensure agents understand test patterns
- Provide examples of existing tests
- Manually fix and add to agent knowledge

## Contributing

To improve the bootstrap system itself:

1. Enhance agent prompts with better context
2. Add more task templates
3. Improve safety mechanisms
4. Share successful bootstrap tasks

## Philosophy

> "A system that can improve itself is the first step toward truly autonomous
> software development."

BuildBureau's bootstrap mode is an experiment in **recursive
self-improvement** - a foundational capability for advanced AI systems. By
enabling BuildBureau to understand and modify its own code, we're exploring the
boundaries of AI-driven software development.

## Resources

- [BuildBureau Architecture](../docs/ARCHITECTURE.md)
- [Agent System](../docs/AGENT_MEMORY.md)
- [Development Guide](../AGENTS.md)
- [Main README](../README.md)

---

**Status**: Experimental - Use with caution and human oversight **Version**:
1.0.0 **Last Updated**: 2026-02-01
