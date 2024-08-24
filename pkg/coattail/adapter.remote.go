package coattail

import "fmt"

/* ====== Type ====== */

type remotePeerAdapter struct {
	details PeerDetails
}

func newRemotePeerAdapter(details PeerDetails) *remotePeerAdapter {
	return &remotePeerAdapter{
		details: details,
	}
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

func (i *remotePeerAdapter) GetPeer(id string) (*Peer, error) {
	return nil, nil
}

func (i *remotePeerAdapter) HasPeer(id string) (bool, error) {
	return false, nil
}

func (i *remotePeerAdapter) Peers() ([]*Peer, error) {
	return []*Peer{}, nil
}
