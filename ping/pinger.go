package ping

import (
	"context"
	"fmt"
	"io"
	"time"
)

const defaultInterval = 30 * time.Second

func doPing(ctx context.Context, writer io.Writer, reset <-chan time.Duration ) {
	var interval time.Duration

	select {
	case <-ctx.Done():
		return
	case interval = <- reset:
	}

	if interval <= 0 {
		interval = defaultInterval
	}

	timer := time.NewTimer(interval)
	defer func() {
		if !timer.Stop() {
			<- timer.C
		}
	}()

	for {
		select {
		case <- ctx.Done():
			return
			case newInterval := <- reset:
				if !timer.Stop() {
					<-timer.C
				}
				if newInterval > 0 {
					interval = newInterval
				}
				case <-timer.C:
					if _, err := writer.Write([]byte("Ping")); err != nil {
						fmt.Println("Ping sent")
						return
					}
		}
		_ = timer.Reset(interval)
	}
}
