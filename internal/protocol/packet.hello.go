package protocol

import (
	"context"
	"encoding/gob"
	"fmt"

	"github.com/nathan-fiscaletti/coattail-go/internal/protocol/protocoltypes"
)

func init() {
	gob.Register(HelloPacket{})
}

type HelloPacket struct {
	Message string `json:"message"`
}

func (h HelloPacket) Handle(ctx context.Context) (protocoltypes.Packet, error) {
	fmt.Printf("handle(Hello): %s\n", h.Message)
	return GoodbyePacket{
		Message: "Goodbye, I am the second functional packet!",
	}, nil
}
