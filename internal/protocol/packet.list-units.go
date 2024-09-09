package protocol

import (
	"context"
	"encoding/gob"
	"fmt"

	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

func init() {
	gob.Register(ListUnitsPacket{})
}

type ListUnitsPacket struct {
	Type coattailtypes.UnitType
}

func (h ListUnitsPacket) Handle(ctx context.Context) (coattailtypes.Packet, error) {
	mgr := GetManager(ctx)

	var vals []string

	switch h.Type {
	case coattailtypes.UnitTypeAction:
		vals, _ = mgr.LocalPeer().Actions(ctx)
	case coattailtypes.UnitTypeReceiver:
		vals, _ = mgr.LocalPeer().Receivers(ctx)
	default:
		return nil, fmt.Errorf("invalid unit type")
	}

	return ListUnitsResponsePacket{
		Values: vals,
	}, nil
}
