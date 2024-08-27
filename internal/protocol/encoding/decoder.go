package encoding

import (
	"encoding/gob"
	"io"
)

type packetDecoder struct {
	*gob.Decoder
}

func newPacketDecoder(r io.Reader) *packetDecoder {
	return &packetDecoder{
		Decoder: gob.NewDecoder(r),
	}
}

func (d *packetDecoder) nextPacket() (EncodedPacket, error) {
	var p EncodedPacket
	err := d.Decode(&p)
	if err != nil {
		return EncodedPacket{}, err
	}

	return p, nil
}
