# Formatutils Package API Documentation

**Package:** `github.com/tripleawwy/onlineliga_discord_bot/internal/formatutils`

**Purpose:** Image generation and result formatting layer

## Overview

The formatutils package generates visual representations of match results using 2D graphics rendering. It handles image
composition, text rendering, badge downloading, and PNG encoding.

## Types

### ImageConfig

```go
type ImageConfig struct {
    FontSize    float64
    Margin      float64
    ColWidth    float64
    RowHeight   float64
    BadgeWidth  float64
    BadgeHeight float64
    R           float64  // Background color: Red (0.0-1.0)
    G           float64  // Background color: Green (0.0-1.0)
    B           float64  // Background color: Blue (0.0-1.0)
}

```

Configuration for image generation parameters.

**Fields:**

- `FontSize`: Text size in points (recommended: 12-16)
- `Margin`: Border spacing in pixels (recommended: 20-40)
- `ColWidth`: Column width in pixels (recommended: 120-180)
- `RowHeight`: Row height in pixels (recommended: 24-32)
- `BadgeWidth`, `BadgeHeight`: Badge dimensions in pixels (recommended: 20-30)
- `R`, `G`, `B`: Background color as RGB (0.0 = black, 1.0 = white)

**Default Configuration (from main.go):**

```go
config := ImageConfig{
    FontSize:    14.0,
    Margin:      28.0,
    ColWidth:    150.0,
    RowHeight:   28.0,
    BadgeWidth:  22.0,
    BadgeHeight: 22.0,
    R:           1.0,  // White background
    G:           1.0,
    B:           1.0,
}

```

## Public Functions

### MatchResultsToImage

```go
func MatchResultsToImage(results []parse.MatchResult, config ImageConfig) (*bytes.Buffer, error)

```

Generates PNG image from match results (main entry point).

**Parameters:**

- `results`: Array of match results to render
- `config`: Image generation configuration

**Returns:**

- `*bytes.Buffer`: PNG image data buffer
- `error`: Generation or encoding error

**Process:**

1. Calculate image dimensions based on result count
1. Create graphics context with background
1. Load font from `FONT_PATH` environment variable
1. Render each result row (text + badge)
1. Encode to PNG format
1. Return buffer

**Error Cases:**

- Empty results array
- Invalid font path
- Badge download failure
- Image encoding error

**Example:**

```go
results := []parse.MatchResult{...}

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

imageBuf, err := formatutils.MatchResultsToImage(results, config)
if err != nil {
    log.Fatalf("Image generation failed: %v", err)
}

// Send to Discord or save to file
os.WriteFile("results.png", imageBuf.Bytes(), 0644)

```

______________________________________________________________________

### SortResults

```go
func SortResults(results []parse.MatchResult, logger *logrus.Logger) []parse.MatchResult

```

Sorts match results for consistent display ordering.

**Parameters:**

- `results`: Unsorted match results
- `logger`: Logger for debugging (currently unused)

**Returns:**

- `[]parse.MatchResult`: Sorted results

**Sort Criteria:**

1. Primary: League level (ascending: Liga 1, Liga 2, ...)
1. Secondary: League position (ascending: 1st, 2nd, ...)

**Example:**

```go
results := scraper.ScrapeMatchResults(userIDs, baseURL)
sortedResults := formatutils.SortResults(results, logger)

// sortedResults now ordered by league level and position

```

**Implementation:**

```go
sort.Slice(results, func(i, j int) bool {
    // Compare league levels first
    levelI, _ := strconv.Atoi(results[i].LeagueLevel[5:])
    levelJ, _ := strconv.Atoi(results[j].LeagueLevel[5:])

    if levelI != levelJ {
        return levelI < levelJ
    }

    // Then compare positions
    posI := extractPosition(results[i].LeaguePosition)
    posJ := extractPosition(results[j].LeaguePosition)

    return posI < posJ
})

```

