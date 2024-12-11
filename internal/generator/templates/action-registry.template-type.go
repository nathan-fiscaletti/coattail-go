package templates

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed action-registry-template/*
var actionRegistryTemplates embed.FS

type ActionRegistryTemplateData struct {
	Actions []ActionTemplateData `yaml:"actions"`

	templates *embed.FS
}

func NewActionRegistryTemplate(data ActionRegistryTemplateData) Template {
	data.templates = &actionRegistryTemplates
	return &data
}

func (d *ActionRegistryTemplateData) Fill(dir string) error {
	modTemplateFs, err := fs.Sub(d.templates, "action-registry-template")
	if err != nil {
		return err
	}

	inputTemplate := "action-registry.go.tmpl"
	filename := "action-registry.go"
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
