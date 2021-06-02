package tftp

import (
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
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		return err
	}

	defer func(){
	_ = conn.Close()
	}()

	fmt.Printf("listening on: %s \n", conn.LocalAddr())

	return s.Serve(conn)
}

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

		_, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			return err
		}

		err = readRequest.UnmarshalBinary(buffer)
		if err != nil {
			fmt.Printf("bad request: %v \n", err)
			continue
		}

		go s.handle(addr.String(), readRequest)
	}
}