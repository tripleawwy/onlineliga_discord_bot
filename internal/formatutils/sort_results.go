package formatutils

import (
	"github.com/sirupsen/logrus"
	"github.com/tripleawwy/onlineliga_discord_bot/internal/parse"
	"sort"
	"strconv"
	"strings"
)

// SortResults sorts the results by league and position
func SortResults(results []parse.MatchResult, logger *logrus.Logger) []parse.MatchResult {
	sort.SliceStable(results, func(i, j int) bool {
		return compareMatchResults(results[i], results[j], logger)
	})
	return results
}

// compareMatchResults compares two MatchResult objects
func compareMatchResults(a, b parse.MatchResult, logger *logrus.Logger) bool {
	// Compare the league
	if a.LeagueLevel != b.LeagueLevel {
		return a.LeagueLevel < b.LeagueLevel
	}

	// Compare the position
	iPosition, err := convertPositionToInt(a.LeaguePosition)
	if err != nil {
		logger.WithError(err).Error("Failed to convert position to int")
		return false
	}
	jPosition, err := convertPositionToInt(b.LeaguePosition)
	if err != nil {
		logger.WithError(err).Error("Failed to convert position to int")
		return false
	}
	if iPosition != jPosition {
		return iPosition < jPosition
	}

	// Compare the home team
	if a.HomeTeam != b.HomeTeam {
		return a.HomeTeam < b.HomeTeam
	}

	// Compare the result
	if a.MatchResult != b.MatchResult {
		return a.MatchResult < b.MatchResult
	}

	// Compare the away team
	return a.AwayTeam < b.AwayTeam
}

// convertPositionToInt converts a league position string to an int
func convertPositionToInt(position string) (int, error) {
	// Remove the leading '#' character
	position = strings.TrimPrefix(position, "#")
	return strconv.Atoi(position)
}
