package peers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Manager struct {
	local *Peer
}

func NewManager() (*Manager, error) {
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

	p.local = newPeer(
		PeerDetails{
			PeerID: LocalPeerId,
		},
		&localPeerAdapter{
			units: []anyUnit{},
			peers: peers,
		},
	)
	return nil
}

func (p *Manager) LocalPeer() *Peer {
	return p.local
}

func (p *Manager) loadPeers() ([]PeerDetails, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	result := []PeerDetails{}

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

			peerDetails := PeerDetails{}
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
