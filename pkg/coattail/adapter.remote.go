package coattail

import (
	"fmt"
	"net"

	"github.com/nathan-fiscaletti/coattail-go/internal/protocol"
)

/* ====== Type ====== */

type remotePeerAdapter struct {
	details PeerDetails

	comm *protocol.Communicator
}

// RunCommunicationTest runs a communication test with the remote peer.
// This is a temporary development function and will be removed in the future.
func (i *remotePeerAdapter) RunCommunicationTest() error {
	comm, err := i.communicator()
	if err != nil {
		return err
	}

	err = comm.WritePacket(protocol.HelloPacketData{
		Message: "Hello, I am the first functional packet!",
	})
	if err != nil {
		return err
	}

	return nil
}

func newRemotePeerAdapter(details PeerDetails) *remotePeerAdapter {
	return &remotePeerAdapter{
		details: details,
	}
}

func (i *remotePeerAdapter) communicator() (*protocol.Communicator, error) {
	if i.comm == nil || i.comm.IsFinished() {
		conn, err := net.Dial("tcp", i.details.Address)
		if err != nil {
			return nil, err
		}

		i.comm = protocol.NewCommunicator(conn)
		i.comm.Start()
	}

	return i.comm, nil
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
