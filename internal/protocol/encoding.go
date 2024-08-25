package protocol

import (
	"encoding/gob"
	"io"
)

type packetEncoder struct {
	*gob.Encoder
}

func newPacketEncoder(w io.Writer) *packetEncoder {
	return &packetEncoder{
		Encoder: gob.NewEncoder(w),
	}
}

func (e *packetEncoder) EncodePacket(p PacketData) error {
	return e.Encode(packet{
		Type: p.Type(),
		Data: p,
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

func (d *packetDecoder) NextPacket() (PacketData, error) {
	var p packet
	err := d.Decode(&p)
	if err != nil {
		return nil, err
	}

	return p.Data.(PacketData), nil
}
