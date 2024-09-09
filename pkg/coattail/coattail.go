// Coattail is a package.
package coattail

import (
	"context"
	"fmt"

	"github.com/nathan-fiscaletti/coattail-go/internal/database"
	"github.com/nathan-fiscaletti/coattail-go/internal/host"
	"github.com/nathan-fiscaletti/coattail-go/internal/protocol"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

// Init initializes the local peer and context for the current process.
func Init() (context.Context, error) {
	ctx := context.Background()

	// Initialize the database
	ctx, err := database.ContextWithDatabase(ctx, database.DatabaseConfig{
		Path: "./data.db",
	})
	if err != nil {
		return nil, err
	}

	ctx, err = protocol.ContextWithManager(ctx)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

// Run starts the local peer. This function will block.
func Run(ctx context.Context) error {
	return host.Run(ctx)
}

// Manage returns the local peer.
func Manage(ctx context.Context) (*coattailtypes.Peer, error) {
	mgr := protocol.GetManager(ctx)
	if mgr == nil {
		return nil, fmt.Errorf("no peer manager found in context")
	}

	return mgr.LocalPeer(), nil
}
