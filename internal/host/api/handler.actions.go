package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nathan-fiscaletti/coattail-go/internal/logging"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

type ActionsHandler struct {
	ctx       context.Context
	localPeer *coattailtypes.Peer
}

func NewActionsHandler(ctx context.Context, localPeer *coattailtypes.Peer) http.Handler {
	return &ActionsHandler{
		ctx:       ctx,
		localPeer: localPeer,
	}
}

func (h *ActionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if logger, _ := logging.GetLogger(h.ctx); logger != nil {
		logger.Printf("GET /actions")
	}

	actions, err := h.localPeer.ListActions(h.ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	actionsData, err := json.Marshal(actions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// disable cors
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(actionsData)
}
