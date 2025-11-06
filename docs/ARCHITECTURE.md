# Architecture Documentation

## Overview

Onlineliga Discord Bot follows a clean, layered architecture pattern that separates concerns across distinct packages.
The design emphasizes modularity, testability, and maintainability.

## System Architecture

### High-Level Architecture

````text

┌──────────────────────────────────────────────────────────┐
│                     Discord Platform                      │
└───────────────────────────┬──────────────────────────────┘
                            │
                            │ Slash Command / Scheduled Task
                            │
┌───────────────────────────▼──────────────────────────────┐
│                    main.go (Entry Point)                  │
│  • Session Management                                     │
│  • Command Registration                                   │
│  • Interaction Handling                                   │
│  • Scheduler Initialization                               │
│  • Signal Handling (Shutdown)                            │
└─────┬──────────────────────────────┬──────────────────────┘
      │                              │
      │ User Command                 │ Cron Schedule
      │                              │
      │                    ┌─────────▼─────────────┐
      │                    │  internal/scheduler   │
      │                    │  • Cron Management    │
      │                    │  • Config Parsing     │
      │                    │  • Task Execution     │
      │                    └─────────┬─────────────┘
      │                              │
      └──────────────┬───────────────┘
                     │
                     │ ScrapeMatchResults()
                     │
┌────────────────────▼─────────────────────────────────────┐
│              internal/scraper (Data Layer)                │
│  • HTTP Request Management                                │
│  • Multi-User Coordination                               │
│  • Error Handling & Resilience                           │
└───────────────────────────┬──────────────────────────────┘
                            │
                            │ GET /apiv1/team/overview
                            │
┌───────────────────────────▼──────────────────────────────┐
│                  Onlineliga API                           │
│  • Match Data                                             │
│  • League Tables                                          │
│  • User Information                                       │
└───────────────────────────┬──────────────────────────────┘
                            │
                            │ JSON Response
                            │
┌───────────────────────────▼──────────────────────────────┐
│              internal/parse (Transform Layer)             │
│  • JSON Unmarshaling                                      │
│  • Data Model Transformation                             │
│  • Business Logic                                         │
└───────────────────────────┬──────────────────────────────┘
                            │
                            │ MatchResult Struct
                            │
┌───────────────────────────▼──────────────────────────────┐
│          internal/formatutils (Presentation Layer)        │
│  • Image Generation                                       │
│  • Badge Downloading                                      │
│  • Text Rendering                                         │
│  • PNG Encoding                                          │
└───────────────────────────┬──────────────────────────────┘
                            │
                            │ PNG Image Buffer
                            │
┌───────────────────────────▼──────────────────────────────┐
│                     Discord Platform                      │
│               (Spoilered Image Attachment)               │
└──────────────────────────────────────────────────────────┘

```text

## Component Details

### 1. Entry Point (main.go)

**Responsibilities:**

- Discord session lifecycle management
- Slash command registration per guild
- Interaction event handling
- Graceful shutdown coordination
- Logging configuration

**Key Functions:**

- `main()`: Bootstrap and lifecycle management
- `onInteractionCreate()`: Handle `/results` command interactions
- `getBaseURL()`: Map location parameter to Onlineliga regional URL

**Dependencies:**

- `github.com/bwmarrin/discordgo`: Discord API client
- `github.com/joho/godotenv`: Environment configuration
- `github.com/sirupsen/logrus`: Structured logging
- `internal/scraper`: Data fetching
- `internal/formatutils`: Result formatting

**Design Patterns:**

- Dependency Injection (logger)
- Event-driven architecture (Discord callbacks)
- Signal handling for graceful shutdown

### 2. Scraper Package (internal/scraper)

**Purpose:** Abstracts HTTP communication with Onlineliga API

**Struct:**

```go
type Scraper struct {
    client *http.Client
    logger *logrus.Logger
    url    string
}

