package protocol

import (
	"context"
	"encoding/gob"

	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

func init() {
	gob.Register(PerformActionResponsePacket{})
}

type PerformActionResponsePacket struct {
	Action               string `json:"action"`
	ResponseData         any    `json:"data"`
	TriggeredPublication bool   `json:"published"`
}

func (h PerformActionResponsePacket) Handle(ctx context.Context) (coattailtypes.Packet, error) {
	return nil, nil
}
