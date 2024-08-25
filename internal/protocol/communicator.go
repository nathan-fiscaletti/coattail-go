package protocol

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/nathan-fiscaletti/coattail-go/internal/protocol/packets"
)

type packetHandler struct {
	packet  packets.Packet
	errChan chan error
}

type Communicator struct {
	conn net.Conn

	wg     sync.WaitGroup
	output chan packetHandler

	finished bool
}

func NewCommunicator(conn net.Conn) *Communicator {
	return &Communicator{
		conn:   conn,
		output: make(chan packetHandler, 100),
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

func (c *Communicator) WritePacket(packet packets.Packet) error {
	errChan := make(chan error)
	c.output <- packetHandler{
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

	for packet := range c.output {
		err := encoder.Encode(packet.packet)
		if err != nil {
			packet.errChan <- err
			continue
		}

		packet.errChan <- nil
	}
}

func (c *Communicator) startInput() {
	defer c.wg.Done()

	// Set the initial read deadline to 10 seconds
	c.conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	decoder := gob.NewDecoder(c.conn)

	for {
		// Read and decode the incoming ProtocolPacket
		var packet packets.Packet
		err := decoder.Decode(&packet)

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

			// If we faild to decode a packet, try again with the next packet
			continue
		}

		// TODO: Consider changing the protocol to function by starting a new
		// TODO: goroutine for each packet received to avoid blocking the connection
		// TODO: but consider that this will not allow for packets that require
		// TODO: a specific order to be processed in the correct order.

		// Process the ProtocolPacket based on the ProtocolOperation
		switch packet.Type {
		case packets.PacketTypeHello:
			// Handle the RunAction operation
			helloPacket := packet.Data.(packets.HelloPacketData)
			fmt.Printf("Received hello packet from %s\n", helloPacket.Message)

			err := c.WritePacket(packets.NewGoodbyePacket(packets.GoodbyePacketData{
				Message: "Goodbye, I am the second functional packet!",
			}))
			if err != nil {
				fmt.Printf("Error sending hello packet: %s\n", err)
			}
		case packets.PacketTypeGoodbye:
			goodByePacket := packet.Data.(packets.GoodbyePacketData)
			fmt.Printf("Received goodbye packet from %s\n", goodByePacket.Message)
		default:
			// Handle other operations or log unknown operations
			// ...
		}

		// Reset the deadline after successful read & process
		c.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	}
}
