package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
)

var channelIDs []string

func main() {
	var token = os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatalln("No $BOT_TOKEN given.")
	}

	flag.Parse()
	channelIDs = parseChannelIDsFromArgs(flag.Args())

	s := session.New("Bot " + token)
	s.AddHandler(func(c *gateway.MessageCreateEvent) {
		// Check if the message is in one of the specified channel IDs.
		if !contains(channelIDs, c.ChannelID.String()) {
			return
		}

		// Check if the message has attachments.
		if c.Message.Attachments == nil {
			return
		}

		// Delete the message.
		if err := s.DeleteMessage(c.ChannelID, c.ID, "No attachments"); err != nil {
			log.Println("Failed to delete message:", err)
		}
	})

	// Add the needed Gateway intents.
	s.AddIntents(gateway.IntentGuildMessages)
	s.AddIntents(gateway.IntentDirectMessages)

	if err := s.Open(context.Background()); err != nil {
		log.Fatalln("Failed to connect:", err)
	}
	defer s.Close()

	u, err := s.Me()
	if err != nil {
		log.Fatalln("Failed to get myself:", err)
	}

	log.Println("Started as", u.Username)

	// Block forever.
	select {}
}

// parseChannelIDsFromArgs parses the channel IDs from the command arguments.
func parseChannelIDsFromArgs(args []string) []string {
	var channelIDs []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "channel:") {
			channelID := strings.TrimPrefix(arg, "channel:")
			channelIDs = append(channelIDs, channelID)
		}
	}
	return channelIDs
}

// contains checks if a string slice contains a specific string.
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}