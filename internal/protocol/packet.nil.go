package protocol

import (
	"context"
	"encoding/gob"
)

func init() {
	gob.Register(NilPacket{})
}

type NilPacket struct{}

func (h NilPacket) Handle(ctx context.Context) (any, error) {
	return NilPacket{}, nil
}

func (h NilPacket) Empty() bool {
	return true
}
