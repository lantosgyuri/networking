package tftp

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"time"
)

type Server struct {
	Payload []byte
	Retires uint8
	Timeout time.Duration
}

func (s *Server) ListenAdnServe(address string) error {
	conn, err := net.ListenPacket("udp", address) // creates the connection
	if err != nil {
		return err
	}
	defer func(){ _ = conn.Close() }()

	fmt.Printf("listening on: %s \n", conn.LocalAddr())

	return s.Serve(conn)
}

// Serve waits for incoming connection than it pass the connection to a goroutine
func (s *Server) Serve(conn net.PacketConn) error {
	if conn == nil {
		return errors.New("no connection provided")
	}

	if s.Payload == nil {
		return errors.New("no payload provided")
	}

	if s.Retires == 0 {
		s.Retires = 10
	}

	if s.Timeout == 0 {
		s.Timeout = 10 * time.Second
	}

	var readRequest ReadReq

	for {
		buffer := make([]byte, DatagramSize)

		_, addr, err := conn.ReadFrom(buffer) // read incoming data -> get client address and request
		if err != nil {
			return err
		}

		err = readRequest.UnmarshalBinary(buffer) // currently only accepting a read request from a client
		if err != nil {
			fmt.Printf("bad request: %v \n", err)
			continue
		}

		go s.handle(addr.String())
	}
}

func (s Server) handle(clientAddress string) {
	conn, err := net.Dial("udp", clientAddress) // make usage of net.Conn (no need to check the sender address)

	if err != nil {
		return
	}
	defer func() { _ = conn.Close() }()

	var (
		ackPkt  Ack
		errPkt  Error
		dataPkt = Data{Payload: bytes.NewReader(s.Payload)}
		buffer  = make([]byte, DatagramSize)
	)

NEXTPACKET:
	for n := DatagramSize; n == DatagramSize; {
		data, err := dataPkt.MarshalBinary() // start creating the data blocks
		if err != nil {
			return
		}

	RETRY:
		for i := s.Retires; i > 0; i-- {
			n, err = conn.Write(data)
			if err != nil {
				return
			}

			_ = conn.SetReadDeadline(time.Now().Add(s.Timeout)) // deadline for the Ack packet

			_, err = conn.Read(buffer)
			if err != nil {
				if nErr, ok := err.(net.Error); ok && nErr.Timeout() {
					continue RETRY
				}
				return
			}

			switch {
			case ackPkt.UnmarshalBinary(buffer) == nil: // check if it is a acknowledgement packet
				if uint16(ackPkt) == dataPkt.Block { // check if the packet is for the same block which last time was sent
					continue NEXTPACKET
				}
			case errPkt.UnMarshalBinary(buffer) == nil: // check for error packet
				return // in case of error just return -> so retry
			}
		}
	}
}