package util

import (
	"context"

	"github.com/nathan-fiscaletti/coattail-go/internal/logging"
)

func CreateServiceContext(ctx context.Context) (context.Context, error) {
	ctx, err := logging.ContextWithLogger(
		ctx,
		"coattail",
	)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}
