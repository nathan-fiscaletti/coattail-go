package coattailtypes

import "context"

type App interface {
	// OnStart is called after the application has been started.
	OnStart(ctx context.Context, local *Peer)
	// LoadUnits is called after the application has been started and should be
	// used to load the units for the Coattail instance.
	LoadUnits(ctx context.Context, local *Peer) error
}

type AppWithName interface {
	Name() string
}

// DefaultApp is a default implementation of the App interface.
type DefaultApp struct{}

func (a *DefaultApp) OnStart(ctx context.Context, local *Peer) {
	// Not implemented by default.
}

func (a *DefaultApp) LoadUnits(ctx context.Context, local *Peer) error {
	return nil
}
