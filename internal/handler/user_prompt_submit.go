package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rxtech-lab/claude-code-telegram-notification/internal/telegram"
	"github.com/rxtech-lab/claude-code-telegram-notification/pkg/templates"
)

type UserPromptSubmitHandler struct {
	telegramClient *telegram.Client
}

type UserPromptSubmitEvent struct {
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message,omitempty"`
	ToolUsed  string    `json:"tool_used,omitempty"`
	FilePath  string    `json:"file_path,omitempty"`
}

type UserPromptSubmitTemplateData struct {
	EventType string
	Timestamp string
	Message   string
	ToolUsed  string
	FilePath  string
}

const userPromptSubmitTemplate = `🤖 *Claude Code Hook: User Prompt Submit*

📝 **Event:** {{.EventType}}
⏰ **Timestamp:** {{.Timestamp}}
{{if .ToolUsed}}🔧 **Tool Used:** {{.ToolUsed}}{{end}}
{{if .Message}}💬 **Message:** {{.Message}}{{end}}
{{if .FilePath}}📁 **File:** {{.FilePath}}{{end}}`

func NewUserPromptSubmitHandler(telegramClient *telegram.Client) *UserPromptSubmitHandler {
	return &UserPromptSubmitHandler{
		telegramClient: telegramClient,
	}
}


func (h *UserPromptSubmitHandler) Handle(ctx context.Context, hookEvent string) error {
	var event UserPromptSubmitEvent
	if err := json.Unmarshal([]byte(hookEvent), &event); err != nil {
		return err
	}

	templateData := UserPromptSubmitTemplateData{
		EventType: event.EventType,
		Timestamp: event.Timestamp.Format("2006-01-02 15:04:05"),
		Message:   event.Message,
		ToolUsed:  event.ToolUsed,
		FilePath:  event.FilePath,
	}

	message, err := templates.RenderTemplate("userPromptSubmit", userPromptSubmitTemplate, templateData)
	if err != nil {
		return err
	}

	return h.telegramClient.SendMessage(message)
}