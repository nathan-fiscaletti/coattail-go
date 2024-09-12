package protocol

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nathan-fiscaletti/coattail-go/internal/host/config"
	"github.com/nathan-fiscaletti/coattail-go/internal/keys"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
	"gopkg.in/yaml.v3"
)

func ContextWithManager(ctx context.Context) (context.Context, error) {
	mgr, err := newManager()
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, keys.PeerManagerKey, mgr), nil
}

func GetManager(ctx context.Context) *Manager {
	if v := ctx.Value(keys.PeerManagerKey); v != nil {
		if m, ok := v.(*Manager); ok {
			return m
		}
	}

	return nil
}

type Manager struct {
	local *coattailtypes.Peer
}

func newManager() (*Manager, error) {
	s := &Manager{}
	err := s.loadLocalPeer()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (p *Manager) loadLocalPeer() error {
	peers, err := p.loadPeers()
	if err != nil {
		return fmt.Errorf("error loading peers: %s", err)
	}

	hostConfig, err := config.GetHostConfig()
	if err != nil {
		return fmt.Errorf("error getting host: %s", err)
	}

	p.local = coattailtypes.NewPeer(
		coattailtypes.PeerDetails{
			PeerID:  coattailtypes.LocalPeerId,
			Address: hostConfig.ServiceAddress,
		},
		&LocalPeerAdapter{
			Units: []coattailtypes.UnitImpl{},
			Peers: peers,
		},
	)
	return nil
}

func (p *Manager) LocalPeer() *coattailtypes.Peer {
	return p.local
}

func (p *Manager) loadPeers() ([]coattailtypes.PeerDetails, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	result := []coattailtypes.PeerDetails{}

	peerDir := filepath.Join(cwd, "peers")
	if _, err := os.Stat(peerDir); os.IsNotExist(err) {
		return result, nil
	}

	files, err := os.ReadDir(peerDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".yaml" {
			f, err := os.Open(filepath.Join(peerDir, file.Name()))
			if err != nil {
				fmt.Printf("error opening peer file %s: %s\n", file.Name(), err)
				continue
			}

			peerDetails := coattailtypes.PeerDetails{}
			err = yaml.NewDecoder(f).Decode(&peerDetails)
			if err != nil {
				fmt.Printf("error decoding peer file %s: %s\n", file.Name(), err)
				continue
			}

			result = append(result, peerDetails)
		}
	}

	return result, nil
}
