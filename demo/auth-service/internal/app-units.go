// Code generated by coattail; DO NOT EDIT.
package internal

import (
    "github.com/nathan-fiscaletti/ct1/internal/actions"
    "context"

	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

func (a *App) LoadUnits(ctx context.Context, local *coattailtypes.Peer) error {
    var err error

    // Register actions
    
    err = local.RegisterAction(ctx, coattailtypes.NewAction(&actions.Authenticate{}))
    if err != nil {
        return err
    }
    
    // Register receivers
    
    return nil
}
