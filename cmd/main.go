package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rxtech-lab/claude-code-telegram-notification/internal/handler"
	"github.com/rxtech-lab/claude-code-telegram-notification/internal/telegram"
)

func main() {
	// Check if user wants to run chat client
	if len(os.Args) > 1 && strings.ToLower(os.Args[1]) == "chat-client" {
		runChatClient()
		return
	}

	// Read environment variables
	botToken := os.Getenv("CLAUDE_CODE_TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("CLAUDE_CODE_TELEGRAM_BOT_TOKEN environment variable is required")
	}

	chatID := os.Getenv("CLAUDE_CODE_TELEGRAM_CHAT_ID")
	if chatID == "" {
		log.Fatal("CLAUDE_CODE_TELEGRAM_CHAT_ID environment variable is required")
	}

	// Parse command line arguments for hook type
	hookType := "generic"
	if len(os.Args) > 1 {
		hookType = strings.ToLower(os.Args[1])
	}

	// Create Telegram client
	telegramClient := telegram.NewClient(botToken, chatID)

	// Create appropriate handler based on hook type
	var hookHandler handler.Handler
	switch hookType {
	case "user-prompt-submit":
		hookHandler = handler.NewUserPromptSubmitHandler(telegramClient)
	case "tool-call":
		hookHandler = handler.NewToolCallHandler(telegramClient)
	case "file-write":
		hookHandler = handler.NewFileWriteHandler(telegramClient)
	case "session-start":
		hookHandler = handler.NewSessionStartHandler(telegramClient)
	case "session-end":
		hookHandler = handler.NewSessionEndHandler(telegramClient)
	default:
		hookHandler = handler.NewGenericHandler(telegramClient)
	}

	// Read hook data from stdin
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading from stdin: %v", err)
	}

	if len(lines) == 0 {
		log.Fatal("No hook data received from stdin")
	}

	// Join all lines to form the complete JSON
	hookData := strings.Join(lines, "\n")

	// Handle the hook event
	ctx := context.Background()
	if err := hookHandler.Handle(ctx, hookData); err != nil {
		log.Fatalf("Error handling hook event: %v", err)
	}

	fmt.Println("Notification sent successfully")
}

func runChatClient() {
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
