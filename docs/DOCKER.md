# 🐳 Docker Deployment Guide

Run BuildBureau with zero dependencies using Docker!

## 🎯 Overview

BuildBureau provides Docker support for easy deployment without installing Go,
SQLite, protoc, or any other dependencies on your local machine. Everything runs
in a container!

## 📦 What's Included

- **Multi-stage Dockerfile**: Optimized for small image size (~50MB final image)
- **Docker Compose**: Easy orchestration with optional Vald integration
- **Build Scripts**: Automated build and run scripts
- **Health Checks**: Built-in container health monitoring
- **Security**: Non-root user, minimal attack surface
- **Persistence**: Data volumes for database and configuration

## 🚀 Quick Start

### Option 1: Docker Compose (Recommended)

The easiest way to run BuildBureau:

```bash
# 1. Set your API key(s)
export GEMINI_API_KEY="your-gemini-key"
# Optional: export OPENAI_API_KEY="your-openai-key"
# Optional: export CLAUDE_API_KEY="your-claude-key"

# 2. Start with Docker Compose
docker-compose up -d

# 3. View logs
docker-compose logs -f

# 4. Stop
docker-compose down
```

### Option 2: Docker CLI

Build and run manually:

```bash
# 1. Build the image
docker build -t buildbureau:latest .

# 2. Run the container
docker run -d \
  --name buildbureau \
  -e GEMINI_API_KEY="your-key" \
  -v buildbureau-data:/app/data \
  -p 8080:8080 \
  buildbureau:latest

# 3. View logs
docker logs -f buildbureau

# 4. Stop
docker stop buildbureau
```

### Option 3: Helper Scripts

Use the provided Makefile targets:

```bash
# Build
make docker-build

# Run
export GEMINI_API_KEY="your-key"
make docker-run
```

## 📋 Requirements

Only requirement:

- Docker 20.10+ or Docker Desktop
- Docker Compose 2.0+ (for docker-compose.yml)

That's it! No Go, no SQLite, no protoc needed!

## 🔧 Configuration

### Environment Variables

Required (at least one):

```bash
GEMINI_API_KEY=your-gemini-key
OPENAI_API_KEY=your-openai-key
CLAUDE_API_KEY=your-claude-key
```

Optional:

```bash
# Model overrides
OPENAI_MODEL=gpt-3.5-turbo
CLAUDE_MODEL=claude-3-haiku-20240307

# Slack integration
SLACK_TOKEN=xoxb-your-token

# Vald configuration (if using)
VALD_HOST=vald
VALD_PORT=8081
```

### Volume Mounts

**Data persistence:**

```yaml
volumes:
  - buildbureau-data:/app/data # SQLite database and memory
```

**Custom configuration:**

```yaml
volumes:
  - ./config.yaml:/app/config/config.yaml:ro
  - ./agents:/app/agents:ro
```

## 🐳 Docker Compose Features

### Basic Usage

```bash
# Start all services
docker-compose up -d

# Start with Vald for semantic search
docker-compose --profile with-vald up -d

# View logs
docker-compose logs -f buildbureau

# Restart services
docker-compose restart

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

### Service Configuration

**BuildBureau Service:**

- Automatic restart unless stopped
- Resource limits: 2 CPU, 2GB RAM
- Health checks every 30s
- Persistent data volume

**Vald Service (Optional):**

- Only starts with `--profile with-vald`
- Provides semantic search capabilities
- Separate data volume
- Accessible on ports 8081 (gRPC) and 8082 (REST)

## 🏗️ Image Details

### Build Stages

**Stage 1: Builder (golang:1.24-alpine)**

- Installs build dependencies (gcc, protobuf)
- Downloads Go modules
- Generates proto code
- Builds static binary

**Stage 2: Runtime (alpine:latest)**

- Minimal base image (~5MB)
- Only runtime dependencies
- Non-root user for security
- Final image: ~50MB

### Security Features

- ✅ Non-root user (UID 1000)
- ✅ Minimal base image
- ✅ No unnecessary dependencies
- ✅ Read-only configuration mounts
- ✅ Health checks
- ✅ Resource limits

## 📊 Resource Usage

**Minimum Requirements:**

- CPU: 1 core
- RAM: 512MB
- Disk: 100MB (excluding data)

**Recommended:**

- CPU: 2 cores
- RAM: 2GB
- Disk: 1GB (with data)

## 🔍 Monitoring & Debugging

### View Logs

```bash
# All logs
docker logs -f buildbureau

