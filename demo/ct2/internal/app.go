package internal

import (
	"context"

	"github.com/nathan-fiscaletti/coattail-go/demo/ct1/pkg/ct1types"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailmodels"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

type CT2 struct{}

func (c *CT2) OnStart(ctx context.Context, local *coattailtypes.Peer) {
	if err := local.RegisterReceiver(ctx, "testReceiver", coattailtypes.NewReceiver[ct1types.Message](&TestReceiver{})); err != nil {
		panic(err)
	}

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
	err = remote.RunAndPublish(ctx, "test", nil)
	if err != nil {
		panic(err)
	}
}
