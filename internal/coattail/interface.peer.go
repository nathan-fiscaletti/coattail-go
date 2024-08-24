package coattail

type RunActionArguments struct {
	Name    string
	Args    interface{}
	Publish bool
}

type PeerAdapter interface {
	/* ====== Actions ====== */
	RunAction(name string, arg interface{}) (interface{}, error)
	RunAndPublishAction(name string, arg interface{}) (interface{}, error)

	Actions() []Action

	AddAction(name string, unit Unit) error
	HasAction(name string) bool

	/* ====== Receivers ====== */
	Receivers() []Receiver
	HasReceiver(name string) bool
	AddReceiver(name string, unit Unit) error

	GetPeer(id string) (*Peer, error)
	HasPeer(id string) (bool, error)
	Peers() ([]*Peer, error)
}

type Peer struct {
	PeerAdapter
}

func newPeer(adapter PeerAdapter) *Peer {
	return &Peer{
		PeerAdapter: adapter,
	}
}
