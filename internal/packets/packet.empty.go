package packets

import (
	"context"
	"encoding/gob"

	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

func init() {
	gob.Register(EmptyPacket{})
}

type EmptyPacket struct{}

func (h EmptyPacket) Handle(ctx context.Context) (coattailtypes.Packet, error) {
	return nil, nil
}
