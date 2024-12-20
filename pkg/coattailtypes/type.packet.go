package coattailtypes

import (
	"context"
)

type Packet interface {
	Handle(ctx context.Context) (Packet, error)
}
