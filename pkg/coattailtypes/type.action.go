package coattailtypes

type Action[
	A any,
	R any,
] interface {
	Execute(args *A) (R, error)
}

type actionUnit[
	A any,
	R any,
] struct {
	action Action[A, R]
}

func (a *actionUnit[A, R]) Execute(args any) (any, error) {
	var argument *A

	if argsAny, ok := args.(A); ok {
		argument = &argsAny
	}

	return a.action.Execute(argument)
}

func NewAction[
	A any,
	R any,
](action Action[A, R]) Unit {
	return &actionUnit[A, R]{
		action: action,
	}
}
