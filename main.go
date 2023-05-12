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
		userIDs := args[1:]
		results, scrapeErr := olScraper.ScrapeResults(userIDs)
		if scrapeErr != nil {
			log.Printf("Error scraping results: %v", scrapeErr)
		}
		// Send a code block with the results separated by newlines.
		_, err := s.ChannelMessageSend(m.ChannelID, "```\n"+strings.Join(results, "\n")+"\n```")
		if err != nil {
			log.Fatalf("Error sending message: %v", err)
		}
	}
}
