package main

import (
	"runtime/debug"

	"github.com/nathan-fiscaletti/coattail-go/internal/generator/templates"
)

func main() {
	if info, ok := debug.ReadBuildInfo(); ok {
		new("c:/git-repos/personal/coattail-go/.test", info.Main.Version)
	}
}

// get package version

func new(path string, ctVersion string) {
	// copy all contents of template directory to specified path
	tmpl := templates.NewAppTemplate(templates.AppTemplateData{
		PackageName:     "github.com/nathan-fiscaletti/coattail-go/coattail_app",
		CoattailVersion: ctVersion,
	})

	// TODO: Fix this after fixing the generators.
	if err := tmpl.Fill(path); err != nil {
		panic(err)
	}
}
