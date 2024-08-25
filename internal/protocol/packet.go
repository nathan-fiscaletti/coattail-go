package protocol

import "encoding/gob"

type OperationType int

type PacketData interface {
	Execute(communicator *Communicator) error
	Type() OperationType
}

type packet struct {
	Type OperationType
	Data interface{}
}

func nextPacket(decoder *gob.Decoder) (PacketData, error) {
	var p packet
	err := decoder.Decode(&p)
	if err != nil {
		return nil, err
	}

	return p.Data.(PacketData), nil
}
