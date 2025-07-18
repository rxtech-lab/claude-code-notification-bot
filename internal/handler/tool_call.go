package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rxtech-lab/claude-code-telegram-notification/internal/telegram"
	"github.com/rxtech-lab/claude-code-telegram-notification/pkg/templates"
)

type ToolCallHandler struct {
	telegramClient *telegram.Client
}

type ToolCallEvent struct {
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	ToolName  string    `json:"tool_name"`
	FilePath  string    `json:"file_path,omitempty"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
}

type ToolCallTemplateData struct {
	EventType string
	Timestamp string
	ToolName  string
	FilePath  string
	Success   bool
	Error     string
}

const toolCallTemplate = `🔨 *Claude Code Hook: Tool Call*

📝 **Event:** {{.EventType}}
🛠️ **Tool:** {{.ToolName}}
⏰ **Timestamp:** {{.Timestamp}}
{{if .FilePath}}📁 **File:** {{.FilePath}}{{end}}
{{if .Success}}✅ **Status:** Success{{else}}❌ **Status:** Failed{{end}}
{{if .Error}}⚠️ **Error:** {{.Error}}{{end}}`

func NewToolCallHandler(telegramClient *telegram.Client) *ToolCallHandler {
	return &ToolCallHandler{
		telegramClient: telegramClient,
	}
}


func (h *ToolCallHandler) Handle(ctx context.Context, hookEvent string) error {
	var event ToolCallEvent
	if err := json.Unmarshal([]byte(hookEvent), &event); err != nil {
		return err
	}

	templateData := ToolCallTemplateData{
		EventType: event.EventType,
		Timestamp: event.Timestamp.Format("2006-01-02 15:04:05"),
		ToolName:  event.ToolName,
		FilePath:  event.FilePath,
		Success:   event.Success,
		Error:     event.Error,
	}

	message, err := templates.RenderTemplate("toolCall", toolCallTemplate, templateData)
	if err != nil {
		return err
	}

	return h.telegramClient.SendMessage(message)
}