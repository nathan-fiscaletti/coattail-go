package coattail

type ActionManager interface {
	RunAction(name string, arg interface{}) (interface{}, error)
	RunAndPublishAction(name string, arg interface{}) (interface{}, error)
	Actions() []Action
	AddAction(name string, unit Unit) error
	HasAction(name string) bool
}

type ReceiverManager interface {
	Receivers() []Receiver
	HasReceiver(name string) bool
	AddReceiver(name string, unit Unit) error
}

type PeerManager interface {
	GetPeer(id string) (*Peer, error)
	HasPeer(id string) (bool, error)
	Peers() ([]*Peer, error)
}

type PeerAdapter interface {
	ActionManager
	ReceiverManager
	PeerManager
}

type PeerDetails struct {
	PeerID string `json:"id"`
}

type Peer struct {
	PeerDetails
	PeerAdapter
}

func newPeer(details PeerDetails, adapter PeerAdapter) *Peer {
	return &Peer{
		PeerDetails: details,
		PeerAdapter: adapter,
	}
}
