package scheduler

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/tripleawwy/onlineliga_discord_bot/internal/formatutils"
	"github.com/tripleawwy/onlineliga_discord_bot/internal/scraper"
)

// Scheduler manages scheduled tasks
type Scheduler struct {
	cron    *cron.Cron
	config  *Config
	session *discordgo.Session
	scraper *scraper.Scraper
	logger  *logrus.Logger
}

// New creates a new Scheduler instance
func New(session *discordgo.Session, scr *scraper.Scraper, logger *logrus.Logger) *Scheduler {
	return &Scheduler{
		cron:    cron.New(cron.WithSeconds()),
		session: session,
		scraper: scr,
		logger:  logger,
	}
}

// LoadAndStart loads the configuration and starts all enabled schedules
func (s *Scheduler) LoadAndStart(configPath string) error {
	config, err := LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	s.config = config

	// Register all enabled schedules
	for _, schedule := range config.Schedules {
		if schedule.Enabled {
			if err := s.registerSchedule(schedule); err != nil {
				s.logger.WithError(err).Errorf("Failed to register schedule: %s", schedule.Name)
				continue
			}
			s.logger.Infof("Registered schedule: %s (cron: %s)", schedule.Name, schedule.Cron)
		} else {
			s.logger.Infof("Skipped disabled schedule: %s", schedule.Name)
		}
	}

	// Start the cron scheduler
	s.cron.Start()
	s.logger.Infof("Scheduler started with %d active schedules", len(s.cron.Entries()))

	return nil
}

// Stop stops the scheduler and waits for running jobs to complete
func (s *Scheduler) Stop() {
	if s.cron != nil {
		ctx := s.cron.Stop()
		<-ctx.Done()
		s.logger.Info("Scheduler stopped")
	}
}

// registerSchedule registers a single schedule with the cron scheduler
func (s *Scheduler) registerSchedule(schedule Schedule) error {
	_, err := s.cron.AddFunc(schedule.Cron, func() {
		s.executeSchedule(schedule)
	})

	return err
}

// executeSchedule executes a scheduled task
func (s *Scheduler) executeSchedule(schedule Schedule) {
	s.logger.Infof("Executing schedule: %s", schedule.Name)

	// Convert user IDs to strings
	userIDs := make([]string, len(schedule.Command.Users))
	for i, userID := range schedule.Command.Users {
		userIDs[i] = strconv.Itoa(userID)
	}

	// Get base URL from location
	baseURL := getBaseURL(schedule.Command.Location)

	// Scrape match results
	results := s.scraper.ScrapeMatchResults(userIDs, baseURL)
	if len(results) == 0 {
		s.logger.Warnf("No results found for schedule: %s", schedule.Name)
		s.sendErrorMessage(schedule.ChannelID, "No results found")
		return
	}

	// Sort and generate image
	results = formatutils.SortResults(results, s.logger)
	imageBuf, err := formatutils.MatchResultsToImage(results, formatutils.ImageConfig{
		FontSize:    14.0,
		Margin:      28.0,
		ColWidth:    150.0,
		RowHeight:   28.0,
		BadgeWidth:  22.0,
		BadgeHeight: 22.0,
		R:           1.0,
		G:           1.0,
		B:           1.0,
	})

	if err != nil {
		s.logger.WithError(err).Errorf("Failed to generate image for schedule: %s", schedule.Name)
		s.sendErrorMessage(schedule.ChannelID, "Failed to generate results image")
		return
	}

	// Send message to Discord
	if err := s.sendResults(schedule, &imageBuf); err != nil {
		s.logger.WithError(err).Errorf("Failed to send results for schedule: %s", schedule.Name)
	} else {
		s.logger.Infof("Successfully executed schedule: %s", schedule.Name)
	}
}

// sendResults sends the results to Discord
func (s *Scheduler) sendResults(schedule Schedule, imageBuf *bytes.Buffer) error {
	// Create command summary
	userIDsStr := make([]string, len(schedule.Command.Users))
	for i, userID := range schedule.Command.Users {
		userIDsStr[i] = strconv.Itoa(userID)
	}

	content := fmt.Sprintf("**%s**\n```/results location: %s users: %s```",
		schedule.Name,
		schedule.Command.Location,
		strings.Join(userIDsStr, " "))

	// Create image reader
	resultAsImage := bytes.NewReader(imageBuf.Bytes())

	// Send message with image attachment
	_, err := s.session.ChannelMessageSendComplex(schedule.ChannelID, &discordgo.MessageSend{
		Content: content,
		Files: []*discordgo.File{
			{
				Name:   "SPOILER_results.png",
				Reader: resultAsImage,
			},
		},
	})

	return err
}

// sendErrorMessage sends an error message to Discord
func (s *Scheduler) sendErrorMessage(channelID, errorMsg string) {
	_, err := s.session.ChannelMessageSend(channelID, fmt.Sprintf("⚠️ Error: %s", errorMsg))
	if err != nil {
		s.logger.WithError(err).Error("Failed to send error message")
	}
}

// getBaseURL converts location to base URL
func getBaseURL(location string) string {
	switch location {
	case ".de":
		return "https://www.onlineliga.de"
	case ".co.uk":
		return "https://www.onlineleague.co.uk"
	case ".at":
		return "https://www.onlineliga.at"
	case ".ch":
		return "https://www.onlineliga.ch"
	default:
		return "https://www.onlineliga.de"
	}
}
