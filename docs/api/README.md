# API Documentation Index

Complete API documentation for Onlineliga Discord Bot internal packages.

## Package Overview

| Package       | Purpose                                | Stability | Documentation          |
| ------------- | -------------------------------------- | --------- | ---------------------- |
| `scraper`     | HTTP data fetching from Onlineliga API | Stable    | [View](scraper.md)     |
| `parse`       | JSON parsing and data transformation   | Stable    | [View](parse.md)       |
| `formatutils` | Image generation and result formatting | Stable    | [View](formatutils.md) |
| `httpclient`  | HTTP client configuration              | Stable    | [View](httpclient.md)  |

## Quick Links

### By Use Case

**Fetching Data:**

- [Scraper Package](scraper.md) - `ScrapeMatchResults()`
- [HTTPClient Package](httpclient.md) - HTTP configuration

**Processing Data:**

- [Parse Package](parse.md) - `ResultObject()`
- [Parse Models](parse.md#data-models) - Data structures

**Generating Images:**

- [Formatutils Package](formatutils.md) - `MatchResultsToImage()`
- [Image Configuration](formatutils.md#imageconfig) - Rendering settings

### By Package

#### [Scraper Package](scraper.md)

**Main Types:**

- `Scraper` - HTTP scraper with logger

**Key Functions:**

- `NewScraper(logger)` - Constructor
- `ScrapeMatchResults(userIDs, baseURL)` - Batch fetching
- `ScrapeMatchResult(userID, baseURL)` - Single fetch

**Use Cases:**

- Fetching match results from Onlineliga
- Multi-user data retrieval
- Error-resilient scraping

______________________________________________________________________

#### [Parse Package](parse.md)

**Main Types:**

- `Root` - API response container
- `MatchResult` - Formatted output
- `User`, `League`, `Match` - Data models

**Key Functions:**

- `ResultObject(body, userID)` - JSON to MatchResult
- `Result(body, userID)` - Legacy string format

**Use Cases:**

- Parsing Onlineliga API responses
- Transforming JSON to structured data
- Match state determination

______________________________________________________________________

#### [Formatutils Package](formatutils.md)

**Main Types:**

- `ImageConfig` - Rendering configuration

**Key Functions:**

- `MatchResultsToImage(results, config)` - PNG generation
- `SortResults(results, logger)` - Result ordering

**Use Cases:**

- Generating match result images
- Customizing image appearance
- Badge rendering

______________________________________________________________________

#### [HTTPClient Package](httpclient.md)

**Main Variables:**

- `DefaultHTTPClient` - Pre-configured client

**Features:**

- Cookie jar support
- Connection pooling
- Session management

**Use Cases:**

- Shared HTTP client
- Consistent API communication
- Cookie-based sessions

## Integration Examples

### Complete Flow

```go
package main

import (
    "github.com/sirupsen/logrus"
    "github.com/tripleawwy/onlineliga_discord_bot/internal/scraper"
    "github.com/tripleawwy/onlineliga_discord_bot/internal/formatutils"
)

func main() {
    logger := logrus.New()

    // 1. Create scraper
    s := scraper.NewScraper(logger)

    // 2. Fetch results
    userIDs := []string{"12345", "67890"}
    baseURL := "https://www.onlineliga.de"
    results := s.ScrapeMatchResults(userIDs, baseURL)

    // 3. Sort results
    sortedResults := formatutils.SortResults(results, logger)

    // 4. Generate image
    config := formatutils.ImageConfig{
        FontSize:    14.0,
        Margin:      28.0,
        ColWidth:    150.0,
        RowHeight:   28.0,
        BadgeWidth:  22.0,
        BadgeHeight: 22.0,
        R:           1.0,
        G:           1.0,
        B:           1.0,
    }

    imageBuf, err := formatutils.MatchResultsToImage(sortedResults, config)
    if err != nil {
        logger.WithError(err).Fatal("Image generation failed")
    }

    // 5. Use image buffer (send to Discord, save to file, etc.)
    _ = imageBuf
}

```

### Error Handling Pattern

```go
// Fetch with error handling
results := scraper.ScrapeMatchResults(userIDs, baseURL)
if len(results) == 0 {
    logger.Warn("No results fetched")
    return
}

// Sort
sortedResults := formatutils.SortResults(results, logger)

// Generate image with fallback
imageBuf, err := formatutils.MatchResultsToImage(sortedResults, config)
if err != nil {
    logger.WithError(err).Error("Image generation failed")
    // Fallback: return text representation
    return fallbackTextResponse(results)
}

// Success
return imageBuf

```

## Data Flow

```text
User IDs
    ↓
Scraper.ScrapeMatchResults()
    ↓
HTTP GET (httpclient.DefaultHTTPClient)
    ↓
parse.ResultObject()
    ↓
formatutils.SortResults()
    ↓
formatutils.MatchResultsToImage()
    ↓
PNG Image Buffer

```

## Testing

### Unit Testing Imports

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"

    "github.com/tripleawwy/onlineliga_discord_bot/internal/scraper"
    "github.com/tripleawwy/onlineliga_discord_bot/internal/parse"
    "github.com/tripleawwy/onlineliga_discord_bot/internal/formatutils"
)

```

### Mock HTTP Responses

```go
import "net/http/httptest"

server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"matchData": {...}}`))
}))
defer server.Close()

// Use server.URL in tests

```

