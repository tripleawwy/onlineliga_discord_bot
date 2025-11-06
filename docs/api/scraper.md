# Scraper Package API Documentation

**Package:** `github.com/tripleawwy/onlineliga_discord_bot/internal/scraper`

**Purpose:** HTTP data fetching layer for Onlineliga API integration

## Overview

The scraper package provides a clean abstraction for fetching match data from the Onlineliga gaming platform API.
It handles HTTP communication, error resilience, and multi-user coordination.

## Types

### Scraper

```go
type Scraper struct {
    client *http.Client
    logger *logrus.Logger
    url    string
}

```

Main struct that encapsulates HTTP client, logger, and API URL.

**Fields:**

- `client`: HTTP client for making requests (from `internal/httpclient`)
- `logger`: Structured logger for debugging and error reporting
- `url`: Currently unused, may be removed or used for base URL caching

## Constructors

### NewScraper

```go
func NewScraper(logger *logrus.Logger) Scraper

```

Creates a new Scraper instance with dependency injection.

**Parameters:**

- `logger`: Logrus logger instance for structured logging

**Returns:**

- `Scraper`: Initialized scraper with default HTTP client

**Example:**

```go
logger := logrus.New()
scraper := scraper.NewScraper(logger)

```

## Methods

### ScrapeResults (Legacy)

```go
func (s *Scraper) ScrapeResults(userIDs []string, baseURL string) [][]string

```

Fetches results for multiple users, returning legacy string slice format.

**Parameters:**

- `userIDs`: Slice of user ID strings
- `baseURL`: Onlineliga base URL (e.g., `https://www.onlineliga.de`)

**Returns:**

- `[][]string`: Two-dimensional array of result strings

**Behavior:**

- Calls `ScrapeResult()` for each user sequentially
- Logs errors but continues with remaining users
- Skips users that fail to fetch

**Example:**

```go
userIDs := []string{"12345", "67890"}
baseURL := "https://www.onlineliga.de"
results := scraper.ScrapeResults(userIDs, baseURL)

```

______________________________________________________________________

### ScrapeResult (Legacy)

```go
func (s *Scraper) ScrapeResult(userID string, baseURL string) ([]string, error)

```

Fetches result for a single user in legacy string format.

**Parameters:**

- `userID`: User ID as string
- `baseURL`: Onlineliga base URL

**Returns:**

- `[]string`: Result data as string slice
- `error`: Error if request or parsing fails

**HTTP Request:**

- Method: `GET`
- Endpoint: `{baseURL}/apiv1/team/overview?userId={userID}`
- Response: JSON payload

**Error Cases:**

- HTTP request failure
- Invalid user ID (not numeric)
- JSON parsing failure
- Network timeout

**Example:**

```go
result, err := scraper.ScrapeResult("12345", "https://www.onlineliga.de")
if err != nil {
    log.Printf("Failed: %v", err)
}

```

______________________________________________________________________

### ScrapeMatchResult

```go
func (s *Scraper) ScrapeMatchResult(userID string, baseURL string) (parse.MatchResult, error)

```

Fetches formatted match result for a single user.

**Parameters:**

- `userID`: User ID as string
- `baseURL`: Onlineliga base URL

**Returns:**

- `parse.MatchResult`: Structured match result data
- `error`: Error if request or parsing fails

**Process:**

1. Construct API URL
1. Make HTTP GET request
1. Read response body
1. Convert userID to integer
1. Parse JSON with `parse.ResultObject()`
1. Return structured MatchResult

**Logging:**

- Logs request URL with user ID field
- Logs cleanup errors for response body
- Errors logged by caller

**Example:**

```go
result, err := scraper.ScrapeMatchResult("12345", "https://www.onlineliga.de")
if err != nil {
    logger.WithError(err).Error("Scraping failed")
    return
}

fmt.Printf("League: %s, Position: %s\n",
    result.LeagueInfo.Name,
    result.LeaguePosition)

```

______________________________________________________________________

### ScrapeMatchResults

```go
func (s *Scraper) ScrapeMatchResults(userIDs []string, baseURL string) []parse.MatchResult

```

Fetches formatted match results for multiple users (current recommended method).

**Parameters:**

- `userIDs`: Slice of user ID strings
- `baseURL`: Onlineliga base URL

**Returns:**

- `[]parse.MatchResult`: Array of successfully fetched results

**Behavior:**

- Processes users **sequentially** (not concurrent)
- Continues on individual failures (graceful degradation)
- Logs errors but doesn't propagate them
- Returns partial results if some users fail

**Error Handling:**

```go
for _, userID := range userIDs {
    result, scrapeErr := s.ScrapeMatchResult(userID, baseURL)
    if scrapeErr != nil {
        // Log error and continue
        s.logger.WithError(scrapeErr).Errorf("...")
        continue
    }
    results = append(results, result)
}

```

**Example:**

