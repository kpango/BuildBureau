# 🚀 Docker Quick Start

Run BuildBureau in 3 simple commands!

## ⚡ Super Quick Start

```bash
# 1. Set API key
export GEMINI_API_KEY="your-key"

# 2. Run with Docker Compose
docker-compose up -d

# 3. Done! View logs:
docker-compose logs -f
```

## 📋 Common Commands

### Start/Stop

```bash
# Start
docker-compose up -d

# Stop
docker-compose down

# Restart
docker-compose restart
```

### Logs & Monitoring

```bash
# View all logs
docker-compose logs -f

# View specific service
docker-compose logs -f buildbureau

# Last 100 lines
docker-compose logs --tail=100 buildbureau
```

### Execute Commands

```bash
# Interactive shell
docker-compose exec buildbureau sh

# Check version
docker-compose exec buildbureau /app/buildbureau --version

# View database
docker-compose exec buildbureau ls -lh /app/data
```

## 🔑 Environment Variables

Set before running:

```bash
# Required (choose at least one)
export GEMINI_API_KEY="your-gemini-key"
export OPENAI_API_KEY="your-openai-key"
export CLAUDE_API_KEY="your-claude-key"

# Optional
export OPENAI_MODEL="gpt-3.5-turbo"
export CLAUDE_MODEL="claude-3-haiku-20240307"
export SLACK_TOKEN="xoxb-your-token"
```

## 🎯 Advanced Options

### With Vald Vector Search

```bash
docker-compose --profile with-vald up -d
```

### Custom Configuration

```bash
# Edit config.yaml, then:
docker-compose restart
```

### Multiple Instances

```bash
# Instance 1
docker-compose -p buildbureau-1 up -d

# Instance 2
docker-compose -p buildbureau-2 up -d
```

## 📦 Build from Source

```bash
# Build image
docker build -t buildbureau:latest .

# Or use Makefile target
make docker-build
```

## 🐛 Troubleshooting

### Container won't start?

```bash
# Check logs
docker-compose logs buildbureau

# Check status
docker-compose ps
```

### Need more memory?

Edit `docker-compose.yml`:

```yaml
resources:
  limits:
    memory: 4G # Increase from 2G
```

### Reset everything?

```bash
# Stop and remove everything
docker-compose down -v

# Start fresh
docker-compose up -d
```

## 📚 Full Documentation

See [Docker Documentation](docs/DOCKER.md) for complete guide.

---

**That's it! BuildBureau is now running with zero dependencies! 🎉**
