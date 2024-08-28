package protocol

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nathan-fiscaletti/coattail-go/internal/keys"
	"github.com/nathan-fiscaletti/coattail-go/internal/protocol/protocoltypes"
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
	local *protocoltypes.Peer
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

	p.local = protocoltypes.NewPeer(
		protocoltypes.PeerDetails{
			PeerID: protocoltypes.LocalPeerId,
		},
		&LocalPeerAdapter{
			Units: []protocoltypes.AnyUnit{},
			Peers: peers,
		},
	)
	return nil
}

func (p *Manager) LocalPeer() *protocoltypes.Peer {
	return p.local
}

func (p *Manager) loadPeers() ([]protocoltypes.PeerDetails, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	result := []protocoltypes.PeerDetails{}

	peerDir := filepath.Join(cwd, "peers")
	if _, err := os.Stat(peerDir); os.IsNotExist(err) {
		return result, nil
	}

	files, err := os.ReadDir(peerDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			f, err := os.Open(filepath.Join(peerDir, file.Name()))
			if err != nil {
				fmt.Printf("error opening peer file %s: %s\n", file.Name(), err)
				continue
			}

			peerDetails := protocoltypes.PeerDetails{}
			err = json.NewDecoder(f).Decode(&peerDetails)
			if err != nil {
				fmt.Printf("error decoding peer file %s: %s\n", file.Name(), err)
				continue
			}

			result = append(result, peerDetails)
		}
	}

	return result, nil
}
