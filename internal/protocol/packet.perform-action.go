package protocol

import (
	"context"
	"encoding/gob"
)

func init() {
	gob.Register(PerformActionPacket{})
}

type PerformActionPacket struct {
	Action string `json:"action"`
	Arg    any    `json:"arg"`
}

func (h PerformActionPacket) Handle(ctx context.Context) (any, error) {
	mgr := GetManager(ctx)

	res, err := mgr.LocalPeer().RunAction(h.Action, h.Arg)

	return PerformActionResponsePacket{
		Action: h.Action,
		Data:   res,
	}, err
}
