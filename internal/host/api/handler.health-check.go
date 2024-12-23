package api

import (
	"context"
	"net/http"

	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

type HealthCheckHandler struct {
	ctx       context.Context
	localPeer *coattailtypes.Peer
}

func NewHealthCheckHandler(ctx context.Context, localPeer *coattailtypes.Peer) http.Handler {
	return &HealthCheckHandler{
		ctx:       ctx,
		localPeer: localPeer,
	}
}

func (h *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
