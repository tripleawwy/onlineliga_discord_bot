package parse

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
)

// Result gets the overview page of a user and returns the result of the last match as a JSON string
func Result(responseBody []byte, userID int) ([]string, error) {
	var rootObject Root
	marshalErr := json.Unmarshal(responseBody, &rootObject)
	if marshalErr != nil {
		return nil, marshalErr
	}

	homeTeam := rootObject.MatchData.LastMatch.UserA.TeamName
	awayTeam := rootObject.MatchData.LastMatch.UserB.TeamName
	leagueLevel := fmt.Sprintf("OL%d", rootObject.User.League.Level)
	leagueTable, filterErr := filterLeagueTableForTeam(rootObject.LeagueTables, userID)
	if filterErr != nil {
		return nil, filterErr
	}
	points := fmt.Sprintf("%d pts", leagueTable.Points)
	leaguePosition := fmt.Sprintf("#%d", leagueTable.Rank)
	matchResult := fmt.Sprintf("%d : %d", rootObject.MatchData.LastMatch.GoalsPlayer1, rootObject.MatchData.LastMatch.GoalsPlayer2)
	matchState := getMatchState(&rootObject.MatchData.LastMatch, userID)

	// Colour the match result by match state
	//matchResult = formatutils.ColourTextByState(matchResult, matchState)
	logrus.Debugf("Match state: %s", matchState)

	badgeURL := rootObject.User.Badge.URL

	// Add all variables to result in this order: league, leaguePosition, homeTeam, matchResult, awayTeam
	result := []string{leagueLevel, leaguePosition, homeTeam, matchResult, awayTeam, points, badgeURL}

	return result, nil
}

// ResultObject gets the overview page of a user and returns the result of the last match as a MatchResult struct
func ResultObject(responseBody []byte, userID int) (MatchResult, error) {
	var rootObject Root
	marshalErr := json.Unmarshal(responseBody, &rootObject)
	if marshalErr != nil {
		return MatchResult{}, marshalErr
	}

	// The league tables contain a league table for each team in the league
	// And we need to filter the league table for the given user id
	leagueTable, filterErr := filterLeagueTableForTeam(rootObject.LeagueTables, userID)
	if filterErr != nil {
		return MatchResult{}, filterErr
	}

	result := MatchResult{
		LeagueInfo:     rootObject.User.League,
		LeagueLevel:    fmt.Sprintf("OL%d", rootObject.User.League.Level),
		BadgeURL:       rootObject.User.Badge.URL,
		LeaguePosition: fmt.Sprintf("#%d", leagueTable.Rank),
		HomeTeam:       rootObject.MatchData.LastMatch.UserA.TeamName,
		MatchResult:    fmt.Sprintf("%d : %d", rootObject.MatchData.LastMatch.GoalsPlayer1, rootObject.MatchData.LastMatch.GoalsPlayer2),
		MatchState:     getMatchState(&rootObject.MatchData.LastMatch, userID),
		AwayTeam:       rootObject.MatchData.LastMatch.UserB.TeamName,
		Points:         fmt.Sprintf("%d pts", leagueTable.Points),
	}

	return result, nil
}

// filterLeagueTableForTeam returns the league table for the given user id
func filterLeagueTableForTeam(leagueTable []LeagueTable, userID int) (LeagueTable, error) {
	for _, team := range leagueTable {
		if team.UserID == userID {
			return team, nil
		}
	}
	return LeagueTable{}, &ResultError{Msg: "Error: No league table found for user " + strconv.Itoa(userID)}
}

// getMatchState determines the state of the last match (WIN, DRAW, LOSS)
func getMatchState(lastMatch *Match, userID int) string {
	// Determine the winnerID based on the goals
	var winnerID int
	switch {
	case lastMatch.GoalsPlayer1 > lastMatch.GoalsPlayer2:
		winnerID = lastMatch.UserA.UID
	case lastMatch.GoalsPlayer1 < lastMatch.GoalsPlayer2:
		winnerID = lastMatch.UserB.UID
	default:
		winnerID = 0
	}

	// Determine the state of the last match
	switch {
	case winnerID == userID:
		return "WIN"
	case winnerID == 0:
		return "DRAW"
	default:
		return "LOSS"
	}
}
