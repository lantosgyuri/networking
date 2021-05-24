package tlv

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

func decode(r io.Reader) (Payload, error) {
	var typ uint8

	err := binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return nil, err
	}

	var payload Payload

	switch typ {
	case BinaryType: payload = new(Binary)
	case StringType: payload = new(String)
	default:
		return nil, errors.New("type is not known")
	}

	_, err = payload.ReadFrom(
		io.MultiReader(bytes.NewReader([]byte{typ}),r)) // concatenate the already written bytes, so the ReadFrom reads the all the bytes as expected
	if err != nil {
		return nil, err
	}

	return payload, nil
}
