# Parse Package API Documentation

**Package:** `github.com/tripleawwy/onlineliga_discord_bot/internal/parse`

**Purpose:** JSON parsing and data transformation layer for Onlineliga API responses

## Overview

The parse package transforms Onlineliga API JSON responses into structured Go types and provides business logic for
match state determination and league filtering.

## Data Models

### Root

```go
type Root struct {
    MatchData    MatchData     `json:"matchData"`
    LeagueTables []LeagueTable `json:"leagueTable"`
    User         User          `json:"user"`
}

```

Top-level container for API response data.

**Fields:**

- `MatchData`: Match information container
- `LeagueTables`: Array of league standings
- `User`: User profile and league information

______________________________________________________________________

### User

```go
type User struct {
    Username string `json:"username"`
    UID      int    `json:"uid"`
    LeagueID int    `json:"leagueId"`
    League   League `json:"league"`
    Badge    Badge  `json:"badgeData"`
}

```

User profile information from the API.

**Fields:**

- `Username`: Display name
- `UID`: Unique user identifier
- `LeagueID`: Current league identifier
- `League`: Nested league details
- `Badge`: User badge data

______________________________________________________________________

### League

```go
type League struct {
    LeagueID    int    `json:"leagueID"`
    Level       int    `json:"level"`
    Name        string `json:"name"`
    Teams       int    `json:"teams"`
    R           int    `json:"r"`  // Color: Red
    G           int    `json:"g"`  // Color: Green
    B           int    `json:"b"`  // Color: Blue
    ID          int    `json:"ID"`
    MatchDay    int    `json:"matchDay"`
    StateId     int    `json:"stateId"`
    CommunityId int    `json:"communityId"`
}

```

League metadata and configuration.

**Fields:**

- `LeagueID`: Unique league identifier
- `Level`: League tier/division
- `Name`: League display name
- `Teams`: Number of teams in league
- `R, G, B`: League color representation (RGB)
- `MatchDay`: Current match day number
- `StateId`: League state identifier
- `CommunityId`: Community/region identifier

______________________________________________________________________

### Badge

```go
type Badge struct {
    URL     string `json:"url"`
    BadgeID int    `json:"badgeId"`
    UserID  int    `json:"userId"`
}

```

User badge/avatar information.

**Fields:**

- `URL`: Badge image URL
- `BadgeID`: Badge identifier
- `UserID`: Owner user ID

______________________________________________________________________

### MatchData

```go
type MatchData struct {
    LastMatch Match `json:"lastMatch"`
}

```

Container for match information.

**Fields:**

- `LastMatch`: Most recent match details

______________________________________________________________________

### Match

```go
type Match struct {
    MatchID             int       `json:"matchId"`
    LeagueID            int       `json:"leagueId"`
    Matchday            int       `json:"matchday"`
    Player1             int       `json:"player1"`
    Player2             int       `json:"player2"`
    GoalsPlayer1        int       `json:"goals_player1"`
    GoalsPlayer2        int       `json:"goals_player2"`
    OlTimestamp         int       `json:"olTimestamp"`
    Timestamp           string    `json:"timestamp"`
    MatchTypeID         int       `json:"matchTypeId"`
    GoalsFirstHalfUser1 int       `json:"goals_first_half_user1"`
    GoalsFirstHalfUser2 int       `json:"goals_first_half_user2"`
    BallPossession1     int       `json:"ballPossession1"`
    BallPossession2     int       `json:"ballPossession2"`
    UserA               MatchUser `json:"userA"`
    UserB               MatchUser `json:"userB"`
}

```

Detailed match information.

**Fields:**

- `MatchID`: Unique match identifier
- `LeagueID`: League where match occurred
- `Matchday`: Match day number
- `Player1`, `Player2`: Player identifiers
- `GoalsPlayer1`, `GoalsPlayer2`: Final score
- `GoalsFirstHalfUser1`, `GoalsFirstHalfUser2`: Half-time score
- `BallPossession1`, `BallPossession2`: Possession percentage
- `Timestamp`: Match completion time (ISO 8601)
- `UserA`, `UserB`: Detailed player information

______________________________________________________________________

### MatchUser

```go
type MatchUser struct {
    UID         int    `json:"uid"`
    Username    string `json:"username"`
    TeamName    string `json:"teamName"`
    CommunityID int    `json:"communityId"`
    LeagueID    int    `json:"leagueId"`
    Badge       Badge  `json:"badge"`
}

```

Player information within a match.

**Fields:**

- `UID`: User identifier
- `Username`: User display name
- `TeamName`: Team/club name
- `CommunityID`: Community identifier
- `LeagueID`: League identifier
- `Badge`: User badge data

