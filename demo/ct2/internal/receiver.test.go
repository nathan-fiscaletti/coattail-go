package internal

import "github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"

type TestReceiver struct{}

func NewTestReceiver() coattailtypes.Unit {
	return coattailtypes.NewReceiver[string](&TestReceiver{})
}

func (t *TestReceiver) Execute(arg *string) error {
	if arg != nil {
		print(*arg)
	}
	return nil
}
