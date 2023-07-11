package main

import (
	"context"
	"flag"
	"log"
	"os"
	"regexp"

	"golang.org/x/exp/slices"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
	"github.com/joho/godotenv"
)

// embedRegex is a regular expression that matches embed URLs.
// https://mathiasbynens.be/demo/url-regex
var embedRegex = regexp.MustCompile(`(https?|ftp):\/\/[^\s\/$.?#].[^\s]*`)

func main() {
	// Load the .env file.
	if err := godotenv.Load(); err != nil {
		log.Println("warning: failed to load .env:", err)
	}

	// Get the bot token.
	var token = os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatalln("No $BOT_TOKEN given.")
	}

	// Parse the command line arguments.
	flag.Parse()
	s := session.New("Bot " + token)
	s.AddHandler(func(c *gateway.MessageCreateEvent) {
		// Check if the message is in one of the specified channel IDs.
		if !slices.Contains(flag.Args(), c.ChannelID.String()) {
			return
		}

		// Check if the message has attachments or embeds.
		if len(c.Message.Attachments) > 0 || containsEmbeds(c) {
			return
		}

		// Send a DM to the user.
		channel, err := s.CreatePrivateChannel(c.Author.ID)
		if err != nil {
			log.Println("Failed to create private channel:", err)
			return
		}

		if _, err := s.SendMessageComplex(channel.ID, api.SendMessageData{
			Content: "Please attach an image or a link to a showcase channel.",
		}); err != nil {
			log.Println("Failed to send message:", err)
		}

		if _, err := s.SendMessageComplex(channel.ID, api.SendMessageData{
			Content: c.Message.Content,
		}); err != nil {
			log.Println("Failed to send message:", err)
		}

		// Delete the message.
		if err := s.DeleteMessage(c.ChannelID, c.ID, "No attachments"); err != nil {
			log.Println("Failed to delete message:", err)
		}
	})

	// Add the needed Gateway intents.
	s.AddIntents(gateway.IntentGuildMessages)
	s.AddIntents(gateway.IntentMessageContent)
	s.AddIntents(gateway.IntentDirectMessages)
	s.AddIntents(gateway.IntentGuildMessages)

	// Get the bot's user.
	u, err := s.Me()
	if err != nil {
		log.Fatalln("Failed to get myself:", err)
	}

	log.Println("Started as", u.Username)

	// Open a connection to Discord.
	if err := s.Connect(context.Background()); err != nil {
		log.Fatalln("Failed to connect:", err)
	}

	log.Println("Shutting down...")
}

// containsEmbeds checks if the message contains embeds regular expressions.
func containsEmbeds(c *gateway.MessageCreateEvent) bool {
	if len(c.Message.Embeds) > 0 {
		return true
	}

	return embedRegex.MatchString(c.Message.Content)
}
