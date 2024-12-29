package coattailtypes

import "context"

// Unit is an interface that defines a unit of work that can be executed.
type Unit interface {
	Execute(context.Context, any) (any, error)
}

// UnitHandler is a function that defines a unit of work that can be executed.
type UnitHandler func(context.Context, any) (any, error)

// NewUnit creates a new Unit from a UnitHandler.
func NewUnit(f UnitHandler) Unit {
	return unitFunc{
		UnitHandler: f,
	}
}

type unitFunc struct {
	UnitHandler
}

func (u unitFunc) Execute(ctx context.Context, args any) (any, error) {
	return u.UnitHandler(ctx, args)
}

type UnitType int

const (
	UnitTypeAction UnitType = iota
	UnitTypeReceiver
)

type UnitImpl struct {
	Unit

	Name     string
	UnitType UnitType
}
