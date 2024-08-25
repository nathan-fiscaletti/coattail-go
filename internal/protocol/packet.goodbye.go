package protocol

import (
	"encoding/gob"
	"fmt"
)

const (
	PacketTypeGoodbye OperationType = 2
)

func init() {
	gob.Register(GoodbyePacketData{})
}

type GoodbyePacketData struct {
	Message string
}

func (g GoodbyePacketData) Type() OperationType {
	return PacketTypeGoodbye
}

func (g GoodbyePacketData) Execute(communicator *Communicator) error {
	fmt.Printf("%s\n", g.Message)
	return nil
}
