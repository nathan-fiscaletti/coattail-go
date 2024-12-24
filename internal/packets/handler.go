package packets

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/nathan-fiscaletti/coattail-go/internal/keys"
	"github.com/nathan-fiscaletti/coattail-go/internal/logging"
	"github.com/nathan-fiscaletti/coattail-go/internal/services/permission"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

// MaxBufferedOperations is the maximum number of operations that can be buffered
// before the PacketHandler will block. If this value is reached, the PacketHandler
// will block until an operation is completed.
const MaxBufferedOperations = 100

type HandlerInputRole int

const (
	InputRoleServer HandlerInputRole = iota
	InputRoleClient
)

// Handler is a handler for incoming and outgoing packets on a connection.
type Handler struct {
	ctx                 context.Context
	inputRole           HandlerInputRole
	conn                net.Conn
	authenticated       bool
	permissions         permission.Permissions
	authenticationError string
	codec               *StreamCodec
	wg                  sync.WaitGroup
	authWg              sync.WaitGroup
	output              chan outputOperation
	connected           bool
}

// NewHandler creates a new PacketHandler with the provided context and
// connection. The PacketHandler will handle incoming and outgoing packets on
// the connection. The context will be used to pass services to the packet
// handlers.
func NewHandler(ctx context.Context, conn net.Conn, inputRole HandlerInputRole) *Handler {
	var ctxWithLogger context.Context = ctx
	if logger, _ := logging.GetLogger(ctx); logger != nil {
		_, port, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err == nil {
			connLogger := log.New(os.Stdout, logger.Prefix()+"[p"+port+"] ", log.LstdFlags)
			ctxWithLogger = context.WithValue(ctxWithLogger, keys.LoggerKey, connLogger)
		}
	}

	if logger, err := logging.GetLogger(ctx); err == nil {
		var role string = "server"
		if inputRole == InputRoleClient {
			role = "client"
		}
		logger.Printf("created connection handler: %s, role: %s\n", conn.RemoteAddr().String(), role)
	}

	return &Handler{
		ctx:           ctxWithLogger,
		inputRole:     inputRole,
		conn:          conn,
		authenticated: inputRole == InputRoleClient,
		codec:         NewStreamCodec(conn),
	}
}

// Context returns the context that was passed to the PacketHandler when it was
// created.
func (c *Handler) Context() context.Context {
	return c.ctx
}

// HandlePackets starts handling incoming and outgoing packets on the connection.
// This function will block until the connection is closed.
func (c *Handler) HandlePackets(logPackets bool) {
	if c.connected {
		panic("attempted to start handling packets on an already connected PacketHandler")
	}

	c.connected = true
	c.wg = sync.WaitGroup{}
	c.output = make(chan outputOperation, MaxBufferedOperations)
	c.wg.Add(2)
	go c.startOutput(logPackets)
	go c.startInput(logPackets)

	c.authWg.Add(1)
	if c.inputRole == InputRoleClient {
		go c.startAuthentication()
	}

	go func() {
		c.wg.Wait()
		c.connected = false
	}()
}

// IsConnected returns true if the PacketHandler is currently connected to a
// remote peer.
func (c *Handler) IsConnected() bool {
	return c.connected
}

// Send sends a packet to the remote peer and returns an error if the packet
// could not be sent. If the remote peer response with a packet, it will be
// automatically handled.
func (c *Handler) Send(packet coattailtypes.Packet) error {
	// Create a new error channel used to return the result of the operation
	errChan := make(chan error)

	// Send the packet to the remote peer
	c.output <- outputOperation{
		callerId: 0,
		packet:   packet,
		errChan:  errChan,
	}

	// Wait for the result of the operation
	return <-errChan
}

// Request is a request to send a packet to the remote peer and wait for a
// response. The response packet will not be automatically handled. You must
// call the Handle method on the response packet to handle it. The context
// passed to the Handle method will be the same context that was passed to the
// PacketHandler when it was created.
type Request struct {
	// Packet to send to the remote peer
	Packet coattailtypes.Packet
	// ResponseTimeout is the amount of time to wait for a response from the
	// remote peer. If the response is not received within this time, an error
	// will be returned. Defaults to 10 seconds.
	ResponseTimeout time.Duration
}

