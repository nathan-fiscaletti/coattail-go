package coattailtypes

import "context"

type App interface {
	// OnStart is called after the application has been started.
	OnStart(ctx context.Context, local *Peer)
}

// DefaultApp is a default implementation of the App interface.
type DefaultApp struct{}

func (a *DefaultApp) OnStart(ctx context.Context, local *Peer) {}
