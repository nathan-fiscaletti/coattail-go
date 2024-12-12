package templates

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed sdk-template/*
var sdkTemplates embed.FS

type SdkTemplateData struct {
	Actions []ActionTemplateData `yaml:"actions"`

	templates *embed.FS
}

func NewSdkTemplate(data SdkTemplateData) Template {
	data.templates = &sdkTemplates
	return &data
}

func (d *SdkTemplateData) Fill(dir string) error {
	modTemplateFs, err := fs.Sub(d.templates, "sdk-template")
	if err != nil {
		return err
	}

	inputTemplate := "sdk.go.tmpl"
	filename := "sdk.go"
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
