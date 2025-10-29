package utils

import (
	"bytes"
	"html/template"
)

func RenderEmailTemplate(filePath string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
