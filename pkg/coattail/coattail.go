// Coattail is a package.
package coattail

import (
	"context"
	"fmt"

	"github.com/nathan-fiscaletti/coattail-go/internal/host"
	"github.com/nathan-fiscaletti/coattail-go/internal/protocol"
	"github.com/nathan-fiscaletti/coattail-go/internal/protocol/protocoltypes"
)

// Init initializes the local peer and context for the current process.
func Init() (context.Context, error) {
	return protocol.ContextWithManager(context.Background())
}

// Run starts the local peer. This function will block.
func Run(ctx context.Context) error {
	return host.Run(host.HostConfig{
		Context: ctx,
		Port:    5244,
	})
}

// Manage returns the local peer.
func Manage(ctx context.Context) (*protocoltypes.Peer, error) {
	mgr := protocol.GetManager(ctx)
	if mgr == nil {
		return nil, fmt.Errorf("no peer manager found in context")
	}

	return mgr.LocalPeer(), nil
}
