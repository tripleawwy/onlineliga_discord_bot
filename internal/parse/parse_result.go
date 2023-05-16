package parse

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/tripleawwy/onlineliga_discord_bot/internal/formatutils"
	"strconv"
	"strings"
)

// Result parses the overview page of a user and returns the result of the last match
func Result(doc *goquery.Document, userID string) ([]string, error) {
	var result []string
	// Search for only the first element with a class of "team-overview-current-matchHistory" and return the text
	selection := doc.Find(".team-overview-matches").First()
	if selection.Length() == 0 {
		return nil, &ResultError{Msg: "No matches found"}
	}

	matchResult, resultErr := parseResult(selection)
	if resultErr != nil {
		return nil, resultErr
	}
	homeTeam, awayTeam, clubNameErr := parseClubNames(selection)
	if clubNameErr != nil {
		return nil, clubNameErr
	}
	league, leagueErr := parseLeague(doc)
	if leagueErr != nil {
		return nil, leagueErr
	}
	leaguePosition, leaguePosErr := parseLeaguePosition(doc)
	if leaguePosErr != nil {
		return nil, leaguePosErr
	}
	matchState, matchStateErr := parseMatchState(doc, userID)
	if matchStateErr != nil {
		return nil, matchStateErr
	}

	// Colour the match result by match state
	matchResult = formatutils.ColourTextByState(matchResult, matchState)

	// Add all variables to result in this order: league, leaguePosition, homeTeam, matchResult, awayTeam
	result = append(result, league)
	result = append(result, leaguePosition)
	result = append(result, homeTeam)
	result = append(result, matchResult)
	result = append(result, awayTeam)

	return result, nil
}

// parseResult parses the result from this match:
func parseResult(selection *goquery.Selection) (string, error) {
	var result string
	// Search for only the first element with a class of "team-overview-current-match" and return the text
	selection = selection.Find(".team-overview-current-match").First()
	if selection.Length() > 0 {
		// Strip all children from the selection
		selection.Children().Remove()
		// Get the text from the selection and strip all whitespaces
		result = selection.Text()
		result = strings.TrimSpace(result)
		result = strings.ReplaceAll(result, " ", "")
	} else {
		// Create an error message and return it
		err := &ResultError{Msg: "Error: No result found"}
		return "", err
	}
	return result, nil
}

// parseClubNames parses the club names from this match:
func parseClubNames(selection *goquery.Selection) (homeTeam string, awayTeam string, err error) {
	// Get all elements with a class of "ol-team-name"
	selection = selection.Find(".ol-team-name")
	// This selection should contain two nodes, one for the home team and one for the away team
	if selection.Length() != 2 {
		// Create an error message and return it
		err = &ResultError{Msg: "Error: Expected two club names, got " + strconv.Itoa(selection.Length()) + " instead"}
		return "", "", err
	} else {
		homeTeam = selection.First().Text()
		homeTeam = strings.TrimSpace(homeTeam)

		awayTeam = selection.Last().Text()
		awayTeam = strings.TrimSpace(awayTeam)
	}
	return homeTeam, awayTeam, nil
}

// parseLeague parses the league from this HTML:
func parseLeague(doc *goquery.Document) (string, error) {
	// Search for only the first element with a class of "ol-tf-league" and return the text
	selection := doc.Find(".ol-tf-league").First()
	if selection.Length() == 0 {
		err := &ResultError{Msg: "Error: No league found"}
		return "", err
	}
	// Get the text from the first child of the selection and strip all whitespaces
	league := selection.Children().First().Text()
	league = strings.TrimSpace(league)
	league = shortenLeagueName(league)

	return league, nil
}

// shortenLeagueName returns the short name of the league
func shortenLeagueName(league string) string {
	// The not shortened league name is something like "2. ONLINELIGA Nord 1"
	// We want to return the league level only "OL2"
	// Split the string by whitespaces
	words := strings.Split(league, " ")
	// The league level is the first word
	leagueLevel := words[0]
	// Remove the dot from the league level
	leagueLevel = strings.ReplaceAll(leagueLevel, ".", "")
	// Prepend "OL" to the league level
	leagueLevel = "OL" + leagueLevel
	return leagueLevel
}

// parseLeaguePosition parses the league position from this HTML:
func parseLeaguePosition(doc *goquery.Document) (string, error) {
	// Search for all elements with a class of "ol-league-table" that have a child with a class of "ol-team-name ol-bold"
	// and return the text of the first child with a class of "ol-table-number"
	selection := doc.Find(".ol-league-table tr:has(.ol-team-name.ol-bold) .ol-table-number")
	if selection.Length() == 0 {
		err := &ResultError{Msg: "Error: No league position found"}
		return "", err
	}
	leaguePosition := selection.Text()
	leaguePosition = strings.TrimSpace(leaguePosition)
	leaguePosition = formatLeaguePosition(leaguePosition)
	return leaguePosition, nil
}

// formatLeaguePosition formats the league position
func formatLeaguePosition(leaguePosition string) string {
	// The league position is something like "1." or "2."
	// We want to return the league position like "P1" or "P2"
	// Remove the dot from the league position
	leaguePosition = strings.ReplaceAll(leaguePosition, ".", "")
	// Prepend "P" to the league position
	leaguePosition = "P" + leaguePosition
	return leaguePosition
}

