package commands

import (
	"github.com/nathan-fiscaletti/coattail-go/internal/api"
	"github.com/spf13/cobra"
)

func NewGenerateCmd() *cobra.Command {
	// Add the generate command
	return &cobra.Command{
		Use:   "generate",
		Short: "Generates actions and receivers from the actions.yaml and receivers.yaml files",
		Run: func(cmd *cobra.Command, args []string) {
			api.RunGeneration(".")
		},
	}
}
