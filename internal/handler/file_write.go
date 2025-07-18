package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rxtech-lab/claude-code-telegram-notification/internal/telegram"
	"github.com/rxtech-lab/claude-code-telegram-notification/pkg/templates"
)

type FileWriteHandler struct {
	telegramClient *telegram.Client
}

type FileWriteEvent struct {
	EventType    string    `json:"event_type"`
	Timestamp    time.Time `json:"timestamp"`
	FilePath     string    `json:"file_path"`
	Success      bool      `json:"success"`
	LinesAdded   int       `json:"lines_added,omitempty"`
	LinesDeleted int       `json:"lines_deleted,omitempty"`
}

func NewFileWriteHandler(telegramClient *telegram.Client) *FileWriteHandler {
	return &FileWriteHandler{
		telegramClient: telegramClient,
	}
}

func (h *FileWriteHandler) Handle(ctx context.Context, hookEvent string) error {
	var event FileWriteEvent
	if err := json.Unmarshal([]byte(hookEvent), &event); err != nil {
		return err
	}

	templateData := templates.TemplateData{
		EventType:    event.EventType,
		Timestamp:    event.Timestamp.Format("2006-01-02 15:04:05"),
		FilePath:     event.FilePath,
		Success:      event.Success,
		LinesAdded:   event.LinesAdded,
		LinesDeleted: event.LinesDeleted,
	}

	message, err := templates.RenderTemplate(templates.FileWriteTemplate, templateData)
	if err != nil {
		return err
	}

	return h.telegramClient.SendMessage(message)
}