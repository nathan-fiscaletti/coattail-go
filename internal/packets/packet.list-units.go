package packets

import (
	"context"
	"encoding/gob"
	"fmt"

	"github.com/nathan-fiscaletti/coattail-go/internal/host"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

func init() {
	gob.Register(ListUnitsPacket{})
}

type ListUnitsPacket struct {
	Type coattailtypes.UnitType
}

func (h ListUnitsPacket) Handle(ctx context.Context) (coattailtypes.Packet, error) {
	ctHost, err := host.GetHost(ctx)
	if err != nil {
		return nil, err
	}

	var values []string

	switch h.Type {
	case coattailtypes.UnitTypeAction:
		values, _ = ctHost.LocalPeer.Actions(ctx)
	case coattailtypes.UnitTypeReceiver:
		values, _ = ctHost.LocalPeer.Receivers(ctx)
	default:
		return nil, fmt.Errorf("invalid unit type")
	}

	return ListUnitsResponsePacket{
		Values: values,
	}, nil
}