```text

**Public API:**

- `NewScraper(logger) Scraper`: Constructor with dependency injection
- `ScrapeResults(userIDs, baseURL) [][]string`: Legacy batch fetcher
- `ScrapeResult(userID, baseURL) ([]string, error)`: Legacy single fetch
- `ScrapeMatchResult(userID, baseURL) (MatchResult, error)`: Single fetch
- `ScrapeMatchResults(userIDs, baseURL) []MatchResult`: Batch fetcher (current)

**Error Handling:**

- Continues processing on individual user failures
- Logs errors with context (user ID, URL)
- Returns partial results on partial success
- Defers cleanup with error logging

**HTTP Communication:**

- Uses shared `httpclient.DefaultHTTPClient`
- Endpoint: `{baseURL}/apiv1/team/overview?userId={id}`
- Response: JSON payload with match and league data

### 3. Parse Package (internal/parse)

**Purpose:** Transform Onlineliga API responses into structured data

**Data Models:**

```go
// API Response Models
Root          // Top-level response container
├── MatchData // Match information container
│   └── Match // Individual match details
├── LeagueTables // League standings array
└── User      // User profile information
    ├── League // League metadata
    └── Badge  // User badge information

// Output Models
MatchResult   // Formatted result for display

```text

**Public API:**

- `Result(jsonBody, userID) ([]string, error)`: Legacy format transformer
- `ResultObject(jsonBody, userID) (MatchResult, error)`: Current transformer

**Business Logic:**

- `filterLeagueTableForTeam()`: Extract user's league position
- `getMatchState()`: Determine Win/Loss/Draw status

**Responsibilities:**

- JSON unmarshaling with error handling
- Data validation and extraction
- Match state determination
- League position filtering
- Badge URL processing

### 4. Formatutils Package (internal/formatutils)

**Purpose:** Generate visual representations of match results

**Configuration:**

```go
type ImageConfig struct {
    FontSize    float64
    Margin      float64
    ColWidth    float64
    RowHeight   float64
    BadgeWidth  float64
    BadgeHeight float64
    R, G, B     float64  // Background color
}

```text

**Public API:**

- `MatchResultsToImage(results, config) (*bytes.Buffer, error)`: Main generator
- `SortResults(results, logger) []MatchResult`: Result sorting

**Internal Operations:**

- `createNewContext()`: Initialize 2D graphics context
- `loadFontFace()`: Load and configure font rendering
- `calculateImageSize()`: Dynamic sizing based on content
- `writeMatchResult()`: Render individual result row
- `drawBadge()`: Download and render user badges
- `encodeToPNG()`: Convert graphics to PNG format

**Graphics Pipeline:**

1. Calculate image dimensions based on result count
2. Create graphics context with fogleman/gg
3. Load font from FONT_PATH environment variable
4. Draw background and text for each result
5. Download and composite user badges
6. Encode final image to PNG buffer

**Dependencies:**

- `github.com/fogleman/gg`: 2D graphics rendering
- `golang.org/x/image`: Image manipulation
- Font: Cascadia Code (embedded in Docker image)

### 5. HTTPClient Package (internal/httpclient)

**Purpose:** Provide configured HTTP client for API requests

**Public API:**

- `DefaultHTTPClient`: Pre-configured `*http.Client` with cookie jar

**Configuration:**

- Cookie handling enabled
- Reusable for connection pooling
- Shared across scraper instances

### 6. Scheduler Package (internal/scheduler)

**Purpose:** Automate periodic execution of result fetching and posting

**Public API:**

- `New(session, scraper, logger)`: Create scheduler instance
- `LoadAndStart(configPath)`: Load schedules and start cron jobs
- `Stop()`: Gracefully stop scheduler and running jobs

**Configuration:**

- YAML-based configuration (`schedules.yaml`)
- Supports multiple schedules with individual cron expressions
- Per-schedule channel targeting
- Enable/disable individual schedules

**Schedule Configuration:**

```yaml
schedules:

  - name: "Daily UK League Results"

    enabled: true
    cron: "0 6 * * *"         # Every day at 6:00 AM UTC
    channel_id: "123456789"   # Target Discord channel
    command:
      location: ".co.uk"
      users: [2597, 77, 138]

