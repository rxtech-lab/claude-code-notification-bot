package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rxtech-lab/claude-code-telegram-notification/internal/telegram"
	"github.com/rxtech-lab/claude-code-telegram-notification/pkg/templates"
)

type SessionStartHandler struct {
	telegramClient *telegram.Client
}

type SessionStartEvent struct {
	EventType        string    `json:"event_type"`
	Timestamp        time.Time `json:"timestamp"`
	WorkingDirectory string    `json:"working_directory"`
	GitRepo          bool      `json:"git_repo"`
}

func NewSessionStartHandler(telegramClient *telegram.Client) *SessionStartHandler {
	return &SessionStartHandler{
		telegramClient: telegramClient,
	}
}

func (h *SessionStartHandler) Handle(ctx context.Context, hookEvent string) error {
	var event SessionStartEvent
	if err := json.Unmarshal([]byte(hookEvent), &event); err != nil {
		return err
	}

	templateData := templates.TemplateData{
		EventType:        event.EventType,
		Timestamp:        event.Timestamp.Format("2006-01-02 15:04:05"),
		WorkingDirectory: event.WorkingDirectory,
		GitRepo:          event.GitRepo,
	}

	message, err := templates.RenderTemplate(templates.SessionStartTemplate, templateData)
	if err != nil {
		return err
	}

	return h.telegramClient.SendMessage(message)
}