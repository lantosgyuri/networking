package proxy

import (
	"net"
	"strings"
	"sync"
	"testing"
)

func TestProxy(t *testing.T) {
	var wg sync.WaitGroup

	destinationServer, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	wg.Add(1)
	defer destinationServer.Close()

	go func() {
		defer wg.Done()
		conn ,err := destinationServer.Accept()
		if err != nil {
			t.Error(err)
			return
		}

		defer conn.Close()

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			t.Error(err)
			return
		}

		if string(buffer[:n]) == "Ping" {
			_, err = conn.Write([]byte("Pong"))
			if err != nil {
				t.Error(err)
				return
			}
			return
		}

		_, err = conn.Write(buffer)
		if err != nil {
			t.Error(err)
			return
		}

	}()

	proxyServer, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	defer proxyServer.Close()

	go func() {
		proxyFrom, err := proxyServer.Accept()
		if err != nil {
			t.Error(err)
			return
		}

		defer proxyFrom.Close()

		proxyTo, err := net.Dial("tcp", destinationServer.Addr().String())
		if err != nil {
			t.Error(err)
			return
		}

		defer proxyTo.Close()

		err = proxy(proxyFrom, proxyTo)
		if err != nil {
			t.Error(err)
			return
		}
	}()

	sender, err := net.Dial("tcp", proxyServer.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	defer sender.Close()

		_, err = sender.Write([]byte("Ping"))

		buffer := make([]byte, 1024)
		_, err = sender.Read(buffer)
		if err != nil {
			t.Error(err)
			return
		}

		actualMessage := string(buffer)

		res := strings.Compare(actualMessage, "Pong")
		if res != 0{
			t.Errorf("message should be Pong, not: %v", actualMessage)
		}

		wg.Wait()
}
