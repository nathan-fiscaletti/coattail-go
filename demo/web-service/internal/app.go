package internal

import (
	"context"
	"fmt"

	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

type CT1 struct{}

func (c *CT1) OnStart(ctx context.Context, local *coattailtypes.Peer) {
	// Register actions from the action registry.
	err := RegisterUnits(ctx, local)
	if err != nil {
		panic(err)
	}

	p, err := local.GetPeer(ctx, "192.168.100.2:5243")
	if err != nil {
		panic(err)
	}

	type Response struct {
		Authenticated bool `json:"authenticated"`
	}

	response, err := p.Run(ctx, "Authenticate", struct {
		Password string `json:"password"`
	}{Password: "password"})
	if err != nil {
		panic(err)
	}

	if res, ok := response.(Response); ok {
		fmt.Println("Authenticated: ", res.Authenticated)
	} else {
		panic("unexpected response type")
	}
}
