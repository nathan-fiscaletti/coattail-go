package protocol

import (
	"context"
	"encoding/gob"
)

func init() {
	gob.Register(PerformActionResponsePacket{})
}

type PerformActionResponsePacket struct {
	Action string `json:"action"`
	Data   any    `json:"data"`
}

func (h PerformActionResponsePacket) Handle(ctx context.Context) (any, error) {
	return NilPacket{}, nil
}
