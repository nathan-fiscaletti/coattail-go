package main

import (
	"github.com/nathan-fiscaletti/coattail-go/demo/ct2/internal"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattail"
)

func main() {
	if err := coattail.Run(&internal.CT2{}); err != nil {
		panic(err)
	}
}
