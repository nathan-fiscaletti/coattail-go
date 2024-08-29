package protocol

import (
	"context"
	"encoding/gob"

	"github.com/nathan-fiscaletti/coattail-go/internal/protocol/protocoltypes"
)

func init() {
	gob.Register(NotifyReceiverPacket{})
}

type NotifyReceiverPacket struct {
	Receiver string `json:"receiver"`
	Data     interface{}
}

func (n NotifyReceiverPacket) Handle(ctx context.Context) (protocoltypes.Packet, error) {
	return nil, GetManager(ctx).LocalPeer().NotifyReceiver(n.Receiver, n.Data)
}
