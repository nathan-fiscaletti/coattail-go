package coattail

import "github.com/nathan-fiscaletti/coattail-go/internal/managers/peers"

type Peer peers.Peer
type Action peers.Action
type PeerDetails peers.PeerDetails
type Receiver peers.Receiver
type Unit peers.Unit
type UnitHandler peers.UnitHandler

const LocalPeerId string = peers.LocalPeerId
