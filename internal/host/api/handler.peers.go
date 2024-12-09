package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nathan-fiscaletti/coattail-go/internal/logging"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

type PeersHandler struct {
	ctx       context.Context
	localPeer *coattailtypes.Peer
}

func NewPeersHandler(ctx context.Context, localPeer *coattailtypes.Peer) http.Handler {
	return &PeersHandler{
		ctx:       ctx,
		localPeer: localPeer,
	}
}

func (h *PeersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if logger, _ := logging.GetLogger(h.ctx); logger != nil {
		logger.Printf("GET /peers")
	}

	peers, err := h.localPeer.ListPeers(h.ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	peerData, err := json.Marshal(peers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// disable cors
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(peerData)
}