```text

**Cron Expression Format:**

- Standard cron format with optional seconds field
- Examples:
  - `"0 6 * * *"` - Daily at 6:00 AM
  - `"0 */2 * * *"` - Every 2 hours
  - `"30 9 * * 1-5"` - Weekdays at 9:30 AM

**Execution Flow:**

1. Parse `schedules.yaml` configuration file
2. Validate schedule configurations
3. Register enabled schedules with cron scheduler
4. On schedule trigger:
   - Execute `ScrapeMatchResults()` with configured users
   - Generate image with `MatchResultsToImage()`
   - Post to configured Discord channel
   - Log execution status

**Error Handling:**

- Individual schedule failures don't affect other schedules
- Failed scrapes send error message to channel
- Logs all execution attempts and outcomes
- Graceful shutdown waits for running jobs

**Dependencies:**

- `github.com/robfig/cron/v3`: Cron job scheduling
- `gopkg.in/yaml.v3`: Configuration parsing

## Data Flow

### Request Flow (Detailed)

1. **User Input** → Discord slash command with location and user IDs

   ```text

   /results location:.de users:12345 67890

   ```text

2. **Discord Event** → `onInteractionCreate()` receives interaction

3. **Initial Response** → Send acknowledgment to Discord

   ```go
   s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
       Type: discordgo.InteractionResponseChannelMessageWithSource,
       Data: &discordgo.InteractionResponseData{
           Content: "content",
       },
   })

   ```text

4. **Data Fetching** → `scraper.ScrapeMatchResults(userIDs, baseURL)`
   - For each user ID:
     - HTTP GET to `/apiv1/team/overview?userId={id}`
     - Receive JSON response
     - Parse with `parse.ResultObject()`
     - Aggregate results

5. **Sorting** → `formatutils.SortResults(results, logger)`
   - Sort by league level and position
   - Ensure consistent ordering

6. **Image Generation** → `formatutils.MatchResultsToImage(results, config)`
   - Calculate image dimensions
   - Render text and badges
   - Encode to PNG

7. **Response Update** → Edit original interaction response

   ```go
   s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
       Content: &content,
       Files: []*discordgo.File{&filePayload},
   })

   ```text

### Error Flow

Errors are handled at each layer with appropriate logging:

```text

Error Occurs
    ↓
Log with Context (logrus)
    ↓
Return Error Up Stack
    ↓
Graceful Degradation (continue with remaining users)
    ↓
Partial Results Returned (if any successful)

```text

## Design Patterns

### Dependency Injection

Logger is injected into scraper to enable testing and configuration:

```go
olScraper := scraper.NewScraper(logger)

```text

### Resource Cleanup

Deferred cleanup ensures resources are released:

```go
defer func(Body io.ReadCloser) {
    if err := Body.Close(); err != nil {
        logger.WithError(err).Warn("cleanup failed")
    }
}(resp.Body)

```text

### Error Handling

Explicit error handling with structured logging:

```go
if scrapeErr != nil {
    s.logger.WithError(scrapeErr).Errorf("Scraping failed for user %s", userID)
    continue  // Graceful degradation
}

```text

### Configuration

Environment-based configuration with sensible defaults:

```go
token := os.Getenv("DISCO_BOT_TOKEN")
baseURL := getBaseURL(location) // Defaults to .de

```text

## Concurrency Model

### Current Implementation

The bot operates **synchronously** for simplicity:

- One interaction handled at a time per goroutine
- Users fetched sequentially in `ScrapeMatchResults()`
- No explicit concurrency primitives

### Goroutines

Discord library handles goroutines internally:

- Separate goroutine per interaction event
- Websocket management in background
- Signal handling in main goroutine

### Potential Improvements

Future optimization opportunities:

- Parallel user fetching with `sync.WaitGroup`
- Request coalescing for duplicate user IDs
- Result caching with TTL

## Deployment Architecture

### Container Structure

```text

┌─────────────────────────────────────┐
│         Builder Stage                │
│  • Go 1.23 base image               │
│  • Dependency download               │
│  • Font installation                 │
│  • Binary compilation                │
└─────────────┬───────────────────────┘
              │
              │ Copy binary + fonts
              │
┌─────────────▼───────────────────────┐
│         Runtime Stage                │
│  • Alpine Linux minimal base        │
│  • Application binary                │
│  • Cascadia Code font                │
│  • No build tools                    │
└─────────────────────────────────────┘

