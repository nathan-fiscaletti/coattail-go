package main

import (
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattail"
	"github.com/nathan-fiscaletti/ct2/internal"
)

func main() {
	if err := coattail.Run(&internal.WebService{}); err != nil {
		panic(err)
	}
}
