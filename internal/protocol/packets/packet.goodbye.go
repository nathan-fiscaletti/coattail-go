package packets

import "encoding/gob"

const (
	PacketTypeGoodbye PacketType = 2
)

func init() {
	gob.Register(GoodbyePacketData{})
}

type GoodbyePacketData struct {
	Message string
}

func NewGoodbyePacket(data GoodbyePacketData) Packet {
	return Packet{
		Type: PacketTypeGoodbye,
		Data: data,
	}
}