## Common Patterns

### Multi-Region Support

```go
regions := map[string]string{
    "de": "https://www.onlineliga.de",
    "uk": "https://www.onlineleague.co.uk",
    "at": "https://www.onlineliga.at",
    "ch": "https://www.onlineliga.ch",
}

for name, url := range regions {
    results := scraper.ScrapeMatchResults(userIDs, url)
    logger.Infof("Region %s: %d results", name, len(results))
}

```

### Concurrent Scraping (Future)

```go
import "sync"

func scrapeMultipleRegions(userIDs []string, regions map[string]string) map[string][]parse.MatchResult {
    var wg sync.WaitGroup
    results := make(map[string][]parse.MatchResult)
    mu := sync.Mutex{}

    for name, url := range regions {
        wg.Add(1)
        go func(n, u string) {
            defer wg.Done()
            r := scraper.ScrapeMatchResults(userIDs, u)
            mu.Lock()
            results[n] = r
            mu.Unlock()
        }(name, url)
    }

    wg.Wait()
    return results
}

```

### Custom Image Themes

```go
// Dark theme
darkTheme := formatutils.ImageConfig{
    FontSize:    14.0,
    Margin:      28.0,
    ColWidth:    150.0,
    RowHeight:   28.0,
    BadgeWidth:  22.0,
    BadgeHeight: 22.0,
    R:           0.1,  // Dark background
    G:           0.1,
    B:           0.1,
}

// Compact theme
compactTheme := formatutils.ImageConfig{
    FontSize:    12.0,
    Margin:      20.0,
    ColWidth:    120.0,
    RowHeight:   24.0,
    BadgeWidth:  18.0,
    BadgeHeight: 18.0,
    R:           1.0,
    G:           1.0,
    B:           1.0,
}

```

## Performance Guidelines

### Memory Efficiency

- Reuse `Scraper` instance (shared HTTP client)
- Don't hold large result sets in memory
- Stream images to Discord/storage
- Close response bodies promptly

### Latency Optimization

- Implement badge caching (future)
- Use concurrent fetching (future)
- Set appropriate HTTP timeouts
- Minimize API requests

## Migration Guides

### From Legacy to Current API

**Old (Legacy):**

```go
results := scraper.ScrapeResults(userIDs, baseURL)  // Returns [][]string

```

**New (Current):**

```go
results := scraper.ScrapeMatchResults(userIDs, baseURL)  // Returns []parse.MatchResult

```

**Benefits:**

- Type safety
- Structured data
- Better error handling
- Easier testing

## Related Documentation

- [Architecture Overview](../ARCHITECTURE.md)
- [Development Guide](../DEVELOPMENT.md)
- [Main README](../../README.md)

## Support

For API questions or issues:

1. Check package-specific documentation
1. Review code examples in docs
1. Examine source code in `internal/`
1. Open GitHub issue if needed

______________________________________________________________________

**Documentation Version:** 2.0
**Last Updated:** 2025-11-05
