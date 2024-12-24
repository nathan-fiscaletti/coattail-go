package coattailtypes

import (
	"context"
	"log"

	"github.com/nathan-fiscaletti/coattail-go/internal/services/authentication"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailmodels"
)

type PeersFile struct {
	Peers []PeerDetails `yaml:"peers"`
}

type PeerDetails struct {
	// Whether or not the peer is local. This should be true if the peer is the
	// local peer, and false if the peer is a remote peer.
	IsLocal bool `yaml:"is_local,omitempty" json:"is_local"`

	// The address of the peer. This is the address that the peer can be reached
	// at. For the local peer, this should be the address of the local machine.
	Address string `yaml:"address" json:"address"`

	// The token of the peer. This is a secret token that is used to authenticate
	// the peer. For the local peer, this should be an empty string.
	Token string `yaml:"token" json:"-"`
}

// Peer represents any coattail peer, whether local or remote.
type Peer struct {
	PeerDetails
	PeerAdapter `json:"-"`
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
	CredentialManager
	LogManager
}

type ActionManager interface {
	// Run runs an action on the peer. The name of the action should be
	// provided as the first argument. The second argument is the data that
	// should be passed to the action. The return value is the result of the
	// action, or an error if the action failed.
	Run(ctx context.Context, name string, arg any) (any, error)

	// Publish publishes data to the peer. The name of the action that produced
	// the data should be provided as the first argument. The second argument
	// is the data that should be published. The return value is an error if
	// the publish failed.
	Publish(ctx context.Context, name string, data any) error

	// RunAndPublish runs an action on the peer and then publishes the result.
	// The name of the action should be provided as the first argument. The second
	// argument is the data that should be passed to the action. The return value
	// is the result of the action, or an error if the action failed.
	RunAndPublish(ctx context.Context, name string, arg any) error

	// ListActions returns a list of all actions that are available on the peer.
	// The return value is a list of action names, or an error if the list could
	// not be retrieved.
	ListActions(ctx context.Context) ([]string, error)

	// RegisterAction adds an action to the peer. The name of the action should be
	// provided as the first argument. The second argument is the unit that
	// should be executed when the action is run. The return value is an error
	// if the action could not be added.
	RegisterAction(ctx context.Context, name string, unit Unit) error

	// HasAction checks if an action is available on the peer. The name of the
	// action should be provided as the first argument. The return value is
	// true if the action is available, or false if it is not.
	HasAction(ctx context.Context, name string) (bool, error)
}

type ReceiverManager interface {
	// ListReceivers returns a list of all receivers that are available on the peer.
	// The return value is a list of receiver names, or an error if the list could
	// not be retrieved.
	ListReceivers(ctx context.Context) ([]string, error)

	// HasReceiver checks if a receiver is available on the peer. The name of the
	// receiver should be provided as the first argument. The return value is
	// true if the receiver is available, or false if it is not.
	HasReceiver(ctx context.Context, name string) (bool, error)

	// RegisterReceiver adds a receiver to the peer. The name of the receiver should be
	// provided as the first argument. The second argument is the unit that should
	// be executed when the receiver is notified. The return value is an error if
	// the receiver could not be added.
	RegisterReceiver(ctx context.Context, name string, unit Unit) error

	// Notify notifies a receiver on the peer. The name of the receiver
	// should be provided as the first argument. The second argument is the data
	// that should be passed to the receiver. The return value is an error if the
	// notification could not be sent.
	Notify(ctx context.Context, name string, arg any) error
}

type PeerManager interface {
	// GetPeer returns a peer by its address. The address of the peer should be
	// provided as the first argument. The return value is the peer, or an error
	// if the peer could not be found.
	GetPeer(ctx context.Context, address string) (*Peer, error)

	// GetPeerBy returns a peer by a predicate. The predicate should be provided
	// as the first argument. The return value is the peer, or an error if the peer
	// could not be found.
	GetPeerBy(ctx context.Context, predicate func(PeerDetails) bool) (*Peer, error)

	// HasPeer checks if a peer is available on the peer. The address of the peer
	// should be provided as the first argument. The return value is true if the
	// peer is available, or false if it is not.
	HasPeer(ctx context.Context, address string) (bool, error)

	// ListPeers returns a list of all peers that are available on the peer. The
	// return value is a list of peers, or an error if the list could not be
	// retrieved.
	ListPeers(ctx context.Context) ([]*Peer, error)

	// Subscribe subscribes to a peer. The subscription details should be provided
	// as the first argument. The return value is an error if the subscription could
	// not be completed.
	Subscribe(ctx context.Context, sub coattailmodels.Subscription) error
}

type CredentialManager interface {
	IssueToken(ctx context.Context, claims authentication.Claims) (*authentication.Token, error)
}

type LogManager interface {
	Logger(ctx context.Context) (*log.Logger, error)
}
