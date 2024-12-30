package commands

import (
	"github.com/nathan-fiscaletti/coattail-go/internal/api"
	"github.com/spf13/cobra"
)

func NewNewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "new <destination> <package-name>", // Sub-command
		Short: "Create a new coattail instance",   // Short description
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			api.CreateCoattailInstance(args[0], args[1])
		},
	}
}
