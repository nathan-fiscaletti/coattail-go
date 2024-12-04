package packets

import (
	"context"
	"encoding/gob"

	"github.com/nathan-fiscaletti/coattail-go/internal/host"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailmodels"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

func init() {
	gob.Register(SubscribePacket{})
}

type SubscribePacket struct {
	Address  string `json:"address"`
	Action   string `json:"action"`
	Receiver string `json:"receiver"`
}

func (h SubscribePacket) Handle(ctx context.Context) (coattailtypes.Packet, error) {
	ctHost, err := host.GetHost(ctx)
	if err != nil {
		return nil, err
	}

	err = ctHost.LocalPeer.Subscribe(ctx, coattailmodels.Subscription{
		Address:  h.Address,
		Action:   h.Action,
		Receiver: h.Receiver,
	})
	if err != nil {
		return nil, err
	}

	return EmptyPacket{}, nil
}
