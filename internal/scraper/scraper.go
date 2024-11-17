package scraper

import (
	"github.com/sirupsen/logrus"
	"github.com/tripleawwy/onlineliga_discord_bot/internal/httpclient"
	"github.com/tripleawwy/onlineliga_discord_bot/internal/parse"
	"io"
	"net/http"
	"strconv"
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
			s.logger.WithError(scrapeErr).Errorf("Scraping results for user %s failed... continuing with next", userID)
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
	overviewURL := strings.Join([]string{baseURL, "/apiv1/team/overview?userId=", userID}, "")
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

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	// Convert UserID to int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return nil, err
	}
	result, parseErr := parse.Result(body, userIDInt)
	if parseErr != nil {
		return nil, parseErr
	}

	return result, err
}

// ScrapeMatchResult scrapes the match result from onlineliga, takes a match id as input and stores it in a MatchResult struct
func (s *Scraper) ScrapeMatchResult(userID string, baseURL string) (parse.MatchResult, error) {
	overviewURL := strings.Join([]string{baseURL, "/apiv1/team/overview?userId=", userID}, "")
	s.logger.WithField("userID", userID).Infof("URL is %s", overviewURL)
	resp, err := s.client.Get(overviewURL)
	if err != nil {
		return parse.MatchResult{}, err
	}
	defer func(Body io.ReadCloser) {
		ReadCloserError := Body.Close()
		if ReadCloserError != nil {
			s.logger.WithError(ReadCloserError).
				Warnf("Closing response body failed with error %s", ReadCloserError.Error())
		}
	}(resp.Body)

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return parse.MatchResult{}, readErr
	}

	// Convert UserID to int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return parse.MatchResult{}, err
	}

	// Get the actual match result
	result, parseErr := parse.ResultObject(body, userIDInt)
	if parseErr != nil {
		return parse.MatchResult{}, parseErr
	}

	return result, nil
}

// ScrapeMatchResults scrapes the match results from onlineliga and takes user ids as input
func (s *Scraper) ScrapeMatchResults(userIDs []string, baseURL string) []parse.MatchResult {
	var results []parse.MatchResult
	for _, userID := range userIDs {
		result, scrapeErr := s.ScrapeMatchResult(userID, baseURL)
		if scrapeErr != nil {
			s.logger.WithError(scrapeErr).Errorf("Scraping match results for user %s failed... continuing with next", userID)
			// Continue with the next user
			continue
		}
		s.logger.WithField("userID", userID).Infof("Result for user %s is %s", userID, result)
		results = append(results, result)
	}
	return results
}
