package udp

import (
	"context"
	"fmt"
	"net"
)

func echoServer(ctx context.Context, address string) (net.Addr, error) {
	s, err := net.ListenPacket("udp", address)
	if err != nil {
		return nil, fmt.Errorf("can not create listener: %v", err)
	}

	go func(){
		// block execution and waiting for cancel sign, cancel context inside a child goroutine, so it cancels for the parent as well
		go func() {
			<- ctx.Done()
			_ = s.Close()
		}()

		buffer := make([]byte, 1024)

		for {
			// read the message
			n, clientAdress, er := s.ReadFrom(buffer)
			if er != nil {
				return
			}

			// echo back the message
			_, err := s.WriteTo(buffer[:n], clientAdress)
			if err != nil {
				return
			}
		}
	}()

	return s.LocalAddr(), nil
}
