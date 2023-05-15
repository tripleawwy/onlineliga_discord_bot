package parse

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/tripleawwy/onlineliga_discord_bot/internal/formatutils"
	"log"
	"strconv"
	"strings"
)

func Result(doc *goquery.Document, userID string) []string {
	var result []string
	// Search for only the first element with a class of "team-overview-current-matchHistory" and return the text
	selection := doc.Find(".team-overview-matches").First()
	matchResult := parseResult(selection)
	homeTeam, awayTeam := parseClubNames(selection)
	league := parseLeague(doc)
	leaguePosition := parseLeaguePosition(doc)
	matchState := determineMatchState(doc, userID)
	matchResult = formatutils.ColourTextByState(matchResult, matchState)

	// add all variables to result in this order: league, leaguePosition, homeTeam, matchResult, awayTeam
	result = append(result, league)
	result = append(result, leaguePosition)
	result = append(result, homeTeam)
	result = append(result, matchResult)
	result = append(result, awayTeam)

	//result = "\u001b[0;32m" + league + " " + leaguePosition + " " + homeTeam + " " + result + " " + awayTeam + "\u001b[0m"
	//result = strings.Join([]string{league, leaguePosition, homeTeam, matchResult, awayTeam}, "\t")
	// Prepend ANSI escape code for green color and append ANSI escape code for reset color
	//result = strings.Join([]string{"\u001b[0;32m", result, "\u001b[0m"}, "")

	return result
}

// parseResult parses the result from this match:
func parseResult(selection *goquery.Selection) string {
	var result string
	// Search for only the first element with a class of "team-overview-current-match" and return the text
	selection = selection.Find(".team-overview-current-match").First()
	if selection.Length() > 0 {
		// Strip all children from the selection
		selection.Children().Remove()
		// Get the text from the selection and strip all whitespaces
		result = selection.Text()
		result = strings.TrimSpace(result)
		// Expected result is something like "0 : 1" but we want "0:1"
		result = strings.ReplaceAll(result, " ", "")
		// Prepend ANSI escape code for green color and append ANSI escape code for reset color
		//result = strings.Join([]string{"\u001b[0;32m", result, "\u001b[0m"}, "")
	}
	return result
}

// parseClubNames parses the club names from this match:
func parseClubNames(selection *goquery.Selection) (homeTeam string, awayTeam string) {
	// Get all elements with a class of "ol-team-name"
	selection = selection.Find(".ol-team-name")
	// This selection should contain two nodes, one for the home team and one for the away team
	if selection.Length() != 2 {
		log.Fatalf("Expected two teams, but got %d", selection.Length())
	} else {
		homeTeam = selection.First().Text()
		// Strip whitespaces from the beginning and end of the string
		homeTeam = strings.TrimSpace(homeTeam)

		awayTeam = selection.Last().Text()
		// Strip whitespaces from the beginning and end of the string
		awayTeam = strings.TrimSpace(awayTeam)

		log.Printf("Home team: %s, away team: %s", homeTeam, awayTeam)
	}
	return homeTeam, awayTeam
}

// parseLeague parses the league from this HTML:
func parseLeague(doc *goquery.Document) string {
	// Search for only the first element with a class of "ol-tf-league" and return the text
	selection := doc.Find(".ol-tf-league").First()
	// Get the text from the first child of the selection and strip all whitespaces
	league := selection.Children().First().Text()
	league = strings.TrimSpace(league)
	league = shortenLeagueName(league)

	return league
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
func parseLeaguePosition(doc *goquery.Document) string {
	// Search for all elements with a class of "ol-league-table" that have a child with a class of "ol-team-name ol-bold"
	// and return the text of the first child with a class of "ol-table-number"
	selection := doc.Find(".ol-league-table tr:has(.ol-team-name.ol-bold) .ol-table-number")
	leaguePosition := selection.Text()
	leaguePosition = strings.TrimSpace(leaguePosition)
	leaguePosition = formatLeaguePosition(leaguePosition)
	return leaguePosition
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

// parseMatchHistory parses a script tag that contains the last 10 matches of a team and
// a dictionary with the match data
func parseMatchHistory(doc *goquery.Document) []map[string]interface{} {
	// Search for all script tags #olTeamOverviewContent > script:last-of-type
	selection := doc.Find("#olTeamOverviewContent > script:last-of-type")
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
		log.Fatal(err)
	}
	return matches
}

// determineMatchState determines the state of the last match (WIN, DRAW, LOSS)
// based on the given user id
func determineMatchState(doc *goquery.Document, userID string) string {
	// Get the match history
	matchHistory := parseMatchHistory(doc)
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
		winner = "draw"
	}

	// Determine the state of the last match
	var state string
	if winner == userID {
		state = "WIN"
	} else if winner == "draw" {
		state = "DRAW"
	} else {
		state = "LOSS"
	}
	return state
}
