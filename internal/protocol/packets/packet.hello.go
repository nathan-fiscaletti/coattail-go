package packets

import (
	"encoding/gob"
	"fmt"
)

func init() {
	gob.Register(HelloPacket{})
}

type HelloPacket struct {
	Message string `json:"message"`
}

func (h HelloPacket) Data() map[string]interface{} {
	return map[string]interface{}{
		"message": h.Message,
	}
}

func (h HelloPacket) Handle() (Packet, error) {
	fmt.Printf("handle(Hello): %s\n", h.Message)
	return GoodbyePacket{
		Message: "Goodbye, I am the second functional packet!",
	}, nil
}
