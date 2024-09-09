package coattailtypes

import (
	"context"

	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailmodels"
)

const (
	// LocalPeerId is the ID of the local peer. This is a unique identifier for
	// the local peer.
	LocalPeerId string = "local"
)

type PeerDetails struct {
	// The ID of the peer. This is a unique identifier for the peer. For the
	// local peer, this should contain LocalPeerId. For remote peers, this is
	// the ID that was assigned to the peer when it was registered.
	PeerID string `yaml:"id"`

	// The address of the peer. This is the address that the peer can be reached
	// at. For the local peer, this should be the address of the local machine.
	Address string `yaml:"address"`

	// The token of the peer. This is a secret token that is used to authenticate
	// the peer. For the local peer, this should be an empty string.
	Token string `yaml:"token"`
}

// Peer represents any coattail peer, whether local or remote.
type Peer struct {
	PeerDetails
	PeerAdapter
}

// NewPeer creates a new peer with the provided details and adapter.
// You should not use this function directly. Instead, use the coattail.Manage()
// function to retrieve the local peer, and the GetPeer() method on the local
// peer to retrieve remote peers.
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
	// RunAction runs an action on the peer. The name of the action should be
	// provided as the first argument. The second argument is the data that
	// should be passed to the action. The return value is the result of the
	// action, or an error if the action failed.
	RunAction(ctx context.Context, name string, arg any) (any, error)

	// Publish publishes data to the peer. The name of the action that produced
	// the data should be provided as the first argument. The second argument
	// is the data that should be published. The return value is an error if
	// the publish failed.
	Publish(ctx context.Context, name string, data any) error

	// Actions returns a list of all actions that are available on the peer.
	// The return value is a list of action names, or an error if the list could
	// not be retrieved.
	Actions(ctx context.Context) ([]string, error)

	// AddAction adds an action to the peer. The name of the action should be
	// provided as the first argument. The second argument is the unit that
	// should be executed when the action is run. The return value is an error
	// if the action could not be added.
	AddAction(ctx context.Context, name string, unit Unit) error

	// HasAction checks if an action is available on the peer. The name of the
	// action should be provided as the first argument. The return value is
	// true if the action is available, or false if it is not.
	HasAction(ctx context.Context, name string) (bool, error)
}

type ReceiverManager interface {
	// Receivers returns a list of all receivers that are available on the peer.
	// The return value is a list of receiver names, or an error if the list could
	// not be retrieved.
	Receivers(ctx context.Context) ([]string, error)
	// HasReceiver checks if a receiver is available on the peer. The name of the
	// receiver should be provided as the first argument. The return value is
	// true if the receiver is available, or false if it is not.
	HasReceiver(ctx context.Context, name string) (bool, error)
	// AddReceiver adds a receiver to the peer. The name of the receiver should be
	// provided as the first argument. The second argument is the unit that should
	// be executed when the receiver is notified. The return value is an error if
	// the receiver could not be added.
	AddReceiver(ctx context.Context, name string, unit Unit) error
	// NotifyReceiver notifies a receiver on the peer. The name of the receiver
	// should be provided as the first argument. The second argument is the data
	// that should be passed to the receiver. The return value is an error if the
	// notification could not be sent.
	NotifyReceiver(ctx context.Context, name string, arg any) error
}

type PeerManager interface {
	// GetPeer returns a peer by its ID. The ID of the peer should be provided as
	// the first argument. The return value is the peer, or an error if the peer
	// could not be found.
	GetPeer(ctx context.Context, id string) (*Peer, error)
	// HasPeer checks if a peer is available on the peer. The ID of the peer should
	// be provided as the first argument. The return value is true if the peer is
	// available, or false if it is not.
	HasPeer(ctx context.Context, id string) (bool, error)
	// ListPeers returns a list of all peers that are available on the peer. The
	// return value is a list of peers, or an error if the list could not be
	// retrieved.
	ListPeers(ctx context.Context) ([]*Peer, error)

	// Subscribe subscribes to a peer. The subscription details should be provided
	// as the first argument. The return value is an error if the subscription could
	// not be completed.
	Subscribe(ctx context.Context, sub coattailmodels.Subscription) error
}
