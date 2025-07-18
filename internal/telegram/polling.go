package telegram

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  *Message `json:"message,omitempty"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	From      *User  `json:"from,omitempty"`
	Chat      *Chat  `json:"chat,omitempty"`
	Date      int64  `json:"date"`
	Text      string `json:"text,omitempty"`
}

type User struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

type Chat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

type GetUpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
	Description string `json:"description,omitempty"`
}

type PollingClient struct {
	botToken   string
	httpClient *http.Client
	offset     int
}

func NewPollingClient(botToken string) *PollingClient {
	return &PollingClient{
		botToken:   botToken,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		offset:     0,
	}
}

func (c *PollingClient) GetUpdates() ([]Update, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d&timeout=10", c.botToken, c.offset)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get updates: %w", err)
	}
	defer resp.Body.Close()
	
	var response GetUpdatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	if !response.Ok {
		return nil, fmt.Errorf("telegram API error: %s", response.Description)
	}
	
	// Update offset to mark messages as processed
	if len(response.Result) > 0 {
		c.offset = response.Result[len(response.Result)-1].UpdateID + 1
	}
	
	return response.Result, nil
}

func (c *PollingClient) StartChatIDDiscovery() {
	fmt.Println("🤖 Telegram Chat ID Discovery")
	fmt.Println("===============================")
	fmt.Println("Send any message to your bot to discover your chat ID...")
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println()
	
	seenChats := make(map[int64]bool)
	
	for {
		updates, err := c.GetUpdates()
		if err != nil {
			log.Printf("Error getting updates: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		
		for _, update := range updates {
			if update.Message != nil && update.Message.Chat != nil {
				chatID := update.Message.Chat.ID
				
				// Skip if we've already seen this chat
				if seenChats[chatID] {
					continue
				}
				seenChats[chatID] = true
				
				// Display chat information
				fmt.Printf("📨 New message received!\n")
				fmt.Printf("   Chat ID: %d\n", chatID)
				fmt.Printf("   Chat Type: %s\n", update.Message.Chat.Type)
				
				if update.Message.From != nil {
					fmt.Printf("   From: %s", update.Message.From.FirstName)
					if update.Message.From.LastName != "" {
						fmt.Printf(" %s", update.Message.From.LastName)
					}
					if update.Message.From.Username != "" {
						fmt.Printf(" (@%s)", update.Message.From.Username)
					}
					fmt.Println()
				}
				
				if update.Message.Chat.Title != "" {
					fmt.Printf("   Group/Channel: %s\n", update.Message.Chat.Title)
				}
				
				fmt.Printf("   Message: %s\n", update.Message.Text)
				fmt.Println()
				
				// Provide copy-paste instructions
				fmt.Printf("💡 To use this chat, set your environment variable:\n")
				fmt.Printf("   export CLAUDE_CODE_TELEGRAM_CHAT_ID=\"%d\"\n", chatID)
				fmt.Println()
				fmt.Println("---")
				fmt.Println()
			}
		}
		
		time.Sleep(1 * time.Second)
	}
}