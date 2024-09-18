package internal

import (
	"context"

	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

type CT1 struct{}

func (c *CT1) OnStart(ctx context.Context, local *coattailtypes.Peer) {
	// Add an action to the local peer.
	err := local.RegisterAction(ctx, "test", NewTestAction())
	if err != nil {
		panic(err)
	}
}
