package tlv

import (
	"net"
	"reflect"
	"testing"
)

func TestPayloads(t *testing.T) {
	bnr := Binary("Test Binary")
	strng := String("Test String")

	payloads := []Payload{&bnr, &strng}

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
			return
		}

		defer conn.Close()
		for _, p := range payloads {
			_, err := p.WriteTo(conn)
			if err != nil {
				t.Error(err)
				break
			}
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	defer conn.Close()

	for i := 0; i < len(payloads); i++ {
		current, err := decode(conn)
		if err != nil {
			t.Fatal(err)
		}

		if expected := payloads[i]; !reflect.DeepEqual(expected, current) {
			t.Errorf("value mismatch: %v != %v", expected, current)
		}
	}
}
