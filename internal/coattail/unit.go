package coattail

type Unit func(interface{}) (interface{}, error)

type UnitType int

const (
	unitTypeAction UnitType = iota
	unitTypeReceiver
)

type AnyUnit struct {
	Unit

	name     string
	unitType UnitType
}

type Action struct {
	Unit

	name string
	peer PeerAdapter
}

type Receiver struct {
	Unit

	name string
	peer PeerAdapter
}
