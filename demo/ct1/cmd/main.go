package main

import (
	"github.com/nathan-fiscaletti/coattail-go/demo/ct1/internal"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattail"
)

func main() {
	if err := coattail.Run(&internal.CT1{}); err != nil {
		panic(err)
	}
}
