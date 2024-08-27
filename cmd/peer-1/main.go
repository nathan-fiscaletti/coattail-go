package main

import (
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattail"
)

func main() {
	// Start the coattail instance..
	err := coattail.RunInstance()
	if err != nil {
		panic(err)
	}
}
