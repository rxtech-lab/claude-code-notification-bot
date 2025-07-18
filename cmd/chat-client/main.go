package main

import (
	"log"
	"os"

	"github.com/rxtech-lab/claude-code-telegram-notification/internal/telegram"
)

func main() {
	// Check if stdin is available
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		// Stdin is a terminal (no piped input), exit
		os.Exit(1)
	}

	// Read bot token from environment
	botToken := os.Getenv("CLAUDE_CODE_TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("CLAUDE_CODE_TELEGRAM_BOT_TOKEN environment variable is required")
	}

	// Create polling client
	client := telegram.NewPollingClient(botToken)

	// Start chat ID discovery
	client.StartChatIDDiscovery()
}
