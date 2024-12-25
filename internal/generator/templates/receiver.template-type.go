package templates

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed receiver-template/*
var receiverTemplates embed.FS

type ReceiverTemplateData struct {
	PackageName string `yaml:"package_name"`
	Name        string `yaml:"name"`
	InputType   string `yaml:"input"`

	templates *embed.FS
}

func NewReceiverTemplate(data ReceiverTemplateData) Template {
	data.templates = &receiverTemplates
	return &data
}

func (d *ReceiverTemplateData) Fill(dir string) error {
	modTemplateFs, err := fs.Sub(d.templates, "receiver-template")
	if err != nil {
		return err
	}

	inputTemplate := "receiver.go.tmpl"
	filename := fmt.Sprintf("receiver.%s.go", strings.ToLower(d.Name))
	outputFile := filepath.Join(dir, filename)

	// Create the target directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return err
	}

	// Read the file content from
	file, err := modTemplateFs.Open(inputTemplate)
	if err != nil {
		return err
	}
	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// Parse the template
	tmpl, err := template.New(filepath.Base(inputTemplate)).Parse(string(content))
	if err != nil {
		return err
	}

	// Create the destination file
	destFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Execute the template into the destination file
	if err := tmpl.Execute(destFile, d); err != nil {
		return err
	}

	return nil
}
