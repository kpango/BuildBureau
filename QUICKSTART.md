# Quick Start Guide

Get started with BuildBureau in 5 minutes

## 1. Installation

### Prerequisites

- Go 1.23 or higher

### Build

```bash
git clone https://github.com/kpango/BuildBureau.git
cd BuildBureau
make deps
make build
```

## 2. Configuration

### Start with minimal configuration

You can use `config.yaml` as is.

### Enable Slack notifications (Optional)

```bash
# Create .env file
cp .env.example .env

# Edit .env and set tokens
export SLACK_BOT_TOKEN="xoxb-your-token"
export SLACK_CHANNEL_ID="C01234567"
```

If you don't need Slack notifications:

```yaml
# Disable in config.yaml
slack:
  enabled: false
```

## 3. Execution

### Run with default configuration

```bash
./bin/buildbureau
```

### Terminal UI Operations

An interactive terminal UI will be displayed when launched:

```
🏢 BuildBureau - Multi-Layer AI Agent System

Requirements Input:
┌──────────────────────────────────────┐
│ Enter project requirements...        │
│                                      │
│                                      │
└──────────────────────────────────────┘

Alt+Enter: Submit | Esc: Exit
```

### Entering Project Requirements

1. Enter project requirements in the text area
2. Press `Alt+Enter` to submit
3. Agents will start processing

## 4. Verification

### Example: Simple Project

```
Please create a contact form for a website.
The following features are required:
- Input fields for name, email address, and message
- Validation
- Submission confirmation
```

### Agent Operations

1. President agent analyzes requirements
2. Director agent divides tasks
3. Section manager agent creates implementation plan
4. Employee agent executes implementation

## 5. Configuration Customization

### Changing the Number of Agents

```yaml
# config.yaml
agents:
  section_manager:
    count: 5  # Increase section managers to 5
  employee:
    count: 20  # Increase employees to 20
```

### Adjusting Timeout

```yaml
agents:
  president:
    timeout: 180  # Extend to 180 seconds
```

### Changing LLM Model

```yaml
agents:
  president:
    model: "gemini-2.5-pro"  # Switch to a more powerful model
```

## Troubleshooting

### Build Errors

```bash
make clean
make deps
make build
```

### Configuration Errors

```bash
# Check YAML syntax
cat config.yaml | grep -E "^\s*-"
```

### Slack Notifications Not Arriving

1. Verify the token is correct
2. Verify the channel ID is correct
3. Verify the Bot is added to the channel

## Next Steps

- Check the [Configuration Guide](docs/CONFIGURATION.md) for detailed configuration methods
- Understand the system architecture in the [Architecture Documentation](docs/ARCHITECTURE.md)
- Learn customization methods in the [Development Guide](docs/DEVELOPMENT.md)

## Frequently Asked Questions

### Q: What is the status of LLM implementation?

A: The current version is a foundation implementation. Integration with Google ADK is planned for future implementation.

### Q: Where do the agents run?

A: Currently they run within a single process. The gRPC interface will enable distributed execution in the future.

### Q: Can I add custom agents?

A: Yes. Please refer to the [Development Guide](docs/DEVELOPMENT.md).

### Q: Can I use this commercially?

A: Please check the license.

## Support

- Bug Reports: [GitHub Issues](https://github.com/kpango/BuildBureau/issues)
- Questions: [GitHub Discussions](https://github.com/kpango/BuildBureau/discussions)
- Contributions: [Contributing Guide](CONTRIBUTING.md)
