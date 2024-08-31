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
	mgr := GetManager(ctx)
	return nil, mgr.LocalPeer().
		NotifyReceiver(n.Receiver, n.Data)
}
