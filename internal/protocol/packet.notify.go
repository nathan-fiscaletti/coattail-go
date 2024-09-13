package protocol

import (
	"context"
	"encoding/gob"

	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

func init() {
	gob.Register(NotifyPacket{})
}

type NotifyPacket struct {
	Receiver string `json:"receiver"`
	Data     interface{}
}

func (n NotifyPacket) Handle(ctx context.Context) (coattailtypes.Packet, error) {
	mgr := GetManager(ctx)
	return nil, mgr.LocalPeer().
		Notify(ctx, n.Receiver, n.Data)
}
