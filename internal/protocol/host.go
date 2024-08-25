package protocol

import "net"

type Host struct {
	host net.Listener
}

func (h *Host) Start() error {
	listener, err := net.Listen("tcp", ":5244")
	if err != nil {
		return err
	}

	h.host = listener

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go h.handleRemotePeer(conn)
	}
}

func (h *Host) handleRemotePeer(conn net.Conn) {
	go NewCommunicator(conn).Start()
}
