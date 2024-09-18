package internal

import "github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"

type TestAction struct{}

func NewTestAction() coattailtypes.Unit {
	return coattailtypes.NewAction[any, string](&TestAction{})
}

func (t *TestAction) Execute(_ *any) (string, error) {
	return "Hello, World!", nil
}
