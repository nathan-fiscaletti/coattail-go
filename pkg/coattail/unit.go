package coattail

type Unit func(interface{}) (interface{}, error)

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
