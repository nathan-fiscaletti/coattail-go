package commands

import (
	"github.com/nathan-fiscaletti/coattail-go/pkg/commands/token"
	"github.com/spf13/cobra"
)

func NewTokenCmd() *cobra.Command {
	tokenCmd := &cobra.Command{
		Use:   "token",
		Short: "Manage tokens",
	}

	tokenCmd.AddCommand(token.NewCreateCommand())

	return tokenCmd
}
