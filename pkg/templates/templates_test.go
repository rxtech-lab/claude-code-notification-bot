package templates

import (
	"strings"
	"testing"
)

func TestRenderTemplate_UserPromptSubmit(t *testing.T) {
	data := TemplateData{
		EventType: "user-prompt-submit",
		Timestamp: "2024-01-01 12:00:00",
		Message:   "Test message",
		ToolUsed:  "Read",
		FilePath:  "/path/to/file.go",
	}
	
	result, err := RenderTemplate(UserPromptSubmitTemplate, data)
	if err != nil {
		t.Fatalf("Error rendering template: %v", err)
	}
	
	if !strings.Contains(result, "user-prompt-submit") {
		t.Error("Result should contain event type")
	}
	
	if !strings.Contains(result, "Test message") {
		t.Error("Result should contain message")
	}
	
	if !strings.Contains(result, "Read") {
		t.Error("Result should contain tool used")
	}
	
	if !strings.Contains(result, "/path/to/file.go") {
		t.Error("Result should contain file path")
	}
}

func TestRenderTemplate_ToolCall(t *testing.T) {
	data := TemplateData{
		EventType: "tool-call",
		Timestamp: "2024-01-01 12:00:00",
		ToolName:  "Bash",
		FilePath:  "/path/to/script.sh",
		Success:   true,
	}
	
	result, err := RenderTemplate(ToolCallTemplate, data)
	if err != nil {
		t.Fatalf("Error rendering template: %v", err)
	}
	
	if !strings.Contains(result, "tool-call") {
		t.Error("Result should contain event type")
	}
	
	if !strings.Contains(result, "Bash") {
		t.Error("Result should contain tool name")
	}
	
	if !strings.Contains(result, "Success") {
		t.Error("Result should contain success status")
	}
}

func TestRenderTemplate_FileWrite(t *testing.T) {
	data := TemplateData{
		EventType:    "file-write",
		Timestamp:    "2024-01-01 12:00:00",
		FilePath:     "/path/to/file.go",
		Success:      true,
		LinesAdded:   10,
		LinesDeleted: 5,
	}
	
	result, err := RenderTemplate(FileWriteTemplate, data)
	if err != nil {
		t.Fatalf("Error rendering template: %v", err)
	}
	
	if !strings.Contains(result, "file-write") {
		t.Error("Result should contain event type")
	}
	
	if !strings.Contains(result, "/path/to/file.go") {
		t.Error("Result should contain file path")
	}
	
	if !strings.Contains(result, "10") {
		t.Error("Result should contain lines added")
	}
	
	if !strings.Contains(result, "5") {
		t.Error("Result should contain lines deleted")
	}
}

func TestRenderTemplate_SessionStart(t *testing.T) {
	data := TemplateData{
		EventType:        "session-start",
		Timestamp:        "2024-01-01 12:00:00",
		WorkingDirectory: "/home/user/project",
		GitRepo:          true,
	}
	
	result, err := RenderTemplate(SessionStartTemplate, data)
	if err != nil {
		t.Fatalf("Error rendering template: %v", err)
	}
	
	if !strings.Contains(result, "session-start") {
		t.Error("Result should contain event type")
	}
	
	if !strings.Contains(result, "/home/user/project") {
		t.Error("Result should contain working directory")
	}
	
	if !strings.Contains(result, "Yes") {
		t.Error("Result should indicate git repository")
	}
}

func TestRenderTemplate_SessionEnd(t *testing.T) {
	data := TemplateData{
		EventType:     "session-end",
		Timestamp:     "2024-01-01 12:00:00",
		Duration:      "30m",
		FilesModified: 5,
	}
	
	result, err := RenderTemplate(SessionEndTemplate, data)
	if err != nil {
		t.Fatalf("Error rendering template: %v", err)
	}
	
	if !strings.Contains(result, "session-end") {
		t.Error("Result should contain event type")
	}
	
	if !strings.Contains(result, "30m") {
		t.Error("Result should contain duration")
	}
	
	if !strings.Contains(result, "5") {
		t.Error("Result should contain files modified count")
	}
}

func TestRenderTemplate_Generic(t *testing.T) {
	data := TemplateData{
		EventType: "custom-event",
		Timestamp: "2024-01-01 12:00:00",
		Data: map[string]interface{}{
			"custom_field": "custom_value",
			"number":       42,
		},
	}
	
	result, err := RenderTemplate(GenericTemplate, data)
	if err != nil {
		t.Fatalf("Error rendering template: %v", err)
	}
	
	if !strings.Contains(result, "custom-event") {
		t.Error("Result should contain event type")
	}
	
	if !strings.Contains(result, "custom_field") {
		t.Error("Result should contain custom field")
	}
	
	if !strings.Contains(result, "custom_value") {
		t.Error("Result should contain custom value")
	}
}

func TestRenderTemplate_InvalidTemplate(t *testing.T) {
	invalidTemplate := "{{.InvalidField"
	data := TemplateData{
		EventType: "test",
		Timestamp: "2024-01-01 12:00:00",
	}
	
	_, err := RenderTemplate(invalidTemplate, data)
	if err == nil {
		t.Error("Expected error for invalid template")
	}
}