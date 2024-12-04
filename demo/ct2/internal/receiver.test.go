package internal

import (
	"github.com/nathan-fiscaletti/coattail-go/demo/ct1/pkg/ct1types"
)

type TestReceiver struct{}

func (t *TestReceiver) Execute(arg *ct1types.Message) error {
	if arg != nil {
		println(arg.Message)
	}
	return nil
}
