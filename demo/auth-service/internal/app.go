package internal

import (
	"context"

	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

type AuthService struct{}

func (c *AuthService) OnStart(ctx context.Context, local *coattailtypes.Peer) {
	// Register actions from the action registry.
	err := RegisterUnits(ctx, local)
	if err != nil {
		panic(err)
	}
}
