package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/nathan-fiscaletti/coattail-go/internal/generator/templates"
)

func main() {
	// get input path from arg[0]
	var outputDir string = "."
	if len(os.Args) != 2 {
		outputDir = os.Args[1]
	}

	var err error
	outputDir, err = filepath.Abs(outputDir)
	if err != nil {
		panic(err)
	}

	fmt.Println(outputDir)

	// outputDir := filepath.Join("C:\\", "git-repos", "coattail-go", ".test")

	if info, ok := debug.ReadBuildInfo(); ok {
		version := info.Main.Version

		modTemplate := templates.ModTemplateData{
			PackageName:     "coattail_app",
			CoattailVersion: version,
			GoVersion:       strings.TrimPrefix(runtime.Version(), "go"),
		}

		if err := templates.NewModTemplate(modTemplate).Fill(outputDir); err != nil {
			panic(err)
		}

		if version == "(devel)" {
			// delete existing go.mod file in outputDir
			err := os.Remove(filepath.Join(outputDir, "go.mod"))
			if err != nil {
				if !os.IsNotExist(err) {
					panic(err)
				}
			}

			// write a new go.mod file
			err = templates.WriteDevModFile(filepath.Join(outputDir, "go.mod"), modTemplate)
			if err != nil {
				panic(err)
			}
		}

		// Run go mod tidy in the output directory
		cmd := exec.Command("go", "mod", "tidy")
		cmd.Dir = outputDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}
}
