package protocol

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type packetResponseHandler struct {
	packet  PacketData
	errChan chan error
}

type Communicator struct {
	conn net.Conn

	wg     sync.WaitGroup
	output chan packetResponseHandler

	finished bool
}

func NewCommunicator(conn net.Conn) *Communicator {
	return &Communicator{
		conn:   conn,
		output: make(chan packetResponseHandler, 100),
	}
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

func (c *Communicator) WritePacket(packet PacketData) error {
	errChan := make(chan error)
	c.output <- packetResponseHandler{
		packet:  packet,
		errChan: errChan,
	}

	err := <-errChan
	if err != nil {
		return err
	}

	return nil
}

func (c *Communicator) startOutput() {
	defer c.wg.Done()

	encoder := gob.NewEncoder(c.conn)

	for packetHandler := range c.output {
		err := encoder.Encode(packet{
			Type: packetHandler.packet.Type(),
			Data: packetHandler.packet,
		})
		if err != nil {
			packetHandler.errChan <- err
			continue
		}

		packetHandler.errChan <- nil
	}
}

func (c *Communicator) startInput() {
	defer c.wg.Done()

	// Set the initial read deadline to 10 seconds
	c.conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	decoder := gob.NewDecoder(c.conn)

	for {
		// Read and decode the incoming ProtocolPacket
		packet, err := nextPacket(decoder)

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
			err := packet.Execute(c)
			if err != nil {
				fmt.Printf("Error executing packet: %s\n", err)
			}
		}()
	}
}
