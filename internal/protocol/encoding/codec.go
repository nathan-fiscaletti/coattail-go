package encoding

import (
	"encoding/gob"
	"io"
	"sync/atomic"

	"github.com/nathan-fiscaletti/coattail-go/internal/protocol/packets"
)

type EncodedPacket struct {
	ID           uint64
	RespondingTo uint64
	Data         interface{}
}

type StreamCodec struct {
	id      atomicId
	encoder *gob.Encoder
	decoder *gob.Decoder
}

func NewStreamCodec(rw io.ReadWriter) *StreamCodec {
	return &StreamCodec{
		id:      newAtomicId(new(uint64)),
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

func (e StreamCodec) Write(callerId uint64, p packets.Packet) (uint64, error) {
	packetId := e.id.next()
	err := e.encoder.Encode(EncodedPacket{
		ID:           packetId,
		RespondingTo: callerId,
		Data:         p,
	})
	if err != nil {
		return 0, err
	}

	return packetId, nil
}

type atomicId struct {
	id *uint64
}

func newAtomicId(id *uint64) atomicId {
	return atomicId{id}
}

func (p atomicId) next() uint64 {
	return atomic.AddUint64(p.id, 1)
}
