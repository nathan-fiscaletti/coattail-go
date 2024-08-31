package logging

import (
	"context"
	"log"
	"os"

	"github.com/nathan-fiscaletti/coattail-go/internal/keys"
)

var logger = log.New(os.Stdout, "[host] ", log.LstdFlags)

func ContextWithLogger(ctx context.Context) context.Context {
	return context.WithValue(ctx, keys.LoggerKey, logger)
}

func GetLogger(ctx context.Context) *log.Logger {
	if v := ctx.Value(keys.LoggerKey); v != nil {
		if l, ok := v.(*log.Logger); ok {
			return l
		}
	}

	return nil
}
