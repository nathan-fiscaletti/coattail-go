package packets

import (
	"context"
	"encoding/gob"
	"errors"
	"net"

	"github.com/nathan-fiscaletti/coattail-go/internal/keys"
	"github.com/nathan-fiscaletti/coattail-go/internal/services/authentication"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

var (
	ErrConnectionNotFound = errors.New("connection not found in context")
)

func init() {
	gob.Register(AuthenticationPacket{})
}

type AuthenticationPacket struct {
	Token string `json:"token"`
}

func (h AuthenticationPacket) Handle(ctx context.Context) (coattailtypes.Packet, error) {
	conn, ok := ctx.Value(keys.ConnectionKey).(net.Conn)
	if !ok {
		return nil, ErrConnectionNotFound
	}

	auth, err := authentication.GetService(ctx)
	if err != nil {
		return nil, err
	}

	host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		return nil, err
	}

	var response AuthenticationResponsePacket

	result, err := auth.Authenticate(ctx, h.Token, net.ParseIP(host))
	response.Authenticated = result.Authenticated
	response.Permitted = result.Token.Permitted
	if err != nil {
		response.Error = err.Error()
	}

	return response, nil
}
