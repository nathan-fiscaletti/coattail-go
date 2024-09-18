package coattailtypes

import "context"

type App interface {
	OnStart(ctx context.Context, local *Peer)
}

type DefaultApp struct{}

func (a *DefaultApp) OnStart(ctx context.Context, local *Peer) {}