// Request sends a packet to the remote peer and waits for a response. Returns the
// response packet and an error if the packet could not be sent or if the
// response could not be received. The response packet will not be automatically
// handled. You must call the Handle method on the response packet to handle it.
// The context passed to the Handle method will be the same context that was
// passed to the PacketHandler when it was created.
func (c *Handler) Request(request Request) (coattailtypes.Packet, error) {
	errChan := make(chan error)
	respChan := make(chan any)
	idChan := make(chan uint64)

	c.output <- outputOperation{
		callerId: 0,
		packet:   request.Packet,
		errChan:  errChan,
		idChan:   idChan,
		respChan: respChan,
	}

	id := <-idChan

	if request.ResponseTimeout == 0 {
		request.ResponseTimeout = 10 * time.Second
	}

	select {
	case <-time.After(request.ResponseTimeout):
		responseHandlers.Delete(id)
		errorHandlers.Delete(id)
		packetName := reflect.TypeOf(request.Packet).Name()
		return nil, fmt.Errorf("timeout waiting for response for packet %v %v", packetName, id)
	case resp := <-respChan:
		if p, ok := resp.(AuthenticationInvalidPacket); ok {
			return nil, errors.New(p.Error)
		}
		return resp.(coattailtypes.Packet), nil
	case err := <-errChan:
		return nil, err
	}
}

var responseHandlers = sync.Map{}
var errorHandlers = sync.Map{}

type outputOperation struct {
	callerId uint64
	packet   coattailtypes.Packet
	idChan   chan uint64
	errChan  chan error
	respChan chan any
}

type response struct {
	CallerID uint64
	Packet   coattailtypes.Packet
}

// respond sends a packet to the remote peer in response to a packet that was
// received from the remote peer. Returns an error if the packet could not be
// sent.
func (c *Handler) respond(resp response) error {
	errChan := make(chan error)
	idChan := make(chan uint64)

	c.output <- outputOperation{
		callerId: resp.CallerID,
		packet:   resp.Packet,
		errChan:  errChan,
		idChan:   idChan,
	}

	id := <-idChan
	errorHandlers.Store(id, errChan)

	return <-errChan
}

func (c *Handler) startAuthentication() {
	handleResponseErr := func(err error) {
		if logger, _ := logging.GetLogger(c.Context()); logger != nil {
			logger.Printf("%s", err)
		}
	}

	resp, err := c.Request(Request{
		Packet: AuthenticationPacket{
			Token: c.ctx.Value(keys.AuthenticationKey).(string),
		},
	})
	if err != nil {
		handleResponseErr(err)
		return
	}

	defer c.authWg.Done()

	respPacket, isRespPacket := resp.(AuthenticationResponsePacket)
	if !isRespPacket {
		handleResponseErr(fmt.Errorf("unexpected response packet of type %v", resp))
		return
	}

	if !respPacket.Authenticated {
		handleResponseErr(fmt.Errorf("authentication failed: %s", respPacket.Error))
		return
	}

	c.authenticated = respPacket.Authenticated
	c.permissions = permission.GetPermissions(respPacket.Permitted)
	c.ctx = permission.ContextWithPermissions(c.ctx, c.permissions)
}

func (c *Handler) startOutput(logPackets bool) {
	defer c.wg.Done()

	for {
		operation, ok := <-c.output
		if !ok {
			break
		}

		// check if output is authentication response
		if authResPacket, isAuthRespPacket := operation.packet.(AuthenticationResponsePacket); isAuthRespPacket {
			c.authenticated = authResPacket.Authenticated
			c.authenticationError = authResPacket.Error
			c.authWg.Done()
		}

		id, err := c.codec.Write(operation.callerId, operation.packet)
		if err != nil {
			operation.errChan <- err
			continue
		}

		if logPackets {
			if logger, _ := logging.GetLogger(c.Context()); logger != nil {
				logger.Printf("(->%s) %T[%d,r%d]{%v}\n", c.conn.RemoteAddr().String(), operation.packet, id, operation.callerId, operation.packet.(any))
			}
		}

		if operation.respChan != nil {
			responseHandlers.Store(id, operation.respChan)
		}

		if operation.errChan != nil {
			errorHandlers.Store(id, operation.errChan)
		}

		if operation.idChan != nil {
			operation.idChan <- id
		}
	}
}

