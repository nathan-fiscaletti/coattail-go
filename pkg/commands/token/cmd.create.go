package token

import (
	"fmt"

	"github.com/nathan-fiscaletti/coattail-go/internal/api"
	"github.com/nathan-fiscaletti/coattail-go/internal/services/permission"
	"github.com/spf13/cobra"
)

func NewCreateCommand() *cobra.Command {
	var keyfile string
	var network string
	var perm int
	var expiry string

	cmd := &cobra.Command{
		Use:   "create -k <keyfile> [-n <network>] [-p <int>] [-e <expiry>]",
		Short: "Create a new token",
		Run: func(cmd *cobra.Command, args []string) {
			api.CreateToken(keyfile, network, perm, expiry)
		},
	}

	allPerms := int(permission.PermissionMask(permission.All))

	// Adding flags
	cmd.Flags().StringVarP(&keyfile, "keyfile", "k", "", "Path to the key file (required)")
	cmd.Flags().StringVarP(&network, "network", "n", "0.0.0.0/0", "Specify the authorized network with CIDR notation")
	cmd.Flags().IntVarP(&perm, "perm", "p", allPerms, fmt.Sprintf("Permission level as an integer [0-%d]", allPerms))
	cmd.Flags().StringVarP(&expiry, "expiry", "e", "", "Expiry time for the token (default 24 hours from now)")

	// Mark `keyfile` as required
	cmd.MarkFlagRequired("keyfile")

	return cmd
}
