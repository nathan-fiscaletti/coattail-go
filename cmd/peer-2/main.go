package main

import (
	"github.com/google/uuid"
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
	local, err := coattail.Manage(ctx)
	if err != nil {
		panic(err)
	}

	err = local.AddReceiver(ctx, "testReceiver", coattailtypes.NewUnit(func(arg any) (any, error) {
		return "Hello, World!", nil
	}))
	if err != nil {
		panic(err)
	}

	// Retrieve the remote peer
	remote, err := local.GetPeer(ctx, "peer-1")
	if err != nil {
		panic(err)
	}

	err = remote.Subscribe(ctx, coattailmodels.Subscription{
		SubscriberID: uuid.New(),
		Action:       "test",
		Receiver:     "testReceiver",
	})
	if err != nil {
		panic(err)
	}

	// Run an action on the remote peer
	// res, err := remote.RunAction(ctx, "test", nil)
	// if err != nil {
	// 	panic(err)
	// }

	// // Print the result
	// if resStr, ok := res.(string); ok {
	// 	println(resStr)
	// }
}
