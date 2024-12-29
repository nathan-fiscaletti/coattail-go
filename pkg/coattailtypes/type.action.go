package coattailtypes

import "context"

type Action[
	A any,
	R any,
] interface {
	Execute(context.Context, *A) (R, error)
}

type actionUnit[
	A any,
	R any,
] struct {
	action Action[A, R]
}

func (a *actionUnit[A, R]) Execute(ctx context.Context, args any) (any, error) {
	var argument *A

	if argsAny, ok := args.(A); ok {
		argument = &argsAny
	}

	return a.action.Execute(ctx, argument)
}

func NewAction[
	A any,
	R any,
](action Action[A, R]) Unit {
	return &actionUnit[A, R]{
		action: action,
	}
}
