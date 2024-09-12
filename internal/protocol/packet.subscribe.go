package protocol

import (
	"context"
	"encoding/gob"

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
	mgr := GetManager(ctx)

	err := mgr.LocalPeer().Subscribe(ctx, coattailmodels.Subscription{
		Address:  h.Address,
		Action:   h.Action,
		Receiver: h.Receiver,
	})
	if err != nil {
		return nil, err
	}

	return GoodbyePacket{}, nil
}
