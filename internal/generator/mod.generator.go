package generator

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

type ModGenerator struct{}

func GenerateNewMod(destination string) error {
	info, ok := debug.ReadBuildInfo()

	if !ok {
		return fmt.Errorf("failed to read build info")
	}

	version := info.Main.Version
	modTemplate := templates.ModTemplateData{
		PackageName:     "coattail_app",
		CoattailVersion: version,
		GoVersion:       strings.TrimPrefix(runtime.Version(), "go"),
	}

	if err := templates.NewModTemplate(modTemplate).Fill(destination); err != nil {
		return fmt.Errorf("failed to generate mod file: %w", err)
	}

	if version == "(devel)" {
		// delete existing go.mod file in outputDir
		err := os.Remove(filepath.Join(destination, "go.mod"))
		if err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("failed to delete existing go.mod file: %w", err)
			}
		}

		// write a new go.mod file
		err = templates.WriteDevModFile(filepath.Join(destination, "go.mod"), modTemplate)
		if err != nil {
			return fmt.Errorf("failed to write go.mod file: %w", err)
		}
	}

	// Run go mod tidy in the output directory
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = destination
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run go mod tidy: %w", err)
	}

	// Run coattail generate in the output directory
	cmd = exec.Command("coattail", "generate")
	cmd.Dir = destination
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run coattail generate: %w", err)
	}

	return nil
}
