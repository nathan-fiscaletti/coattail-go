// Coattail is a package.
package coattail

import (
	"github.com/nathan-fiscaletti/coattail-go/internal/host"
	"github.com/nathan-fiscaletti/coattail-go/internal/managers/peers"
)

// Manage returns the local peer.
func Manage() *Peer {
	return (*Peer)(getManager().LocalPeer())
}

// RunInstance starts the local peer. This function will block.
func RunInstance() error {
	return host.Run(host.HostConfig{
		Port: 5244,
	})
}

var peersManager *peers.Manager

func getManager() *peers.Manager {
	if peersManager == nil {
		m, err := peers.NewManager()
		if err != nil {
			panic(err)
		}
		peersManager = m
	}

	return peersManager
}
