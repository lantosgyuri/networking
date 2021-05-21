package ping

import (
	"context"
	"fmt"
	"io"
	"time"
)

func RunPinger() {
	ctx, cancel := context.WithCancel(context.Background())
	reader, write := io.Pipe()
	resetInterval := make(chan time.Duration, 1)
	resetInterval <- time.Second
	done := make(chan struct{})

	go func() {
		doPing(ctx, write, resetInterval)
		close(done)
	}()

	receivePing := func(d time.Duration, r io.Reader) {
		if d >= 0 {
			resetInterval <- d
		}

		buf := make([]byte, 1024)
		n, err := r.Read(buf)
		if err != nil {
			fmt.Printf("can not read buffer: %v", err)
		}

		fmt.Printf(" recieved is %q", buf[:n])
	}

	for _, v := range []int{0,2,3,14000,0} {
		receivePing(time.Duration(v) * time.Millisecond, reader)
	}

	cancel()
	<-done
}


