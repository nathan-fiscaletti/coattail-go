package packets

import (
	"context"
	"encoding/gob"

	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

func init() {
	gob.Register(AuthenticationInvalidPacket{})
}

type AuthenticationInvalidPacket struct {
	Error string `json:"error"`
}

func (h AuthenticationInvalidPacket) Handle(ctx context.Context) (coattailtypes.Packet, error) {
	return nil, nil
}
