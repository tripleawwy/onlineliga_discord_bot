package parse_test

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/tripleawwy/onlineliga_discord_bot/internal/parse"
	"strings"
	"testing"
)

func TestResult(t *testing.T) {
	html := `
		<div class="team-overview-matches">
			<div class="team-overview-current-match-result" onclick="olAnchorNavigation.load('/match', { season : 28, matchId : 8203 });">
				<div class="team-overview-current-match">0 : 1
					<div class="mobile-matchdaytable-halftime-result">( 0 : 0 )</div>
				</div>
			</div>
		</div>
	`
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	expected := "0 : 1"
	result := parse.Result(doc)
	if result != expected {
		t.Errorf("Unexpected result: expected %q, but got %q", expected, result)
	}
}
