package packets

import (
	"context"
	"encoding/gob"
	"fmt"
)

func init() {
	gob.Register(GoodbyePacket{})
}

type GoodbyePacket struct {
	Message string `json:"message"`
}

func (g GoodbyePacket) Handle(ctx context.Context) (Packet, error) {
	fmt.Printf("handle(Goodbye): %s\n", g.Message)
	return nil, nil
}
