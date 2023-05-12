package parse

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
)

func Result(doc *goquery.Document) string {
	var result string
	// Search for only the first element with a class of "team-overview-current-match" and return the text
	selection := doc.Find(".team-overview-matches").First()
	result = parseResult(selection)
	homeTeam, awayTeam := parseClubNames(selection)
	league := parseLeague(doc)
	leaguePosition := parseLeaguePosition(doc)
	result = league + " " + leaguePosition + " " + homeTeam + " " + result + " " + awayTeam
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
	// The unshortened league name is something like "2. ONLINELIGA Nord 1"
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
