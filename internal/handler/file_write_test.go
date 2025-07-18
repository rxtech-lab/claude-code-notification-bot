package handler

import (
	"context"
	"testing"

	"github.com/rxtech-lab/claude-code-telegram-notification/internal/telegram"
)

func TestFileWriteHandler_Handle(t *testing.T) {
	realClient := telegram.NewClient("test-token", "test-chat")
	handler := NewFileWriteHandler(realClient)
	
	// Test with valid JSON
	hookEvent := `{
		"event_type": "file-write",
		"timestamp": "2024-01-01T12:00:00Z",
		"file_path": "/path/to/file.go",
		"success": true,
		"lines_added": 10,
		"lines_deleted": 5
	}`
	
	ctx := context.Background()
	err := handler.Handle(ctx, hookEvent)
	
	// We expect this to fail because we don't have valid telegram credentials
	if err == nil {
		t.Error("Expected error due to invalid telegram credentials, but got none")
	}
}

func TestFileWriteHandler_Handle_InvalidJSON(t *testing.T) {
	realClient := telegram.NewClient("test-token", "test-chat")
	handler := NewFileWriteHandler(realClient)
	
	// Test with invalid JSON
	hookEvent := `invalid json`
	
	ctx := context.Background()
	err := handler.Handle(ctx, hookEvent)
	
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestNewFileWriteHandler(t *testing.T) {
	client := telegram.NewClient("test-token", "test-chat")
	handler := NewFileWriteHandler(client)
	
	if handler == nil {
		t.Error("Expected handler to be created")
	}
	
	if handler.telegramClient != client {
		t.Error("Expected telegram client to be set correctly")
	}
}