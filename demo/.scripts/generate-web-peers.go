package main

import (
	"os"
	"path/filepath"

	"github.com/nathan-fiscaletti/coattail-go/internal/api"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
	"gopkg.in/yaml.v3"
)

func main() {
	peers := coattailtypes.PeersFile{
		Peers: []coattailtypes.PeerDetails{{
			Address: "192.168.100.2:5243",
			Token:   api.CreateToken(filepath.Join(".", "auth-service", "secret.key"), "0.0.0.0/0", 7, ""),
		}},
	}

	// Write the peers file to the web-service directory
	peersPath := filepath.Join(".", "web-service", "peers.yaml")
	yamlData, err := yaml.Marshal(peers)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(peersPath, yamlData, 0644)
	if err != nil {
		panic(err)
	}
}
