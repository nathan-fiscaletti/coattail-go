package coattailtypes

import "github.com/invopop/jsonschema"

type Receiver[A any] interface {
	Execute(*A) error
}

type ReceiverWithInputSchema interface {
	InputSchema() *jsonschema.Schema
}

type receiverUnit[A any] struct {
	receiver Receiver[A]
}

func (a *receiverUnit[A]) Execute(args any) (any, error) {
	var argument *A

	if argsAny, ok := args.(A); ok {
		argument = &argsAny
	}

	err := a.receiver.Execute(argument)
	return nil, err
}

func (a *receiverUnit[A]) InputSchema() *jsonschema.Schema {
	if receiver, ok := a.receiver.(ReceiverWithInputSchema); ok {
		return receiver.InputSchema()
	}

	return nil
}

func (a *receiverUnit[A]) OutputSchema() *jsonschema.Schema {
	return nil
}

func NewReceiver[A any](receiver Receiver[A]) Unit {
	return &receiverUnit[A]{
		receiver: receiver,
	}
}
