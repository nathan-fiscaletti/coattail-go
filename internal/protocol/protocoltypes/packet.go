package protocoltypes

import (
	"context"
)

type Packet interface {
	Handle(ctx context.Context) (any, error)
}

type EmptyPacket interface {
	Empty() bool
}
