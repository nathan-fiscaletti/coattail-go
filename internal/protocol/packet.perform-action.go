package protocol

import (
	"context"
	"encoding/gob"

	"github.com/nathan-fiscaletti/coattail-go/internal/protocol/protocoltypes"
)

func init() {
	gob.Register(PerformActionPacket{})
}

type PerformActionPacket struct {
	Action  string `json:"action"`
	Arg     any    `json:"arg"`
	Publish bool   `json:"publish"`
}

func (h PerformActionPacket) Handle(ctx context.Context) (protocoltypes.Packet, error) {
	mgr := GetManager(ctx)

	res, err := mgr.LocalPeer().RunAction(protocoltypes.RunActionArguments{
		Name:    h.Action,
		Arg:     h.Arg,
		Publish: h.Publish,
	})
	if err != nil {
		return nil, err
	}

	var published bool
	var publishedError error

	if h.Publish {
		published = true
		publishedError = mgr.LocalPeer().PublishActionResult(h.Action, res)
		if publishedError != nil {
			published = false
		}
	}

	return PerformActionResponsePacket{
		Action:         h.Action,
		Data:           res,
		Published:      published,
		PublishedError: publishedError,
	}, nil
}