## Internal Functions

### createNewContext

```go
func createNewContext(width, height float64, config ImageConfig) *gg.Context

```

Initializes 2D graphics context with background.

**Parameters:**

- `width`: Image width in pixels
- `height`: Image height in pixels
- `config`: Image configuration for background color

**Returns:**

- `*gg.Context`: Initialized graphics context

**Process:**

1. Create context with dimensions
1. Set background color from config (R, G, B)
1. Fill entire canvas
1. Reset to default color (black)

______________________________________________________________________

### loadFontFace

```go
func loadFontFace(dc *gg.Context, fontPath string, fontSize float64) error

```

Loads TrueType font into graphics context.

**Parameters:**

- `dc`: Graphics context
- `fontPath`: Absolute path to .ttf font file
- `fontSize`: Font size in points

**Returns:**

- `error`: Font loading error

**Font Source:**

- Read from `FONT_PATH` environment variable
- Default in Docker: Cascadia Code Bold
- Must be valid TrueType font (.ttf)

**Error Cases:**

- File not found
- Invalid font format
- Permission denied

______________________________________________________________________

### calculateImageSize

```go
func calculateImageSize(results []parse.MatchResult, config ImageConfig) (float64, float64)

```

Computes image dimensions based on content.

**Parameters:**

- `results`: Match results to render
- `config`: Configuration with margins and row height

**Returns:**

- `width`: Image width in pixels
- `height`: Image height in pixels

**Formula:**

```go
numRows := len(results)
width = calculateImageWidth(results, config)
height = config.Margin*2 + float64(numRows)*config.RowHeight

```

______________________________________________________________________

### calculateImageWidth

```go
func calculateImageWidth(results []parse.MatchResult, config ImageConfig) float64

```

Calculates required image width for results.

**Parameters:**

- `results`: Match results
- `config`: Configuration with column widths

**Returns:**

- `float64`: Image width in pixels

**Logic:**

```go
// Fixed columns: Level, Badge, Position, Home, Score, State, Away, Points
numColumns := 8
return config.Margin*2 + float64(numColumns)*config.ColWidth

```

______________________________________________________________________

### calculateColumnWidth

```go
func calculateColumnWidth(text string, dc *gg.Context) float64

```

Measures text width for dynamic column sizing.

**Parameters:**

- `text`: Text to measure
- `dc`: Graphics context with loaded font

**Returns:**

- `float64`: Text width in pixels

**Usage:**

```go
textWidth := calculateColumnWidth("Liga 1", dc)

```

______________________________________________________________________

### calculateTextWidth

```go
func calculateTextWidth(text string, dc *gg.Context) float64

```

Alias for `calculateColumnWidth` (measures text width).

______________________________________________________________________

### writeMatchResult

```go
func writeMatchResult(dc *gg.Context, result parse.MatchResult, y float64, config ImageConfig) error

```

Renders single match result row.

**Parameters:**

- `dc`: Graphics context
- `result`: Match result data
- `y`: Vertical position (top of row)
- `config`: Rendering configuration

**Returns:**

- `error`: Rendering error (badge download, drawing)

**Columns Rendered:**

1. League Level (e.g., "Liga 1")
1. Badge (downloaded image)
1. League Position (e.g., "1.")
1. Home Team
1. Match Score (e.g., "3:1")
1. Match State (WIN/LOSS/DRAW)
1. Away Team
1. Points (+3, +1, 0)

**Text Positioning:**

```go
x := config.Margin
dc.DrawString(result.LeagueLevel, x, y+config.RowHeight/2)
x += config.ColWidth
// ... continue for each column

```

______________________________________________________________________

### drawBadge

```go
func drawBadge(dc *gg.Context, badgeURL string, x, y float64, config ImageConfig) error

```

Downloads and renders user badge at position.

**Parameters:**

- `dc`: Graphics context
- `badgeURL`: Badge image URL
- `x`, `y`: Top-left position
- `config`: Badge dimensions

