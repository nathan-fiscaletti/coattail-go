package actions

import (
	"coattail_app/pkg/types"
)

type Authenticate struct{}

func (a *Authenticate) Execute(arg *types.Request) (types.Response, error) {
	return types.Response{
		Authenticated: arg.Password == "password",
	}, nil
}
