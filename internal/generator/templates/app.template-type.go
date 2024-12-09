package templates

import (
	"embed"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed app-template/**/*
var appTemplates embed.FS

type AppTemplateData struct {
	PackageName string

	templates *embed.FS
}

func NewAppTemplate(data AppTemplateData) Template {
	data.templates = &appTemplates
	return &data
}

func (d *AppTemplateData) Fill(dir string) error {
	// TODO: this logic should be abstracted somewhere or put into a
	// TODO: common function so that each generator can use it.
	// loop through each directory in the tmplFs and copy it to the dir
	// each file in the directory with the extension ".tmpl" should be fille in using the template data
	// if a file does not have the .tmpl extension, it should be copied as-is

	// loop through each directory in the tmplFs and copy it to the dir
	err := fs.WalkDir(d.templates, ".", func(path string, _d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if _d.IsDir() {
			return nil
		}

		// Check if the file ends with .tmpl
		if !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		// Determine the relative path without the .tmpl suffix
		relativePath := strings.TrimSuffix(path, ".tmpl")

		// Create the target directory if it doesn't exist
		targetPath := filepath.Join(dir, relativePath)
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		// Read the file content from
		content, err := d.templates.ReadFile(path)
		if err != nil {
			return err
		}

		// Parse the template
		tmpl, err := template.New(filepath.Base(path)).Parse(string(content))
		if err != nil {
			return err
		}

		// Create the destination file
		destFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		// Execute the template into the destination file
		if err := tmpl.Execute(destFile, d); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
