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
	Action  string `json:"action"`
	Arg     any    `json:"arg"`
	Publish bool   `json:"publish"`
}

func (h PerformActionPacket) Handle(ctx context.Context) (coattailtypes.Packet, error) {
	mgr := GetManager(ctx)

	var res any
	var err error

	var publishFunc func(context.Context, string, any) (any, error) = mgr.LocalPeer().Run
	if h.Publish {
		publishFunc = mgr.LocalPeer().RunAndPublish
	}

	res, err = publishFunc(ctx, h.Action, h.Arg)
	if err != nil {
		return nil, err
	}

	return PerformActionResponsePacket{
		Action:               h.Action,
		ResponseData:         res,
		TriggeredPublication: h.Publish,
	}, nil
}
