package coattailtypes

type Action[A, R any] interface {
	Execute(*A) (R, error)
}

type actionUnit[A, R any] struct {
	action Action[A, R]
}

func (a *actionUnit[A, R]) Execute(args any) (any, error) {
	var arguments *A

	if argsAny, ok := args.(A); ok {
		arguments = &argsAny
	}

	return a.action.Execute(arguments)
}

func NewAction[A, R any](action Action[A, R]) Unit {
	return &actionUnit[A, R]{
		action: action,
	}
}
