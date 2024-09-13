package main

import (
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattail"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailmodels"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

func main() {
	// Initialize coattail
	ctx, err := coattail.Init()
	if err != nil {
		panic(err)
	}

	// Retrieve the local peer
	local := coattail.Manage(ctx)

	if err := local.AddReceiver(ctx, "testReceiver", coattailtypes.NewUnit(func(arg any) (any, error) {
		if arg != nil {
			if argStr, ok := arg.(string); ok {
				println(argStr)
			} else {
				println("Received an invalid argument")
			}
		} else {
			println("Received no argument")
		}

		return nil, nil
	})); err != nil {
		panic(err)
	}

	if err := coattail.Run(ctx, func() {
		// Retrieve the remote peer
		remote, err := local.GetPeer(ctx, "127.0.0.1:5243")
		if err != nil {
			panic(err)
		}

		// Subscribe to the "test" action on the remote peer
		// with the "testReceiver" receiver registered.
		err = remote.Subscribe(ctx, coattailmodels.Subscription{
			Address:  local.Address,
			Action:   "test",
			Receiver: "testReceiver",
		})
		if err != nil {
			panic(err)
		}

		// Run an action on the remote peer and publish it's actions to it's subscribers.
		_, err = remote.RunAndPublish(ctx, "test", nil)
		if err != nil {
			panic(err)
		}
	}); err != nil {
		panic(err)
	}
}
