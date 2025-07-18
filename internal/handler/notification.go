package handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/rxtech-lab/claude-code-telegram-notification/internal/telegram"
	"github.com/rxtech-lab/claude-code-telegram-notification/pkg/templates"
)

type NotificationHandler struct {
	telegramClient *telegram.Client
}

type NotificationEvent struct {
	EventType string                 `json:"event_type"`
	Timestamp time.Time              `json:"timestamp"`
	Message   string                 `json:"message,omitempty"`
	Level     string                 `json:"level,omitempty"`
	Source    string                 `json:"source,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

func NewNotificationHandler(telegramClient *telegram.Client) *NotificationHandler {
	return &NotificationHandler{
		telegramClient: telegramClient,
	}
}

func (h *NotificationHandler) Handle(ctx context.Context, hookEvent string) error {
	var event NotificationEvent
	if err := json.Unmarshal([]byte(hookEvent), &event); err != nil {
		return err
	}

	log.Println("Received notification event:", event)
	templateData := templates.TemplateData{
		EventType: event.EventType,
		Timestamp: event.Timestamp.Format("2006-01-02 15:04:05"),
		Message:   event.Message,
		Data:      event.Data,
	}

	message, err := templates.RenderTemplate(templates.NotificationTemplate, templateData)
	if err != nil {
		return err
	}

	return h.telegramClient.SendMessage(message)
}
