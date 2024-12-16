package packets

import (
	"context"
	"encoding/gob"

	"github.com/nathan-fiscaletti/coattail-go/internal/logging"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

func init() {
	gob.Register(AuthenticationResponsePacket{})
}

type AuthenticationResponsePacket struct {
	Authenticated bool   `json:"authenticated"`
	Permitted     int32  `json:"permitted"`
	Error         string `json:"error"`
}

func (h AuthenticationResponsePacket) Handle(ctx context.Context) (coattailtypes.Packet, error) {
	if !h.Authenticated {
		if logger, _ := logging.GetLogger(ctx); logger != nil {
			logger.Printf("authentication failed: %s", h.Error)
		}
		return nil, nil
	}

	if logger, _ := logging.GetLogger(ctx); logger != nil {
		logger.Printf("permissions granted: %d", h.Permitted)
	}

	return nil, nil
}
