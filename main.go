package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/tripleawwy/onlineliga_discord_bot/internal/scraper"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	// Read Token from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	token := os.Getenv("DISCO_BOT_TOKEN")

	// Create a new Discord session using the provided bot token.
	discord, err := discordgo.New("Bot " + string(token))
	if err != nil {
		// Log error
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// Register messageCreate as a callback for the messageCreate events.
	discord.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		// Log error
		log.Fatalf("Error opening Discord session: %v", err)
	}

	// Wait here until interrupted.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Clean up resources.
	err = discord.Close()
	if err != nil {
		log.Fatalf("Error closing Discord session: %v", err)
	}

}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	olScraper := scraper.NewScraper()
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

		userId := args[1]
		// Convert string to int

		result, err := olScraper.ScrapeResult(userId)
		if err != nil {
			log.Fatalf("Error scraping result: %v", err)
		}
		// Send the result as code block to the channel
		_, _ = s.ChannelMessageSend(m.ChannelID, "```"+result+"```")

	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	if m.Content == "555" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Nase")
	}

	if m.Content == "Nase" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "555")
	}
}
