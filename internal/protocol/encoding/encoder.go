package encoding

import (
	"encoding/gob"
	"io"

	"github.com/nathan-fiscaletti/coattail-go/internal/protocol/packets"
)

type packetEncoder struct {
	*gob.Encoder
	packetIdentifier
}

func newPacketEncoder(pi packetIdentifier, w io.Writer) *packetEncoder {
	return &packetEncoder{
		Encoder:          gob.NewEncoder(w),
		packetIdentifier: pi,
	}
}

func (e *packetEncoder) encodePacket(callerId uint64, p packets.Packet) (uint64, error) {
	packetId := e.packetIdentifier.next()
	return packetId, e.Encode(EncodedPacket{
		ID:           packetId,
		RespondingTo: callerId,
		Data:         p,
	})
}
