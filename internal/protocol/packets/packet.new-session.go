package packets

import (
	"context"
	"encoding/gob"

	"github.com/nathan-fiscaletti/coattail-go/internal/services/authentication"
)

func init() {
	gob.Register(NewSessionPacket{})
}

type NewSessionPacket struct {
	AuthenticationToken string `json:"authentication_token"`
}

func (h NewSessionPacket) Handle(ctx context.Context) (Packet, error) {
	authService := authentication.GetService(ctx)

	res := authService.Authenticate(h.AuthenticationToken)

	return GoodbyePacket{
		Message: res,
	}, nil
}
