package packets

import "encoding/gob"

const (
	PacketTypeHello PacketType = 1
)

func init() {
	gob.Register(HelloPacketData{})
}

type HelloPacketData struct {
	Message string
}

func NewHelloPacket(data HelloPacketData) Packet {
	return Packet{
		Type: PacketTypeHello,
		Data: data,
	}
}
