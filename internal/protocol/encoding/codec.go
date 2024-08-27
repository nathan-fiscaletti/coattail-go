package encoding

import (
	"io"

	"github.com/nathan-fiscaletti/coattail-go/internal/protocol/packets"
)

type EncodedPacket struct {
	ID           uint64
	RespondingTo uint64
	Data         interface{}
}

type Codec struct {
	packetIdentifier

	encoder *packetEncoder
	decoder *packetDecoder
}

func NewCodec(rw io.ReadWriter) *Codec {
	res := Codec{
		packetIdentifier: newPacketIdentifier(new(uint64)),
	}

	res.encoder = newPacketEncoder(res.packetIdentifier, rw)
	res.decoder = newPacketDecoder(rw)

	return &res
}

func (e Codec) Read() (EncodedPacket, error) {
	return e.decoder.nextPacket()
}

func (e Codec) Write(callerId uint64, p packets.Packet) (uint64, error) {
	return e.encoder.encodePacket(callerId, p)
}
