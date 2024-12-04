package internal

import (
	"context"

	"github.com/nathan-fiscaletti/coattail-go/demo/ct1/pkg/ct1types"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

type CT1 struct{}

func (c *CT1) OnStart(ctx context.Context, local *coattailtypes.Peer) {
	// Add an action to the local peer.
	err := local.RegisterAction(ctx, "test", coattailtypes.NewAction[any, ct1types.Message](&TestAction{}))
	if err != nil {
		panic(err)
	}
}
