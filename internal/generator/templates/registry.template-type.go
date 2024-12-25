package templates

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed registry-template/*
var registryTemplates embed.FS

type RegistryTemplateData struct {
	Actions     []ActionTemplateData   `yaml:"actions"`
	Receivers   []ReceiverTemplateData `yaml:"receivers"`
	PackageName string                 `yaml:"package_name"`

	templates *embed.FS
}

func NewRegistryTemplate(data RegistryTemplateData) Template {
	data.templates = &registryTemplates
	return &data
}

func (d *RegistryTemplateData) Fill(dir string) error {
	modTemplateFs, err := fs.Sub(d.templates, "registry-template")
	if err != nil {
		return err
	}

	inputTemplate := "registry.go.tmpl"
	filename := "registry.go"
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
