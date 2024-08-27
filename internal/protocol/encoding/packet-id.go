package encoding

import "sync/atomic"

type packetIdentifier struct {
	id *uint64
}

func newPacketIdentifier(id *uint64) packetIdentifier {
	return packetIdentifier{id}
}

func (p packetIdentifier) next() uint64 {
	return atomic.AddUint64(p.id, 1)
}
