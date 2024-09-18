package packets

import (
	"context"
	"encoding/gob"

	"github.com/nathan-fiscaletti/coattail-go/internal/host"
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
	ctHost, err := host.GetHost(ctx)
	if err != nil {
		return nil, err
	}

	return nil, ctHost.LocalPeer.
		Notify(ctx, n.Receiver, n.Data)
}
