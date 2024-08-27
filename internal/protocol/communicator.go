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

type Communicator struct {
	ctx      context.Context
	conn     net.Conn
	codec    *encoding.Codec
	wg       sync.WaitGroup
	output   chan outputOperation
	finished bool
}

func NewCommunicator(ctx context.Context, conn net.Conn) *Communicator {
	// Initialize services.
	ctx = authentication.ContextWithService(ctx)

	return &Communicator{
		ctx:    ctx,
		conn:   conn,
		codec:  encoding.NewCodec(conn),
		output: make(chan outputOperation, 100),
	}
}

func (c *Communicator) Context() context.Context {
	return c.ctx
}

func (c *Communicator) Start() {
	c.wg.Add(2)
	go c.startOutput()
	go c.startInput()
	go c.Wait()
}

func (c *Communicator) Wait() {
	c.wg.Wait()
	c.finished = true
}

func (c *Communicator) IsFinished() bool {
	return c.finished
}

// Say sends a packet to the remote peer and returns an error if the packet
// could not be sent. If the remote peer response with a packet, it will be
// automatically handled.
func (c *Communicator) Say(packet packets.Packet) error {
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

type Question struct {
	// Packet to send to the remote peer
	Packet packets.Packet
	// ResponseTimeout is the amount of time to wait for a response from the
	// remote peer. If the response is not received within this time, an error
	// will be returned. Defaults to 10 seconds.
	ResponseTimeout time.Duration
}

// Ask sends a packet to the remote peer and waits for a response. Returns the
// response packet and an error if the packet could not be sent or if the
// response could not be received. The response packet will not be automatically
// handled. You must call the Handle method on the response packet to handle it.
// The context passed to the Handle method will be the same context that was
// passed to the Communicator when it was created.
func (c *Communicator) Ask(question Question) (packets.Packet, error) {
	errChan := make(chan error)
	respChan := make(chan interface{})
	idChan := make(chan uint64)

	c.output <- outputOperation{
		callerId: 0,
		packet:   question.Packet,
		errChan:  errChan,
		idChan:   idChan,
		respChan: respChan,
	}

	err := <-errChan
	if err != nil {
		return nil, err
	}

	id := <-idChan

	if question.ResponseTimeout == 0 {
		question.ResponseTimeout = 10 * time.Second
	}

	select {
	case <-time.After(question.ResponseTimeout):
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
func (c *Communicator) respond(resp response) error {
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

func (c *Communicator) startOutput() {
	defer c.wg.Done()

	for {
		packetHandler, ok := <-c.output
		if !ok {
			break
		}

		id, err := c.codec.Write(packetHandler.callerId, packetHandler.packet)
		if err != nil {
			packetHandler.errChan <- err
			continue
		}

		if packetHandler.respChan != nil {
			responseHandlers.Store(id, packetHandler.respChan)
		}

		packetHandler.errChan <- nil

		if packetHandler.idChan != nil {
			packetHandler.idChan <- id
		}
	}
}

func (c *Communicator) startInput() {
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
			continue
		}

		// Reset the deadline after successful read & process
		c.conn.SetReadDeadline(time.Now().Add(10 * time.Second))

		// Process the ProtocolPacket based on the ProtocolOperation
		go func() {
			if packet.RespondingTo != 0 {
				if respChan, ok := responseHandlers.Load(packet.RespondingTo); ok {
					respChan.(chan interface{}) <- packet.Data
					responseHandlers.Delete(packet.RespondingTo)
					return
				}
			}

			resp, err := packet.Data.(packets.Packet).Handle(c.ctx)
			if err != nil {
				fmt.Printf("Error executing packet: %s\n", err)
			}

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
