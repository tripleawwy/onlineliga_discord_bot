package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/tripleawwy/onlineliga_discord_bot/internal/httpclient"
	"github.com/tripleawwy/onlineliga_discord_bot/internal/parse"
	"log"
	"net/http"
)

// Scraper is the interface for the scraper
type Scraper struct {
	client *http.Client
}

// NewScraper returns a new Scraper
func NewScraper() Scraper {
	return Scraper{
		client: httpclient.DefaultHTTPClient,
	}
}

// ScrapeResults scrapes the results from onlineliga and takes user ids as input
func (s *Scraper) ScrapeResults(userIDs []string) {
	for _, userID := range userIDs {
		s.ScrapeResult(userID)
	}
}

// ScrapeResult scrapes the results from onlineliga and takes a user id as input
func (s *Scraper) ScrapeResult(userID string) (string, error) {
	overviewURL := "https://www.onlineliga.de/team/overview?userId=" + userID
	log.Printf("URL is %s", overviewURL)
	resp, err := s.client.Get(overviewURL)
	if err != nil {
		return "Requesting results for user " + userID + " failed", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "Parsing results for user " + userID + " failed", err
	}
	//log.Print(doc.Text())
	result := parse.Result(doc)
	return result, err
}
