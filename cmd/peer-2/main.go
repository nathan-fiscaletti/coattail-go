package main

import (
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattail"
)

func main() {
	// Initialize coattail
	ctx, err := coattail.Init()
	if err != nil {
		panic(err)
	}

	// Retrieve the local peer
	local, err := coattail.Manage(ctx)
	if err != nil {
		panic(err)
	}

	// Retrieve the remote peer
	remote, err := local.GetPeer("peer-1")
	if err != nil {
		panic(err)
	}

	// Run an action on the remote peer
	res, err := remote.RunAction("test", nil)
	if err != nil {
		panic(err)
	}

	// Print the result
	if resStr, ok := res.(string); ok {
		println(resStr)
	}
}
