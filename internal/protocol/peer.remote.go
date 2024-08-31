package protocol

import (
	"context"
	"fmt"
	"net"

	"github.com/nathan-fiscaletti/coattail-go/internal/protocol/protocoltypes"
)

/* ====== Type ====== */

type RemotePeerAdapter struct {
	details protocoltypes.PeerDetails

	packetHandler *PacketHandler
}

func newRemotePeerAdapter(details protocoltypes.PeerDetails) *RemotePeerAdapter {
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

func (i *RemotePeerAdapter) RunAction(arg protocoltypes.RunActionArguments) (any, error) {
	ph, err := i.getPacketHandler()
	if err != nil {
		return nil, err
	}

	packet, err := ph.Request(Request{
		Packet: PerformActionPacket{
			Action:  arg.Name,
			Arg:     arg.Arg,
			Publish: arg.Publish,
		},
	})
	if err != nil {
		return nil, err
	}

	respPacket, isRespPacket := packet.(PerformActionResponsePacket)
	if !isRespPacket {
		return nil, fmt.Errorf("unexpected response packet")
	}

	return respPacket.Data, nil
}

func (i *RemotePeerAdapter) PublishActionResult(name string, data any) error {
	return nil
}

func (i *RemotePeerAdapter) Actions() []protocoltypes.Action {
	return []protocoltypes.Action{}
}

func (i *RemotePeerAdapter) HasAction(name string) bool {
	// TODO: implement
	return false
}

func (i *RemotePeerAdapter) AddAction(name string, unit protocoltypes.Unit) error {
	return fmt.Errorf("cannot add action to remote peer")
}

/* ====== Receivers ====== */

func (i *RemotePeerAdapter) Receivers() []protocoltypes.Receiver {
	return []protocoltypes.Receiver{}
}

func (i *RemotePeerAdapter) HasReceiver(name string) bool {
	return false
}

func (i *RemotePeerAdapter) AddReceiver(name string, unit protocoltypes.Unit) error {
	return fmt.Errorf("cannot add receiver to remote peer")
}

func (i *RemotePeerAdapter) NotifyReceiver(name string, arg any) error {
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

func (i *RemotePeerAdapter) GetPeer(id string) (*protocoltypes.Peer, error) {
	return nil, nil
}

func (i *RemotePeerAdapter) HasPeer(id string) (bool, error) {
	return false, nil
}

func (i *RemotePeerAdapter) ListPeers() ([]*protocoltypes.Peer, error) {
	return []*protocoltypes.Peer{}, nil
}