**Returns:**

- `error`: Download or render error

**Process:**

1. Check if URL is valid image URL
1. Download image via `downloadImage()`
1. Resize to badge dimensions
1. Composite onto canvas at position

**Supported Formats:**

- PNG
- JPEG
- GIF

______________________________________________________________________

### downloadImage

```go
func downloadImage(url string) (image.Image, error)

```

Fetches image from URL.

**Parameters:**

- `url`: Image URL

**Returns:**

- `image.Image`: Downloaded image
- `error`: Download or decode error

**HTTP Client:**

- Uses standard `http.DefaultClient`
- No timeout configured (uses default)
- Follows redirects

**Error Cases:**

- Network failure
- HTTP error (404, 500, etc.)
- Unsupported image format
- Malformed image data

______________________________________________________________________

### resizeImage

```go
func resizeImage(img image.Image, width, height float64) image.Image

```

Resizes image to target dimensions.

**Parameters:**

- `img`: Source image
- `width`, `height`: Target dimensions

**Returns:**

- `image.Image`: Resized image

**Algorithm:**

- Uses `github.com/golang/freetype` or `golang.org/x/image`
- Maintains aspect ratio? (Implementation detail)
- Quality: Depends on resizing method

______________________________________________________________________

### isImageURL

```go
func isImageURL(url string) bool

```

Validates if URL points to image resource.

**Parameters:**

- `url`: URL to check

**Returns:**

- `bool`: True if valid image URL

**Validation:**

```go
// Common image extensions
extensions := []string{".png", ".jpg", ".jpeg", ".gif"}
for _, ext := range extensions {
    if strings.HasSuffix(strings.ToLower(url), ext) {
        return true
    }
}
return false

```

______________________________________________________________________

### setDrawingColor

```go
func setDrawingColor(dc *gg.Context, r, g, b float64)

```

Sets current drawing color.

**Parameters:**

- `dc`: Graphics context
- `r`, `g`, `b`: RGB values (0.0-1.0)

**Usage:**

```go
setDrawingColor(dc, 1.0, 0.0, 0.0) // Red
setDrawingColor(dc, 0.0, 0.0, 0.0) // Black

```

______________________________________________________________________

### encodeToPNG

```go
func encodeToPNG(dc *gg.Context) (*bytes.Buffer, error)

```

Encodes graphics context to PNG buffer.

**Parameters:**

- `dc`: Graphics context with rendered content

**Returns:**

- `*bytes.Buffer`: PNG image data
- `error`: Encoding error

**Format:**

- PNG with default compression
- No transparency unless explicitly set
- RGB color mode

## Usage Examples

### Basic Image Generation

```go
package main

import (
    "os"
    "github.com/tripleawwy/onlineliga_discord_bot/internal/formatutils"
    "github.com/tripleawwy/onlineliga_discord_bot/internal/parse"
)

func main() {
    results := []parse.MatchResult{
        {
            LeagueLevel:    "Liga 1",
            LeaguePosition: "1.",
            BadgeURL:       "https://example.com/badge.png",
            HomeTeam:       "FC Example",
            MatchResult:    "3:1",
            MatchState:     "WIN",
            AwayTeam:       "Test FC",
            Points:         "+3",
        },
    }

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

    imgBuf, err := formatutils.MatchResultsToImage(results, config)
    if err != nil {
        panic(err)
    }

    os.WriteFile("output.png", imgBuf.Bytes(), 0644)
}

```

### Custom Styling

```go
// Dark mode configuration
darkConfig := ImageConfig{
    FontSize:    16.0,
    Margin:      32.0,
    ColWidth:    160.0,
    RowHeight:   32.0,
    BadgeWidth:  24.0,
    BadgeHeight: 24.0,
    R:           0.1,  // Dark background
    G:           0.1,
    B:           0.1,
}

// Compact configuration
compactConfig := ImageConfig{
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

### Error Handling

```go
results := scraper.ScrapeMatchResults(userIDs, baseURL)
if len(results) == 0 {
    log.Error("No results to render")
    return
}

