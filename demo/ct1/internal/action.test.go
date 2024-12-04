package internal

import (
	"encoding/gob"

	"github.com/nathan-fiscaletti/coattail-go/demo/ct1/pkg/ct1types"
)

type TestAction struct{}

func init() {
	gob.Register(ct1types.Message{})
}

func (t *TestAction) Execute(_ *any) (ct1types.Message, error) {
	return ct1types.Message{
		Message: "Hello, world!",
	}, nil
}
