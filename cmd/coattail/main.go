package main

import (
	"fmt"
	"os"

	"github.com/nathan-fiscaletti/coattail-go/internal/commands"
	"github.com/spf13/cobra"
)

const longDescription = `The Coattail CLI can be used to create a new Coattail instance or manage an existing one.

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

	rootCmd.AddCommand(commands.NewNewCmd())
	rootCmd.AddCommand(commands.NewGenerateCmd())
	rootCmd.AddCommand(commands.NewTokenCmd())

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