func (c *Handler) startInput(logPackets bool) {
	defer c.wg.Done()
	defer close(c.output)

	// Set the initial read deadline to 10 seconds
	c.conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	for {
		// Read and decode the incoming ProtocolPacket
		packet, err := c.codec.Read()

		// Handle timeout error or EOF
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				if logger, _ := logging.GetLogger(c.Context()); logger != nil {
					logger.Println("Connection timed out due to inactivity.")
				}
				return
			}
			if err == io.EOF || strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host") {
				if logger, _ := logging.GetLogger(c.Context()); logger != nil {
					logger.Println("Connection closed by peer.")
				}
				return
			}

			if logger, _ := logging.GetLogger(c.Context()); logger != nil {
				logger.Printf("Error reading packet: %s\n", err)
			}

			// If we failed to decode a packet, try again with the next packet
			// TODO: Implement rate limiting.
			continue
		}

		// Reset the deadline after successful read & process
		c.conn.SetReadDeadline(time.Now().Add(10 * time.Second))

		// print the packet
		if logPackets {
			if logger, _ := logging.GetLogger(c.Context()); logger != nil {
				logger.Printf("(<-%s) %T[%d,r%d]{%v}\n", c.conn.RemoteAddr().String(), packet.Data, packet.ID, packet.RespondingTo, packet.Data)
			}
		}

		// Process the Packet in a new goroutine
		go func() {
			// Make sure that either the connection is authenticated, or that
			// the packet is an authentication packet.
			if c.inputRole == InputRoleServer {
				if !c.authenticated { // Should only be run on the host since authenticated is always true on the clients
					if _, isAuthPacket := packet.Data.(AuthenticationPacket); !isAuthPacket {
						c.authWg.Wait()
						if !c.authenticated {
							logger, _ := logging.GetLogger(c.Context())
							packetName := reflect.TypeOf(packet.Data).Name()
							logger.Printf("Authentication failed for packet %v (responding to: %v)\n", packetName, packet.RespondingTo)
							err = c.respond(response{
								CallerID: packet.ID,
								Packet: AuthenticationInvalidPacket{
									Error: fmt.Sprintf("authentication failed: %s", c.authenticationError),
								},
							})
							if err != nil {
								if logger, _ := logging.GetLogger(c.Context()); logger != nil {
									logger.Printf("Error writing response packet: %s\n", err)
								}
							}
							return
						}
					}
				}
			}

			// Should only have an impact on the client since the client doesn't send this packet type.
			if c.inputRole == InputRoleClient {
				if authInvalidPacket, isAuthInvalidPacket := packet.Data.(AuthenticationInvalidPacket); isAuthInvalidPacket {
					if _, ok := responseHandlers.Load(packet.RespondingTo); !ok {
						if errChan, ok := errorHandlers.Load(packet.RespondingTo); ok {
							errChan.(chan error) <- errors.New(authInvalidPacket.Error)

							// Clean up the response handlers if any
							if _, hasResponseHandler := responseHandlers.Load(packet.RespondingTo); hasResponseHandler {
								responseHandlers.Delete(packet.RespondingTo)
							}

							// Delete the error handler after it has been used
							errorHandlers.Delete(packet.RespondingTo)
						}
					}
				}
			}

			// If this is a response packet, load the response handler for the
			// original packet (if any) and send it the response.
			if packet.RespondingTo != 0 {
				if respChan, ok := responseHandlers.Load(packet.RespondingTo); ok {
					respChan.(chan any) <- packet.Data
					// Delete the response handler after it has been used
					responseHandlers.Delete(packet.RespondingTo)
					return
				}
			}

			// Handle the packet.
			resp, err := packet.Data.(coattailtypes.Packet).Handle(
				context.WithValue(c.ctx, keys.ConnectionKey, c.conn),
			)
			if err != nil {
				if logger, _ := logging.GetLogger(c.Context()); logger != nil {
					logger.Printf("Error executing packet: %s\n", err)
				}
			}

			// If the packet handler returned a response, send it back to the
			// remote peer.
			if resp != nil {
				err = c.respond(response{
					CallerID: packet.ID,
					Packet:   resp,
				})
				if err != nil {
					if logger, _ := logging.GetLogger(c.Context()); logger != nil {
						logger.Printf("Error writing response packet: %s\n", err)
					}
				}
			}
		}()
	}
}
