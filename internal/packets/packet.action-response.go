package packets

import (
	"context"
	"encoding/gob"

	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

func init() {
	gob.Register(ActionResponsePacket{})
}

type ActionResponsePacket struct {
	Action       string `json:"action"`
	ResponseData any    `json:"data"`
}

func (h ActionResponsePacket) Handle(ctx context.Context) (coattailtypes.Packet, error) {
	return nil, nil
}