```go
userIDs := []string{"12345", "67890", "11223"}
baseURL := "https://www.onlineliga.de"

results := scraper.ScrapeMatchResults(userIDs, baseURL)

for _, result := range results {
    fmt.Printf("%s: %s\n", result.HomeTeam, result.MatchResult)
}

```

## Usage Patterns

### Basic Usage

```go
package main

import (
    "github.com/sirupsen/logrus"
    "github.com/tripleawwy/onlineliga_discord_bot/internal/scraper"
)

func main() {
    logger := logrus.New()
    logger.Level = logrus.InfoLevel

    s := scraper.NewScraper(logger)

    userIDs := []string{"12345", "67890"}
    baseURL := "https://www.onlineliga.de"

    results := s.ScrapeMatchResults(userIDs, baseURL)

    // Process results...
}

```

### With Error Handling

```go
result, err := scraper.ScrapeMatchResult("12345", baseURL)
if err != nil {
    logger.WithFields(logrus.Fields{
        "userID": "12345",
        "error": err.Error(),
    }).Error("Failed to scrape result")
    return
}

// Use result...

```

### Multiple Regions

```go
type Region struct {
    Name string
    URL  string
}

regions := []Region{
    {"Germany", "https://www.onlineliga.de"},
    {"UK", "https://www.onlineleague.co.uk"},
    {"Austria", "https://www.onlineliga.at"},
    {"Switzerland", "https://www.onlineliga.ch"},
}

for _, region := range regions {
    results := scraper.ScrapeMatchResults(userIDs, region.URL)
    logger.Infof("Region %s: %d results", region.Name, len(results))
}

```

## Implementation Details

### Resource Management

Response bodies are properly cleaned up:

```go
defer func(Body io.ReadCloser) {
    ReadCloserError := Body.Close()
    if ReadCloserError != nil {
        s.logger.WithError(ReadCloserError).
            Warnf("Closing response body failed")
    }
}(resp.Body)

```

### User ID Conversion

```go
userIDInt, err := strconv.Atoi(userID)
if err != nil {
    return parse.MatchResult{}, err
}

```

User IDs must be numeric strings.

## Dependencies

**External:**

- `github.com/sirupsen/logrus`: Structured logging
- Standard library: `io`, `net/http`, `strconv`, `strings`

**Internal:**

- `internal/httpclient`: Provides `DefaultHTTPClient`
- `internal/parse`: Parses JSON responses

## Testing Recommendations

### Unit Tests

Mock HTTP responses:

```go
func TestScrapeMatchResult(t *testing.T) {
    // Create mock HTTP server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"matchData": {...}}`))
    }))
    defer server.Close()

    logger := logrus.New()
    s := scraper.NewScraper(logger)

    result, err := s.ScrapeMatchResult("12345", server.URL)
    assert.NoError(t, err)
    assert.NotNil(t, result)
}

```

### Integration Tests

Test against real API (staging):

```go
func TestScrapeMatchResult_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    logger := logrus.New()
    s := scraper.NewScraper(logger)

    result, err := s.ScrapeMatchResult("12345", "https://www.onlineliga.de")
    assert.NoError(t, err)
}

```

## Performance Considerations

### Current Limitations

- **Sequential fetching**: Users processed one at a time
- **No caching**: Every request hits the API
- **No rate limiting**: Could overwhelm API with many users

### Optimization Opportunities

**Concurrent Fetching:**

```go
func (s *Scraper) ScrapeMatchResultsConcurrent(userIDs []string, baseURL string) []parse.MatchResult {
    var wg sync.WaitGroup
    resultsChan := make(chan parse.MatchResult, len(userIDs))

    for _, userID := range userIDs {
        wg.Add(1)
        go func(id string) {
            defer wg.Done()
            result, err := s.ScrapeMatchResult(id, baseURL)
            if err == nil {
                resultsChan <- result
            }
        }(userID)
    }

    wg.Wait()
    close(resultsChan)

    results := []parse.MatchResult{}
    for result := range resultsChan {
        results = append(results, result)
    }
    return results
}

```

**Request Deduplication:**

```go
// If same userID appears multiple times, fetch once
uniqueUserIDs := deduplicate(userIDs)

```

## Error Scenarios

| Scenario        | Behavior       | Recovery             |
| --------------- | -------------- | -------------------- |
| Invalid user ID | Error returned | Caller handles       |
| Network timeout | Error returned | Retry logic needed   |
| API returns 404 | Parsing fails  | Graceful degradation |
| Malformed JSON  | Parse error    | Logged, continues    |
| API rate limit  | HTTP 429       | No retry, fails      |

## Related Documentation

- [Parse Package](parse.md) - JSON parsing
- [HTTPClient Package](httpclient.md) - HTTP client configuration
- [ARCHITECTURE.md](../ARCHITECTURE.md) - System architecture

______________________________________________________________________

**Package Version:** 2.0
**API Stability:** Stable
**Last Updated:** 2025-11-05
