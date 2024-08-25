package protocol

import (
	"encoding/gob"
	"fmt"
)

const (
	PacketTypeHello OperationType = 1
)

func init() {
	gob.Register(HelloPacketData{})
}

type HelloPacketData struct {
	Message string
}

func (h HelloPacketData) Type() OperationType {
	return PacketTypeHello
}

func (h HelloPacketData) Execute(communicator *Communicator) error {
	fmt.Printf("%s\n", h.Message)
	return communicator.WritePacket(GoodbyePacketData{
		Message: "Goodbye, I am the second functional packet!",
	})
}