sortedResults := formatutils.SortResults(results, logger)

imageBuf, err := formatutils.MatchResultsToImage(sortedResults, config)
if err != nil {
    log.WithError(err).Error("Image generation failed")
    // Fallback: send text response instead
    return
}

// Success: use imageBuf

```

## Dependencies

**External:**

- `github.com/fogleman/gg`: 2D graphics rendering
- `golang.org/x/image`: Image manipulation
- `github.com/golang/freetype`: Font rendering

**Standard Library:**

- `bytes`: Buffer management
- `image`: Image types and operations
- `image/png`, `image/jpeg`: Image encoding/decoding
- `net/http`: Badge downloading
- `os`: Environment variables (FONT_PATH)

## Performance Characteristics

### Generation Time

| Operation           | Time           | Notes                   |
| ------------------- | -------------- | ----------------------- |
| Image creation      | ~10-30ms       | Context setup           |
| Font loading        | ~5-15ms        | Cached after first load |
| Text rendering      | ~1-5ms per row | Depends on text length  |
| Badge download      | ~50-200ms each | Network dependent       |
| PNG encoding        | ~10-50ms       | Depends on size         |
| **Total (3 users)** | ~200-500ms     | With badges             |

### Memory Usage

| Component          | Memory         | Notes                 |
| ------------------ | -------------- | --------------------- |
| Graphics context   | ~1-5MB         | Depends on dimensions |
| Font data          | ~500KB-1MB     | Loaded once           |
| Downloaded badges  | ~50-200KB each | Temporary             |
| PNG buffer         | ~100KB-1MB     | Final output          |
| **Peak (3 users)** | ~3-10MB        | Typical               |

## Optimization Opportunities

### Badge Caching

```go
var badgeCache = make(map[string]image.Image)
var cacheMutex sync.RWMutex

func getCachedBadge(url string) (image.Image, error) {
    cacheMutex.RLock()
    if cached, ok := badgeCache[url]; ok {
        cacheMutex.RUnlock()
        return cached, nil
    }
    cacheMutex.RUnlock()

    img, err := downloadImage(url)
    if err != nil {
        return nil, err
    }

    cacheMutex.Lock()
    badgeCache[url] = img
    cacheMutex.Unlock()

    return img, nil
}

```

### Parallel Badge Download

```go
type badgeResult struct {
    index int
    img   image.Image
    err   error
}

func downloadBadgesParallel(urls []string) []image.Image {
    results := make(chan badgeResult, len(urls))

    for i, url := range urls {
        go func(idx int, u string) {
            img, err := downloadImage(u)
            results <- badgeResult{idx, img, err}
        }(i, url)
    }

    badges := make([]image.Image, len(urls))
    for range urls {
        r := <-results
        if r.err == nil {
            badges[r.index] = r.img
        }
    }

    return badges
}

```

## Error Scenarios

| Scenario             | Impact            | Recovery                |
| -------------------- | ----------------- | ----------------------- |
| FONT_PATH not set    | Generation fails  | Set env var, restart    |
| Invalid font file    | Generation fails  | Verify font path        |
| Badge download fails | Row missing badge | Continue without badge  |
| Network timeout      | Slow generation   | Set HTTP timeout        |
| Invalid badge URL    | Skip badge        | Log warning, continue   |
| PNG encoding fails   | No output         | Check memory/disk space |

## Related Documentation

- [Parse Package](parse.md) - Data structures
- [Scraper Package](scraper.md) - Data fetching
- [ARCHITECTURE.md](../ARCHITECTURE.md) - Image generation pipeline

______________________________________________________________________

**Package Version:** 2.0
**API Stability:** Stable
**Last Updated:** 2025-11-05
