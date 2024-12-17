package internal

import (
	"context"

	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

type CT1 struct{}

func (c *CT1) OnStart(ctx context.Context, local *coattailtypes.Peer) {
	// Register actions from the action registry.
	err := RegisterUnits(ctx, local)
	if err != nil {
		panic(err)
	}
}
