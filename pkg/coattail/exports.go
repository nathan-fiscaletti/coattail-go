package coattail

import "github.com/nathan-fiscaletti/coattail-go/internal/protocol/protocoltypes"

const LocalPeerId string = protocoltypes.LocalPeerId

type PeerDetails protocoltypes.PeerDetails

type Unit protocoltypes.Unit
type UnitHandler protocoltypes.UnitHandler

var NewUnit = protocoltypes.NewUnit

type Action protocoltypes.Action
type Receiver protocoltypes.Receiver
