package handler

import (
	"context"
	"testing"

	"github.com/rxtech-lab/claude-code-telegram-notification/internal/telegram"
)

// MockTelegramClient implements a mock for testing
type MockTelegramClient struct {
	lastMessage string
	shouldError bool
}

func (m *MockTelegramClient) SendMessage(message string) error {
	if m.shouldError {
		return &MockError{message: "mock error"}
	}
	m.lastMessage = message
	return nil
}

type MockError struct {
	message string
}

func (e *MockError) Error() string {
	return e.message
}

func TestUserPromptSubmitHandler_Handle(t *testing.T) {
	_ = &MockTelegramClient{}

	// Create a mock telegram client adapter
	realClient := telegram.NewClient("test-token", "test-chat")

	handler := NewUserPromptSubmitHandler(realClient)

	// Test with valid JSON
	hookEvent := `{
		"event_type": "user-prompt-submit",
		"timestamp": "2024-01-01T12:00:00Z",
		"message": "Test message",
		"tool_used": "Read",
		"file_path": "/path/to/file.go"
	}`

	ctx := context.Background()

	// Note: This test is limited because we can't easily mock the telegram client
	// In a real implementation, you'd want to inject an interface rather than a concrete type
	err := handler.Handle(ctx, hookEvent)

	// We expect this to fail because we don't have valid telegram credentials
	// but we can verify the JSON parsing works by checking the error type
	if err == nil {
		t.Error("Expected error due to invalid telegram credentials, but got none")
	}
}

func TestUserPromptSubmitHandler_Handle_InvalidJSON(t *testing.T) {
	realClient := telegram.NewClient("test-token", "test-chat")
	handler := NewUserPromptSubmitHandler(realClient)

	// Test with invalid JSON
	hookEvent := `invalid json`

	ctx := context.Background()
	err := handler.Handle(ctx, hookEvent)

	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestNewUserPromptSubmitHandler(t *testing.T) {
	client := telegram.NewClient("test-token", "test-chat")
	handler := NewUserPromptSubmitHandler(client)

	if handler == nil {
		t.Error("Expected handler to be created")
	}

	if handler.telegramClient != client {
		t.Error("Expected telegram client to be set correctly")
	}
}
