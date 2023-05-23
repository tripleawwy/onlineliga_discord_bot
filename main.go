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
		// If there is no .env file, we assume that the token is set as an environment variable
		logger.WithError(err).Warnf("Error loading .env file")
	}
	token := os.Getenv("DISCO_BOT_TOKEN")

	// Create a new Discord session using the provided bot token.
	discord, err := discordgo.New("Bot " + string(token))
	if err != nil {
		logger.WithError(err).Fatal("Error creating Discord session")
	}

	// Register messageCreate as a callback for the messageCreate events.
	//discord.AddHandler(messageCreate)
	discord.AddHandler(onInteractionCreate)

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		logger.WithError(err).Fatal("Error opening connection")
	}

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "results",
			Description: "Get the results of a user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "location",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  ".de",
							Value: ".de",
						},
						{
							Name:  ".co.uk",
							Value: ".co.uk",
						},
						{
							Name:  ".at",
							Value: ".at",
						},
						{
							Name:  ".ch",
							Value: ".ch",
						},
					},
					Description: "The location of the results",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "users",
					Description: "The user ids to get the results for",
					Required:    true,
				},
			},
		},
	}

	// Bulk register the commands for all guilds
	for _, guild := range discord.State.Guilds {
		_, err = discord.ApplicationCommandBulkOverwrite(discord.State.User.ID, guild.ID, commands)
		if err != nil {
			logger.WithError(err).Errorf("Error registering commands for guild %s", guild.Name)
		}
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

	baseURL := getBaseURL(".de")
	if args[0] == prefix {
		userIDs := args[1:]
		results := olScraper.ScrapeResults(userIDs, baseURL)
		sortedResults := formatutils.SortResults(results, logger)
		resultsAsTable := formatutils.ResultsToTable(sortedResults)

		imageBuf, imageErr := formatutils.ResultsToImage(sortedResults)
		if imageErr != nil {
			logger.WithError(imageErr).Error("Error creating image")
		}
		// Read the image buffer into a reader
		resultAsImage := strings.NewReader(string(imageBuf.Bytes()))

		filePayload := discordgo.File{
			Name:   "SPOILER_results.png",
			Reader: resultAsImage,
		}

		// Send the message to the channel in a code block with syntax highlighting.
		_, sendErr := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "||```ansi\n" + resultsAsTable + "\n```||",
			Files: []*discordgo.File{
				&filePayload,
			},
		})
		if sendErr != nil {
			return
		}
	}
	return
}

func onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	olScraper := scraper.NewScraper(logger)

	if i.ApplicationCommandData().Name == "results" {
		location := i.ApplicationCommandData().Options[0].StringValue()
		logger.Infof("Location: %s", location)
		userIDs := strings.Split(i.ApplicationCommandData().Options[1].StringValue(), " ")

		var content string
		content = "```/results location: " + i.ApplicationCommandData().Options[0].StringValue() + " users: " + i.ApplicationCommandData().Options[1].StringValue() + "```"

		// Send initial response
		interactionErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "content",
			},
		})
		if interactionErr != nil {
			logger.WithError(interactionErr).Error("Error responding initially to interaction")
		}

		baseURL := getBaseURL(location)

		results := olScraper.ScrapeResults(userIDs, baseURL)
		sortedResults := formatutils.SortResults(results, logger)

		imageBuf, imageErr := formatutils.ResultsToImage(sortedResults)
		if imageErr != nil {
			logger.WithError(imageErr).Error("Error creating image")
		}
		// Read the image buffer into a reader
		resultAsImage := strings.NewReader(string(imageBuf.Bytes()))

		filePayload := discordgo.File{
			Name:   "SPOILER_results.png",
			Reader: resultAsImage,
		}

		_, interactionEditErr := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &content,
			Files: []*discordgo.File{
				&filePayload,
			},
		})

		if interactionEditErr != nil {
			logger.WithError(interactionEditErr).Error("Error responding to interaction")
		}
	}
	return
}

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
