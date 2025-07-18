package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rxtech-lab/claude-code-telegram-notification/internal/telegram"
	"github.com/rxtech-lab/claude-code-telegram-notification/pkg/templates"
)

type GenericHandler struct {
	telegramClient *telegram.Client
}

type GenericEvent struct {
	EventType string                 `json:"event_type"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

type GenericTemplateData struct {
	EventType string
	Timestamp string
	Data      map[string]interface{}
}

const genericTemplate = `📋 *Claude Code Hook: {{.EventType}}*

⏰ **Timestamp:** {{.Timestamp}}
{{range $key, $value := .Data}}
**{{$key}}:** {{$value}}
{{end}}`

func NewGenericHandler(telegramClient *telegram.Client) *GenericHandler {
	return &GenericHandler{
		telegramClient: telegramClient,
	}
}

func (h *GenericHandler) Handle(ctx context.Context, hookEvent string) error {
	var event GenericEvent
	if err := json.Unmarshal([]byte(hookEvent), &event); err != nil {
		return err
	}

	templateData := GenericTemplateData{
		EventType: event.EventType,
		Timestamp: event.Timestamp.Format("2006-01-02 15:04:05"),
		Data:      event.Data,
	}

	message, err := templates.RenderTemplate("generic", genericTemplate, templateData)
	if err != nil {
		return err
	}

	return h.telegramClient.SendMessage(message)
}
