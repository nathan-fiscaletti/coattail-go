package main

import (
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattail"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
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

	// Add an action to the local peer.
	err = local.AddAction(ctx, "test", coattailtypes.NewUnit(func(arg any) (any, error) {
		return "Hello, World!", nil
	}))
	if err != nil {
		panic(err)
	}

	// Start the local peer
	err = coattail.Run(ctx)
	if err != nil {
		panic(err)
	}
}
