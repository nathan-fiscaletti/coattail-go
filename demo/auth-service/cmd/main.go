package main

import (
	"coattail_app/internal"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattail"
)

func main() {
	if err := coattail.Run(&internal.CT1{}); err != nil {
		panic(err)
	}
}
