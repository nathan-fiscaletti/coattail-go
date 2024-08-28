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

func (i *RemotePeerAdapter) RunAction(name string, arg any) (any, error) {
	ph, err := i.getPacketHandler()
	if err != nil {
		return nil, err
	}

	respPacket, err := ph.Request(Request{
		Packet: PerformActionPacket{
			Action: name,
			Arg:    arg,
		},
	})
	if err != nil {
		return nil, err
	}

	return respPacket.(PerformActionResponsePacket).Data, nil
}

func (i *RemotePeerAdapter) RunAndPublishAction(name string, arg any) (any, error) {
	return nil, nil
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
