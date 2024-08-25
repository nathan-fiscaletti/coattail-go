package protocol

type PacketType int

type Packet interface {
	Execute(communicator *Communicator) error
	Type() PacketType
}
