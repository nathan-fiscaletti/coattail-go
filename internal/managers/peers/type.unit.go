package peers

// Unit is an interface that defines a unit of work that can be executed.
type Unit interface {
	Execute(interface{}) (interface{}, error)
}

// UnitHandler is a function that defines a unit of work that can be executed.
type UnitHandler func(interface{}) (interface{}, error)

// NewUnit creates a new Unit from a UnitHandler.
func NewUnit(f UnitHandler) Unit {
	return unitFunc{
		UnitHandler: f,
	}
}

type unitFunc struct {
	UnitHandler
}

func (u unitFunc) Execute(args interface{}) (interface{}, error) {
	return u.UnitHandler(args)
}

type unitType int

const (
	unitTypeAction unitType = iota
	unitTypeReceiver
)

type anyUnit struct {
	Unit

	name     string
	unitType unitType
}
