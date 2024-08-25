package protocol

import (
	"encoding/gob"
	"fmt"
)

const (
	PacketTypeHello PacketType = 1
)

func init() {
	gob.Register(HelloPacket{})
}

type HelloPacket struct {
	Message string
}

func (h HelloPacket) Type() PacketType {
	return PacketTypeHello
}

func (h HelloPacket) Execute(communicator *Communicator) error {
	fmt.Printf("%s\n", h.Message)
	return communicator.WritePacket(GoodbyePacket{
		Message: "Goodbye, I am the second functional packet!",
	})
}
