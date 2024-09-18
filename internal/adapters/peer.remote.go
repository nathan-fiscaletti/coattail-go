package adapters

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/nathan-fiscaletti/coattail-go/internal/packets"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailmodels"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
	"github.com/samber/lo"
)

/* ====== Type ====== */

var (
	ErrAccessDenied = errors.New("access denied")
)

type RemotePeerAdapter struct {
	details coattailtypes.PeerDetails

	handler *packets.Handler
}

func newRemotePeerAdapter(details coattailtypes.PeerDetails) *RemotePeerAdapter {
	return &RemotePeerAdapter{
		details: details,
	}
}

func (i *RemotePeerAdapter) getHandler() (*packets.Handler, error) {
	if i.handler == nil || !i.handler.IsConnected() {
		conn, err := net.Dial("tcp", i.details.Address)
		if err != nil {
			return nil, err
		}

		ctx := context.Background()
		i.handler = packets.NewHandler(ctx, conn)
		i.handler.HandlePackets()
	}

	return i.handler, nil
}

/* ====== Actions ====== */

func (i *RemotePeerAdapter) Run(ctx context.Context, name string, arg any) (any, error) {
	ph, err := i.getHandler()
	if err != nil {
		return nil, err
	}

	packet, err := ph.Request(packets.Request{
		Packet: packets.ActionPacket{
			Type:   packets.ActionPacketTypePerform,
			Action: name,
			Arg:    arg,
		},
	})
	if err != nil {
		return nil, err
	}

	respPacket, isRespPacket := packet.(packets.ActionResponsePacket)
	if !isRespPacket {
		return nil, fmt.Errorf("unexpected response packet")
	}

	return respPacket.ResponseData, nil
}

func (i *RemotePeerAdapter) Publish(ctx context.Context, name string, data any) error {
	ph, err := i.getHandler()
	if err != nil {
		return err
	}

	// Run as a request to block until the publish is complete
	_, err = ph.Request(packets.Request{
		Packet: packets.ActionPacket{
			Type:   packets.ActionPacketTypePublish,
			Action: name,
			Arg:    data,
		},
	})

	return err
}

func (i *RemotePeerAdapter) RunAndPublish(ctx context.Context, name string, arg any) (any, error) {
	ph, err := i.getHandler()
	if err != nil {
		return nil, err
	}

	packet, err := ph.Request(packets.Request{
		Packet: packets.ActionPacket{
			Type:   packets.ActionPacketTypePerformAndPublish,
			Action: name,
			Arg:    arg,
		},
	})
	if err != nil {
		return nil, err
	}

	respPacket, isRespPacket := packet.(packets.ActionResponsePacket)
	if !isRespPacket {
		return nil, fmt.Errorf("unexpected response packet")
	}

	return respPacket.ResponseData, nil
}

func (i *RemotePeerAdapter) Actions(ctx context.Context) ([]string, error) {
	ph, err := i.getHandler()
	if err != nil {
		return nil, err
	}

	packet, err := ph.Request(packets.Request{
		Packet: packets.ListUnitsPacket{
			Type: coattailtypes.UnitTypeAction,
		},
	})
	if err != nil {
		return nil, err
	}

	respPacket, isRespPacket := packet.(packets.ListUnitsResponsePacket)
	if !isRespPacket {
		return nil, fmt.Errorf("unexpected response packet")
	}

	return respPacket.Values, nil
}

func (i *RemotePeerAdapter) HasAction(ctx context.Context, name string) (bool, error) {
	actions, err := i.Actions(ctx)
	if err != nil {
		return false, err
	}

	return lo.Contains(actions, name), nil
}

func (i *RemotePeerAdapter) RegisterAction(ctx context.Context, name string, unit coattailtypes.Unit) error {
	return fmt.Errorf("cannot add action to remote peer")
}

/* ====== Receivers ====== */

func (i *RemotePeerAdapter) Receivers(ctx context.Context) ([]string, error) {
	ph, err := i.getHandler()
	if err != nil {
		return nil, err
	}

	packet, err := ph.Request(packets.Request{
		Packet: packets.ListUnitsPacket{
			Type: coattailtypes.UnitTypeReceiver,
		},
	})
	if err != nil {
		return nil, err
	}

	respPacket, isRespPacket := packet.(packets.ListUnitsResponsePacket)
	if !isRespPacket {
		return nil, fmt.Errorf("unexpected response packet")
	}

	return respPacket.Values, nil
}

func (i *RemotePeerAdapter) HasReceiver(ctx context.Context, name string) (bool, error) {
	receivers, err := i.Receivers(ctx)
	if err != nil {
		return false, err
	}

	return lo.Contains(receivers, name), nil
}

func (i *RemotePeerAdapter) RegisterReceiver(ctx context.Context, name string, unit coattailtypes.Unit) error {
	return fmt.Errorf("cannot add receiver to remote peer")
}

func (i *RemotePeerAdapter) Notify(ctx context.Context, name string, arg any) error {
	ph, err := i.getHandler()
	if err != nil {
		return err
	}

	err = ph.Send(packets.NotifyPacket{
		Receiver: name,
		Data:     arg,
	})

	return err
}

/* ====== Peers ====== */

func (i *RemotePeerAdapter) GetPeer(ctx context.Context, id string) (*coattailtypes.Peer, error) {
	return nil, ErrAccessDenied
}

func (i *RemotePeerAdapter) GetPeerBy(ctx context.Context, predicate func(coattailtypes.PeerDetails) bool) (*coattailtypes.Peer, error) {
	return nil, ErrAccessDenied
}

func (i *RemotePeerAdapter) HasPeer(ctx context.Context, id string) (bool, error) {
	return false, ErrAccessDenied
}

func (i *RemotePeerAdapter) ListPeers(ctx context.Context) ([]*coattailtypes.Peer, error) {
	return nil, ErrAccessDenied
}

func (i *RemotePeerAdapter) Subscribe(ctx context.Context, sub coattailmodels.Subscription) error {
	ph, err := i.getHandler()
	if err != nil {
		return err
	}

	// Should use Request here to block until the subscription is complete
	_, err = ph.Request(packets.Request{
		Packet: packets.SubscribePacket{
			Address:  sub.Address,
			Action:   sub.Action,
			Receiver: sub.Receiver,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
