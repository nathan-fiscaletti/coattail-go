package coattailtypes

import (
	"context"
	"reflect"
)

type Action[
	A any,
	R any,
] interface {
	Execute(context.Context, *A) (R, error)
}

type ActionUnit[
	A any,
	R any,
] struct {
	name   string
	action Action[A, R]
}

func (a *ActionUnit[A, R]) Execute(ctx context.Context, args any) (any, error) {
	var argument *A

	if argsAny, ok := args.(A); ok {
		argument = &argsAny
	}

	return a.action.Execute(ctx, argument)
}

func (a *ActionUnit[A, R]) Name() string {
	return a.name
}

func NewAction[
	A any,
	R any,
](action Action[A, R]) Unit {
	// get name of action using reflection, considering that it might be a pointer
	// if it is a pointer, we need to get the name of the underlying type
	actionType := reflect.TypeOf(action)
	if actionType.Kind() == reflect.Ptr {
		actionType = actionType.Elem()
	}
	name := actionType.Name()

	return &ActionUnit[A, R]{
		name:   name,
		action: action,
	}
}
