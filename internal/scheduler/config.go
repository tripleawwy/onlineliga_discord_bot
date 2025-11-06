// Package scheduler provides cron-based task scheduling for automated Discord bot actions.
package scheduler

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the entire schedules configuration
type Config struct {
	Schedules []Schedule `yaml:"schedules"`
}

// Schedule represents a single scheduled task
type Schedule struct {
	Name      string  `yaml:"name"`
	Enabled   bool    `yaml:"enabled"`
	Cron      string  `yaml:"cron"`
	ChannelID string  `yaml:"channel_id"`
	Command   Command `yaml:"command"`
}

// Command represents the command to execute
type Command struct {
	Location string `yaml:"location"`
	Users    []int  `yaml:"users"`
}

// LoadConfig reads and parses the schedules.yaml file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if len(c.Schedules) == 0 {
		return fmt.Errorf("no schedules defined")
	}

	for i, schedule := range c.Schedules {
		if err := schedule.Validate(); err != nil {
			return fmt.Errorf("schedule %d (%s): %w", i, schedule.Name, err)
		}
	}

	return nil
}

// Validate checks if a single schedule is valid
func (s *Schedule) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("name is required")
	}

	if s.Cron == "" {
		return fmt.Errorf("cron expression is required")
	}

	if s.ChannelID == "" {
		return fmt.Errorf("channel_id is required")
	}

	if s.Command.Location == "" {
		return fmt.Errorf("command.location is required")
	}

	if len(s.Command.Users) == 0 {
		return fmt.Errorf("command.users cannot be empty")
	}

	// Validate location format
	validLocations := map[string]bool{
		".de":    true,
		".co.uk": true,
		".at":    true,
		".ch":    true,
	}

	if !validLocations[s.Command.Location] {
		return fmt.Errorf("invalid location %s (must be .de, .co.uk, .at, or .ch)", s.Command.Location)
	}

	return nil
}
