package protocol

import (
	"context"
	"encoding/gob"

	"github.com/nathan-fiscaletti/coattail-go/internal/protocol/protocoltypes"
)

func init() {
	gob.Register(PerformActionResponsePacket{})
}

type PerformActionResponsePacket struct {
	Action string `json:"action"`
	Data   any    `json:"data"`
}

func (h PerformActionResponsePacket) Handle(ctx context.Context) (protocoltypes.Packet, error) {
	return nil, nil
}
