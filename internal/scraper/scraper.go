package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/tripleawwy/onlineliga_discord_bot/internal/httpclient"
	"github.com/tripleawwy/onlineliga_discord_bot/internal/parse"
	"io"
	"net/http"
	"strings"
)

// Scraper is the interface for the scraper
type Scraper struct {
	client *http.Client
	logger *logrus.Logger
	url    string
}

// NewScraper returns a new Scraper
func NewScraper(logger *logrus.Logger) Scraper {
	return Scraper{
		client: httpclient.DefaultHTTPClient,
		logger: logger,
	}
}

// ScrapeResults scrapes the results from onlineliga and takes user ids as input
func (s *Scraper) ScrapeResults(userIDs []string, baseURL string) [][]string {
	var results [][]string
	for _, userID := range userIDs {
		result, scrapeErr := s.ScrapeResult(userID, baseURL)
		if scrapeErr != nil {
			s.logger.WithError(scrapeErr).Errorf("Scraping results for user %s failed", userID)
			// Continue with the next user
			continue
		}
		s.logger.WithField("userID", userID).Infof("Result for user %s is %s", userID, result)
		results = append(results, result)
	}
	return results
}

// ScrapeResult scrapes the results from onlineliga and takes a user id as input
func (s *Scraper) ScrapeResult(userID string, baseURL string) ([]string, error) {
	overviewURL := strings.Join([]string{baseURL, "/team/overview?userId=", userID}, "")
	s.logger.WithField("userID", userID).Infof("URL is %s", overviewURL)
	resp, err := s.client.Get(overviewURL)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		ReadCloserError := Body.Close()
		if ReadCloserError != nil {
			s.logger.WithError(ReadCloserError).
				Warnf("Closing response body failed with error %s", ReadCloserError.Error())
		}
	}(resp.Body)

	// Parse the document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	s.logger.WithField("userID", userID).Debugf("Parsed document for user %s is %s", userID, doc.Text())
	if err != nil {
		return nil, err
	}

	// Parse the result
	result, parseErr := parse.Result(doc, userID)
	if parseErr != nil {
		return nil, parseErr
	}

	return result, err
}
