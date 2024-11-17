package parse_test

import (
	"github.com/tripleawwy/onlineliga_discord_bot/internal/parse"
	"os"
	"testing"
)

func TestResult(t *testing.T) {
	// open file which is located in files path
	path := "files/parse_test.json"
	jsonData, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("Failed to read file: %v", err)
	}

	userID := 8315
	result, parseErr := parse.ResultObject(jsonData, userID)
	if parseErr != nil {
		t.Errorf("Failed to parse result: %v", parseErr)
	}

	expectedResult := parse.MatchResult{
		LeagueLevel:    "OL1",
		BadgeURL:       "https://bla.xyz/image/1.png",
		LeaguePosition: "#11",
		HomeTeam:       "Test A",
		MatchResult:    "4 : 3",
		MatchState:     "WIN",
		AwayTeam:       "Test B",
		Points:         "3 pts",
	}

	if result != expectedResult {
		t.Errorf("Expected %v, got %v", expectedResult, result)
	}

}
