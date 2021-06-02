package tftp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
)

type Error struct {
	Error ErrorCode
	Message string
}

func (e *Error) MarshalBinary() ([]byte, error) {
	capacity := 2 + 2 + len(e.Message) + 1

	buffer := new(bytes.Buffer)
	buffer.Grow(capacity)

	err := binary.Write(buffer, binary.BigEndian, OpErr) // write the error operation code
	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, e.Error) // write the error code
	if err != nil {
		return nil, err
	}

	_, err = buffer.WriteString(e.Message) // write message
	if err != nil {
		return nil, err
	}

	err = buffer.WriteByte(0)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (e *Error) UnMarshalBinary(p []byte) error {
	buffer := bytes.NewBuffer(p)

	var code OperationCode

	err := binary.Read(buffer, binary.BigEndian, &code) // read operation code
	if err != nil || code != OpErr {
		return errors.New("invalid operation code")
	}

	err = binary.Read(buffer, binary.BigEndian, &e.Error) // read error code
	if err != nil {
		return err
	}

	e.Message, err = buffer.ReadString(0)
	e.Message = strings.TrimRight(e.Message, "\x00")

	return err
}
