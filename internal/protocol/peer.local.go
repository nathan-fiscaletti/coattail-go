package protocol

import (
	"context"
	"fmt"

	"github.com/nathan-fiscaletti/coattail-go/internal/database"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailmodels"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
	"github.com/samber/lo"
)

/* ====== Type ====== */

type LocalPeerAdapter struct {
	Units []coattailtypes.UnitImpl
	Peers []coattailtypes.PeerDetails
}

/* ====== Units ====== */

func (i *LocalPeerAdapter) getUnit(hType coattailtypes.UnitType, name string) (coattailtypes.UnitImpl, error) {
	h, ok := lo.Find(i.Units, func(h coattailtypes.UnitImpl) bool {
		return h.UnitType == hType && h.Name == name
	})
	if !ok {
		return coattailtypes.UnitImpl{}, fmt.Errorf("handler %s not found", name)
	}

	return h, nil
}

type runUnitArguments struct {
	Type coattailtypes.UnitType
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

func (i *LocalPeerAdapter) RunAction(ctx context.Context, name string, arg any) (any, error) {
	return i.runUnit(runUnitArguments{
		Type: coattailtypes.UnitTypeAction,
		Name: name,
		Args: arg,
	})
}

func (i *LocalPeerAdapter) Publish(ctx context.Context, name string, data any) error {
	db, err := database.GetDatabase(ctx)
	if err != nil {
		return err
	}

	for _, unit := range i.Units {
		if unit.UnitType == coattailtypes.UnitTypeAction && unit.Name == name {
			var subscriptions []coattailmodels.Subscription
			if err := db.Where("action = ?", name).Find(&subscriptions).Error; err != nil {
				return err
			}

			for _, sub := range subscriptions {
				peer, err := i.GetPeer(ctx, sub.SubscriberID.String())
				if err != nil {
					return err
				}

				if err := peer.NotifyReceiver(ctx, sub.Receiver, data); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (i *LocalPeerAdapter) Actions(ctx context.Context) ([]string, error) {
	return lo.Map(lo.Filter(i.Units, func(h coattailtypes.UnitImpl, _ int) bool {
		return h.UnitType == coattailtypes.UnitTypeAction
	}), func(h coattailtypes.UnitImpl, _ int) string {
		return h.Name
	}), nil
}

func (i *LocalPeerAdapter) HasAction(ctx context.Context, name string) (bool, error) {
	return lo.ContainsBy(i.Units, func(h coattailtypes.UnitImpl) bool {
		return h.UnitType == coattailtypes.UnitTypeAction && h.Name == name
	}), nil
}

func (i *LocalPeerAdapter) AddAction(ctx context.Context, name string, unit coattailtypes.Unit) error {
	exists, _ := i.HasAction(ctx, name)
	if exists {
		return fmt.Errorf("action %s already exists", name)
	}

	i.Units = append(i.Units, coattailtypes.UnitImpl{
		Unit:     unit,
		Name:     name,
		UnitType: coattailtypes.UnitTypeAction,
	})

	return nil
}

/* ====== Receivers ====== */

func (i *LocalPeerAdapter) Receivers(ctx context.Context) ([]string, error) {
	return lo.Map(lo.Filter(i.Units, func(h coattailtypes.UnitImpl, _ int) bool {
		return h.UnitType == coattailtypes.UnitTypeReceiver
	}), func(h coattailtypes.UnitImpl, _ int) string {
		return h.Name
	}), nil
}

func (i *LocalPeerAdapter) HasReceiver(ctx context.Context, name string) (bool, error) {
	return lo.ContainsBy(i.Units, func(h coattailtypes.UnitImpl) bool {
		return h.UnitType == coattailtypes.UnitTypeReceiver && h.Name == name
	}), nil
}

func (i *LocalPeerAdapter) AddReceiver(ctx context.Context, name string, unit coattailtypes.Unit) error {
	exists, _ := i.HasReceiver(ctx, name)
	if exists {
		return fmt.Errorf("receiver %s already exists", name)
	}

	i.Units = append(i.Units, coattailtypes.UnitImpl{
		Unit:     unit,
		Name:     name,
		UnitType: coattailtypes.UnitTypeReceiver,
	})

	return nil
}

func (i *LocalPeerAdapter) NotifyReceiver(ctx context.Context, name string, arg any) error {
	_, err := i.runUnit(runUnitArguments{
		Type: coattailtypes.UnitTypeReceiver,
		Name: name,
		Args: arg,
	})

	return err
}

/* ====== Peers ====== */

func (i *LocalPeerAdapter) GetPeer(ctx context.Context, id string) (*coattailtypes.Peer, error) {
	for _, peerDetails := range i.Peers {
		if peerDetails.PeerID == id {
			return coattailtypes.NewPeer(peerDetails, newRemotePeerAdapter(peerDetails)), nil
		}
	}

	return nil, fmt.Errorf("peer %s not found", id)
}

func (i *LocalPeerAdapter) HasPeer(ctx context.Context, id string) (bool, error) {
	return lo.ContainsBy(i.Peers, func(peerDetails coattailtypes.PeerDetails) bool {
		return peerDetails.PeerID == id
	}), nil
}

func (i *LocalPeerAdapter) ListPeers(ctx context.Context) ([]*coattailtypes.Peer, error) {
	return lo.Map(i.Peers, func(peerDetails coattailtypes.PeerDetails, _ int) *coattailtypes.Peer {
		return coattailtypes.NewPeer(peerDetails, newRemotePeerAdapter(peerDetails))
	}), nil
}

func (i *LocalPeerAdapter) Subscribe(ctx context.Context, sub coattailmodels.Subscription) error {
	db, err := database.GetDatabase(ctx)
	if err != nil {
		return err
	}

	return db.Create(&sub).Error
}
