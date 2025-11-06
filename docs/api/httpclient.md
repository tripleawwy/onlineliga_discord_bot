# HTTPClient Package API Documentation

**Package:** `github.com/tripleawwy/onlineliga_discord_bot/internal/httpclient`

**Purpose:** HTTP client configuration and management

## Overview

The httpclient package provides a pre-configured HTTP client with cookie jar support for consistent
API communication across the application.

## Public Variables

### DefaultHTTPClient

```go
var DefaultHTTPClient *http.Client

```

Pre-configured HTTP client with cookie jar support.

**Configuration:**

- Cookie jar enabled for session management
- Default timeout (inherited from `http.DefaultClient`)
- Connection pooling enabled
- Follows redirects (default behavior)

**Usage:**

```go
import "github.com/tripleawwy/onlineliga_discord_bot/internal/httpclient"

resp, err := httpclient.DefaultHTTPClient.Get("https://www.onlineliga.de/api")

```

## Initialization

The client is initialized in `init()` function:

```go
func init() {
    jar, err := cookiejar.New(nil)
    if err != nil {
        log.Fatal(err)
    }

    DefaultHTTPClient = &http.Client{
        Jar: jar,
    }
}

```

**Features:**

- Automatic cookie storage and management
- Session persistence across requests
- Thread-safe cookie jar

## Cookie Management

### Cookie Jar

The `cookiejar.Jar` implementation stores cookies across requests:

**Automatic Behaviors:**

- Saves cookies from `Set-Cookie` headers
- Sends appropriate cookies with subsequent requests
- Handles cookie expiration
- Domain and path matching

**Example:**

```go
// First request sets cookies
resp1, _ := httpclient.DefaultHTTPClient.Get("https://example.com/login")

// Second request automatically includes cookies
resp2, _ := httpclient.DefaultHTTPClient.Get("https://example.com/api")

```

## Usage Patterns

### Basic HTTP GET

```go
package main

import (
    "io"
    "log"
    "github.com/tripleawwy/onlineliga_discord_bot/internal/httpclient"
)

func main() {
    resp, err := httpclient.DefaultHTTPClient.Get("https://api.example.com/data")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    // Process response...
}

```

### HTTP POST

```go
import (
    "bytes"
    "net/http"
)

func postData(url string, data []byte) error {
    req, err := http.NewRequest("POST", url, bytes.NewReader(data))
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/json")

    resp, err := httpclient.DefaultHTTPClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    return nil
}

```

### With Context (Timeout)

```go
import (
    "context"
    "net/http"
    "time"
)

func getWithTimeout(url string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return err
    }

    resp, err := httpclient.DefaultHTTPClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    return nil
}

```

### Custom Headers

```go
func getWithHeaders(url string) error {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return err
    }

    req.Header.Set("User-Agent", "OnlineligaBot/2.0")
    req.Header.Set("Accept", "application/json")

    resp, err := httpclient.DefaultHTTPClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    return nil
}

```

## Configuration Customization

### Custom Client with Timeout

```go
import (
    "net/http"
    "net/http/cookiejar"
    "time"
)

func createCustomClient() *http.Client {
    jar, _ := cookiejar.New(nil)

    return &http.Client{
        Jar:     jar,
        Timeout: 30 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns:        100,
            MaxIdleConnsPerHost: 10,
            IdleConnTimeout:     90 * time.Second,
        },
    }
}

```

### With Retry Logic

```go
import (
    "net/http"
    "time"
)

type RetryTransport struct {
    Transport http.RoundTripper
    MaxRetries int
}

func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    var resp *http.Response
    var err error

    for i := 0; i < t.MaxRetries; i++ {
        resp, err = t.Transport.RoundTrip(req)
        if err == nil && resp.StatusCode < 500 {
            return resp, nil
        }

        if i < t.MaxRetries-1 {
            time.Sleep(time.Second * time.Duration(i+1))
        }
    }

    return resp, err
}

func createRetryClient() *http.Client {
    jar, _ := cookiejar.New(nil)

    return &http.Client{
        Jar: jar,
        Transport: &RetryTransport{
            Transport:  http.DefaultTransport,
            MaxRetries: 3,
        },
    }
}

```

## Integration with Scraper

The scraper package uses `DefaultHTTPClient`:

```go
// internal/scraper/scraper.go
type Scraper struct {
    client *http.Client  // Set to httpclient.DefaultHTTPClient
    logger *logrus.Logger
}

func NewScraper(logger *logrus.Logger) Scraper {
    return Scraper{
        client: httpclient.DefaultHTTPClient,
        logger: logger,
    }
}

```

**Benefits:**

- Shared cookie jar across all scraper instances
- Consistent HTTP configuration
- Connection pooling and reuse
- Session management

