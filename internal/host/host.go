package host

import (
	"context"
	"fmt"
	"net"

	"github.com/nathan-fiscaletti/coattail-go/internal/protocol"
)

type HostConfig struct {
	Port int
}

type host struct {
	config HostConfig
}

func newHost(config HostConfig) *host {
	return &host{config}
}

func (h *host) start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", h.config.Port))
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go h.handleRemotePeer(conn)
	}
}

func (h *host) handleRemotePeer(conn net.Conn) {
	go protocol.NewCommunicator(context.Background(), conn).Start()
}
