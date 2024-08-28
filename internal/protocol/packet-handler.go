package protocol

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/nathan-fiscaletti/coattail-go/internal/protocol/encoding"
	"github.com/nathan-fiscaletti/coattail-go/internal/protocol/packets"
	"github.com/nathan-fiscaletti/coattail-go/internal/services/authentication"
)

type PacketHandler struct {
	ctx       context.Context
	conn      net.Conn
	codec     *encoding.StreamCodec
	wg        sync.WaitGroup
	output    chan outputOperation
	connected bool
}

func NewPacketHandler(ctx context.Context, conn net.Conn) *PacketHandler {
	// Initialize services.
	ctx = authentication.ContextWithService(ctx)

	return &PacketHandler{
		ctx:    ctx,
		conn:   conn,
		codec:  encoding.NewStreamCodec(conn),
		output: make(chan outputOperation, 100),
	}
}

func (c *PacketHandler) Context() context.Context {
	return c.ctx
}

func (c *PacketHandler) HandlePackets() {
	c.connected = true
	c.wg.Add(2)
	go c.startOutput()
	go c.startInput()
	go func() {
		c.wg.Wait()
		c.connected = false
	}()
}

func (c *PacketHandler) IsConnected() bool {
	return c.connected
}

// Send sends a packet to the remote peer and returns an error if the packet
// could not be sent. If the remote peer response with a packet, it will be
// automatically handled.
func (c *PacketHandler) Send(packet packets.Packet) error {
	errChan := make(chan error)

	c.output <- outputOperation{
		callerId: 0,
		packet:   packet,
		errChan:  errChan,
	}

	err := <-errChan
	if err != nil {
		return err
	}

	return nil
}

type Request struct {
	// Packet to send to the remote peer
	Packet packets.Packet
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
func (c *PacketHandler) Request(request Request) (packets.Packet, error) {
	errChan := make(chan error)
	respChan := make(chan interface{})
	idChan := make(chan uint64)

	c.output <- outputOperation{
		callerId: 0,
		packet:   request.Packet,
		errChan:  errChan,
		idChan:   idChan,
		respChan: respChan,
	}

	err := <-errChan
	if err != nil {
		return nil, err
	}

	id := <-idChan

	if request.ResponseTimeout == 0 {
		request.ResponseTimeout = 10 * time.Second
	}

	select {
	case <-time.After(request.ResponseTimeout):
		responseHandlers.Delete(id)
		return nil, fmt.Errorf("timeout waiting for response")
	case resp := <-respChan:
		return resp.(packets.Packet), nil
	}
}

var responseHandlers = sync.Map{}

type outputOperation struct {
	callerId uint64
	packet   packets.Packet
	idChan   chan uint64
	errChan  chan error
	respChan chan interface{}
}

type response struct {
	CallerID uint64
	Packet   packets.Packet
}

// respond sends a packet to the remote peer in response to a packet that was
// received from the remote peer. Returns an error if the packet could not be
// sent.
func (c *PacketHandler) respond(resp response) error {
	errChan := make(chan error)

	c.output <- outputOperation{
		callerId: resp.CallerID,
		packet:   resp.Packet,
		errChan:  errChan,
	}

	err := <-errChan
	if err != nil {
		return err
	}

	return nil
}

func (c *PacketHandler) startOutput() {
	defer c.wg.Done()

	for {
		operation, ok := <-c.output
		if !ok {
			break
		}

		id, err := c.codec.Write(operation.callerId, operation.packet)
		if err != nil {
			operation.errChan <- err
			continue
		}

		if operation.respChan != nil {
			responseHandlers.Store(id, operation.respChan)
		}

		operation.errChan <- nil

		if operation.idChan != nil {
			operation.idChan <- id
		}
	}
}

func (c *PacketHandler) startInput() {
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
				fmt.Println("Connection timed out due to inactivity.")
				return
			}
			if err == io.EOF {
				fmt.Println("Connection closed by peer.")
				return
			}

			// If we failed to decode a packet, try again with the next packet
			// TODO: Implement rate limiting.
			continue
		}

		// Reset the deadline after successful read & process
		c.conn.SetReadDeadline(time.Now().Add(10 * time.Second))

		// Process the Packet in a new goroutine
		go func() {
			// If this is a response packet, load the response handler for the
			// original packet (if any) and send it the response.
			if packet.RespondingTo != 0 {
				if respChan, ok := responseHandlers.Load(packet.RespondingTo); ok {
					respChan.(chan interface{}) <- packet.Data
					// Delete the response handler after it has been used
					responseHandlers.Delete(packet.RespondingTo)
					return
				}
			}

			// Handle the packet.
			resp, err := packet.Data.(packets.Packet).Handle(c.ctx)
			if err != nil {
				fmt.Printf("Error executing packet: %s\n", err)
			}

			// If the packet handler returned a response, send it back to the
			// remote peer.
			if resp != nil {
				err = c.respond(response{
					CallerID: packet.ID,
					Packet:   resp,
				})
				if err != nil {
					fmt.Printf("Error writing response packet: %s\n", err)
				}
			}
		}()
	}
}
