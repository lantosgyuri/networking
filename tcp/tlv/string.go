package tlv

import (
"encoding/binary"
"errors"
"io"
)

type String string

func (s *String) Bytes() []byte {
	return []byte(*s)
}

func (s *String) String() string {
	return string(*s)
}

func (s *String) WriteTo (w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, StringType) // the type of the message
	if err != nil {
		return 0, err
	}

	var n int64 = 1

	err = binary.Write(w, binary.BigEndian, uint32(len(*s))) // the length of the message
	if err != nil {
		return n, err
	}

	n += 4

	o, err := w.Write([]byte(*s)) // the message

	return n + int64(o), err
}

func (s *String) ReadFrom(r io.Reader) (int64, error) {
	var typ uint8 // 1 byte

	err := binary.Read(r, binary.BigEndian, &typ) // read type
	if err != nil {
		return 0, nil
	}

	var n int64 = 1

	if typ != StringType {
		return n, errors.New("invalid type")
	}

	var size uint32 // 4 bytes
	err = binary.Read(r, binary.BigEndian, &size)
	if err != nil {
		return n, err
	}

	n += 4

	if size > MayPayloadSize {
		return n, ErrMayPayloadExceeded
	}

	buffer := make([]byte, size)
	o, err := r.Read(buffer)
	if err != nil {
		return n, err
	}

	*s = String(buffer)

	return n + int64(o), err
}
