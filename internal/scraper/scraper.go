package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/tripleawwy/onlineliga_discord_bot/internal/httpclient"
	"github.com/tripleawwy/onlineliga_discord_bot/internal/parse"
	"io"
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
func (s *Scraper) ScrapeResults(userIDs []string) ([]string, error) {
	var results []string
	var err error
	for _, userID := range userIDs {
		result, scrapeErr := s.ScrapeResult(userID)
		if scrapeErr != nil {
			log.Printf("Scraping results for user %s failed with error %s", userID, scrapeErr.Error())
			err = scrapeErr
			break // stop scraping if there is an error for a user
		}
		log.Printf("Result for user %s is %s", userID, result)
		results = append(results, result)
	}
	return results, err
}

// ScrapeResult scrapes the results from onlineliga and takes a user id as input
func (s *Scraper) ScrapeResult(userID string) (string, error) {
	overviewURL := "https://www.onlineliga.de/team/overview?userId=" + userID
	log.Printf("URL is %s", overviewURL)
	resp, err := s.client.Get(overviewURL)
	if err != nil {
		return "Requesting results for user " + userID + " failed", err
	}
	defer func(Body io.ReadCloser) {
		ReadCloserError := Body.Close()
		if ReadCloserError != nil {
			log.Printf("Closing response body failed with error %s", ReadCloserError.Error())
		}
	}(resp.Body)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "Parsing results for user " + userID + " failed", err
	}
	//log.Print(doc.Text())
	result := parse.Result(doc)
	return result, err
}
