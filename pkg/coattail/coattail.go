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

var peersService *peers.Service

func getManager() *peers.Service {
	if peersService == nil {
		m, err := peers.NewService()
		if err != nil {
			panic(err)
		}
		peersService = m
	}

	return peersService
}
