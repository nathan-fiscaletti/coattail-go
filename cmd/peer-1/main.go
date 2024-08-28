package main

import (
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattail"
)

func main() {
	ctx, err := coattail.Init()
	if err != nil {
		panic(err)
	}

	mgr, err := coattail.Manage(ctx)
	if err != nil {
		panic(err)
	}

	err = mgr.LocalPeer().AddAction("test", coattail.NewUnit(func(arg any) (any, error) {
		return "Hello, World!", nil
	}))
	if err != nil {
		panic(err)
	}

	err = coattail.RunInstance(ctx)
	if err != nil {
		panic(err)
	}
}
