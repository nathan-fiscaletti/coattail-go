package coattailtypes

type Receiver[A any] interface {
	Execute(*A) error
}

type receiverUnit[A any] struct {
	receiver Receiver[A]
}

func (a *receiverUnit[A]) Execute(args any) (any, error) {
	var arguments *A

	if argsAny, ok := args.(A); ok {
		arguments = &argsAny
	}

	err := a.receiver.Execute(arguments)
	return nil, err
}

func NewReceiver[A any](receiver Receiver[A]) Unit {
	return &receiverUnit[A]{
		receiver: receiver,
	}
}
