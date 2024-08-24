package coattail

import (
	"fmt"

	"github.com/samber/lo"
)

/* ====== Type ====== */

type localPeerAdapter struct {
	units []anyUnit
	peers []PeerDetails
}

/* ====== Units ====== */

func (i *localPeerAdapter) getUnit(hType unitType, name string) (anyUnit, error) {
	h, ok := lo.Find(i.units, func(h anyUnit) bool {
		return h.unitType == hType && h.name == name
	})
	if !ok {
		return anyUnit{}, fmt.Errorf("handler %s not found", name)
	}

	return h, nil
}

type runUnitArguments struct {
	Type unitType
	Name string
	Args interface{}
}

func (i *localPeerAdapter) runUnit(arg runUnitArguments) (interface{}, error) {
	h, err := i.getUnit(arg.Type, arg.Name)
	if err != nil {
		return nil, err
	}

	return h.Unit(arg.Args)
}

/* ====== Actions ====== */

func (i *localPeerAdapter) RunAction(name string, arg interface{}) (interface{}, error) {
	return i.runUnit(runUnitArguments{
		Type: unitTypeAction,
		Name: name,
		Args: arg,
	})
}

func (i *localPeerAdapter) RunAndPublishAction(name string, arg interface{}) (interface{}, error) {
	result, err := i.RunAction(name, arg)
	if err != nil {
		return nil, err
	}

	err = i.publishActionResult(name, arg, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// TODO: implement
func (i *localPeerAdapter) publishActionResult(_ string, _ interface{}, _ interface{}) error {
	return nil
}

func (i *localPeerAdapter) Actions() []Action {
	return lo.Map(lo.Filter(i.units, func(h anyUnit, _ int) bool {
		return h.unitType == unitTypeAction
	}), func(h anyUnit, _ int) Action {
		return Action{
			name: h.name,
			Unit: h.Unit,
			peer: i,
		}
	})
}

func (i *localPeerAdapter) HasAction(name string) bool {
	return lo.ContainsBy(i.units, func(h anyUnit) bool {
		return h.unitType == unitTypeAction && h.name == name
	})
}

func (i *localPeerAdapter) AddAction(name string, unit Unit) error {
	if i.HasAction(name) {
		return fmt.Errorf("action %s already exists", name)
	}

	i.units = append(i.units, anyUnit{
		Unit:     unit,
		name:     name,
		unitType: unitTypeAction,
	})

	return nil
}

/* ====== Receivers ====== */

func (i *localPeerAdapter) Receivers() []Receiver {
	return lo.Map(lo.Filter(i.units, func(h anyUnit, _ int) bool {
		return h.unitType == unitTypeReceiver
	}), func(h anyUnit, _ int) Receiver {
		return Receiver{
			name: h.name,
			Unit: h.Unit,

			peer: i,
		}
	})
}

func (i *localPeerAdapter) HasReceiver(name string) bool {
	return lo.ContainsBy(i.units, func(h anyUnit) bool {
		return h.unitType == unitTypeReceiver && h.name == name
	})
}

func (i *localPeerAdapter) AddReceiver(name string, unit Unit) error {
	if i.HasReceiver(name) {
		return fmt.Errorf("receiver %s already exists", name)
	}

	i.units = append(i.units, anyUnit{
		Unit:     unit,
		name:     name,
		unitType: unitTypeReceiver,
	})

	return nil
}

func (i *localPeerAdapter) GetPeer(id string) (*Peer, error) {
	for _, peerDetails := range i.peers {
		if peerDetails.PeerID == id {
			return newPeer(peerDetails, newRemotePeerAdapter(peerDetails)), nil
		}
	}

	return nil, fmt.Errorf("peer %s not found", id)
}

func (i *localPeerAdapter) HasPeer(id string) (bool, error) {
	return lo.ContainsBy(i.peers, func(peerDetails PeerDetails) bool {
		return peerDetails.PeerID == id
	}), nil
}

func (i *localPeerAdapter) Peers() ([]*Peer, error) {
	return lo.Map(i.peers, func(peerDetails PeerDetails, _ int) *Peer {
		return newPeer(peerDetails, newRemotePeerAdapter(peerDetails))
	}), nil
}
