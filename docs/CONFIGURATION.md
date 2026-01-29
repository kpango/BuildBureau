# Configuration Guide

Detailed configuration guide for BuildBureau

## Configuration File Structure

BuildBureau manages all settings in `config.yaml`.

## Agent Configuration (agents)

You can configure the following parameters for each agent type:

### Common Parameters

```yaml
agents:
  <agent_type>:
    count: <number>           # Number of agents
    model: <string>           # LLM model name to use
    instruction: <string>     # System prompt
    allowTools: <boolean>     # Allow tool usage
    tools: [<strings>]        # List of available tools
    timeout: <seconds>        # Timeout duration
    retryCount: <number>      # Number of retries
```

### Agent Types

#### 1. president (President)
```yaml
president:
  count: 1
  model: "gemini-2.0-flash-exp"
  instruction: |
    You are the president who oversees the entire project and makes strategic decisions.
    Understand the client's requirements and create an overall project plan.
  allowTools: true
  tools:
    - web_search      # Web search
    - knowledge_base  # Knowledge base access
  timeout: 120
  retryCount: 3
```

#### 2. president_secretary (President's Secretary)
```yaml
president_secretary:
  count: 1
  model: "gemini-2.0-flash-exp"
  instruction: |
    You are the president's secretary. Record requirements based on the president's instructions,
    and update the internal knowledge base.
  allowTools: true
  tools:
    - knowledge_base
    - document_manager
  timeout: 60
  retryCount: 3
```

#### 3. department_manager (Department Manager)
```yaml
department_manager:
  count: 1
  model: "gemini-2.0-flash-exp"
  instruction: |
    You are the department manager responsible for dividing the entire project into section manager units.
  allowTools: true
  tools:
    - web_search
    - knowledge_base
  timeout: 120
  retryCount: 3
```

#### 4. section_manager (Section Manager)
```yaml
section_manager:
  count: 3  # Multiple assignments possible
  model: "gemini-2.0-flash-exp"
  instruction: |
    You are the section manager responsible for creating detailed implementation plans and final specifications.
  allowTools: true
  tools:
    - code_analyzer
    - knowledge_base
  timeout: 90
  retryCount: 3
```

#### 5. employee (Employee)
```yaml
employee:
  count: 6  # Multiple assignments possible
  model: "gemini-2.0-flash-exp"
  instruction: |
    You are an engineer who implements based on given specifications.
  allowTools: true
  tools:
    - code_execution
    - file_operations
    - knowledge_base
  timeout: 180
  retryCount: 3
```

## LLM Configuration (llm)

```yaml
llm:
  provider: "google"                                    # Provider name
  apiEndpoint: "https://generativelanguage.googleapis.com"  # API endpoint
  defaultModel: "gemini-2.0-flash-exp"                 # Default model
  maxTokens: 8192                                       # Maximum tokens
  temperature: 0.7                                      # Temperature parameter (0.0-1.0)
  topP: 0.95                                           # Top-P sampling
```

### Provider Configuration

Currently planned supported providers:
- `google`: Google AI (Gemini)
- `openai`: OpenAI (GPT-4, etc.)
- `anthropic`: Anthropic (Claude)

### Model Selection

Recommended models:
- High-speed processing: `gemini-2.0-flash-exp`
- High quality: `gemini-2.5-pro`
- Balanced: `gemini-2.0-flash-exp`

## gRPC Configuration (grpc)

```yaml
grpc:
  port: 50051                  # gRPC server port
  maxMessageSize: 10485760     # Maximum message size (bytes)
  timeout: 300                 # Timeout (seconds)
  enableReflection: true       # Enable reflection
```

### Port Configuration

- Default: `50051`
- You may need to open this port in your firewall

### Message Size

- Default: 10MB
- Increase this when handling large files

## Slack Notification Configuration (slack)

```yaml
slack:
  enabled: true                      # Enable Slack notifications
  token: "${SLACK_BOT_TOKEN}"        # Bot token (environment variable)
  channelID: "${SLACK_CHANNEL_ID}"   # Channel ID (environment variable)
  retryCount: 3                      # Number of retries
  timeout: 10                        # Timeout (seconds)
  
  notifications:
    projectStart:
      enabled: true
      message: "🚀 Project \"{{.ProjectName}}\" has started"
    
    taskComplete:
      enabled: true
      message: "✅ Task \"{{.TaskName}}\" has been completed ({{.Agent}})"
    
    error:
      enabled: true
      message: "❌ An error occurred: {{.ErrorMessage}} ({{.Agent}})"
    
    projectComplete:
      enabled: true
      message: "🎉 Project \"{{.ProjectName}}\" has been completed!"
```

