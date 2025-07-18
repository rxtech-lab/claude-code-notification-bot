package handler

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/rxtech-lab/claude-code-telegram-notification/internal/telegram"
	"github.com/rxtech-lab/claude-code-telegram-notification/pkg/templates"
)

type NotificationHandler struct {
	telegramClient *telegram.Client
}

type NotificationEvent struct {
	SessionID      string `json:"session_id"`
	TranscriptPath string `json:"transcript_path"`
	Cwd            string `json:"cwd"`
	HookEventName  string `json:"hook_event_name"`
	Message        string `json:"message"`
}

type NotificationTemplateData struct {
	EventType      string
	Timestamp      string
	Message        string
	SessionID      string
	TranscriptPath string
	Cwd            string
}

const notificationTemplate = `🔔 *Claude Code Notification*
{{.Message}}

📝 **Event:** {{.EventType}}
⏰ **Timestamp:** {{.Timestamp}}
📁 **Working Directory:** {{.Cwd}}`

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

	// skip the message if it's a waiting for input message
	if strings.Contains(event.Message, "Claude is waiting for your input") {
		return nil
	}

	log.Println("Received notification event:", event)
	templateData := NotificationTemplateData{
		EventType:      event.HookEventName,
		Timestamp:      time.Now().Format("2006-01-02 15:04:05"),
		Message:        event.Message,
		SessionID:      event.SessionID,
		TranscriptPath: event.TranscriptPath,
		Cwd:            event.Cwd,
	}

	message, err := templates.RenderTemplate("notification", notificationTemplate, templateData)
	if err != nil {
		return err
	}

	return h.telegramClient.SendMessage(message)
}
