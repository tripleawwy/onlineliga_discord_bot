package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/tripleawwy/onlineliga_discord_bot/internal/formatutils"
	"github.com/tripleawwy/onlineliga_discord_bot/internal/scraper"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// Set up global logger
var logger = logrus.New()

func main() {
	logger.Level = logrus.InfoLevel
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	// Read Token from .env file
	err := godotenv.Load()
	if err != nil {
		logger.WithError(err).Fatal("Error loading .env file")
	}
	token := os.Getenv("DISCO_BOT_TOKEN")

	// Create a new Discord session using the provided bot token.
	discord, err := discordgo.New("Bot " + string(token))
	if err != nil {
		logger.WithError(err).Fatal("Error creating Discord session")
	}

	// Register messageCreate as a callback for the messageCreate events.
	discord.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		logger.WithError(err).Fatal("Error opening connection")
	}

	// Wait here until interrupted.
	logger.Infoln("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Clean up resources.
	err = discord.Close()
	if err != nil {
		logger.WithError(err).Fatal("Error closing connection")
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	olScraper := scraper.NewScraper(logger)
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	const prefix = "!results"

	// Split the message into a slice of words.
	// We expect the first word to be the command name.
	// The rest of the words are the arguments.
	args := strings.Split(m.Content, " ")

	if args[0] == prefix {
		userIDs := args[1:]
		results := olScraper.ScrapeResults(userIDs)
		resultsAsTable := formatutils.ResultsToTable(results)
		// Send the message to the channel in a code block with syntax highlighting.
		_, sendErr := s.ChannelMessageSend(m.ChannelID, "```ansi\n"+resultsAsTable+"\n```")

		if sendErr != nil {
			return
		}
	}
	return
}
