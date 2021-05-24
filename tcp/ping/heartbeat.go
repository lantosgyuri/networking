package ping

import (
	"context"
	"fmt"
	"net"
	"time"
)

// The connection deadline extended with 5 second with a heartbeat

func Heartbeat() {
	done := make(chan struct{})
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		fmt.Printf("can not listen: %v", err)
	}

	go func(){
		defer func() {
			close(done)
		}()

		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("can not accept connection: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer func() {
			cancel()
			if err := conn.Close(); err != nil {
				fmt.Printf("can not close connection: %v", err)
			}
		}()

		resetInterval := make(chan time.Duration)
		resetInterval <- time.Second

		go doPing(ctx, conn, resetInterval)

		err = conn.SetDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			fmt.Printf("can not set deadline: %v", err)
			return
		}

		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Printf("can not read buffer: %v", err)
			}

			fmt.Printf("Recieved: %s", buf[:n])

			resetInterval <- 0
			err = conn.SetDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				fmt.Printf("can not set deadline: %v", err)
				return
			}
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		fmt.Printf("can not dial listener: %v", err)
	}

	defer conn.Close()

	buf := make([]byte, 1024)
	for i := 0; i < 4; i++ {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("can not read buffer: %v", err)
		}
		fmt.Printf("Recieved is: %s ", buf[:n])
	}

	_, err = conn.Write([]byte("Pong"))
	if err != nil {
		fmt.Printf("can not write to connection: %v", err)
	}

	for i := 0; i < 4; i++ {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("can not read buffer: %v", err)
		}
		fmt.Printf("Recieved is: %s ", buf[:n])
	}

	<-done

}
