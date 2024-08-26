package protocol

import (
	"encoding/gob"
	"io"
	"sync/atomic"

	"github.com/nathan-fiscaletti/coattail-go/internal/protocol/packets"
)

var packetId uint64

func nextPacketId() uint64 {
	return atomic.AddUint64(&packetId, 1)
}

type encodedPacket struct {
	ID           uint64
	RespondingTo uint64
	Data         interface{}
}

type packetEncoder struct {
	*gob.Encoder
}

func newPacketEncoder(w io.Writer) *packetEncoder {
	return &packetEncoder{
		Encoder: gob.NewEncoder(w),
	}
}

func (e *packetEncoder) EncodePacket(callerId uint64, p packets.Packet) (uint64, error) {
	packetId := nextPacketId()
	return packetId, e.Encode(encodedPacket{
		ID:           packetId,
		RespondingTo: callerId,
		Data:         p,
	})
}

type packetDecoder struct {
	*gob.Decoder
}

func newPacketDecoder(r io.Reader) *packetDecoder {
	return &packetDecoder{
		Decoder: gob.NewDecoder(r),
	}
}

func (d *packetDecoder) NextPacket() (encodedPacket, error) {
	var p encodedPacket
	err := d.Decode(&p)
	if err != nil {
		return encodedPacket{}, err
	}

	return p, nil
}
