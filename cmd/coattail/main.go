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
		Use:   "coattail", // Command name
		Short: "Coattail CLI",
		Long: `The Coattail CLI can be used to create a new Coattail instance from the command line.

For more information, please visit: 

    https://github.com/nathan-fiscaletti/coattail-go`,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	}

	// Add the 'new' command
	var newCmd = &cobra.Command{
		Use:   "new [destination]",              // Sub-command
		Short: "Create a new coattail instance", // Short description
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var destination string = "."
			if len(os.Args) > 0 {
				destination = args[0]
			}
			generate(destination)
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

func generate(destination string) {
	logger := log.Default()

	var err error
	destination, err = filepath.Abs(destination)
	if err != nil {
		panic(err)
	}

	logger.Printf("Creating new Coattail instance at: %v\n", destination)

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

	if err := templates.NewModTemplate(modTemplate).Fill(destination); err != nil {
		panic(err)
	}

	if version == "(devel)" {
		// delete existing go.mod file in outputDir
		err := os.Remove(filepath.Join(destination, "go.mod"))
		if err != nil {
			if !os.IsNotExist(err) {
				panic(err)
			}
		}

		// write a new go.mod file
		err = templates.WriteDevModFile(filepath.Join(destination, "go.mod"), modTemplate)
		if err != nil {
			panic(err)
		}
	}

	// Run go mod tidy in the output directory
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = destination
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
