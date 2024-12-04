package packets

import (
	"context"
	"encoding/gob"

	"github.com/nathan-fiscaletti/coattail-go/internal/host"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

func init() {
	gob.Register(ActionPacket{})
}

type ActionPacketType int

const (
	ActionPacketTypePerformAndPublish ActionPacketType = iota
	ActionPacketTypePerform
	ActionPacketTypePublish
)

type ActionPacket struct {
	Action string           `json:"action"`
	Arg    any              `json:"arg"`
	Type   ActionPacketType `json:"type"`
}

func (h ActionPacket) Handle(ctx context.Context) (coattailtypes.Packet, error) {
	ctHost, err := host.GetHost(ctx)
	if err != nil {
		return nil, err
	}

	var resp any

	switch h.Type {
	case ActionPacketTypePerformAndPublish:
		err = ctHost.LocalPeer.RunAndPublish(ctx, h.Action, h.Arg)
	case ActionPacketTypePerform:
		resp, err = ctHost.LocalPeer.Run(ctx, h.Action, h.Arg)
	case ActionPacketTypePublish:
		err = ctHost.LocalPeer.Publish(ctx, h.Action, h.Arg)
	}
	if err != nil {
		return nil, err
	}

	if h.Type == ActionPacketTypePerform {
		return ActionResponsePacket{
			Action:       h.Action,
			ResponseData: resp,
		}, nil
	}

	return nil, nil
}
