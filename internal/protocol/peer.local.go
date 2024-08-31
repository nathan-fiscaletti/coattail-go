package protocol

import (
	"fmt"

	"github.com/nathan-fiscaletti/coattail-go/internal/protocol/protocoltypes"
	"github.com/samber/lo"
)

/* ====== Type ====== */

type LocalPeerAdapter struct {
	Units []protocoltypes.AnyUnit
	Peers []protocoltypes.PeerDetails
}

/* ====== Units ====== */

func (i *LocalPeerAdapter) getUnit(hType protocoltypes.UnitType, name string) (protocoltypes.AnyUnit, error) {
	h, ok := lo.Find(i.Units, func(h protocoltypes.AnyUnit) bool {
		return h.UnitType == hType && h.Name == name
	})
	if !ok {
		return protocoltypes.AnyUnit{}, fmt.Errorf("handler %s not found", name)
	}

	return h, nil
}

type runUnitArguments struct {
	Type protocoltypes.UnitType
	Name string
	Args any
}

func (i *LocalPeerAdapter) runUnit(arg runUnitArguments) (any, error) {
	h, err := i.getUnit(arg.Type, arg.Name)
	if err != nil {
		return nil, err
	}

	return h.Execute(arg.Args)
}

/* ====== Actions ====== */

func (i *LocalPeerAdapter) RunAction(arg protocoltypes.RunActionArguments) (any, error) {
	return i.runUnit(runUnitArguments{
		Type: protocoltypes.UnitTypeAction,
		Name: arg.Name,
		Args: arg.Arg,
	})
}

func (i *LocalPeerAdapter) PublishActionResult(name string, data any) error {
	return nil
}

func (i *LocalPeerAdapter) Actions() []protocoltypes.Action {
	return lo.Map(lo.Filter(i.Units, func(h protocoltypes.AnyUnit, _ int) bool {
		return h.UnitType == protocoltypes.UnitTypeAction
	}), func(h protocoltypes.AnyUnit, _ int) protocoltypes.Action {
		return protocoltypes.Action{
			Unit: h.Unit,
			Name: h.Name,
		}
	})
}

func (i *LocalPeerAdapter) HasAction(name string) bool {
	return lo.ContainsBy(i.Units, func(h protocoltypes.AnyUnit) bool {
		return h.UnitType == protocoltypes.UnitTypeAction && h.Name == name
	})
}

func (i *LocalPeerAdapter) AddAction(name string, unit protocoltypes.Unit) error {
	if i.HasAction(name) {
		return fmt.Errorf("action %s already exists", name)
	}

	i.Units = append(i.Units, protocoltypes.AnyUnit{
		Unit:     unit,
		Name:     name,
		UnitType: protocoltypes.UnitTypeAction,
	})

	return nil
}

/* ====== Receivers ====== */

func (i *LocalPeerAdapter) Receivers() []protocoltypes.Receiver {
	return lo.Map(lo.Filter(i.Units, func(h protocoltypes.AnyUnit, _ int) bool {
		return h.UnitType == protocoltypes.UnitTypeReceiver
	}), func(h protocoltypes.AnyUnit, _ int) protocoltypes.Receiver {
		return protocoltypes.Receiver{
			Unit: h.Unit,
			Name: h.Name,
		}
	})
}

func (i *LocalPeerAdapter) HasReceiver(name string) bool {
	return lo.ContainsBy(i.Units, func(h protocoltypes.AnyUnit) bool {
		return h.UnitType == protocoltypes.UnitTypeReceiver && h.Name == name
	})
}

func (i *LocalPeerAdapter) AddReceiver(name string, unit protocoltypes.Unit) error {
	if i.HasReceiver(name) {
		return fmt.Errorf("receiver %s already exists", name)
	}

	i.Units = append(i.Units, protocoltypes.AnyUnit{
		Unit:     unit,
		Name:     name,
		UnitType: protocoltypes.UnitTypeReceiver,
	})

	return nil
}

func (i *LocalPeerAdapter) NotifyReceiver(name string, arg any) error {
	_, err := i.runUnit(runUnitArguments{
		Type: protocoltypes.UnitTypeReceiver,
		Name: name,
		Args: arg,
	})

	return err
}

/* ====== Peers ====== */

func (i *LocalPeerAdapter) GetPeer(id string) (*protocoltypes.Peer, error) {
	for _, peerDetails := range i.Peers {
		if peerDetails.PeerID == id {
			return protocoltypes.NewPeer(peerDetails, newRemotePeerAdapter(peerDetails)), nil
		}
	}

	return nil, fmt.Errorf("peer %s not found", id)
}

func (i *LocalPeerAdapter) HasPeer(id string) (bool, error) {
	return lo.ContainsBy(i.Peers, func(peerDetails protocoltypes.PeerDetails) bool {
		return peerDetails.PeerID == id
	}), nil
}

func (i *LocalPeerAdapter) ListPeers() ([]*protocoltypes.Peer, error) {
	return lo.Map(i.Peers, func(peerDetails protocoltypes.PeerDetails, _ int) *protocoltypes.Peer {
		return protocoltypes.NewPeer(peerDetails, newRemotePeerAdapter(peerDetails))
	}), nil
}
