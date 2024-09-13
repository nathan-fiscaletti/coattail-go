package protocol

import (
	"context"
	"encoding/gob"

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
	mgr := GetManager(ctx)

	var resp any
	var err error

	switch h.Type {
	case ActionPacketTypePerformAndPublish:
		resp, err = mgr.LocalPeer().RunAndPublish(ctx, h.Action, h.Arg)
	case ActionPacketTypePerform:
		resp, err = mgr.LocalPeer().Run(ctx, h.Action, h.Arg)
	case ActionPacketTypePublish:
		err = mgr.LocalPeer().Publish(ctx, h.Action, h.Arg)
	}
	if err != nil {
		return nil, err
	}

	return ActionResponsePacket{
		Action:       h.Action,
		ResponseData: resp,
	}, nil
}
