package templates

import (
	"bytes"
	"text/template"
)

func RenderTemplate[T any](templateName string, templateStr string, data T) (string, error) {
	tmpl, err := template.New(templateName).Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
