package packets

type PacketType int

type Packet struct {
	Type PacketType
	Data interface{}
}
