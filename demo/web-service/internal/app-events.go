package internal

import (
	"context"
	"log"

	"github.com/nathan-fiscaletti/ct1/pkg/sdk"
	"github.com/nathan-fiscaletti/ct1/pkg/types"

	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

func (a *App) OnStart(ctx context.Context, local *coattailtypes.Peer) {
	authPeer, err := local.GetPeer(ctx, "192.168.100.2:5243")
	if err != nil {
		panic(err)
	}

	authSdk := sdk.NewSdk(authPeer)
	response, err := authSdk.Authenticate(ctx, types.Request{Password: "password"})
	if err != nil {
		panic(err)
	}

	log.Default().Println("Authenticated: ", response.Authenticated)
	log.Default().Println("If you see this message, the demo has completed successfully! 🎉")
}
