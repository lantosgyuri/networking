package tftp

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Ack uint16

func (a *Ack) MarshalBinary() ([]byte, error) {
	capacity := 2 + 2 // the ack is only an OpCode and a BlockNumber, so 4 bytes

	buffer := new(bytes.Buffer)
	buffer.Grow(capacity)

	err := binary.Write(buffer, binary.BigEndian, OpAck) // write the operation ode
	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, a) // write the block number
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (a *Ack) UnmarshalBinary(p []byte) error {
	var code OperationCode

	r := bytes.NewReader(p)

	err := binary.Read(r, binary.BigEndian, &code)
	if err != nil || code != OpAck {
		return errors.New("invalid ack")
	}

	return binary.Read(r, binary.BigEndian, a)
}
