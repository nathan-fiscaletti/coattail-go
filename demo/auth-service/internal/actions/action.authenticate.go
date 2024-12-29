package actions

import (
	"context"

	"github.com/nathan-fiscaletti/ct1/pkg/types"
)

type Authenticate struct{}

func (a *Authenticate) Execute(ctx context.Context, arg *types.Request) (types.Response, error) {
	return types.Response{
		Authenticated: arg.Password == "password",
	}, nil
}
