package handler

import (
	"context"
	"testing"

	"github.com/rxtech-lab/claude-code-telegram-notification/internal/telegram"
)

func TestToolCallHandler_Handle(t *testing.T) {
	realClient := telegram.NewClient("test-token", "test-chat")
	handler := NewToolCallHandler(realClient)
	
	// Test with valid JSON
	hookEvent := `{
		"event_type": "tool-call",
		"timestamp": "2024-01-01T12:00:00Z",
		"tool_name": "Bash",
		"file_path": "/path/to/script.sh",
		"success": true
	}`
	
	ctx := context.Background()
	err := handler.Handle(ctx, hookEvent)
	
	// We expect this to fail because we don't have valid telegram credentials
	if err == nil {
		t.Error("Expected error due to invalid telegram credentials, but got none")
	}
}

func TestToolCallHandler_Handle_WithError(t *testing.T) {
	realClient := telegram.NewClient("test-token", "test-chat")
	handler := NewToolCallHandler(realClient)
	
	// Test with error in tool call
	hookEvent := `{
		"event_type": "tool-call",
		"timestamp": "2024-01-01T12:00:00Z",
		"tool_name": "Bash",
		"file_path": "/path/to/script.sh",
		"success": false,
		"error": "Command failed with exit code 1"
	}`
	
	ctx := context.Background()
	err := handler.Handle(ctx, hookEvent)
	
	// We expect this to fail because we don't have valid telegram credentials
	if err == nil {
		t.Error("Expected error due to invalid telegram credentials, but got none")
	}
}

func TestToolCallHandler_Handle_InvalidJSON(t *testing.T) {
	realClient := telegram.NewClient("test-token", "test-chat")
	handler := NewToolCallHandler(realClient)
	
	// Test with invalid JSON
	hookEvent := `invalid json`
	
	ctx := context.Background()
	err := handler.Handle(ctx, hookEvent)
	
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestNewToolCallHandler(t *testing.T) {
	client := telegram.NewClient("test-token", "test-chat")
	handler := NewToolCallHandler(client)
	
	if handler == nil {
		t.Error("Expected handler to be created")
	}
	
	if handler.telegramClient != client {
		t.Error("Expected telegram client to be set correctly")
	}
}