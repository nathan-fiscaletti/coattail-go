package peers

import (
	"context"
	"fmt"
	"net"

	"github.com/nathan-fiscaletti/coattail-go/internal/protocol"
	"github.com/nathan-fiscaletti/coattail-go/internal/protocol/packets"
)

/* ====== Type ====== */

type remotePeerAdapter struct {
	details PeerDetails

	packetHandler *protocol.PacketHandler
}

// RunCommunicationTest runs a communication test with the remote peer.
// This is a temporary development function and will be removed in the future.
func (i *remotePeerAdapter) RunCommunicationTest() error {
	packetHandler, err := i.getPacketHandler()
	if err != nil {
		return err
	}

	resp, err := packetHandler.Request(protocol.Request{
		Packet: packets.HelloPacket{
			Message: "Hello, I am the first functional packet!",
		},
	})
	if err != nil {
		return err
	}

	fmt.Printf("response: %v\n", resp)
	resp.Handle(packetHandler.Context())

	return nil
}

func newRemotePeerAdapter(details PeerDetails) *remotePeerAdapter {
	return &remotePeerAdapter{
		details: details,
	}
}

func (i *remotePeerAdapter) getPacketHandler() (*protocol.PacketHandler, error) {
	if i.packetHandler == nil || !i.packetHandler.IsConnected() {
		conn, err := net.Dial("tcp", i.details.Address)
		if err != nil {
			return nil, err
		}

		ctx := context.Background()
		i.packetHandler = protocol.NewPacketHandler(ctx, conn)
		i.packetHandler.HandlePackets()
	}

	return i.packetHandler, nil
}

/* ====== Actions ====== */

func (i *remotePeerAdapter) RunAction(name string, arg interface{}) (interface{}, error) {
	return nil, nil
}

func (i *remotePeerAdapter) RunAndPublishAction(name string, arg interface{}) (interface{}, error) {
	return nil, nil
}

func (i *remotePeerAdapter) Actions() []Action {
	return []Action{}
}

func (i *remotePeerAdapter) HasAction(name string) bool {
	// TODO: implement
	return false
}

func (i *remotePeerAdapter) AddAction(name string, unit Unit) error {
	return fmt.Errorf("cannot add action to remote peer")
}

/* ====== Receivers ====== */

func (i *remotePeerAdapter) Receivers() []Receiver {
	return []Receiver{}
}

func (i *remotePeerAdapter) HasReceiver(name string) bool {
	return false
}

func (i *remotePeerAdapter) AddReceiver(name string, unit Unit) error {
	return fmt.Errorf("cannot add receiver to remote peer")
}

/* ====== Peers ====== */

func (i *remotePeerAdapter) GetPeer(id string) (*Peer, error) {
	return nil, nil
}

func (i *remotePeerAdapter) HasPeer(id string) (bool, error) {
	return false, nil
}

func (i *remotePeerAdapter) Peers() ([]*Peer, error) {
	return []*Peer{}, nil
}
