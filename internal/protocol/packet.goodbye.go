package protocol

import (
	"context"
	"encoding/gob"

	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

func init() {
	gob.Register(GoodbyePacket{})
}

type GoodbyePacket struct {
	Message string `json:"message"`
}

func (g GoodbyePacket) Handle(ctx context.Context) (coattailtypes.Packet, error) {
	return nil, nil
}
