package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/nathan-fiscaletti/coattail-go/internal/generator/templates"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "ct", // Command name
		Short: "Coattail CLI",
	}

	// Add the 'new' command
	var newCmd = &cobra.Command{
		Use:   "new [path]",                     // Sub-command
		Short: "Create a new coattail instance", // Short description
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var outputDir string = "."
			if len(os.Args) > 0 {
				outputDir = args[0]
			}
			generate(outputDir)
		},
	}

	// Attach the command to the root
	rootCmd.AddCommand(newCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func generate(outputDir string) {
	logger := log.Default()

	var err error
	outputDir, err = filepath.Abs(outputDir)
	if err != nil {
		panic(err)
	}

	logger.Printf("Creating new Coattail instance at: %v\n", outputDir)

	info, ok := debug.ReadBuildInfo()

	if !ok {
		logger.Println("Failed to read build info")
		os.Exit(1)
	}

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
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
