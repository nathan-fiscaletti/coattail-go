package commands

import (
	"github.com/nathan-fiscaletti/coattail-go/internal/api"
	"github.com/spf13/cobra"
)

func NewNewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "new [destination]",              // Sub-command
		Short: "Create a new coattail instance", // Short description
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			api.CreateCoattailInstance(args[0])
		},
	}
}
