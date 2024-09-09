package protocol

import (
	"context"
	"encoding/gob"
	"fmt"

	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

func init() {
	gob.Register(HelloPacket{})
}

type HelloPacket struct {
	Message string `json:"message"`
}

func (h HelloPacket) Handle(ctx context.Context) (coattailtypes.Packet, error) {
	fmt.Printf("handle(Hello): %s\n", h.Message)
	return GoodbyePacket{
		Message: "Goodbye, I am the second functional packet!",
	}, nil
}
