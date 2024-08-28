package protocoltypes

const (
	LocalPeerId string = "local"
)

type PeerDetails struct {
	// The ID of the peer. This is a unique identifier for the peer. For the
	// local peer, this should contain coattail.LocalPeerId. For remote peers,
	// this is the ID that was assigned to the peer when it was registered.
	PeerID string `json:"id"`

	// The address of the peer. This is the address that the peer can be reached
	// at. For the local peer, this should be the address of the local machine.
	Address string `json:"address"`

	// The token of the peer. This is a secret token that is used to authenticate
	// the peer. For the local peer, this should be an empty string.
	Token string `json:"token"`
}

// Peer represents any coattail peer, whether local or remote.
type Peer struct {
	PeerDetails
	PeerAdapter
}

func NewPeer(details PeerDetails, adapter PeerAdapter) *Peer {
	return &Peer{
		PeerDetails: details,
		PeerAdapter: adapter,
	}
}

// PeerAdapter is a unified interface for both local and remote peers.
type PeerAdapter interface {
	ActionManager
	ReceiverManager
	PeerManager
}

type ActionManager interface {
	RunAction(name string, arg any) (any, error)
	RunAndPublishAction(name string, arg any) (any, error)
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
	ListPeers() ([]*Peer, error)
}
