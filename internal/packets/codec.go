package packets

import (
	"encoding/gob"
	"io"

	"github.com/nathan-fiscaletti/coattail-go/internal/util/atomicid"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

type EncodedPacket struct {
	ID           uint64
	RespondingTo uint64
	Data         any
}

type StreamCodec struct {
	id      *atomicid.AtomicId
	encoder *gob.Encoder
	decoder *gob.Decoder
}

func NewStreamCodec(rw io.ReadWriter) *StreamCodec {
	return &StreamCodec{
		id:      atomicid.New(new(uint64)),
		encoder: gob.NewEncoder(rw),
		decoder: gob.NewDecoder(rw),
	}
}

func (e StreamCodec) Read() (EncodedPacket, error) {
	var p EncodedPacket
	err := e.decoder.Decode(&p)
	if err != nil {
		return EncodedPacket{}, err
	}

	return p, nil
}

func (e StreamCodec) Write(callerId uint64, p coattailtypes.Packet) (uint64, error) {
	packetId := e.id.Next()
	err := e.encoder.Encode(EncodedPacket{
		ID:           packetId,
		RespondingTo: callerId,
		Data:         p,
	})
	if err != nil {
		return packetId, err
	}

	return packetId, nil
}
