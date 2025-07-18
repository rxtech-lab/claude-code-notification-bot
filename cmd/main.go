package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/rxtech-lab/claude-code-telegram-notification/internal/config"
	"github.com/rxtech-lab/claude-code-telegram-notification/internal/handler"
	"github.com/rxtech-lab/claude-code-telegram-notification/internal/telegram"
)

func main() {
	// Setup file logging
	setupFileLogging()

	// Parse command line arguments
	var (
		token  = flag.String("token", "", "Telegram bot token (for initialization)")
		chatID = flag.String("chatid", "", "Telegram chat ID (for initialization)")
		help   = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	// Show help if requested
	if *help {
		showHelp()
		return
	}

	// Handle initialization with --token and --chatid
	if *token != "" && *chatID != "" {
		initializeConfig(*token, *chatID)
		return
	}

	// Check if user wants to run chat client (non-flag arguments)
	args := flag.Args()
	if len(args) > 0 && strings.ToLower(args[0]) == "chat-client" {
		runChatClient()
		return
	}

	// Load configuration from file
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Parse command line arguments for hook type
	hookType := "generic"
	if len(args) > 0 {
		hookType = strings.ToLower(args[0])
	}

	// Create Telegram client
	telegramClient := telegram.NewClient(cfg.BotToken, cfg.ChatID)

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
	case "notification":
		hookHandler = handler.NewNotificationHandler(telegramClient)
	case "stop":
		hookHandler = handler.NewStopHandler(telegramClient)
	default:
		log.Fatalf("Invalid hook type: %s", hookType)
	}

	log.Println("Hook handler created:", hookType)
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

func setupFileLogging() {
	// Get log file path from environment variable, default to ~/claude-code-notification/notification.log
	logFilePath := os.Getenv("CLAUDE_CODE_LOG_FILE")
	if logFilePath == "" {
		configDir, err := config.GetConfigDir()
		if err != nil {
			log.Printf("Warning: Could not get config directory: %v", err)
			logFilePath = "notification.log"
		} else {
			// Create config directory if it doesn't exist
			if err := os.MkdirAll(configDir, 0755); err != nil {
				log.Printf("Warning: Could not create config directory: %v", err)
				logFilePath = "notification.log"
			} else {
				logFilePath = configDir + "/notification.log"
			}
		}
	}

	// Create or open log file
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Warning: Could not open log file %s: %v", logFilePath, err)
		return
	}

	// Set log output to both file and stdout
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("Logging to file: %s", logFilePath)
}

func initializeConfig(token, chatID string) {
	cfg := &config.Config{
		BotToken: token,
		ChatID:   chatID,
	}

	if err := config.SaveConfig(cfg); err != nil {
		log.Fatalf("Failed to save config: %v", err)
	}

	configFile, _ := config.GetConfigFilePath()
	fmt.Printf("Configuration saved to: %s\n", configFile)
	fmt.Println("You can now run the application without --token and --chatid arguments.")
}

func showHelp() {
	fmt.Println("Claude Code Telegram Notification")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  Initialize configuration:")
	fmt.Println("    ./claude-code-telegram-notification --token=<bot_token> --chatid=<chat_id>")
	fmt.Println()
	fmt.Println("  Run notification hooks:")
	fmt.Println("    ./claude-code-telegram-notification [hook_type]")
	fmt.Println("    Hook types: user-prompt-submit, tool-call, file-write, session-start, session-end, notification, stop")
	fmt.Println()
	fmt.Println("  Run chat client:")
	fmt.Println("    ./claude-code-telegram-notification chat-client")
	fmt.Println()
	fmt.Println("Environment variables:")
	fmt.Println("  CLAUDE_CODE_LOG_FILE - Custom log file path (default: ~/claude-code-notification/notification.log)")
}

func runChatClient() {
	// Setup file logging for chat client too
	setupFileLogging()

	// Load configuration from file
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create polling client
	client := telegram.NewPollingClient(cfg.BotToken)

	// Start chat ID discovery
	client.StartChatIDDiscovery()
}
