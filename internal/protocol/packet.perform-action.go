package protocol

import (
	"context"
	"encoding/gob"

	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

func init() {
	gob.Register(PerformActionPacket{})
}

type PerformActionPacket struct {
	Action string `json:"action"`
	Arg    any    `json:"arg"`
}

func (h PerformActionPacket) Handle(ctx context.Context) (coattailtypes.Packet, error) {
	mgr := GetManager(ctx)

	res, err := mgr.LocalPeer().RunAction(ctx, h.Action, h.Arg)
	if err != nil {
		return nil, err
	}

	return PerformActionResponsePacket{
		Action:       h.Action,
		ResponseData: res,
	}, nil
}
