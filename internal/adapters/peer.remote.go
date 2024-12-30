package adapters

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/nathan-fiscaletti/coattail-go/internal/keys"
	"github.com/nathan-fiscaletti/coattail-go/internal/logging"
	"github.com/nathan-fiscaletti/coattail-go/internal/packets"
	"github.com/nathan-fiscaletti/coattail-go/internal/services/authentication"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailmodels"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
	"github.com/samber/lo"
)

/* ====== Type ====== */

var (
	ErrAccessDenied = errors.New("access denied")
)

type RemotePeerAdapter struct {
	details coattailtypes.PeerDetails

	handler *packets.Handler
}

func newRemotePeerAdapter(details coattailtypes.PeerDetails) *RemotePeerAdapter {
	return &RemotePeerAdapter{
		details: details,
	}
}

func (i *RemotePeerAdapter) getHandler(ctx context.Context) (*packets.Handler, error) {
	if i.handler == nil || !i.handler.IsConnected() {
		conn, err := net.Dial("tcp", i.details.Address)
		if err != nil {
			return nil, err
		}

		tlsConfig := &tls.Config{
			// InsecureSkipVerify skips certificate validation (not recommended in production)
			// Set this to true only for testing or if you're using a self-signed certificate
			InsecureSkipVerify: true,
		}

		tlsConn := tls.Client(conn, tlsConfig)

		// Perform the TLS handshake
		err = tlsConn.Handshake()
		if err != nil {
			return nil, fmt.Errorf("failed to perform TLS handshake: %w", err)
		}

		if logger, _ := logging.GetLogger(ctx); logger != nil {
			state := tlsConn.ConnectionState()
			logger.Printf("TLS Connection established with %s\n", state.ServerName)
		}

		ctxWithAuthKey := context.WithValue(ctx, keys.AuthenticationKey, i.details.Token)
		i.handler = packets.NewHandler(ctxWithAuthKey, tlsConn, packets.InputRoleClient)
		i.handler.HandlePackets(false)
	}

	return i.handler, nil
}

/* ====== Actions ====== */

func (i *RemotePeerAdapter) Run(ctx context.Context, name string, arg any) (any, error) {
	ph, err := i.getHandler(ctx)
	if err != nil {
		return nil, err
	}

	packet, err := ph.Request(packets.Request{
		Packet: packets.ActionPacket{
			Type:   packets.ActionPacketTypePerform,
			Action: name,
			Arg:    arg,
		},
	})
	if err != nil {
		return nil, err
	}

	respPacket, isRespPacket := packet.(packets.ActionResponsePacket)
	if !isRespPacket {
		return nil, fmt.Errorf("unexpected response packet")
	}

	return respPacket.ResponseData, nil
}

func (i *RemotePeerAdapter) Publish(ctx context.Context, name string, data any) error {
	ph, err := i.getHandler(ctx)
	if err != nil {
		return err
	}

	// Run as a request to block until the publish is complete
	_, err = ph.Request(packets.Request{
		Packet: packets.ActionPacket{
			Type:   packets.ActionPacketTypePublish,
			Action: name,
			Arg:    data,
		},
	})

	return err
}

func (i *RemotePeerAdapter) RunAndPublish(ctx context.Context, name string, arg any) error {
	ph, err := i.getHandler(ctx)
	if err != nil {
		return err
	}

	err = ph.Send(packets.ActionPacket{
		Type:   packets.ActionPacketTypePerformAndPublish,
		Action: name,
		Arg:    arg,
	})
	if err != nil {
		return err
	}

	return nil
}

func (i *RemotePeerAdapter) ListActions(ctx context.Context) ([]string, error) {
	ph, err := i.getHandler(ctx)
	if err != nil {
		return nil, err
	}

	packet, err := ph.Request(packets.Request{
		Packet: packets.ListUnitsPacket{
			Type: coattailtypes.UnitTypeAction,
		},
	})
	if err != nil {
		return nil, err
	}

	respPacket, isRespPacket := packet.(packets.ListUnitsResponsePacket)
	if !isRespPacket {
		return nil, fmt.Errorf("unexpected response packet")
	}

	return respPacket.Values, nil
}

func (i *RemotePeerAdapter) HasAction(ctx context.Context, name string) (bool, error) {
	actions, err := i.ListActions(ctx)
	if err != nil {
		return false, err
	}

	return lo.Contains(actions, name), nil
}

func (i *RemotePeerAdapter) RegisterAction(ctx context.Context, unit coattailtypes.Unit) error {
	return fmt.Errorf("cannot add action to remote peer")
}

/* ====== Receivers ====== */

func (i *RemotePeerAdapter) ListReceivers(ctx context.Context) ([]string, error) {
	ph, err := i.getHandler(ctx)
	if err != nil {
		return nil, err
	}

	packet, err := ph.Request(packets.Request{
		Packet: packets.ListUnitsPacket{
			Type: coattailtypes.UnitTypeReceiver,
		},
	})
	if err != nil {
		return nil, err
	}

	respPacket, isRespPacket := packet.(packets.ListUnitsResponsePacket)
	if !isRespPacket {
		return nil, fmt.Errorf("unexpected response packet")
	}

	return respPacket.Values, nil
}

func (i *RemotePeerAdapter) HasReceiver(ctx context.Context, name string) (bool, error) {
	receivers, err := i.ListReceivers(ctx)
	if err != nil {
		return false, err
	}

	return lo.Contains(receivers, name), nil
}

func (i *RemotePeerAdapter) RegisterReceiver(ctx context.Context, unit coattailtypes.Unit) error {
	return fmt.Errorf("cannot add receiver to remote peer")
}

func (i *RemotePeerAdapter) Notify(ctx context.Context, name string, arg any) error {
	ph, err := i.getHandler(ctx)
	if err != nil {
		return err
	}

	err = ph.Send(packets.NotifyPacket{
		Receiver: name,
		Data:     arg,
	})

	return err
}

/* ====== Peers ====== */

func (i *RemotePeerAdapter) GetPeer(ctx context.Context, id string) (*coattailtypes.Peer, error) {
	return nil, ErrAccessDenied
}

func (i *RemotePeerAdapter) GetPeerBy(ctx context.Context, predicate func(coattailtypes.PeerDetails) bool) (*coattailtypes.Peer, error) {
	return nil, ErrAccessDenied
}

func (i *RemotePeerAdapter) HasPeer(ctx context.Context, id string) (bool, error) {
	return false, ErrAccessDenied
}

func (i *RemotePeerAdapter) ListPeers(ctx context.Context) ([]*coattailtypes.Peer, error) {
	return nil, ErrAccessDenied
}

func (i *RemotePeerAdapter) Subscribe(ctx context.Context, sub coattailmodels.Subscription) error {
	ph, err := i.getHandler(ctx)
	if err != nil {
		return err
	}

	// Should use Request here to block until the subscription is complete
	_, err = ph.Request(packets.Request{
		Packet: packets.SubscribePacket{
			Address:  sub.Address,
			Action:   sub.Action,
			Receiver: sub.Receiver,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

/* ====== Credentials ====== */

func (i *RemotePeerAdapter) IssueToken(ctx context.Context, claims authentication.Claims) (*authentication.Token, error) {
	return nil, ErrAccessDenied
}

/* ====== Logger ====== */

func (i *RemotePeerAdapter) Logger(ctx context.Context) (*log.Logger, error) {
	return nil, ErrAccessDenied
}
