package main

import (
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattail"
)

func main() {
	// Start the coattail application.
	err := coattail.Run()
	if err != nil {
		panic(err)
	}
}
