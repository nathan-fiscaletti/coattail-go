package protocol

type OperationType int

type PacketData interface {
	Execute(communicator *Communicator) error
	Type() OperationType
}

type packet struct {
	Type OperationType
	Data interface{}
}