## Performance Characteristics

### Connection Pooling

Default `http.Client` maintains connection pool:

**Settings (Default):**

- `MaxIdleConns`: 100 (total idle connections)
- `MaxIdleConnsPerHost`: 2 (idle per host)
- `IdleConnTimeout`: 90 seconds

**Benefits:**

- Reuses TCP connections
- Reduces handshake overhead
- Improves throughput

### Memory Usage

**Per Client:**

- Cookie jar: ~1-10KB (depends on cookies)
- Connection pool: ~10-100KB
- Total: ~10-110KB

**Multiple Requests:**

- Connections reused (no multiplication)
- Cookies accumulate (minimal overhead)

## Error Handling

### Common Errors

**Network Errors:**

```go
resp, err := httpclient.DefaultHTTPClient.Get(url)
if err != nil {
    if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
        log.Error("Request timed out")
    } else {
        log.Errorf("Network error: %v", err)
    }
    return err
}

```

**HTTP Status Errors:**

```go
if resp.StatusCode >= 400 {
    log.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
    return fmt.Errorf("HTTP error: %d", resp.StatusCode)
}

```

**Cookie Errors:**

```go
// Cookie jar errors are rare (only on initialization)
jar, err := cookiejar.New(nil)
if err != nil {
    log.Fatal("Failed to create cookie jar")
}

```

## Testing

### Mock HTTP Client

```go
import (
    "net/http"
    "net/http/httptest"
)

func TestHTTPRequest(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status": "ok"}`))
    }))
    defer server.Close()

    resp, err := httpclient.DefaultHTTPClient.Get(server.URL)
    if err != nil {
        t.Fatal(err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected 200, got %d", resp.StatusCode)
    }
}

```

### Cookie Testing

```go
func TestCookieJar(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        http.SetCookie(w, &http.Cookie{
            Name:  "session",
            Value: "abc123",
        })
        w.WriteHeader(http.StatusOK)
    }))
    defer server.Close()

    // First request sets cookie
    resp1, _ := httpclient.DefaultHTTPClient.Get(server.URL)
    resp1.Body.Close()

    // Check if cookie was stored
    url, _ := url.Parse(server.URL)
    cookies := httpclient.DefaultHTTPClient.Jar.Cookies(url)

    if len(cookies) == 0 {
        t.Error("Cookie not stored")
    }
}

```

## Best Practices

### Resource Cleanup

Always close response bodies:

```go
resp, err := httpclient.DefaultHTTPClient.Get(url)
if err != nil {
    return err
}
defer resp.Body.Close()

// Read body...

```

### Timeout Configuration

Set timeouts for production:

```go
// Per-request timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
resp, err := httpclient.DefaultHTTPClient.Do(req)

```

### Error Checking

Check both error and status code:

```go
resp, err := httpclient.DefaultHTTPClient.Get(url)
if err != nil {
    return fmt.Errorf("request failed: %w", err)
}
defer resp.Body.Close()

if resp.StatusCode >= 400 {
    return fmt.Errorf("HTTP %d", resp.StatusCode)
}

```

## Security Considerations

### HTTPS Enforcement

Always use HTTPS URLs:

```go
// Good
url := "https://www.onlineliga.de/api"

// Bad (insecure)
url := "http://www.onlineliga.de/api"

```

### Certificate Validation

Default client validates certificates:

```go
// Certificate validation enabled by default
// Do NOT disable in production

```

### Cookie Security

Cookies handled securely:

- Respects `Secure` flag
- Honors `HttpOnly` attribute
- Domain and path restrictions enforced

## Troubleshooting

### Connection Refused

```text
Error: connection refused

```

**Causes:**

- Server is down
- Firewall blocking
- Incorrect URL

**Solution:**

- Verify server is running
- Check network connectivity
- Validate URL

### Timeout Errors

```text
Error: context deadline exceeded

```

**Causes:**

- Server slow to respond
- Network latency
- No timeout set

**Solution:**

- Increase timeout duration
- Check server performance
- Verify network quality

### Cookie Issues

```text
Cookies not being sent

```

**Causes:**

- Domain mismatch
- Secure flag on HTTP
- Path restrictions

**Solution:**

- Verify cookie domain
- Use HTTPS if Secure flag set
- Check cookie path

## Related Documentation

- [Scraper Package](scraper.md) - Main HTTP client consumer
- [ARCHITECTURE.md](../ARCHITECTURE.md) - HTTP layer design
- [Go net/http Documentation](https://pkg.go.dev/net/http)
- [Go http.Client Documentation](https://pkg.go.dev/net/http#Client)

______________________________________________________________________

**Package Version:** 1.0
**API Stability:** Stable
**Last Updated:** 2025-11-05