______________________________________________________________________

### LeagueTable

```go
type LeagueTable struct {
    UserID         int    `json:"userId"`
    LeagueID       int    `json:"leagueId"`
    Matchday       int    `json:"matchday"`
    Type           string `json:"type"`
    Rank           int    `json:"rank"`
    Points         int    `json:"points"`
    ConcedingGoals int    `json:"concedingGoals"`
    ScoredGoals    int    `json:"scoredGoals"`
    Win            int    `json:"win"`
    Lost           int    `json:"lost"`
    Draw           int    `json:"draw"`
    MatchCount     int    `json:"matchCount"`
    TeamName       string `json:"teamName"`
    BadgeID        int    `json:"badgeId"`
    UID            int    `json:"uid"`
}

```

League standings entry.

**Fields:**

- `Rank`: Position in league table (1 = first place)
- `Points`: Total points earned
- `Win`, `Draw`, `Lost`: Match outcomes
- `ScoredGoals`, `ConcedingGoals`: Goal statistics
- `MatchCount`: Matches played
- `TeamName`: Team display name

______________________________________________________________________

### MatchResult (Output Model)

```go
type MatchResult struct {
    LeagueInfo     League `json:"leagueInfo"`
    LeagueLevel    string `json:"leagueLevel"`
    BadgeURL       string `json:"badgeURL"`
    LeaguePosition string `json:"leaguePosition"`
    HomeTeam       string `json:"homeTeam"`
    MatchResult    string `json:"matchResult"`
    MatchState     string `json:"matchState"`
    AwayTeam       string `json:"awayTeam"`
    Points         string `json:"points"`
}

```

Formatted match result for display/rendering.

**Fields:**

- `LeagueInfo`: Complete league metadata
- `LeagueLevel`: Formatted level string (e.g., "Liga 1")
- `BadgeURL`: User badge image URL
- `LeaguePosition`: Formatted position (e.g., "1st")
- `HomeTeam`: Home team name
- `MatchResult`: Score (e.g., "3:1")
- `MatchState`: "WIN", "LOSS", or "DRAW"
- `AwayTeam`: Away team name
- `Points`: Points string (e.g., "+3", "0")

## Public Functions

### Result (Legacy)

```go
func Result(body []byte, userID int) ([]string, error)

```

Parses API response and returns legacy string slice format.

**Parameters:**

- `body`: Raw JSON response body
- `userID`: User ID to extract data for

**Returns:**

- `[]string`: Result data as string array
- `error`: Parse or validation error

**Deprecated:** Use `ResultObject()` instead for structured data.

______________________________________________________________________

### ResultObject

```go
func ResultObject(body []byte, userID int) (MatchResult, error)

```

Parses API response and returns structured MatchResult (recommended).

**Parameters:**

- `body`: Raw JSON response from Onlineliga API
- `userID`: User ID to extract match data for

**Returns:**

- `MatchResult`: Formatted result structure
- `error`: Parse, validation, or logic error

**Process:**

1. Unmarshal JSON to `Root` struct
1. Filter league table for user's team
1. Determine match state (Win/Loss/Draw)
1. Format all fields for display
1. Return structured `MatchResult`

**Error Cases:**

- Invalid JSON format
- Missing required fields
- User not found in league table
- Match data not available

**Example:**

```go
apiResponse := []byte(`{
    "matchData": {...},
    "leagueTable": [...],
    "user": {...}
}`)

result, err := parse.ResultObject(apiResponse, 12345)
if err != nil {
    log.Fatalf("Parse failed: %v", err)
}

fmt.Printf("League: %s, Position: %s, Result: %s (%s)\n",
    result.LeagueInfo.Name,
    result.LeaguePosition,
    result.MatchResult,
    result.MatchState)

```

## Internal Functions

### filterLeagueTableForTeam

```go
func filterLeagueTableForTeam(leagueTables []LeagueTable, teamID int) LeagueTable

```

Extracts user's entry from league table array.

**Parameters:**

- `leagueTables`: Array of all league standings
- `teamID`: User ID to find

**Returns:**

- `LeagueTable`: User's league table entry

**Logic:**

```go
for _, table := range leagueTables {
    if table.UserID == teamID {
        return table
    }
}
return LeagueTable{} // Empty if not found

```

______________________________________________________________________

### getMatchState

```go
func getMatchState(match Match, userID int) string

```

Determines match outcome from user's perspective.

**Parameters:**

- `match`: Match details
- `userID`: User to determine result for

**Returns:**

- `"WIN"`: User scored more goals
- `"LOSS"`: User scored fewer goals
- `"DRAW"`: Equal goals scored

