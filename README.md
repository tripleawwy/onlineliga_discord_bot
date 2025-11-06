# Onlineliga Discord Bot

A Discord bot that fetches football match statistics from the Onlineliga gaming platform and presents them as
beautifully formatted images.

[![Docker](https://img.shields.io/badge/docker-ready-blue.svg)](https://github.com/tripleawwy/onlineliga_discord_bot/pkgs/container/onlineliga_discord_bot)
[![Go Version](https://img.shields.io/badge/go-1.23-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## Features

- 🎮 Fetch match results from Onlineliga gaming platform
- 🌍 Multi-region support (Germany, UK, Austria, Switzerland)
- 🖼️ Generate formatted PNG images with match statistics
- 🏆 Display league information, positions, and badges
- ⏰ Automated scheduled tasks with cron-based scheduler
- 🚀 Multi-platform Docker support (amd64/arm64)
- 📊 Structured logging with detailed error handling

## Quick Start

### Using Docker Compose (Recommended)

1. Clone the repository:

```bash
git clone https://github.com/tripleawwy/onlineliga_discord_bot.git
cd onlineliga_discord_bot
```

1. Create a `.env` file with your Discord bot token:

```env
DISCO_BOT_TOKEN=your_discord_bot_token_here
FONT_PATH=/usr/share/fonts/ttf/static/CascadiaCode-Bold.ttf
```

1. Start the bot:

```bash
docker compose up -d
```

### Using Docker

Pull and run the pre-built image from GitHub Container Registry:

```bash
docker pull ghcr.io/tripleawwy/onlineliga_discord_bot:latest

docker run -d \
  -e DISCO_BOT_TOKEN=your_token_here \
  -e FONT_PATH=/usr/share/fonts/ttf/static/CascadiaCode-Bold.ttf \
  ghcr.io/tripleawwy/onlineliga_discord_bot:latest
```

### Building from Source

#### Prerequisites

- Go 1.23 or later
- Git

#### Build Steps

```bash
# Clone the repository
git clone https://github.com/tripleawwy/onlineliga_discord_bot.git
cd onlineliga_discord_bot

# Download dependencies
go mod download

# Build the binary
go build -o bin/ol_discord_bot .

# Run the bot
DISCO_BOT_TOKEN=your_token_here ./bin/ol_discord_bot
```

#### Cross-Platform Build

Use the provided build script for cross-platform compilation:

```bash
# Build for Linux AMD64
./scripts/build.sh linux/amd64

# Build for Linux ARM64
./scripts/build.sh linux/arm64
```

## Usage

### Discord Commands

The bot provides a single slash command:

```text
/results location:<region> users:<user_ids>
```

**Parameters:**

- `location`: Onlineliga region (`.de`, `.co.uk`, `.at`, `.ch`)
- `users`: Space-separated list of Onlineliga user IDs

**Example:**

```text
/results location:.de users:12345 67890 11223
```

**Response:**
The bot will generate a formatted PNG image showing:

- League information and level
- User badges
- League positions
- Match results (home team, score, away team)
- Match state (Win, Loss, Draw)
- Points earned

The image is sent as a spoilered attachment to prevent accidental reveals.

## Configuration

### Environment Variables

| Variable          | Required | Default | Description                           |
| ----------------- | -------- | ------- | ------------------------------------- |
| `DISCO_BOT_TOKEN` | Yes      | -       | Discord bot authentication token      |
| `FONT_PATH`       | Yes\*    | -       | Path to font file for image rendering |

\*Required in Docker, optional for local development if system fonts are available.

### Supported Regions

| Region   | Base URL                         | Description    |
| -------- | -------------------------------- | -------------- |
| `.de`    | <https://www.onlineliga.de>      | Germany        |
| `.co.uk` | <https://www.onlineleague.co.uk> | United Kingdom |
| `.at`    | <https://www.onlineliga.at>      | Austria        |
| `.ch`    | <https://www.onlineliga.ch>      | Switzerland    |

### Scheduler Configuration (Optional)

The bot supports automated scheduled tasks via a `schedules.yaml` configuration file. This allows you to
automatically post results to specific channels at scheduled times.

**Example `schedules.yaml`:**

```yaml
schedules:
  - name: "Daily UK League Results"
    enabled: true
    cron: "0 6 * * *"  # Every day at 6:00 AM UTC
    channel_id: "YOUR_CHANNEL_ID_HERE"
    command:
      location: ".co.uk"
      users: [2597, 77, 138, 5236, 8247]
```

**Configuration:**

1. Create `schedules.yaml` in the project root
1. Add your schedule configuration(s)
1. Replace `YOUR_CHANNEL_ID_HERE` with your Discord channel ID
1. Optionally set `SCHEDULES_PATH` environment variable to use a custom path

**Cron Expression Format:**

- `"0 6 * * *"` - Every day at 6:00 AM UTC
- `"0 */2 * * *"` - Every 2 hours
- `"30 9 * * 1-5"` - 9:30 AM on weekdays
- `"0 0 * * 0"` - Midnight every Sunday

See the included `schedules.yaml` for more examples and documentation.

## Architecture

The bot is built with a clean, layered architecture:

```text
┌─────────────────┐     ┌─────────────────┐
│   Discord API   │────▶│    Scheduler    │ (internal/scheduler)
│    (main.go)    │     │  (Cron-based)   │
└────────┬────────┘     └────────┬────────┘
         │                       │
         └───────────┬───────────┘
                     │
         ┌───────────▼───────────┐
         │       Scraper         │ (internal/scraper)
         └───────────┬───────────┘
                     │
         ┌───────────▼───────────┐
         │        Parser         │ (internal/parse)
         └───────────┬───────────┘
                     │
         ┌───────────▼───────────┐
         │      Formatutils      │ (internal/formatutils)
         └───────────────────────┘
```

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for detailed design documentation.

## Development

For development setup, testing, and contribution guidelines, see:

- [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) - Developer guide
- [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md) - Deployment and operations
- [docs/api/](docs/api/) - API documentation for internal packages

### Project Structure

```text
onlineliga_discord_bot/
├── main.go                 # Entry point, Discord integration
├── schedules.yaml          # Scheduler configuration (optional)
├── internal/
│   ├── scraper/           # HTTP data fetching layer
│   ├── parse/             # JSON parsing & transformation
│   ├── formatutils/       # Image generation
│   ├── scheduler/         # Cron-based task scheduler
│   └── httpclient/        # HTTP client configuration
├── scripts/
│   └── build.sh           # Cross-platform build script
├── Dockerfile             # Multi-stage container build
└── compose.yml            # Docker Compose configuration
```

## Deployment

### Docker (Production)

The project includes CI/CD via GitHub Actions that automatically builds and publishes multi-platform Docker images
to GitHub Container Registry on every push to main.

**Supported Platforms:**

- `linux/amd64`
- `linux/arm64`

See [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md) for detailed deployment instructions.

## API Integration

The bot uses the Onlineliga API endpoint:

```text
/apiv1/team/overview?userId=<user_id>
```

This API provides:

- Match data (last match details)
- League table information
- User details and badges
- Team statistics

## Troubleshooting

### Bot doesn't respond to commands

1. Verify the bot token is correct
1. Ensure the bot has proper permissions in your Discord server
1. Check that commands are registered (they register per-guild on startup)
1. Review logs for error messages

### Font rendering issues

If images don't render properly, ensure `FONT_PATH` points to a valid TrueType font file. The Docker image includes
Cascadia Code by default.

### API rate limiting

The Onlineliga API may rate limit requests. The bot continues with the next user if one fails, ensuring partial
results are still returned.

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

### Development Guidelines

1. Follow Go best practices and project conventions
1. Add tests for new functionality
1. Update documentation for significant changes
1. Ensure Docker builds successfully
1. Test with multiple user IDs and regions

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Built with [DiscordGo](https://github.com/bwmarrin/discordgo)
- Image rendering by [fogleman/gg](https://github.com/fogleman/gg)
- Font: [Cascadia Code](https://github.com/microsoft/cascadia-code)

## Support

For issues or questions:

- Open an issue on GitHub
- Check existing documentation in the `docs/` folder

______________________________________________________________________

## Made with ❤️ for the Onlineliga community
