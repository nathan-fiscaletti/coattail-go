package internal

import (
	"coattail_app/pkg/types"
	"context"
	"log"

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

	log.Default().Println("Running action Authenticate")
	response, err := p.Run(ctx, "Authenticate", types.Request{Password: "password"})
	if err != nil {
		panic(err)
	}
	log.Default().Println("Authenticate responded")

	if res, ok := response.(types.Response); ok {
		log.Default().Println("Authenticated: ", res.Authenticated)
	} else {
		panic("unexpected response type")
	}
}
