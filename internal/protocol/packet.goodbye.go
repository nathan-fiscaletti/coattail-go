package protocol

import (
	"context"
	"encoding/gob"
	"fmt"

	"github.com/nathan-fiscaletti/coattail-go/internal/protocol/protocoltypes"
)

func init() {
	gob.Register(GoodbyePacket{})
}

type GoodbyePacket struct {
	Message string `json:"message"`
}

func (g GoodbyePacket) Handle(ctx context.Context) (protocoltypes.Packet, error) {
	fmt.Printf("handle(Goodbye): %s\n", g.Message)
	return nil, nil
}