# Last 100 lines
docker logs --tail 100 buildbureau

# With timestamps
docker logs -t buildbureau
```

### Execute Commands

```bash
# Interactive shell
docker exec -it buildbureau sh

# Run specific command
docker exec buildbureau /app/buildbureau --version

# Check database
docker exec buildbureau ls -lh /app/data
```

### Health Check

```bash
# Check container health
docker inspect --format='{{.State.Health.Status}}' buildbureau

# View health check logs
docker inspect --format='{{json .State.Health}}' buildbureau | jq
```

## 🔄 Updates & Maintenance

### Update Image

```bash
# Rebuild with latest code
docker-compose build --no-cache

# Pull latest from registry (if published)
docker-compose pull

# Restart with new image
docker-compose up -d
```

### Backup Data

```bash
# Backup database
docker cp buildbureau:/app/data ./backup/

# Or use volume backup
docker run --rm \
  -v buildbureau-data:/data \
  -v $(pwd)/backup:/backup \
  alpine tar czf /backup/buildbureau-data.tar.gz -C /data .
```

### Restore Data

```bash
# Restore from backup
docker run --rm \
  -v buildbureau-data:/data \
  -v $(pwd)/backup:/backup \
  alpine tar xzf /backup/buildbureau-data.tar.gz -C /data
```

## 🐛 Troubleshooting

### Container Won't Start

```bash
# Check logs for errors
docker logs buildbureau

# Verify environment variables
docker inspect buildbureau | grep -A 20 "Env"

# Check port conflicts
docker ps -a | grep 8080
```

### Out of Memory

```bash
# Increase memory limit in docker-compose.yml
deploy:
  resources:
    limits:
      memory: 4G  # Increase from 2G
```

### Permission Issues

```bash
# Fix volume permissions
docker run --rm \
  -v buildbureau-data:/data \
  alpine chown -R 1000:1000 /data
```

### Build Fails

```bash
# Build without cache
docker build --no-cache -t buildbureau:latest .

# Check Docker version
docker version

# Ensure enough disk space
docker system df
```

## 🌐 Network Configuration

### Expose Additional Ports

```yaml
ports:
  - "8080:8080" # gRPC
  - "8081:8081" # Metrics (if implemented)
  - "9090:9090" # Debug port
```

### Connect to Other Containers

```yaml
networks:
  - buildbureau-network

networks:
  buildbureau-network:
    driver: bridge
```

## 📝 Examples

### Example 1: Development Setup

```bash
# docker-compose.dev.yml
version: '3.8'
services:
  buildbureau:
    build: .
    volumes:
      - ./config.yaml:/app/config/config.yaml
      - ./agents:/app/agents
      - ./data:/app/data
    environment:
      - GEMINI_API_KEY=${GEMINI_API_KEY}
    ports:
      - "8080:8080"
```

### Example 2: Production with All Providers

```yaml
environment:
  - GEMINI_API_KEY=${GEMINI_API_KEY}
  - OPENAI_API_KEY=${OPENAI_API_KEY}
  - CLAUDE_API_KEY=${CLAUDE_API_KEY}
  - SLACK_TOKEN=${SLACK_TOKEN}
```

### Example 3: With External Vald

```yaml
environment:
  - VALD_HOST=vald.production.local
  - VALD_PORT=8081
```

## 🎓 Best Practices

1. **Use Docker Compose**: Easier management and configuration
2. **Set Resource Limits**: Prevent resource exhaustion
3. **Use Volumes**: Persist data across restarts
4. **Enable Health Checks**: Monitor container health
5. **Backup Regularly**: Protect your data
6. **Use .env Files**: Keep secrets out of compose files
7. **Tag Images**: Version your builds
8. **Monitor Logs**: Watch for issues
9. **Update Regularly**: Get security patches
10. **Test First**: Try in dev before production

## 📚 Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [BuildBureau Documentation](../README.md)
- [Configuration Guide](../README.md#configuration)

## 🆘 Support

If you encounter issues:

1. Check container logs: `docker logs buildbureau`
2. Verify environment variables
3. Check resource usage: `docker stats`
4. Review health checks: `docker inspect buildbureau`
5. Consult troubleshooting section above

---

**Zero dependencies. Maximum convenience. That's Docker! 🐳**
