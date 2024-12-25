package actions

import (
	"github.com/nathan-fiscaletti/ct1/pkg/types"
)

type Authenticate struct{}

func (a *Authenticate) Execute(arg *types.Request) (types.Response, error) {
	return types.Response{
		Authenticated: arg.Password == "password",
	}, nil
}
