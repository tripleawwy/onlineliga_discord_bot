package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/tripleawwy/onlineliga_discord_bot/internal/httpclient"
	"github.com/tripleawwy/onlineliga_discord_bot/internal/parse"
	"io"
	"net/http"
)

// Scraper is the interface for the scraper
type Scraper struct {
	client *http.Client
	logger *logrus.Logger
}

// NewScraper returns a new Scraper
func NewScraper(logger *logrus.Logger) Scraper {
	return Scraper{
		client: httpclient.DefaultHTTPClient,
		logger: logger,
	}
}

// ScrapeResults scrapes the results from onlineliga and takes user ids as input
func (s *Scraper) ScrapeResults(userIDs []string) ([]string, error) {
	var results []string
	for _, userID := range userIDs {
		result, scrapeErr := s.ScrapeResult(userID)
		if scrapeErr != nil {
			s.logger.WithError(scrapeErr).Errorf("Scraping results for user %s failed", userID)

			// Stop scraping if there is an error for a user
			break
		}
		s.logger.WithField("userID", userID).Infof("Result for user %s is %s", userID, result)
		results = append(results, result)
	}
	return results
}

// ScrapeResult scrapes the results from onlineliga and takes a user id as input
func (s *Scraper) ScrapeResult(userID string) (string, error) {
	overviewURL := "https://www.onlineliga.de/team/overview?userId=" + userID
	s.logger.WithField("userID", userID).Debugf("URL is %s", overviewURL)
	resp, err := s.client.Get(overviewURL)
	if err != nil {
		return "Requesting results for user " + userID + " failed", err
	}
	defer func(Body io.ReadCloser) {
		ReadCloserError := Body.Close()
		if ReadCloserError != nil {
			s.logger.WithError(ReadCloserError).
				Warnf("Closing response body failed with error %s", ReadCloserError.Error())
		}
	}(resp.Body)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "Parsing results for user " + userID + " failed", err
	}
	s.logger.WithField("userID", userID).Debugf("Parsed document for user %s is %s", userID, doc.Text())
	result := parse.Result(doc)
	return result, err
}