// parseMatchState determines the state of the last match (WIN, DRAW, LOSS)
// based on the given user id
func parseMatchState(doc *goquery.Document, userID string) (string, error) {
	// Get the match history
	matchHistory, matchHistoryErr := parseMatchHistory(doc)
	if matchHistoryErr != nil {
		return "", matchHistoryErr
	}
	// Get the last match
	match := matchHistory[len(matchHistory)-1]
	// The winner is the player with the most goals
	// Get the goals of the player
	goalsPlayer1 := match["goals_player1"].(float64)
	goalsPlayer2 := match["goals_player2"].(float64)
	// Determine the winner
	var winner string
	if goalsPlayer1 > goalsPlayer2 {
		winnerId := match["player1"].(float64)
		// Convert the winner id to a string
		winner = strconv.FormatFloat(winnerId, 'f', 0, 64)

	} else if goalsPlayer1 < goalsPlayer2 {
		winnerId := match["player2"].(float64)
		// Convert the winner id to a string
		winner = strconv.FormatFloat(winnerId, 'f', 0, 64)
	} else {
		winner = "no winner"
	}

	return determineState(winner, userID), nil
}

// parseMatchHistory parses a script tag that contains the last 10 matches of a team and
// a dictionary with the match data
func parseMatchHistory(doc *goquery.Document) ([]map[string]interface{}, error) {
	// Search for all script tags #olTeamOverviewContent > script:last-of-type
	selection := doc.Find("#olTeamOverviewContent > script:last-of-type")
	if selection.Length() == 0 {
		err := &ResultError{Msg: "Error: No script tag found. Thus no match history found"}
		return nil, err
	}
	// Get the text from the selection
	script := selection.Text()
	//$(document).ready(function()
	//{
	//	drawMatchHistory("matchHistory1", 71606, '[{"matchId":44457,"leagueId":125,"matchday":13,"player1":71606,"player2":10188,"goals_player1":3,"goals_player2":0,"goals_first_half_user1":1,"goals_first_half_user2":0,"state":0,"season":28},{"matchId":44465,"leagueId":125,"matchday":14,"player1":21222,"player2":71606,"goals_player1":0,"goals_player2":3,"goals_first_half_user1":0,"goals_first_half_user2":3,"state":0,"season":28},{"matchId":44473,"leagueId":125,"matchday":15,"player1":71606,"player2":42965,"goals_player1":3,"goals_player2":1,"goals_first_half_user1":2,"goals_first_half_user2":1,"state":0,"season":28},{"matchId":44481,"leagueId":125,"matchday":16,"player1":13160,"player2":71606,"goals_player1":1,"goals_player2":2,"goals_first_half_user1":0,"goals_first_half_user2":0,"state":0,"season":28},{"matchId":44491,"leagueId":125,"matchday":17,"player1":23757,"player2":71606,"goals_player1":0,"goals_player2":2,"goals_first_half_user1":0,"goals_first_half_user2":1,"state":0,"season":28},{"matchId":44501,"leagueId":125,"matchday":18,"player1":44236,"player2":71606,"goals_player1":1,"goals_player2":3,"goals_first_half_user1":0,"goals_first_half_user2":2,"state":0,"season":28},{"matchId":44511,"leagueId":125,"matchday":19,"player1":71606,"player2":15116,"goals_player1":1,"goals_player2":1,"goals_first_half_user1":1,"goals_first_half_user2":1,"state":0,"season":28},{"matchId":44521,"leagueId":125,"matchday":20,"player1":61247,"player2":71606,"goals_player1":0,"goals_player2":5,"goals_first_half_user1":0,"goals_first_half_user2":2,"state":0,"season":28},{"matchId":44531,"leagueId":125,"matchday":21,"player1":71606,"player2":25826,"goals_player1":5,"goals_player2":0,"goals_first_half_user1":1,"goals_first_half_user2":0,"state":0,"season":28},{"matchId":44541,"leagueId":125,"matchday":22,"player1":28153,"player2":71606,"goals_player1":0,"goals_player2":1,"goals_first_half_user1":0,"goals_first_half_user2":0,"state":0,"season":28}]', 15, 0.2);
	//	olGUI.setToggleButtonActiveToActive($('#toggleButtonInfo'));
	//});
	// Extract the JSON string from the script
	// The JSON string is between the first occurrence of "[" and the last occurrence of "]"
	start := strings.Index(script, "[")
	end := strings.LastIndex(script, "]")
	script = script[start : end+1]

	// Get Last Match of the Match History
	// The last match is the last element of the JSON array
	// Convert the JSON string to a dictionary
	var matches []map[string]interface{}
	err := json.Unmarshal([]byte(script), &matches)
	if err != nil {
		return nil, err
	}
	return matches, nil
}

// determineState determines the state of the last match (WIN, DRAW, LOSS)
func determineState(winner string, userID string) string {
	// Determine the state of the last match
	var state string
	if winner == userID {
		state = "WIN"
	} else if winner == "no winner" {
		state = "DRAW"
	} else {
		state = "LOSS"
	}
	return state
}
