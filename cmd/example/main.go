package main

import (
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattail"
)

func main() {
	ct := coattail.Manage()

	err := ct.AddAction("test", func(args interface{}) (interface{}, error) {
		return args, nil
	})
	if err != nil {
		panic(err)
	}

	res, err := ct.RunAction("test", "Hello, World!")
	if err != nil {
		panic(err)
	}

	println(res.(string))
}
