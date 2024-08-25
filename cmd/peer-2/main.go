package main

import (
	"time"

	"github.com/nathan-fiscaletti/coattail-go/pkg/coattail"
)

func main() {
	ct := coattail.Local()

	peer, err := ct.GetPeer("peer-1")
	if err != nil {
		panic(err)
	}

	_, err = peer.RunAction("test", "Hello, world!")
	if err != nil {
		panic(err)
	}

	time.Sleep(20 * time.Minute)

	// err := ct.AddAction("test", coattail.NewUnit(func(args interface{}) (interface{}, error) {
	// 	return args, nil
	// }))
	// if err != nil {
	// 	panic(err)
	// }

	// res, err := ct.RunAction("test", "Hello, World!")
	// if err != nil {
	// 	panic(err)
	// }

	// println(res.(string))
}
