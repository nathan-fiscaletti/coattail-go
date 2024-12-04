package logging

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/nathan-fiscaletti/coattail-go/internal/keys"
)

var (
	ErrLoggerNotFound = errors.New("logger not found in context")
)

func ContextWithLogger(ctx context.Context, appName string) (context.Context, error) {
	logger := log.New(os.Stdout, "["+appName+"] ", log.LstdFlags)
	return context.WithValue(ctx, keys.LoggerKey, logger), nil
}

func GetLogger(ctx context.Context) (*log.Logger, error) {
	if v := ctx.Value(keys.LoggerKey); v != nil {
		if l, ok := v.(*log.Logger); ok {
			return l, nil
		}
	}

	return nil, ErrLoggerNotFound
}
