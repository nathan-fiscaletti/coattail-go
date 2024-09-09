package protocol

import (
	"context"
	"fmt"
	"net"

	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailmodels"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
	"github.com/samber/lo"
)

/* ====== Type ====== */

type RemotePeerAdapter struct {
	details coattailtypes.PeerDetails

	packetHandler *PacketHandler
}

func newRemotePeerAdapter(details coattailtypes.PeerDetails) *RemotePeerAdapter {
	return &RemotePeerAdapter{
		details: details,
	}
}

func (i *RemotePeerAdapter) getPacketHandler() (*PacketHandler, error) {
	if i.packetHandler == nil || !i.packetHandler.IsConnected() {
		conn, err := net.Dial("tcp", i.details.Address)
		if err != nil {
			return nil, err
		}

		ctx := context.Background()
		i.packetHandler = NewPacketHandler(ctx, conn)
		i.packetHandler.HandlePackets()
	}

	return i.packetHandler, nil
}

/* ====== Actions ====== */

func (i *RemotePeerAdapter) RunAction(ctx context.Context, name string, arg any) (any, error) {
	ph, err := i.getPacketHandler()
	if err != nil {
		return nil, err
	}

	packet, err := ph.Request(Request{
		Packet: PerformActionPacket{
			Action: name,
			Arg:    arg,
		},
	})
	if err != nil {
		return nil, err
	}

	respPacket, isRespPacket := packet.(PerformActionResponsePacket)
	if !isRespPacket {
		return nil, fmt.Errorf("unexpected response packet")
	}

	return respPacket.ResponseData, nil
}

func (i *RemotePeerAdapter) Publish(ctx context.Context, name string, data any) error {
	return nil
}

func (i *RemotePeerAdapter) Actions(ctx context.Context) ([]string, error) {
	ph, err := i.getPacketHandler()
	if err != nil {
		return nil, err
	}

	packet, err := ph.Request(Request{
		Packet: ListUnitsPacket{
			Type: coattailtypes.UnitTypeAction,
		},
	})
	if err != nil {
		return nil, err
	}

	respPacket, isRespPacket := packet.(ListUnitsResponsePacket)
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

func (i *RemotePeerAdapter) AddAction(ctx context.Context, name string, unit coattailtypes.Unit) error {
	return fmt.Errorf("cannot add action to remote peer")
}

/* ====== Receivers ====== */

func (i *RemotePeerAdapter) Receivers(ctx context.Context) ([]string, error) {
	ph, err := i.getPacketHandler()
	if err != nil {
		return nil, err
	}

	packet, err := ph.Request(Request{
		Packet: ListUnitsPacket{
			Type: coattailtypes.UnitTypeReceiver,
		},
	})
	if err != nil {
		return nil, err
	}

	respPacket, isRespPacket := packet.(ListUnitsResponsePacket)
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

func (i *RemotePeerAdapter) AddReceiver(ctx context.Context, name string, unit coattailtypes.Unit) error {
	return fmt.Errorf("cannot add receiver to remote peer")
}

func (i *RemotePeerAdapter) NotifyReceiver(ctx context.Context, name string, arg any) error {
	ph, err := i.getPacketHandler()
	if err != nil {
		return err
	}

	err = ph.Send(NotifyReceiverPacket{
		Receiver: name,
		Data:     arg,
	})

	return err
}

/* ====== Peers ====== */

func (i *RemotePeerAdapter) GetPeer(ctx context.Context, id string) (*coattailtypes.Peer, error) {
	return nil, nil
}

func (i *RemotePeerAdapter) HasPeer(ctx context.Context, id string) (bool, error) {
	return false, nil
}

func (i *RemotePeerAdapter) ListPeers(ctx context.Context) ([]*coattailtypes.Peer, error) {
	return []*coattailtypes.Peer{}, nil
}

func (i *RemotePeerAdapter) Subscribe(ctx context.Context, sub coattailmodels.Subscription) error {
	ph, err := i.getPacketHandler()
	if err != nil {
		return err
	}

	return ph.Send(SubscribePacket{
		SubscriberID: sub.SubscriberID,
		Action:       sub.Action,
		Receiver:     sub.Receiver,
	})
}
