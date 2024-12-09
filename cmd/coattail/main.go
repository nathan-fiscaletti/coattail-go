package main

import (
	"github.com/nathan-fiscaletti/coattail-go/internal/generator/templates"
)

func main() {
	new("c:/git-repos/coattail-go/.test")
}

func new(path string) {
	// copy all contents of template directory to specified path
	tmpl := templates.NewAppTemplate(templates.AppTemplateData{
		PackageName: "github.com/nathan-fiscaletti/coattail-go/coattail_app",
	})

	// TODO: Fix this after fixing the generators.
	if err := tmpl.Fill(path); err != nil {
		panic(err)
	}
}
