package main

import (
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattail"
)

func main() {
	ct := coattail.Manage()

	peer, err := ct.GetPeer("peer-1")
	if err != nil {
		panic(err)
	}

	// We're currently using a temporary test function just
	// to test the network connection between the two peers.
	err = peer.RunCommunicationTest()
	if err != nil {
		panic(err)
	}

	// Keep the program running
	select {}
}
