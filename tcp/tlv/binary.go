package tlv

import (
	"encoding/binary"
	"errors"
	"io"
)

type Binary []byte

func (b *Binary) Bytes() []byte {
	return *b
}

func (b *Binary) String() string {
	return string(*b)
}

func (b *Binary) WriteTo (w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, BinaryType) // the type of the message
	if err != nil {
		return 0, err
	}

	var n int64 = 1

	err = binary.Write(w, binary.BigEndian, uint32(len(*b))) // the length of the message
	if err != nil {
		return n, err
	}

	n += 4

	o, err := w.Write(*b) // the message

	return n + int64(o), err
}

func (b *Binary) ReadFrom(r io.Reader) (int64, error) {
	var typ uint8 // 1 byte

	err := binary.Read(r, binary.BigEndian, &typ) // read type
	if err != nil {
		return 0, nil
	}

	var n int64 = 1

	if typ != BinaryType {
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

	*b = make([]byte, size)
	o, err := r.Read(*b)

	return n + int64(o), err
}