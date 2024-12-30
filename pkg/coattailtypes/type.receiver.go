package coattailtypes

import (
	"context"
	"reflect"

	"github.com/invopop/jsonschema"
)

type Receiver[A any] interface {
	Execute(context.Context, *A) error
}

type ReceiverWithInputSchema interface {
	InputSchema() *jsonschema.Schema
}

type receiverUnit[A any] struct {
	name     string
	receiver Receiver[A]
}

func (a *receiverUnit[A]) Execute(ctx context.Context, args any) (any, error) {
	var argument *A

	if argsAny, ok := args.(A); ok {
		argument = &argsAny
	}

	err := a.receiver.Execute(ctx, argument)
	return nil, err
}

func (a *receiverUnit[A]) Name() string {
	return a.name
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
	// get name of action using reflection, considering that it might be a pointer
	// if it is a pointer, we need to get the name of the underlying type
	receiverType := reflect.TypeOf(receiver)
	if receiverType.Kind() == reflect.Ptr {
		receiverType = receiverType.Elem()
	}
	name := receiverType.Name()

	return &receiverUnit[A]{
		name:     name,
		receiver: receiver,
	}
}
