# Development Guide

This guide covers everything you need to contribute to the Onlineliga Discord Bot project.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Development Setup](#development-setup)
  - [Pre-commit Hooks](#5-set-up-pre-commit-hooks-recommended)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Code Style](#code-style)
- [Debugging](#debugging)
- [Common Tasks](#common-tasks)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### Required Tools

- **Go 1.23 or later** - [Install Go](https://go.dev/doc/install)
- **Git** - Version control
- **Docker** (optional) - For container testing
- **Make** (optional) - For convenience commands

### Recommended Tools

- **VS Code** with Go extension
- **GoLand** or other Go IDE
- **golangci-lint** - For code linting
- **go-fumpt** - For code formatting

### Discord Setup

1. Create a Discord application at [Discord Developer Portal](https://discord.com/developers/applications)
1. Create a bot user and copy the token
1. Enable necessary bot permissions:
   - Send Messages
   - Use Slash Commands
   - Attach Files
1. Invite bot to your test server

## Development Setup

### 1. Clone Repository

```bash
git clone https://github.com/tripleawwy/onlineliga_discord_bot.git
cd onlineliga_discord_bot
```

### 2. Install Dependencies

```bash
go mod download
go mod verify
```

### 3. Configure Environment

Create a `.env` file in the project root:

```env
DISCO_BOT_TOKEN=your_discord_bot_token_here
FONT_PATH=/path/to/your/font.ttf
```

**Font Setup (Local Development):**

Download Cascadia Code font:

```bash
mkdir -p ~/.fonts
cd ~/.fonts
curl -LO https://github.com/microsoft/cascadia-code/releases/download/v2111.01/CascadiaCode-2111.01.zip
unzip CascadiaCode-2111.01.zip
fc-cache -f -v
```

Update `.env`:

```env
FONT_PATH=$HOME/.fonts/ttf/static/CascadiaCode-Bold.ttf
```

### 4. Verify Setup

```bash
# Run tests
go test ./...

# Build the binary
go build -o bin/ol_discord_bot .

# Run the bot
./bin/ol_discord_bot
```

You should see:

```text
INFO Bot is now running. Press CTRL-C to exit.
```

### 5. Set Up Pre-commit Hooks (Recommended)

Pre-commit hooks ensure code quality and consistency before commits.

**Install pre-commit:**

```bash
# Using pip
pip install pre-commit

# Using homebrew (macOS)
brew install pre-commit

# Using apt (Ubuntu/Debian)
sudo apt install pre-commit
```

**Install hooks:**

```bash
# Install the git hooks
pre-commit install

# Run hooks manually on all files
pre-commit run --all-files
```

**What's included:**

- **File validation**: End-of-file fixer, trailing whitespace, line ending normalization
- **Security**: Detect private keys, check for large files
- **Go linting**: golangci-lint with comprehensive checks
- **Go formatting**: gofmt and go mod tidy
- **Markdown**: mdformat for formatting, markdownlint for linting

**Configuration files:**

- `.pre-commit-config.yaml` - Hook configuration
- `.golangci.yml` - Go linting rules
- `.markdownlint.yaml` - Markdown linting rules

**Skip hooks (if needed):**

```bash
# Skip all hooks for a commit
git commit --no-verify -m "message"

# Skip specific hook
SKIP=golangci-lint git commit -m "message"
```

## Project Structure

```text
onlineliga_discord_bot/
├── main.go                       # Entry point
├── go.mod                        # Go module definition
├── go.sum                        # Dependency checksums
├── .env                          # Environment configuration (gitignored)
├── .editorconfig                 # Editor configuration
├── .gitignore                    # Git ignore rules
│
├── internal/                     # Private packages
│   ├── scraper/                  # HTTP data fetching
│   │   └── scraper.go
│   ├── parse/                    # JSON parsing & transformation
│   │   ├── models.go             # Data models
│   │   ├── parse_result.go       # Parsing logic
│   │   ├── parse_result_test.go  # Unit tests
│   │   ├── parse_error.go        # Error handling
│   │   └── files/                # Test data
│   │       └── parse_test.json
│   ├── formatutils/              # Image generation
│   │   ├── format_results.go     # Main renderer
│   │   ├── format_error.go       # Error formatting
│   │   └── sort_results.go       # Result sorting
│   └── httpclient/               # HTTP client config
│       ├── client.go             # Default client
│       └── cookie.go             # Cookie handling
│
├── scripts/                      # Build and utility scripts
│   └── build.sh                  # Cross-platform build
│
├── docs/                         # Documentation
│   ├── api/                      # API documentation
│   └── images/                   # Diagrams and images
│
├── .github/                      # GitHub configuration
│   └── workflows/                # CI/CD workflows
│       └── docker_push.yml       # Docker build pipeline
│
├── Dockerfile                    # Container image definition
├── compose.yml                   # Docker Compose config
│
├── README.md                     # Project overview
├── ARCHITECTURE.md               # System design docs
├── DEVELOPMENT.md                # This file
└── DEPLOYMENT.md                 # Operations guide
```

## Development Workflow

### Git Workflow

This project follows a standard Git flow:

1. **Create Feature Branch**

   ```bash
   git checkout -b feature/your-feature-name
   ```

1. **Make Changes**

   - Write code
   - Add tests
   - Update documentation

1. **Commit Changes**

   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

   **Commit Message Format:**

   ```text
   <type>: <description>

   [optional body]
   ```

   Types: `feat`, `fix`, `docs`, `test`, `refactor`, `chore`

1. **Push and Create PR**

   ```bash
   git push origin feature/your-feature-name
   ```

### Code Review Process

- PRs require review before merging
- CI must pass (Docker build)
- Code should follow style guidelines
- Tests should pass locally and in CI

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/parse

# Run with verbose output
go test -v ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Writing Tests

Place test files next to the code they test:

```text
internal/parse/
├── parse_result.go
└── parse_result_test.go
```

**Test File Template:**

```go
package parse

import (
    "testing"
)

func TestResultObject(t *testing.T) {
    tests := []struct {
        name     string
        input    []byte
        userID   int
        expected MatchResult
        wantErr  bool
    }{
        {
            name:    "valid match result",
            input:   []byte(`{...}`),
            userID:  12345,
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := ResultObject(tt.input, tt.userID)
            if (err != nil) != tt.wantErr {
                t.Errorf("ResultObject() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            // Add assertions
        })
    }
}
```

### Test Data

Store test fixtures in `internal/*/files/`:

```go
testData, err := os.ReadFile("files/parse_test.json")
if err != nil {
    t.Fatal(err)
}
```

## Code Style

### Go Formatting

```bash
# Format all Go files
go fmt ./...

# Or use gofumpt (stricter formatting)
gofumpt -w .
```

### Linting

```bash
# Run go vet
go vet ./...

# Run golangci-lint (if installed)
golangci-lint run

# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Code Conventions

**Naming:**

- Exported: `PascalCase`
- Unexported: `camelCase`
- Files: `snake_case.go`
- Test files: `*_test.go`

**Comments:**

```go
// Package parse provides JSON parsing for Onlineliga API responses.
package parse

// MatchResult represents a formatted match result for display.
type MatchResult struct {
    // HomeTeam is the name of the home team
    HomeTeam string `json:"homeTeam"`
}

// ResultObject parses API response and returns a MatchResult.
// It returns an error if parsing fails.
func ResultObject(body []byte, userID int) (MatchResult, error) {
    // Implementation
}
```

**Error Handling:**

```go
if err != nil {
    logger.WithError(err).Error("descriptive context")
    return nil, fmt.Errorf("operation failed: %w", err)
}
```

**Resource Cleanup:**

```go
defer func(Body io.ReadCloser) {
    if err := Body.Close(); err != nil {
        logger.WithError(err).Warn("cleanup failed")
    }
}(resp.Body)
```

## Debugging

### Logging Levels

Adjust log level in `main.go`:

```go
logger.Level = logrus.DebugLevel  // Verbose
logger.Level = logrus.InfoLevel   // Default
logger.Level = logrus.WarnLevel   // Quiet
```

### Debug Mode

Run with environment variable:

```bash
LOG_LEVEL=debug go run main.go
```

### VS Code Launch Configuration

Create `.vscode/launch.json`:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Bot",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/main.go",
            "env": {
                "DISCO_BOT_TOKEN": "your_token",
                "FONT_PATH": "/path/to/font.ttf"
            }
        }
    ]
}
```

### Common Debugging Scenarios

**Discord Command Not Working:**

1. Check bot is online in Discord
1. Verify token is correct
1. Ensure bot has permissions
1. Check logs for registration errors

**API Request Failing:**

```go
// Add detailed logging
logger.WithFields(logrus.Fields{
    "userID": userID,
    "url": overviewURL,
}).Info("Making API request")
```

**Image Not Rendering:**

1. Verify FONT_PATH is correct
1. Check font file exists and is readable
1. Examine image generation logs
1. Test with minimal ImageConfig

## Common Tasks

### Adding a New Command

1. **Define Command Structure** (main.go):

   ```go
   commands := []*discordgo.ApplicationCommand{
       {
           Name: "newcommand",
           Description: "Description of new command",
           Options: []*discordgo.ApplicationCommandOption{
               // Define options
           },
       },
   }
   ```

1. **Handle Command** (main.go):

   ```go
   func onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
       switch i.ApplicationCommandData().Name {
       case "results":
           // Existing handler
       case "newcommand":
           // New handler
       }
   }
   ```

1. **Test Command:**

   - Restart bot
   - Run `/newcommand` in Discord
   - Check logs for issues

### Adding New Package

1. **Create Package Directory:**

   ```bash
   mkdir -p internal/newpackage
   ```

1. **Create Package File:**

   ```go
   // internal/newpackage/newpackage.go
   package newpackage

   // Public API
   func NewFunction() error {
       return nil
   }
   ```

1. **Add Tests:**

   ```bash
   touch internal/newpackage/newpackage_test.go
   ```

1. **Import in Main:**

   ```go
   import "github.com/tripleawwy/onlineliga_discord_bot/internal/newpackage"
   ```

### Updating Dependencies

```bash
# Update specific dependency
go get -u github.com/bwmarrin/discordgo

# Update all dependencies (minor/patch)
go get -u ./...

# Update all dependencies (including major)
go get -u=patch ./...

# Tidy module files
go mod tidy

# Verify dependencies
go mod verify
```

### Building for Production

```bash
# Build with optimizations
go build -ldflags="-w -s" -o bin/ol_discord_bot .

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o bin/ol_discord_bot_linux_amd64 .

# Use build script
./scripts/build.sh linux/amd64
```

### Docker Development

```bash
# Build Docker image locally
docker build -t onlineliga-bot:dev .

# Run container with .env
docker run --env-file .env onlineliga-bot:dev

# Build multi-platform
docker buildx build --platform linux/amd64,linux/arm64 -t onlineliga-bot:latest .

# Use Docker Compose
docker compose up --build
```

## Troubleshooting

### Build Issues

#### "cannot find package"

```bash
go mod download
go mod tidy
```

#### "undefined: X"

- Check import statements
- Verify function/variable is exported (capitalized)
- Run `go mod tidy`

### Runtime Issues

**Bot doesn't connect:**

- Verify DISCO_BOT_TOKEN is valid
- Check internet connectivity
- Review Discord API status

**Commands don't appear:**

- Commands register on bot startup
- Restart bot to refresh commands
- Check guild ID is correct
- Allow time for Discord to propagate

**Font errors:**

- Ensure FONT_PATH points to valid .ttf file
- Check file permissions
- Try absolute path

### Performance Issues

**Slow responses:**

- Check network latency to Onlineliga API
- Review logs for timeout errors
- Consider adding caching

**High memory usage:**

- Profile with `go tool pprof`
- Check for goroutine leaks
- Review image generation efficiency

## IDE Configuration

### VS Code Settings

Create `.vscode/settings.json`:

```json
{
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.lintOnSave": "workspace",
    "go.formatTool": "gofumpt",
    "go.testFlags": ["-v"],
    "editor.formatOnSave": true,
    "[go]": {
        "editor.codeActionsOnSave": {
            "source.organizeImports": true
        }
    }
}
```

### GoLand Configuration

1. Enable Go modules: Preferences → Go → Go Modules
1. Set GOROOT to Go 1.23+
1. Enable gofmt on save
1. Configure run configuration with env vars

## Resources

### Documentation

- [Go Documentation](https://go.dev/doc/)
- [DiscordGo Documentation](https://pkg.go.dev/github.com/bwmarrin/discordgo)
- [Discord Developer Docs](https://discord.com/developers/docs)

### Tools

- [Go Playground](https://go.dev/play/) - Test code snippets
- [Regex101](https://regex101.com/) - Test regex patterns
- [JSON Formatter](https://jsonformatter.org/) - Format API responses

### Community

- [Go Discord Server](https://discord.gg/golang)
- [DiscordGo GitHub](https://github.com/bwmarrin/discordgo)

## Getting Help

If you encounter issues:

1. **Check Documentation** - README, ARCHITECTURE, this guide
1. **Search Issues** - GitHub issues for similar problems
1. **Enable Debug Logging** - Set `LOG_LEVEL=debug`
1. **Create Issue** - With logs, steps to reproduce, environment details

______________________________________________________________________

**Happy Coding!** 🚀

For deployment information, see [DEPLOYMENT.md](DEPLOYMENT.md).
For architecture details, see [ARCHITECTURE.md](ARCHITECTURE.md).
