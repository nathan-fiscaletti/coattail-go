// Coattail is a package.
package coattail

import (
	"context"
	"net"
	"reflect"

	"github.com/nathan-fiscaletti/coattail-go/internal/adapters"
	"github.com/nathan-fiscaletti/coattail-go/internal/database"
	"github.com/nathan-fiscaletti/coattail-go/internal/host"
	"github.com/nathan-fiscaletti/coattail-go/internal/logging"
	"github.com/nathan-fiscaletti/coattail-go/internal/packets"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

// Run starts the local peer and runs the main function. This function will block forever.
func Run(app coattailtypes.App) error {
	if app == nil {
		app = &coattailtypes.DefaultApp{}
	}

	ctx, err := createContext(app)
	if err != nil {
		return err
	}

	h, err := host.GetHost(ctx)
	if err != nil {
		return err
	}

	// Initialize the local peer in memory for the host.
	if err := adapters.InitLocalPeer(h); err != nil {
		return err
	}

	if err := h.Start(ctx, func(ctx context.Context, conn net.Conn) {
		go packets.NewHandler(ctx, conn).HandlePackets()
	}); err != nil {
		return err
	}

	app.OnStart(ctx, h.LocalPeer)

	// Block forever
	select {}
}

func createContext(app coattailtypes.App) (context.Context, error) {
	ctx := logging.ContextWithLogger(
		context.Background(),
		reflect.TypeOf(app).String(),
	)

	// Initialize the database
	ctx, err := database.ContextWithDatabase(ctx, database.DatabaseConfig{
		Path: "./data.db",
	})
	if err != nil {
		return nil, err
	}

	ctx, err = host.ContextWithHost(ctx)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}
