package internal

import (
	"context"

	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

type TestAction struct{}

func (t *TestAction) Execute(arg any) (any, error) {
	return "Hello, World!", nil
}

type CT1 struct{}

func (c *CT1) OnStart(ctx context.Context, local *coattailtypes.Peer) {
	// Add an action to the local peer.
	err := local.AddAction(ctx, "test", &TestAction{})
	if err != nil {
		panic(err)
	}
}
