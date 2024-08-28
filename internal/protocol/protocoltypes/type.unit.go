package protocoltypes

// Unit is an interface that defines a unit of work that can be executed.
type Unit interface {
	Execute(any) (any, error)
}

// UnitHandler is a function that defines a unit of work that can be executed.
type UnitHandler func(any) (any, error)

// NewUnit creates a new Unit from a UnitHandler.
func NewUnit(f UnitHandler) Unit {
	return unitFunc{
		UnitHandler: f,
	}
}

type unitFunc struct {
	UnitHandler
}

func (u unitFunc) Execute(args any) (any, error) {
	return u.UnitHandler(args)
}

type UnitType int

const (
	UnitTypeAction UnitType = iota
	UnitTypeReceiver
)

type AnyUnit struct {
	Unit

	Name     string
	UnitType UnitType
}
