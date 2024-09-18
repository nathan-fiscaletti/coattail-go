package adapters

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nathan-fiscaletti/coattail-go/internal/database"
	"github.com/nathan-fiscaletti/coattail-go/internal/host"
	"github.com/nathan-fiscaletti/coattail-go/internal/logging"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailmodels"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

/* ====== Local Peer Initialization ====== */

func InitLocalPeer(host *host.Host) error {
	peers, err := loadPeers()
	if err != nil {
		return fmt.Errorf("error loading peers: %s", err)
	}

	host.LocalPeer = coattailtypes.NewPeer(
		coattailtypes.PeerDetails{
			IsLocal: true,
			Address: host.Config.ServiceAddress,
		},
		&LocalPeerAdapter{
			Units: []coattailtypes.UnitImpl{},
			Peers: peers,
		},
	)
	return nil
}

func loadPeers() ([]coattailtypes.PeerDetails, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	result := []coattailtypes.PeerDetails{}

	fmt.Printf("cwd: %s\n", cwd)

	peersFile := filepath.Join(cwd, "peers.yaml")
	if _, err := os.Stat(peersFile); os.IsNotExist(err) {
		return result, nil
	}

	f, err := os.Open(peersFile)
	if err != nil {
		return nil, err
	}

	peers := struct {
		Peers []coattailtypes.PeerDetails `yaml:"peers"`
	}{}
	err = yaml.NewDecoder(f).Decode(&peers)
	if err != nil {
		return nil, err
	}

	return peers.Peers, nil
}

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

func (i *LocalPeerAdapter) Run(ctx context.Context, name string, arg any) (any, error) {
	logging.GetLogger(ctx).Printf("running action: %s", name)

	return i.runUnit(runUnitArguments{
		Type: coattailtypes.UnitTypeAction,
		Name: name,
		Args: arg,
	})
}

func (i *LocalPeerAdapter) Publish(ctx context.Context, name string, data any) error {
	logging.GetLogger(ctx).Printf("publishing action: %s", name)

	db, err := database.GetDatabase(ctx)
	if err != nil {
		return err
	}

	var action *coattailtypes.UnitImpl

	for _, unit := range i.Units {
		if unit.UnitType == coattailtypes.UnitTypeAction && unit.Name == name {
			action = &unit
			break
		}
	}

	if action == nil {
		return fmt.Errorf("action %s not found", name)
	}

	var subscriptions []coattailmodels.Subscription
	if err := db.Where("action = ?", action.Name).Find(&subscriptions).Error; err != nil {
		return err
	}

	for _, sub := range subscriptions {
		peer, err := i.GetPeerBy(ctx, func(details coattailtypes.PeerDetails) bool {
			return details.Address == sub.Address
		})
		if err != nil {
			return err
		}

		if err := peer.Notify(ctx, sub.Receiver, data); err != nil {
			return err
		}
	}

	return nil
}

func (i *LocalPeerAdapter) RunAndPublish(ctx context.Context, name string, arg any) (any, error) {
	res, err := i.Run(ctx, name, arg)
	if err != nil {
		return nil, err
	}

	if err := i.Publish(ctx, name, res); err != nil {
		return nil, err
	}

	return res, nil
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

func (i *LocalPeerAdapter) Notify(ctx context.Context, name string, arg any) error {
	logging.GetLogger(ctx).Printf("notifying receiver: %s", name)

	_, err := i.runUnit(runUnitArguments{
		Type: coattailtypes.UnitTypeReceiver,
		Name: name,
		Args: arg,
	})

	return err
}

/* ====== Peers ====== */

func (i *LocalPeerAdapter) GetPeer(ctx context.Context, address string) (*coattailtypes.Peer, error) {
	for _, peerDetails := range i.Peers {
		if peerDetails.Address == address {
			return coattailtypes.NewPeer(peerDetails, newRemotePeerAdapter(peerDetails)), nil
		}
	}

	return nil, fmt.Errorf("peer %s not found", address)
}

func (i *LocalPeerAdapter) GetPeerBy(ctx context.Context, predicate func(coattailtypes.PeerDetails) bool) (*coattailtypes.Peer, error) {
	for _, peerDetails := range i.Peers {
		if predicate(peerDetails) {
			return coattailtypes.NewPeer(peerDetails, newRemotePeerAdapter(peerDetails)), nil
		}
	}

	return nil, fmt.Errorf("peer not found")
}

func (i *LocalPeerAdapter) HasPeer(ctx context.Context, address string) (bool, error) {
	return lo.ContainsBy(i.Peers, func(peerDetails coattailtypes.PeerDetails) bool {
		return peerDetails.Address == address
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

	if db.Find(&sub).RowsAffected > 0 {
		return nil
	}

	logging.GetLogger(ctx).Printf("registering subscriber: %s", sub.String())
	return db.Create(&sub).Error
}
