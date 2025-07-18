package handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/rxtech-lab/claude-code-telegram-notification/internal/telegram"
	"github.com/rxtech-lab/claude-code-telegram-notification/pkg/templates"
)

type StopHandler struct {
	telegramClient *telegram.Client
}

type StopEvent struct {
	SessionID      string `json:"session_id"`
	TranscriptPath string `json:"transcript_path"`
	Cwd            string `json:"cwd"`
	HookEventName  string `json:"hook_event_name"`
	Message        string `json:"message"`
}

type StopTemplateData struct {
	EventType      string
	Timestamp      string
	Message        string
	SessionID      string
	TranscriptPath string
	Cwd            string
}

const stopTemplate = `🛑 *Claude Code Session Stopped*

📝 **Event:** {{.EventType}}
⏰ **Timestamp:** {{.Timestamp}}
💬 **Message:** {{.Message}}
📁 **Working Directory:** {{.Cwd}}
🔗 **Session ID:** {{.SessionID}}`

func NewStopHandler(telegramClient *telegram.Client) *StopHandler {
	return &StopHandler{
		telegramClient: telegramClient,
	}
}

func (h *StopHandler) Handle(ctx context.Context, hookEvent string) error {
	var event StopEvent
	if err := json.Unmarshal([]byte(hookEvent), &event); err != nil {
		return err
	}

	log.Println("Received stop event:", event)
	templateData := StopTemplateData{
		EventType:      event.HookEventName,
		Timestamp:      time.Now().Format("2006-01-02 15:04:05"),
		Message:        event.Message,
		SessionID:      event.SessionID,
		TranscriptPath: event.TranscriptPath,
		Cwd:            event.Cwd,
	}

	message, err := templates.RenderTemplate("stop", stopTemplate, templateData)
	if err != nil {
		return err
	}

	return h.telegramClient.SendMessage(message)
}