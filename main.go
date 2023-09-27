package main

import (
	"context"
	"flag"
	"log"
	"os"
	"regexp"

	"golang.org/x/exp/slices"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
	"github.com/joho/godotenv"
)

// embedRegex is a regular expression that matches embed URLs.
// https://mathiasbynens.be/demo/url-regex
var embedRegex = regexp.MustCompile(`https?:\/\/[^\s\/$.?#].[^\s]*`)

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

	// Add a handler for the message create and update events.
	handle := func(m discord.Message) {
		// Check if the message is in one of the specified channel IDs.
		if !slices.Contains(flag.Args(), m.ChannelID.String()) {
			return
		}

		// Check if the message has attachments or embeds.
		if len(m.Attachments) > 0 || containsEmbeds(m) {
			return
		}

		// Send a DM to the user.
		channel, err := s.CreatePrivateChannel(m.Author.ID)
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
			Content: m.Content,
		}); err != nil {
			log.Println("Failed to send message:", err)
		}

		// Delete the message.
		if err := s.DeleteMessage(m.ChannelID, m.ID, "No attachments"); err != nil {
			log.Println("Failed to delete message:", err)
		}
	}

	s.AddHandler(func(c *gateway.MessageCreateEvent) {
		handle(c.Message)
	})

	s.AddHandler(func(c *gateway.MessageUpdateEvent) {
		handle(c.Message)
	})

	// Add the needed Gateway intents.
	s.AddIntents(gateway.IntentGuildMessages)
	s.AddIntents(gateway.IntentMessageContent)

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
func containsEmbeds(m discord.Message) bool {
	if len(m.Embeds) > 0 {
		return true
	}

	return embedRegex.MatchString(m.Content)
}
