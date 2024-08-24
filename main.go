package main

import (
	"github.com/nathan-fiscaletti/coattail-go/internal/coattail"
)

func main() {
	local := coattail.Manage()

	err := local.AddAction("test", func(args interface{}) (interface{}, error) {
		return args, nil
	})
	if err != nil {
		panic(err)
	}

	res, err := local.RunAction("test", "Hello, World!")
	if err != nil {
		panic(err)
	}

	println(res.(string))
}
