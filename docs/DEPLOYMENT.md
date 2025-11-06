# Deployment Guide

This guide covers production deployment, operations, and maintenance of the Onlineliga Discord Bot.

## Table of Contents

- [Deployment Options](#deployment-options)
- [Docker Deployment](#docker-deployment)
- [Docker Compose Deployment](#docker-compose-deployment)
- [Kubernetes Deployment](#kubernetes-deployment)
- [Binary Deployment](#binary-deployment)
- [CI/CD Pipeline](#cicd-pipeline)
- [Configuration](#configuration)
- [Monitoring](#monitoring)
- [Maintenance](#maintenance)
- [Troubleshooting](#troubleshooting)
- [Security](#security)

## Deployment Options

| Method         | Use Case                          | Complexity | Scalability |
| -------------- | --------------------------------- | ---------- | ----------- |
| Docker Compose | Small deployments, development    | Low        | Low         |
| Docker         | Single instance, simple setup     | Low        | Medium      |
| Kubernetes     | Multi-instance, high availability | High       | High        |
| Binary         | Bare metal, VM                    | Medium     | Low         |

## Docker Deployment

### Using Pre-built Images (Recommended)

Images are automatically built and published to GitHub Container Registry:

````bash
# Pull latest image
docker pull ghcr.io/tripleawwy/onlineliga_discord_bot:latest

# Run container
docker run -d \
  --name ol-discord-bot \
  --restart unless-stopped \
  -e DISCO_BOT_TOKEN=your_token_here \
  -e FONT_PATH=/usr/share/fonts/ttf/static/CascadiaCode-Bold.ttf \
  ghcr.io/tripleawwy/onlineliga_discord_bot:latest

```text

### Building Custom Image

```bash
# Clone repository
git clone https://github.com/tripleawwy/onlineliga_discord_bot.git
cd onlineliga_discord_bot

# Build image
docker build -t onlineliga-bot:custom .

# Run custom image
docker run -d \
  --name ol-discord-bot \
  --restart unless-stopped \
  --env-file .env \
  onlineliga-bot:custom

```text

### Multi-Platform Build

Build for multiple architectures:

```bash
# Set up Docker Buildx
docker buildx create --use

# Build and push multi-platform image
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t your-registry/onlineliga-bot:latest \
  --push .

```text

### Docker Run Options

```bash
docker run -d \
  --name ol-discord-bot \
  --restart unless-stopped \              # Auto-restart on failure
  --memory="256m" \                       # Memory limit
  --cpus="0.5" \                          # CPU limit
  --log-driver json-file \                # Logging driver
  --log-opt max-size=10m \                # Log rotation
  --log-opt max-file=3 \                  # Keep 3 log files
  -e DISCO_BOT_TOKEN=token \              # Bot token
  -e FONT_PATH=/usr/share/fonts/... \     # Font path
  ghcr.io/tripleawwy/onlineliga_discord_bot:latest

```text

## Docker Compose Deployment

### Basic Setup

1. **Create Project Directory:**

   ```bash
   mkdir onlineliga-bot
   cd onlineliga-bot

   ```text

2. **Create `compose.yml`:**

   ```yaml
   services:
     disco_bot:
       image: ghcr.io/tripleawwy/onlineliga_discord_bot:latest
       restart: unless-stopped
       environment:

         - DISCO_BOT_TOKEN=${DISCO_BOT_TOKEN:?error}
         - FONT_PATH=/usr/share/fonts/ttf/static/CascadiaCode-Bold.ttf

       deploy:
         resources:
           limits:
             memory: 256M
             cpus: '0.5'
       logging:
         driver: json-file
         options:
           max-size: "10m"
           max-file: "3"

   ```text

3. **Create `.env` File:**

   ```env
   DISCO_BOT_TOKEN=your_discord_bot_token_here

   ```text

4. **Deploy:**

   ```bash
   docker compose up -d

   ```text

### Management Commands

```bash
# Start bot
docker compose up -d

# View logs
docker compose logs -f

# Stop bot
docker compose stop

# Restart bot
docker compose restart

# Remove bot
docker compose down

# Update to latest image
docker compose pull
docker compose up -d

```text

### Production Configuration

```yaml
services:
  disco_bot:
    image: ghcr.io/tripleawwy/onlineliga_discord_bot:latest
    restart: unless-stopped

    environment:

      - DISCO_BOT_TOKEN=${DISCO_BOT_TOKEN:?error}
      - FONT_PATH=/usr/share/fonts/ttf/static/CascadiaCode-Bold.ttf

    deploy:
      resources:
        limits:
          memory: 512M
          cpus: '1.0'
        reservations:
          memory: 128M
          cpus: '0.25'

    logging:
      driver: json-file
      options:
        max-size: "20m"
        max-file: "5"

    healthcheck:
      test: ["CMD", "pgrep", "ol_discord_bot"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s

```text

## Kubernetes Deployment

### Basic Deployment

**deployment.yaml:**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: onlineliga-bot
  namespace: discord-bots
spec:
  replicas: 1
  selector:
    matchLabels:
      app: onlineliga-bot
  template:
    metadata:
      labels:
        app: onlineliga-bot
    spec:
      containers:

      - name: bot

        image: ghcr.io/tripleawwy/onlineliga_discord_bot:latest
        imagePullPolicy: Always
        env:

        - name: DISCO_BOT_TOKEN

          valueFrom:
            secretKeyRef:
              name: discord-secrets
              key: bot-token

        - name: FONT_PATH

          value: /usr/share/fonts/ttf/static/CascadiaCode-Bold.ttf
        resources:
          requests:
            memory: "128Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "1000m"

```text

**secret.yaml:**

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: discord-secrets
  namespace: discord-bots
type: Opaque
stringData:
  bot-token: your_discord_bot_token_here

```text

### Deployment Commands

```bash
# Create namespace
kubectl create namespace discord-bots

# Apply secret
kubectl apply -f secret.yaml

# Apply deployment
kubectl apply -f deployment.yaml

# Check status
kubectl get pods -n discord-bots

# View logs
kubectl logs -f deployment/onlineliga-bot -n discord-bots

# Scale (if needed)
kubectl scale deployment/onlineliga-bot --replicas=2 -n discord-bots

# Delete deployment
kubectl delete -f deployment.yaml

```text

## Binary Deployment

### Building Binary

```bash
# Build for current platform
go build -ldflags="-w -s" -o ol_discord_bot .

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ol_discord_bot .

```text

### Systemd Service

**Create Service File:** `/etc/systemd/system/onlineliga-bot.service`

```ini
[Unit]
Description=Onlineliga Discord Bot
After=network.target

[Service]
Type=simple
User=discord
Group=discord
WorkingDirectory=/opt/onlineliga-bot
Environment="DISCO_BOT_TOKEN=your_token"
Environment="FONT_PATH=/usr/share/fonts/truetype/cascadia/CascadiaCode-Bold.ttf"
ExecStart=/opt/onlineliga-bot/ol_discord_bot
Restart=on-failure
RestartSec=10s

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/onlineliga-bot/logs

[Install]
WantedBy=multi-user.target

```text

**Manage Service:**

```bash
# Enable service
sudo systemctl enable onlineliga-bot

# Start service
sudo systemctl start onlineliga-bot

# Check status
sudo systemctl status onlineliga-bot

# View logs
sudo journalctl -u onlineliga-bot -f

# Restart service
sudo systemctl restart onlineliga-bot

# Stop service
sudo systemctl stop onlineliga-bot

```text

## CI/CD Pipeline

### GitHub Actions Workflow

The project includes automated CI/CD via `.github/workflows/docker_push.yml`:

**Workflow Triggers:**

- Push to any branch
- Pull request creation

**Pipeline Steps:**

1. Checkout code
2. Set up Docker Buildx (multi-platform support)
3. Authenticate with GitHub Container Registry
4. Extract metadata (tags, labels)
5. Build Docker images for linux/amd64 and linux/arm64
6. Push images to ghcr.io

**Image Tags:**

- `main` branch → `ghcr.io/tripleawwy/onlineliga_discord_bot:latest`
- Other branches → `ghcr.io/tripleawwy/onlineliga_discord_bot:<branch>`
- Tags → `ghcr.io/tripleawwy/onlineliga_discord_bot:<tag>`

### Manual Workflow Trigger

```bash
# Push tag for release
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# Workflow automatically builds and pushes image tagged v1.0.0

```text

### Deploying Specific Version

```bash
# Pull specific version
docker pull ghcr.io/tripleawwy/onlineliga_discord_bot:v1.0.0

# Run specific version
docker run -d \
  -e DISCO_BOT_TOKEN=token \
  -e FONT_PATH=/usr/share/fonts/... \
  ghcr.io/tripleawwy/onlineliga_discord_bot:v1.0.0

```text

## Configuration

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DISCO_BOT_TOKEN` | Yes | - | Discord bot authentication token |
| `FONT_PATH` | Yes* | - | Absolute path to TrueType font file |

*Required in Docker; may use system fonts in binary deployment

### Font Configuration

**Docker (default):**

```text

FONT_PATH=/usr/share/fonts/ttf/static/CascadiaCode-Bold.ttf

```text

**Custom font in Docker:**

```bash
docker run -d \
  -v /path/to/your/font.ttf:/fonts/custom.ttf:ro \
  -e FONT_PATH=/fonts/custom.ttf \
  ghcr.io/tripleawwy/onlineliga_discord_bot:latest

```text

**System font (binary):**

```bash
# Linux
FONT_PATH=/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf

# macOS
FONT_PATH=/System/Library/Fonts/Helvetica.ttc

```text

### Logging Configuration

Adjust log level in code or via environment:

```go
// main.go
logger.Level = logrus.InfoLevel  // Options: Debug, Info, Warn, Error

```text

## Monitoring

### Docker Logs

```bash
# Follow logs
docker logs -f ol-discord-bot

# Last 100 lines
docker logs --tail 100 ol-discord-bot

# Logs since timestamp
docker logs --since 2024-01-01T00:00:00 ol-discord-bot

# Export logs
docker logs ol-discord-bot > bot-logs.txt

```text

### Docker Stats

```bash
# Real-time stats
docker stats ol-discord-bot

# One-time snapshot
docker stats --no-stream ol-discord-bot

```text

### Health Checks

**Docker health check:**

```bash
docker inspect --format='{{json .State.Health}}' ol-discord-bot | jq

```text

**Manual health check:**

```bash
# Check if process is running
docker exec ol-discord-bot ps aux | grep ol_discord_bot

# Check Discord connection
docker logs ol-discord-bot | grep "Bot is now running"

```text

### Metrics (Future Enhancement)

Consider adding Prometheus metrics:

- Request count
- Response time
- Error rate
- Active connections

## Maintenance

### Updates

**Docker:**

```bash
# Pull latest image
docker pull ghcr.io/tripleawwy/onlineliga_discord_bot:latest

# Stop current container
docker stop ol-discord-bot
docker rm ol-discord-bot

# Start with new image
docker run -d \
  --name ol-discord-bot \
  --restart unless-stopped \
  --env-file .env \
  ghcr.io/tripleawwy/onlineliga_discord_bot:latest

```text

**Docker Compose:**

```bash
docker compose pull
docker compose up -d

```text

**Kubernetes:**

```bash
kubectl rollout restart deployment/onlineliga-bot -n discord-bots

```text

### Backups

**Configuration Backup:**

```bash
# Backup .env and compose files
tar -czf onlineliga-bot-backup.tar.gz .env compose.yml

# Backup to remote
scp onlineliga-bot-backup.tar.gz user@backup-server:/backups/

```text

**No database backup needed** - bot is stateless

### Log Rotation

**Docker:**

```yaml
logging:
  driver: json-file
  options:
    max-size: "10m"
    max-file: "3"

```text

**Systemd:**

```bash
# Logs auto-rotate via journald
journalctl --vacuum-time=7d  # Keep 7 days

```text

## Troubleshooting

### Bot Not Starting

**Check logs:**

```bash
docker logs ol-discord-bot

```text

**Common issues:**

- Invalid bot token → Verify `DISCO_BOT_TOKEN`
- Font file missing → Check `FONT_PATH`
- Network issues → Verify internet connectivity
- Permission errors → Check file permissions

### High Memory Usage

**Check resource usage:**

```bash
docker stats ol-discord-bot

```text

**Solutions:**

- Set memory limits in docker-compose.yml
- Restart bot periodically
- Profile with `pprof` if persistent

### Command Registration Issues

**Symptoms:**

- Commands don't appear in Discord
- Bot responds to old commands

**Solutions:**

1. Restart bot to refresh commands
2. Verify bot is in guild
3. Check bot permissions
4. Wait for Discord propagation (~1 hour max)

### API Request Failures

**Check logs for:**

```text

Error responding to interaction
Scraping failed for user

```text

**Solutions:**

- Verify Onlineliga API is accessible
- Check for rate limiting
- Validate user IDs are correct
- Ensure network connectivity

## Security

### Token Security

**Never:**

- Commit `.env` to version control
- Share bot token publicly
- Log bot token

**Best Practices:**

- Use environment variables
- Rotate tokens periodically
- Use secrets management (Kubernetes, AWS Secrets Manager)

### Container Security

**Image Security:**

- Uses Alpine Linux (minimal attack surface)
- Multi-stage build (no build tools in final image)
- Static binary (`CGO_ENABLED=0`)
- Regular updates via CI/CD

**Runtime Security:**

- Run as non-root user (if needed)
- Read-only filesystem where possible
- Resource limits (memory, CPU)
- Network segmentation

### Access Control

**Docker:**

```bash
# Limit Docker socket access
# Only trusted users in docker group

```text

**Kubernetes:**

```yaml
# Use RBAC
# Namespace isolation
# Network policies

```text

## Performance Tuning

### Resource Allocation

**Small deployment (1-10 guilds):**

- Memory: 128-256MB
- CPU: 0.25-0.5 cores

**Medium deployment (10-100 guilds):**

- Memory: 256-512MB
- CPU: 0.5-1.0 cores

**Large deployment (100+ guilds):**

- Memory: 512MB-1GB
- CPU: 1.0-2.0 cores

### Optimization Tips

1. **Caching** - Cache API responses (future)
2. **Concurrency** - Parallel user fetching (future)
3. **Connection pooling** - Reuse HTTP connections (done)
4. **Image optimization** - Reduce image size/quality if needed

## Rollback Procedures

### Docker

```bash
# Pull previous version
docker pull ghcr.io/tripleawwy/onlineliga_discord_bot:v1.0.0

# Stop current
docker stop ol-discord-bot
docker rm ol-discord-bot

# Start previous version
docker run -d \
  --name ol-discord-bot \
  --restart unless-stopped \
  --env-file .env \
  ghcr.io/tripleawwy/onlineliga_discord_bot:v1.0.0

```text

### Kubernetes

```bash
# Rollback to previous revision
kubectl rollout undo deployment/onlineliga-bot -n discord-bots

# Rollback to specific revision
kubectl rollout undo deployment/onlineliga-bot --to-revision=2 -n discord-bots

# Check rollout history
kubectl rollout history deployment/onlineliga-bot -n discord-bots

```text

## Support

For deployment issues:

1. Check this guide
2. Review [DEVELOPMENT.md](DEVELOPMENT.md) for debugging
3. Check GitHub Issues
4. Open new issue with:
   - Deployment method
   - Environment details
   - Error logs
   - Steps to reproduce

---

**Deployment Version:** 2.0
**Last Updated:** 2025-11-05
````
