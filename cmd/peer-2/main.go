package main

import (
	"time"

	"github.com/nathan-fiscaletti/coattail-go/pkg/coattail"
)

func main() {
	ct := coattail.Local()

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
	time.Sleep(20 * time.Minute)
}
