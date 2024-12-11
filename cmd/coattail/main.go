package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nathan-fiscaletti/coattail-go/internal/generator"
	"github.com/spf13/cobra"
)

const longDescription = `The Coattail CLI can be used to create a new Coattail instance from the command line.

For more information, please visit: 

    https://github.com/nathan-fiscaletti/coattail-go`

func main() {
	var rootCmd = &cobra.Command{
		Use:   "coattail",
		Short: "Coattail CLI",
		Long:  longDescription,
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

	// Add the create-action command
	var generateActions = &cobra.Command{
		Use:   "generate-actions",
		Short: "Generates actions from the actions.yaml file",
		Run: func(cmd *cobra.Command, args []string) {
			generateActions("./internal/actions")
		},
	}

	// Attach the command to the root
	rootCmd.AddCommand(generateActions)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func generateActions(destination string) {
	logger := log.Default()

	var err error
	destination, err = filepath.Abs(destination)
	if err != nil {
		panic(err)
	}

	err = generator.GenerateActions(destination, "./actions.yaml")
	if err != nil {
		logger.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	logger.Printf("Successfully generated actions.\n")
}

func generate(destination string) {
	logger := log.Default()

	var err error
	destination, err = filepath.Abs(destination)
	if err != nil {
		panic(err)
	}

	logger.Printf("Creating new Coattail instance at: %v\n", destination)

	err = generator.GenerateNewMod(destination)
	if err != nil {
		logger.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	logger.Printf("Successfully created new Coattail instance at: %v\n", destination)
}
