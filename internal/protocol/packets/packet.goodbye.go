package packets

import (
	"encoding/gob"
	"fmt"
)

func init() {
	gob.Register(GoodbyePacket{})
}

type GoodbyePacket struct {
	Message string `json:"message"`
}

func (g GoodbyePacket) Handle() (Packet, error) {
	fmt.Printf("handle(Goodbye): %s\n", g.Message)
	return nil, nil
}
