package main

import (
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattail"
	"github.com/nathan-fiscaletti/ct1/internal"
)

func main() {
	if err := coattail.Run(&internal.AuthService{}); err != nil {
		panic(err)
	}
}
