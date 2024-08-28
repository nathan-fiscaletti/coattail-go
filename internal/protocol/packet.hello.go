package protocol

import (
	"context"
	"encoding/gob"
	"fmt"
)

func init() {
	gob.Register(HelloPacket{})
}

type HelloPacket struct {
	Message string `json:"message"`
}

func (h HelloPacket) Handle(ctx context.Context) (any, error) {
	fmt.Printf("handle(Hello): %s\n", h.Message)
	return GoodbyePacket{
		Message: "Goodbye, I am the second functional packet!",
	}, nil
}