### Slack Bot Setup Instructions

1. Create an app at [Slack API](https://api.slack.com/apps)
2. Add the following to Bot Token Scopes:
   - `chat:write`
   - `chat:write.public`
3. Install to workspace
4. Obtain Bot User OAuth Token
5. Set environment variables:
   ```bash
   export SLACK_BOT_TOKEN="xoxb-your-token"
   export SLACK_CHANNEL_ID="C01234567"
   ```

### Message Templates

Available variables:
- `{{.ProjectName}}`: Project name
- `{{.TaskName}}`: Task name
- `{{.Agent}}`: Agent ID
- `{{.ErrorMessage}}`: Error message
- `{{.Timestamp}}`: Timestamp

## UI Configuration (ui)

```yaml
ui:
  enableTUI: true        # Enable Terminal UI
  refreshRate: 100       # Refresh interval (milliseconds)
  theme: "default"       # Theme
  showProgress: true     # Show progress
  logLevel: "info"       # Log level
```

### Log Levels

- `debug`: All logs including debug information
- `info`: Normal information logs
- `warn`: Warnings only
- `error`: Errors only

### Themes

Currently available themes:
- `default`: Default theme

## System Configuration (system)

```yaml
system:
  workDir: "./work"              # Working directory
  logDir: "./logs"               # Log directory
  cacheDir: "./cache"            # Cache directory
  maxConcurrentTasks: 10         # Maximum concurrent tasks
  globalTimeout: 3600            # Global timeout (seconds)
```

### Directory Structure

```
BuildBureau/
├── work/      # Temporary working files
├── logs/      # Log files
└── cache/     # Cache files
```

## Environment Variables

### Required Environment Variables

When using Slack notifications:
```bash
export SLACK_BOT_TOKEN="xoxb-..."
export SLACK_CHANNEL_ID="C..."
```

When using Google AI API:
```bash
export GOOGLE_AI_API_KEY="..."
```

### Optional Environment Variables

```bash
# Custom configuration file path
export CONFIG_PATH="/path/to/custom/config.yaml"

# Override log level
export LOG_LEVEL="debug"
```

## Configuration Examples

### Development Environment Configuration

```yaml
agents:
  president:
    count: 1
    timeout: 60
  # ... other agents (with shorter timeouts)

slack:
  enabled: false  # Disable notifications during development

ui:
  enableTUI: true
  logLevel: "debug"  # Enable debug logs
```

### Production Environment Configuration

```yaml
agents:
  president:
    count: 1
    timeout: 180
  section_manager:
    count: 5  # Scale up
  employee:
    count: 20  # Scale up

slack:
  enabled: true  # Enable notifications

system:
  maxConcurrentTasks: 20  # Increase parallelism

ui:
  logLevel: "info"  # Information logs only
```

### High-Load Environment Configuration

```yaml
grpc:
  maxMessageSize: 52428800  # 50MB

system:
  maxConcurrentTasks: 50
  globalTimeout: 7200  # 2 hours

agents:
  employee:
    count: 50
    timeout: 300
```

## Troubleshooting

### Configuration File Validation

Check configuration file syntax:
```bash
# YAML syntax check
yamllint config.yaml

# Validate with BuildBureau
./bin/buildbureau --validate-config  # (not yet implemented)
```

### Common Errors

1. **"failed to load config"**
   - Check for YAML syntax errors
   - Verify that indentation is correct

2. **"Slack token is required"**
   - Set the `SLACK_BOT_TOKEN` environment variable
   - Or set `slack.enabled: false`

3. **"president agent count must be at least 1"**
   - Check the count of required agents

## Best Practices

1. **Sensitive Information Management**
   - Manage tokens via environment variables
   - Add `.env` file to `.gitignore`

2. **Timeout Configuration**
   - Set appropriately according to agent roles
   - Set longer timeouts for tasks involving implementation

3. **Retry Count**
   - Set higher in unstable network environments
   - Set an upper limit to avoid infinite loops

4. **Log Level**
   - Use `debug` during development
   - Use `info` or `warn` in production
