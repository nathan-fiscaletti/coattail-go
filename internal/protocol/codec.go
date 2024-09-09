package protocol

import (
	"encoding/gob"
	"io"
	"sync/atomic"

	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

type EncodedPacket struct {
	ID           uint64
	RespondingTo uint64
	Data         any
}

type StreamCodec struct {
	id      *atomicId
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

func (e StreamCodec) Write(callerId uint64, p coattailtypes.Packet) (uint64, error) {
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

type atomicId uint64

func newAtomicId(id *uint64) *atomicId {
	return (*atomicId)(id)
}

func (p *atomicId) next() uint64 {
	return atomic.AddUint64((*uint64)(p), 1)
}

func (p *atomicId) current() uint64 {
	return atomic.LoadUint64((*uint64)(p))
}
