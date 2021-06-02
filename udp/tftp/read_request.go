package tftp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
)

type ReadReq struct {
	FileName string
	Mode string
}

func (r *ReadReq) MarshalBinary() ([]byte, error) {
	mode := "octet"
	if r.Mode != "" {
		mode = r.Mode
	}

	capacity := 2 + 2 + len(r.FileName) + 1 + len(r.Mode) + 1

	b := new(bytes.Buffer)
	b.Grow(capacity)

	// fill buffer
	err := binary.Write(b, binary.BigEndian, OpRRQ)
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(r.FileName)
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(0)
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(mode)
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(0)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (r *ReadReq) UnmarshalBinary(p []byte) error {
	buffer := bytes.NewBuffer(p)

	var code OperationCode

	err := binary.Read(buffer, binary.BigEndian, &code) // get operation code
	if err != nil {
		return err
	}

	if code != OpRRQ {
		return errors.New("request should be a read request")
	}

	r.FileName, err = buffer.ReadString(0) // read file name, read until the 0 byte
	if err != nil {
		return errors.New("invalid read request")
	}

	r.FileName = strings.TrimRight(r.FileName, "\x00") // remove the 0 byte
	if len(r.FileName) == 0 {
		return errors.New("no file name was provided")
	}

	r.Mode, err = buffer.ReadString(0)
	if err != nil {
		return errors.New("invalid read request")
	}

	r.Mode = strings.TrimRight(r.Mode, "\x00")
	if len(r.Mode) == 0 {
		return errors.New("no read mode was provided")
	}

	if strings.ToLower(r.Mode) != "octet" {
		return errors.New("no read should be octet")
	}

	return nil
}
