package packets

import (
	"context"
	"encoding/gob"

	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

func init() {
	gob.Register(ListUnitsResponsePacket{})
}

type ListUnitsResponsePacket struct {
	Values []string
}

func (h ListUnitsResponsePacket) Handle(ctx context.Context) (coattailtypes.Packet, error) {
	return nil, nil
}
