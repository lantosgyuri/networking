package tftp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type Data struct {
	Block uint16
	Payload io.Reader
}

func (d *Data) MarshalBinary() ([]byte, error) {
	buffer := new(bytes.Buffer)
	buffer.Grow(DatagramSize)

	d.Block++ // increment the block identifier

	err := binary.Write(buffer, binary.BigEndian, OpData) // write data operation code
	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, d.Block) // write block number
	if err != nil {
		return nil, err
	}

	_, err = io.CopyN(buffer, d.Payload, BlockSize) // write a block from the payload
	if err != nil && err != io.EOF {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (d *Data) UnmarshalBinary(p []byte) error {
	if size := len(p); size < 4 || size > DatagramSize { // if the received data is smaller than the header(4 bytes) or bigger than max block size
		return errors.New("invalid data")
	}

	var opCode OperationCode

	err := binary.Read(bytes.NewReader(p[:2]), binary.BigEndian, &opCode) // read operation code
	if err != nil || opCode != OpData {
		return errors.New("operation code should be: OpData")
	}

	err = binary.Read(bytes.NewReader(p[2:4]), binary.BigEndian, &d.Block) // read block number
	if err != nil {
		return errors.New("invalid data")
	}

	d.Payload = bytes.NewBuffer(p[4:]) // add bytes after the first 4 byte

	return nil
}