```text

### Multi-Platform Support

Build process supports multiple architectures:

- Platform detection via `BUILDPLATFORM` and `TARGETPLATFORM`
- Cross-compilation with `GOOS` and `GOARCH`
- Unified Docker image with manifest list

### CI/CD Pipeline

GitHub Actions workflow:

1. Checkout code
2. Set up Docker Buildx (multi-platform)
3. Authenticate with GHCR
4. Build for linux/amd64 and linux/arm64
5. Push tagged images to ghcr.io

## Security Considerations

### Token Management

- Bot token stored in environment variable (never hardcoded)
- Not logged or exposed in error messages
- Loaded from `.env` file or environment

### API Security

- HTTPS-only communication with Onlineliga
- No user authentication required (public API)
- Rate limiting handled gracefully

### Container Security

- Multi-stage build minimizes attack surface
- No build tools in final image
- Alpine Linux base for minimal footprint
- `CGO_ENABLED=0` for static binary (no C dependencies)

### Input Validation

- Discord validates slash command parameters
- User IDs converted to integers with error handling
- Location parameter validated against known regions

## Performance Characteristics

### Latency

Typical response time breakdown:

- API fetching: ~100-500ms per user
- Image generation: ~50-200ms
- Discord API: ~100-300ms
- **Total: ~250-1000ms for 1-3 users**

### Resource Usage

- Memory: ~20-50MB per bot instance
- CPU: Minimal (mostly I/O bound)
- Network: ~10-50KB per user query
- Disk: Static binary (~10-20MB)

### Scalability

Current limitations:

- Sequential user fetching limits throughput
- No caching (every request hits API)
- Single bot instance per deployment

Horizontal scaling:

- Stateless design enables multiple instances
- Each instance handles independent guilds
- No shared state or coordination required

## Testing Strategy

### Current Coverage

- `internal/parse` package has unit tests
- Test data in `internal/parse/files/parse_test.json`
- Other packages lack explicit tests

### Testing Approaches

**Unit Testing:**

- Mock HTTP responses for scraper tests
- Parse logic with various JSON payloads
- Image generation with different result counts

**Integration Testing:**

- End-to-end command flow
- Real API integration (staging environment)
- Docker container functionality

**Manual Testing:**

- Deploy to test Discord server
- Verify image rendering across platforms
- Test error scenarios (invalid IDs, API failures)

## Future Architecture Considerations

### Potential Enhancements

1. **Caching Layer**
   - Cache API responses with TTL
   - Reduce Onlineliga API load
   - Improve response time

2. **Async Processing**
   - Parallel user fetching
   - Background image generation
   - Non-blocking Discord responses

3. **Database Integration**
   - User preferences
   - Historical match data
   - Analytics and metrics

4. **Observability**
   - Prometheus metrics
   - Distributed tracing
   - Performance dashboards

5. **Extended Features**
   - More commands (leaderboards, statistics)
   - Customizable image themes
   - Notification system for new matches

### Migration Paths

**From HTML Scraping to API** (Completed)

- Old: Parse HTML structure
- New: Use official API endpoints
- Benefits: Reliability, performance, maintainability

#### Future: GraphQL API

- Could reduce over-fetching
- Single query for all users
- Better type safety

## Appendix

### Technology Stack Summary

| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| Language | Go | 1.23.3 | Core implementation |
| Discord SDK | discordgo | 0.28.1 | Discord integration |
| Graphics | fogleman/gg | 1.3.0 | Image generation |
| Logging | logrus | 1.9.3 | Structured logging |
| Config | godotenv | 1.5.1 | Environment management |
| Container | Docker | - | Deployment |
| CI/CD | GitHub Actions | - | Automation |
| Registry | GHCR | - | Image hosting |

### API Endpoints Reference

**Onlineliga API:**

- Base URL: `https://www.onlineliga.{tld}`
- Endpoint: `/apiv1/team/overview`
- Method: `GET`
- Query: `?userId={integer}`
- Response: `application/json`

**Discord API:**

- Used via discordgo library
- Slash commands API
- Interaction responses
- File attachments

---

**Last Updated:** 2025-11-05
**Architecture Version:** 2.0 (API-based)
````