**Logic:**

```go
if match.UserA.UID == userID {
    if match.GoalsPlayer1 > match.GoalsPlayer2 {
        return "WIN"
    } else if match.GoalsPlayer1 < match.GoalsPlayer2 {
        return "LOSS"
    }
}
// Similar for UserB
return "DRAW"

```

## Usage Examples

### Basic Parsing

```go
package main

import (
    "fmt"
    "io"
    "net/http"
    "github.com/tripleawwy/onlineliga_discord_bot/internal/parse"
)

func main() {
    apiURL := "https://www.onlineliga.de/apiv1/team/overview?userId=12345"
    resp, err := http.Get(apiURL)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

    result, err := parse.ResultObject(body, 12345)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Result: %s %s %s\n",
        result.HomeTeam,
        result.MatchResult,
        result.AwayTeam)
}

```

### Error Handling

```go
result, err := parse.ResultObject(body, userID)
if err != nil {
    switch {
    case errors.Is(err, json.SyntaxError{}):
        log.Error("Invalid JSON from API")
    default:
        log.Errorf("Parse failed: %v", err)
    }
    return
}

```

### Batch Processing

```go
type Result struct {
    UserID int
    Data   parse.MatchResult
    Error  error
}

func parseMultipleUsers(responses map[int][]byte) []Result {
    results := make([]Result, 0, len(responses))

    for userID, body := range responses {
        data, err := parse.ResultObject(body, userID)
        results = append(results, Result{
            UserID: userID,
            Data:   data,
            Error:  err,
        })
    }

    return results
}

```

## Field Formatting

### League Level

```go
LeagueLevel = fmt.Sprintf("Liga %d", league.Level)
// Output: "Liga 1", "Liga 2", etc.

```

### League Position

```go
LeaguePosition = fmt.Sprintf("%d.", leagueTable.Rank)
// Output: "1.", "2.", "10.", etc.

```

### Match Result

```go
MatchResult = fmt.Sprintf("%d:%d",
    match.GoalsPlayer1,
    match.GoalsPlayer2)
// Output: "3:1", "0:0", "2:4", etc.

```

### Points

```go
if matchState == "WIN" {
    Points = "+3"
} else if matchState == "DRAW" {
    Points = "+1"
} else {
    Points = "0"
}

```

## Testing

### Unit Test Example

```go
func TestResultObject(t *testing.T) {
    testData, err := os.ReadFile("files/parse_test.json")
    if err != nil {
        t.Fatal(err)
    }

    result, err := ResultObject(testData, 12345)
    assert.NoError(t, err)
    assert.Equal(t, "WIN", result.MatchState)
    assert.Equal(t, "+3", result.Points)
}

```

### Test Data Structure

Store test JSON in `internal/parse/files/parse_test.json`:

```json
{
    "matchData": {
        "lastMatch": {
            "matchId": 123,
            "goals_player1": 3,
            "goals_player2": 1,
            "userA": {
                "uid": 12345,
                "username": "TestUser"
            },
            "userB": {
                "uid": 67890,
                "username": "Opponent"
            }
        }
    },
    "leagueTable": [
        {
            "userId": 12345,
            "rank": 1,
            "points": 45
        }
    ],
    "user": {
        "uid": 12345,
        "username": "TestUser"
    }
}

```

## Performance Considerations

- **JSON Unmarshaling:** ~1-5ms for typical API response
- **Memory Allocation:** ~10-50KB per parsed result
- **No caching:** Each call parses from scratch

### Optimization Opportunities

**Struct Tags:**

```go
// Use omitempty for optional fields
type Match struct {
    MatchID int `json:"matchId,omitempty"`
}

```

**Streaming Parser:**

```go
// For large responses, use json.Decoder
decoder := json.NewDecoder(resp.Body)
var root Root
err := decoder.Decode(&root)

```

## Error Handling Patterns

| Error Type                | Cause                    | Recovery                 |
| ------------------------- | ------------------------ | ------------------------ |
| `json.SyntaxError`        | Invalid JSON             | Log and skip             |
| `json.UnmarshalTypeError` | Type mismatch            | Validate API schema      |
| Empty result              | User not in league table | Return empty MatchResult |
| Missing match data        | No recent match          | Handle in caller         |

## Related Documentation

- [Scraper Package](scraper.md) - HTTP fetching
- [Formatutils Package](formatutils.md) - Result rendering
- [ARCHITECTURE.md](../ARCHITECTURE.md) - Data flow

______________________________________________________________________

**Package Version:** 2.0
**API Stability:** Stable
**Last Updated:** 2025-11-05
