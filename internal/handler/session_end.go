package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rxtech-lab/claude-code-telegram-notification/internal/telegram"
	"github.com/rxtech-lab/claude-code-telegram-notification/pkg/templates"
)

type SessionEndHandler struct {
	telegramClient *telegram.Client
}

type SessionEndEvent struct {
	EventType     string    `json:"event_type"`
	Timestamp     time.Time `json:"timestamp"`
	Duration      string    `json:"duration"`
	FilesModified int       `json:"files_modified,omitempty"`
}

func NewSessionEndHandler(telegramClient *telegram.Client) *SessionEndHandler {
	return &SessionEndHandler{
		telegramClient: telegramClient,
	}
}

func (h *SessionEndHandler) Handle(ctx context.Context, hookEvent string) error {
	var event SessionEndEvent
	if err := json.Unmarshal([]byte(hookEvent), &event); err != nil {
		return err
	}

	templateData := templates.TemplateData{
		EventType:     event.EventType,
		Timestamp:     event.Timestamp.Format("2006-01-02 15:04:05"),
		Duration:      event.Duration,
		FilesModified: event.FilesModified,
	}

	message, err := templates.RenderTemplate(templates.SessionEndTemplate, templateData)
	if err != nil {
		return err
	}

	return h.telegramClient.SendMessage(message)
}