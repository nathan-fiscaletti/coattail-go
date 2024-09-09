package protocol

import (
	"context"
	"encoding/gob"

	"github.com/google/uuid"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailmodels"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

func init() {
	gob.Register(SubscribePacket{})
}

type SubscribePacket struct {
	SubscriberID uuid.UUID `json:"subscriber_id"`
	Action       string    `json:"action"`
	Receiver     string    `json:"receiver"`
}

func (h SubscribePacket) Handle(ctx context.Context) (coattailtypes.Packet, error) {
	mgr := GetManager(ctx)

	err := mgr.LocalPeer().Subscribe(ctx, coattailmodels.Subscription{
		SubscriberID: h.SubscriberID,
		Action:       h.Action,
		Receiver:     h.Receiver,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
