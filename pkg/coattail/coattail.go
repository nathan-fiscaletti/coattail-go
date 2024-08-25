package coattail

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nathan-fiscaletti/coattail-go/internal/protocol"
)

var local *Peer

// Local returns the local peer.
func Local() *Peer {
	if local == nil {
		local = initializeLocalPeer()
	}

	return local
}

// Run starts the local peer. This function will block.
func Run() error {
	return Local().PeerAdapter.(*localPeerAdapter).host.Start()
}

func initializeLocalPeer() *Peer {
	peers, err := loadPeers()
	if err != nil {
		panic(fmt.Sprintf("error loading peers: %s", err))
	}

	return newPeer(
		PeerDetails{
			PeerID: LocalPeerId,
		},
		&localPeerAdapter{
			units: []anyUnit{},
			peers: peers,
			host:  &protocol.Host{},
		},
	)
}

func loadPeers() ([]PeerDetails, error) {
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
