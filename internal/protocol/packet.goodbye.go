package protocol

import (
	"encoding/gob"
	"fmt"
)

const (
	PacketTypeGoodbye PacketType = 2
)

func init() {
	gob.Register(GoodbyePacket{})
}

type GoodbyePacket struct {
	Message string
}

func (g GoodbyePacket) Type() PacketType {
	return PacketTypeGoodbye
}

func (g GoodbyePacket) Execute(communicator *Communicator) error {
	fmt.Printf("%s\n", g.Message)
	return nil
}
