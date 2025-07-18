package templates

import (
	"bytes"
	"text/template"
)

var (
	// Template for user prompt submit hook
	UserPromptSubmitTemplate = `🤖 *Claude Code Hook: User Prompt Submit*

📝 **Event:** {{.EventType}}
⏰ **Timestamp:** {{.Timestamp}}
{{if .ToolUsed}}🔧 **Tool Used:** {{.ToolUsed}}{{end}}
{{if .Message}}💬 **Message:** {{.Message}}{{end}}
{{if .FilePath}}📁 **File:** {{.FilePath}}{{end}}`

	// Template for tool call hook
	ToolCallTemplate = `🔨 *Claude Code Hook: Tool Call*

📝 **Event:** {{.EventType}}
🛠️ **Tool:** {{.ToolName}}
⏰ **Timestamp:** {{.Timestamp}}
{{if .FilePath}}📁 **File:** {{.FilePath}}{{end}}
{{if .Success}}✅ **Status:** Success{{else}}❌ **Status:** Failed{{end}}
{{if .Error}}⚠️ **Error:** {{.Error}}{{end}}`

	// Template for file write hook
	FileWriteTemplate = `💾 *Claude Code Hook: File Write*

📝 **Event:** {{.EventType}}
📁 **File:** {{.FilePath}}
⏰ **Timestamp:** {{.Timestamp}}
{{if .LinesAdded}}➕ **Lines Added:** {{.LinesAdded}}{{end}}
{{if .LinesDeleted}}➖ **Lines Deleted:** {{.LinesDeleted}}{{end}}
{{if .Success}}✅ **Status:** Success{{else}}❌ **Status:** Failed{{end}}`

	// Template for session start hook
	SessionStartTemplate = `🚀 *Claude Code Hook: Session Start*

📝 **Event:** {{.EventType}}
⏰ **Timestamp:** {{.Timestamp}}
🗂️ **Working Directory:** {{.WorkingDirectory}}
{{if .GitRepo}}🔗 **Git Repository:** Yes{{else}}🔗 **Git Repository:** No{{end}}`

	// Template for session end hook
	SessionEndTemplate = `🏁 *Claude Code Hook: Session End*

📝 **Event:** {{.EventType}}
⏰ **Timestamp:** {{.Timestamp}}
⌛ **Duration:** {{.Duration}}
{{if .FilesModified}}📝 **Files Modified:** {{.FilesModified}}{{end}}`

	// Generic template for unknown hook types
	GenericTemplate = `📋 *Claude Code Hook: {{.EventType}}*

⏰ **Timestamp:** {{.Timestamp}}
{{range $key, $value := .Data}}
**{{$key}}:** {{$value}}
{{end}}`
)

type TemplateData struct {
	EventType        string
	Timestamp        string
	ToolUsed         string
	ToolName         string
	Message          string
	FilePath         string
	Success          bool
	Error            string
	LinesAdded       int
	LinesDeleted     int
	WorkingDirectory string
	GitRepo          bool
	Duration         string
	FilesModified    int
	Data             map[string]interface{}
}

func RenderTemplate(templateStr string, data TemplateData) (string, error) {
	tmpl, err := template.New("hook").Parse(templateStr)
	if err != nil {
		return "", err
	}
	
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	
	return buf.String(), nil
}