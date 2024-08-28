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

	remote, err := mgr.LocalPeer().GetPeer("peer-1")
	if err != nil {
		panic(err)
	}

	res, err := remote.RunAction("test", nil)
	if err != nil {
		panic(err)
	}

	if resStr, ok := res.(string); ok {
		println(resStr)
	}
}
