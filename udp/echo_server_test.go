package udp

import (
	"bytes"
	"context"
	"net"
	"testing"
)

func TestEchoServer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server, err := echoServer(ctx, "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	client, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _= client.Close() }()

	message := []byte("Test")
	_, err = client.WriteTo(message, server)
	if err != nil {
		t.Fatal(err)
	}

	buffer := make([]byte, 1024)
	n, addr, err := client.ReadFrom(buffer)
	if err != nil {
		t.Fatal(err)
	}

	if addr.String() != server.String() {
		t.Fatalf("the sender adress: %s, should be the server adress: %s", addr.String(), server.String())
	}

	if !bytes.Equal(message, buffer[:n]) {
		t.Errorf("response :%q should be: %q", message, buffer[:n])
	}
}