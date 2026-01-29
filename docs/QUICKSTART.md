# BuildBureau Quick Start Guide

This guide will help you get BuildBureau up and running in minutes.

## Prerequisites

- Go 1.21 or later installed
- Google Gemini API key (get one from [Google AI Studio](https://makersuite.google.com/app/apikey))
- (Optional) Slack Bot Token for notifications

## Step 1: Clone the Repository

```bash
git clone https://github.com/kpango/BuildBureau.git
cd BuildBureau
```

## Step 2: Set Up Environment Variables

Copy the example environment file:

```bash
cp .env.example .env
```

Edit `.env` and add your API key:

```bash
# Required
GOOGLE_API_KEY=your-actual-gemini-api-key-here

# Optional: Slack integration
SLACK_BOT_TOKEN=xoxb-your-slack-bot-token
SLACK_CHANNEL_MAIN=C01234567
```

Load the environment variables:

```bash
export $(cat .env | xargs)
```

Or on Windows PowerShell:
```powershell
Get-Content .env | ForEach-Object {
    $name, $value = $_.split('=')
    Set-Item -Path "Env:$name" -Value $value
}
```

## Step 3: Build the Application

```bash
make build
```

Or manually:

```bash
mkdir -p bin
go build -o bin/buildbureau ./cmd/buildbureau
```

## Step 4: Run BuildBureau

```bash
make run
```

Or directly:

```bash
./bin/buildbureau
```

## Using BuildBureau

### 1. Enter a Client Request

When the Terminal UI starts, you'll see a prompt:

```
🏢 BuildBureau - Multi-Agent AI System

┌─ Conversation ─────────────────────────────────────────┐
│                                                         │
│ Waiting for client request...                          │
│                                                         │
└─────────────────────────────────────────────────────────┘

Enter your request:
┌─────────────────────────────────────────────────────────┐
│ |                                                       │
└─────────────────────────────────────────────────────────┘
Press Enter to submit • Ctrl+C to quit
```

Type your request, for example:

```
We need to build an e-commerce website with user authentication, 
product catalog, shopping cart, and payment processing.
```

Press **Enter** to submit.

### 2. Watch the Agents Work

You'll see the agents process your request hierarchically:

```
09:30:15 [CEO]: Analyzing client requirements for e-commerce platform...
  09:30:18 [CEO Secretary]: Recording project requirements to knowledge base...
  
    09:30:22 [Department Head]: Breaking down into technical categories...
      09:30:25 [DeptHead Secretary]: Researching e-commerce best practices...
      
        09:30:30 [Manager (Frontend)]: Creating UI/UX specifications...
          09:30:32 [Manager Secretary]: Researching React frameworks...
          
            09:30:35 [Worker 1-1-1]: Implementing authentication UI...
            09:30:40 [Worker 1-1-2]: Implementing product catalog UI...
        
        09:30:45 [Manager (Backend)]: Creating API specifications...
          09:30:47 [Manager Secretary]: Researching payment gateways...
          
            09:30:50 [Worker 1-2-1]: Implementing authentication API...
            09:30:55 [Worker 1-2-2]: Implementing payment processing...
```

### 3. Check Slack Notifications (if configured)

If you configured Slack, you'll receive notifications for:
- Project start
- Task assignments
- Task completions
- Any errors

## Configuration

### Customizing Agent Behavior

Edit `configs/config.yaml` to customize:

**Agent Instructions:**
```yaml
agents:
  ceo:
    instruction: |
      You are an experienced CEO with 20 years in software...
      [customize the CEO's behavior here]
```

**Organizational Structure:**
```yaml
hierarchy:
  departments: 1
  managers_per_department: 3
  manager_specialties:
    - Frontend Development
    - Backend Development  
    - DevOps  # Changed from QA
  workers_per_manager: 3  # Increased from 2
```

**Notification Settings:**
```yaml
slack:
  enabled: true  # Set to false to disable Slack
  notify_on:
    task_completed:
      enabled: true
      roles: ["Manager", "Worker"]  # Add/remove roles
      channel: "development"
```

### Changing AI Models

Update the model in `configs/config.yaml`:

```yaml
agents:
  ceo:
    model: "gemini-1.5-pro"  # Changed from gemini-2.0-flash-exp
    temperature: 0.8  # Adjust creativity (0.0 - 1.0)
```

Available models:
- `gemini-2.0-flash-exp` (Fast, experimental)
- `gemini-1.5-pro` (More powerful)
- `gemini-1.5-flash` (Balanced)

## Common Issues

### Issue: "failed to create Gemini client"

**Solution:** Make sure your API key is set correctly:

```bash
echo $GOOGLE_API_KEY  # Should print your API key
```

If empty, export it:

```bash
export GOOGLE_API_KEY="your-api-key-here"
```

### Issue: "failed to authenticate with Slack"

**Solution:** Check your Slack token has the required scopes:
- `chat:write` - Send messages
- `chat:write.public` - Post to channels

### Issue: Build fails with "package not found"

**Solution:** Run `go mod download`:

```bash
go mod tidy
go mod download
```

### Issue: "terminal not supported"

**Solution:** Make sure you're running in a proper terminal:
- On Linux/Mac: Use Terminal.app, iTerm2, or similar
- On Windows: Use Windows Terminal (not Command Prompt)

## Advanced Usage

### Running Without UI (Headless Mode)

Edit `configs/config.yaml`:

```yaml
system:
  ui:
    enabled: false
```

This is useful for:
- Server deployments
- Automated testing
- CI/CD pipelines

### Custom Configuration File

Use a different config file:

```bash
./bin/buildbureau -config path/to/custom-config.yaml
```

### Viewing Logs

Enable file logging in `configs/config.yaml`:

```yaml
system:
  logging:
    enable_file_logging: true
    file: "logs/buildbureau.log"
    level: "DEBUG"  # DEBUG, INFO, WARN, ERROR
```

View logs:

```bash
tail -f logs/buildbureau.log
```

## Next Steps

- Read the [Architecture Documentation](docs/ARCHITECTURE.md)
- Explore [example configurations](configs/)
- Customize agent prompts for your use case
- Set up Slack notifications
- Integrate with your development workflow

## Getting Help

- Check the [README](README.md) for detailed documentation
- Review the [Architecture Guide](docs/ARCHITECTURE.md)
- Open an issue on GitHub for bugs or feature requests

## Example Use Cases

### 1. Software Project Planning

Request:
```
Plan a mobile app for task management with offline sync, 
push notifications, and team collaboration features.
```

### 2. Code Review Process Design

Request:
```
Design a code review process for our team including 
automated checks, peer review workflow, and approval gates.
```

### 3. API Design

Request:
```
Design a RESTful API for a social media platform with 
posts, comments, likes, and user profiles.
```

### 4. Testing Strategy

Request:
```
Create a comprehensive testing strategy for a microservices 
architecture including unit, integration, and e2e tests.
```

## Tips for Better Results

1. **Be Specific**: Provide clear requirements and constraints
2. **Include Context**: Mention technology preferences, constraints, team size
3. **Ask Follow-ups**: Agents can refine based on additional information
4. **Review Configuration**: Adjust agent instructions for your domain
5. **Monitor Logs**: Check logs for detailed agent reasoning

---

Happy building with BuildBureau! 🏢✨
